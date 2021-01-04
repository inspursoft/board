package service

import (
	"encoding/json"
	"errors"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"
	"net/http"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

// message from web socket client
type xtermMessage struct {
	MsgType string `json:"type"`
	Input   string `json:"input"`
	Rows    uint16 `json:"rows"`
	Cols    uint16 `json:"cols"`
}

type WSStreamHandler struct {
	conn        *websocket.Conn
	resizeEvent chan model.TerminalSize

	sync.Mutex
	cond *sync.Cond
	rbuf []byte
}

// Run start a loop to fetch from ws client and store the data in byte buffer
func (h *WSStreamHandler) Run() error {
	h.conn.SetReadDeadline(time.Now().Add(600 * time.Second))
	for {
		t, p, err := h.conn.ReadMessage()
		if err != nil {
			logs.Warn("ws ReadMessage err: %+v", err)
			return err
		}
		if t == websocket.CloseMessage {
			logs.Info("recieve close type ws message")
			return nil
		}
		xmsg := xtermMessage{}
		if err := json.Unmarshal(p, &xmsg); err != nil {
			logs.Warn("unmarshal message error: %+v", err)
			continue
		}

		switch xmsg.MsgType {
		case "input":
			{
				h.Lock()
				h.rbuf = append(h.rbuf, xmsg.Input...)
				h.cond.Signal()
				h.Unlock()
			}
		case "resize":
			{
				ev := model.TerminalSize{
					Width:  xmsg.Cols,
					Height: xmsg.Rows}
				logs.Info("recieve resize event %+v", ev)
				h.resizeEvent <- ev
			}
		default:
			logs.Info("other message type %s: not input or resize.", xmsg.MsgType)
		}
	}
}

func (h *WSStreamHandler) Read(b []byte) (int, error) {
	h.Lock()
	for len(h.rbuf) == 0 {
		h.cond.Wait()
	}
	size := copy(b, h.rbuf)
	h.rbuf = h.rbuf[size:]
	h.Unlock()
	return size, nil
}

func (h *WSStreamHandler) Write(b []byte) (int, error) {
	return len(b), h.conn.WriteMessage(websocket.TextMessage, b)
}

func (h *WSStreamHandler) Next() *model.TerminalSize {
	ret := <-h.resizeEvent
	return &ret
}

func PodShell(namespace, pod, container string, w http.ResponseWriter, r *http.Request) error {
	// upgrade to websocket.
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		return err
	}
	defer conn.Close()

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	// maybe someday cmd args will be setting to `/bin/bash, cmd or powershell` base on container information.
	cmd := []string{"/bin/sh"}
	ptyHandler := &WSStreamHandler{
		conn:        conn,
		resizeEvent: make(chan model.TerminalSize),
	}
	ptyHandler.cond = sync.NewCond(ptyHandler)
	// run loop to fetch data from ws client
	go ptyHandler.Run()
	logs.Info("invoke kubernetes %s/%s pod exec in container %s.", namespace, pod, container)
	// check Container Privileged
	modelPod, err := k8sclient.AppV1().Pod(namespace).Get(pod)
	if err != nil {
		return err
	}
	if modelPod.Spec.Containers[0].SecurityContext != nil &&
		modelPod.Spec.Containers[0].SecurityContext.Privileged != nil &&
		*modelPod.Spec.Containers[0].SecurityContext.Privileged {
		err = errors.New("the container privileged is true, cannot be connected")
		ptyHandler.Write([]byte(err.Error()))
		return err
	}
	err = k8sclient.AppV1().Pod(namespace).ShellExec(pod, container, cmd, ptyHandler)
	ptyHandler.Write([]byte(err.Error()))
	return err
}

func CopyFromPod(namespace, podName, container, src, dest string) error {
	if len(src) == 0 || len(dest) == 0 {
		return errors.New("filepath can not be empty")
	}

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	cmd := []string{"tar", "cf", "-", src}
	logs.Info("Copy kubernetes pod %s/%s:%s in container %s to %s on host.", namespace, podName, src, container, dest)
	return k8sclient.AppV1().Pod(namespace).CopyFromPod(podName, container, src, dest, cmd)
}

func CopyToPod(namespace, podName, container, src, dest string) error {
	if len(src) == 0 || len(dest) == 0 {
		return errors.New("filepath can not be empty")
	}

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	// check Container Privileged
	modelPod, err := k8sclient.AppV1().Pod(namespace).Get(podName)
	if err != nil {
		return err
	}
	if modelPod.Spec.Containers[0].SecurityContext != nil &&
		modelPod.Spec.Containers[0].SecurityContext.Privileged != nil &&
		*modelPod.Spec.Containers[0].SecurityContext.Privileged {
		return errors.New("the container privileged is true, cannot be connected")
	}
	logs.Info("Copying the content of '%s' to '%s/%s/%s:%s'", src, namespace, podName, container, dest)
	return k8sclient.AppV1().Pod(namespace).CopyToPod(podName, container, src, dest)
}

func GetPodsByLabelSelector(projectName string, labelSelector *model.LabelSelector) (*model.PodList, error) {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	var opts model.ListOptions
	opts.LabelSelector = types.LabelSelectorToString(labelSelector)
	return k8sclient.AppV1().Pod(projectName).List(opts)
}
