package apps

import (
	"fmt"
	"io"

	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

type pods struct {
	k8sClient kubernetes.Interface
	cfg       *types.Config
	namespace string
	pod       v1.PodInterface
}

func (p *pods) Create(pod *model.Pod) (*model.Pod, error) {
	k8sPod := types.ToK8sPod(pod)
	k8sPod, err := p.pod.Create(k8sPod)
	if err != nil {
		logs.Error("Create pod of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) Update(pod *model.Pod) (*model.Pod, error) {
	k8sPod := types.ToK8sPod(pod)
	k8sPod, err := p.pod.Update(k8sPod)
	if err != nil {
		logs.Error("Update pod of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) UpdateStatus(pod *model.Pod) (*model.Pod, error) {
	k8sPod := types.ToK8sPod(pod)
	k8sPod, err := p.pod.UpdateStatus(k8sPod)
	if err != nil {
		logs.Error("Create pod status of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) Delete(name string) error {
	err := p.pod.Delete(name, nil)
	if err != nil {
		logs.Error("delete pod of %s/%s failed. Err:%+v", name, p.namespace, err)
	}
	return err
}

func (p *pods) Get(name string) (*model.Pod, error) {
	pod, err := p.pod.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("get pod of %s/%s failed. Err:%+v", name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(pod)
	return modelPod, nil
}

func (p *pods) List(opts model.ListOptions) (*model.PodList, error) {
	podList, err := p.pod.List(types.ToK8sListOptions(opts))
	if err != nil {
		logs.Error("list pods failed. Err:%+v", err)
		return nil, err
	}

	modelPodList := types.FromK8sPodList(podList)
	return modelPodList, nil
}

func (p *pods) GetLogs(name string, opts *model.PodLogOptions) (io.ReadCloser, error) {
	request := p.pod.GetLogs(name, types.ToK8sPodLogOptions(opts))
	if request == nil {
		err := fmt.Errorf("get pod of %s/%s logs failed, request client is null", name, p.namespace)
		logs.Error("%+v", err)
		return nil, err
	}
	return request.Stream()
}

func (p *pods) Exec(podName, containerName string, cmd []string, ptyHandler model.PtyHandler) error {
	req := p.k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(p.namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.cfg, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: types.ToK8sTerminalSizeQueue(ptyHandler),
		Tty:               true,
	})
	if err != nil {
		return err
	}

	return nil
}

func NewPods(k8sClient kubernetes.Interface, cfg *types.Config, namespace string, pod v1.PodInterface) *pods {
	return &pods{
		k8sClient: k8sClient,
		cfg:       cfg,
		namespace: namespace,
		pod:       pod,
	}
}
