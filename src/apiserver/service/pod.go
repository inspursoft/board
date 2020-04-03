package service

import (
	"encoding/json"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
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
	defer func(){
		logs.Info("exitrun..............")	
	}
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
		logs.Info("recieve xtermmessage %s", string(p))
		xmsg := xtermMessage{}
		if err := json.Unmarshal(p, &xmsg); err != nil {
			logs.Warn("json.Unmarshal err: ", err)
			return err
		}

		switch xmsg.MsgType {
		case "input":
			{
				h.Lock()
				h.rbuf = append(h.rbuf, xmsg.Input...)
				logs.Info("recieve data event %s", string(xmsg.Input))
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
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
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
	logs.Info("invoke kubernetes %s/%s pod exec in container %s.", namespace, pod, container)
	err = k8sclient.AppV1().Pod(namespace).Exec(pod, container, cmd, ptyHandler)
	if err != nil {
		return err
	}
	// run loop to fetch data from ws client
	return ptyHandler.Run()
}
