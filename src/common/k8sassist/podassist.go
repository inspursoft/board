package k8sassist

import (
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

type pods struct {
	namespace string
	pod       v1.PodInterface
}

func (p *pods) Create(pod *model.Pod) (*model.Pod, error) {
	k8sPod := toK8sPod(pod)
	k8sPod, err := p.pod.Create(k8sPod)
	if err != nil {
		logs.Error("Create pod of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := fromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) Update(pod *model.Pod) (*model.Pod, error) {
	k8sPod := toK8sPod(pod)
	k8sPod, err := p.pod.Update(k8sPod)
	if err != nil {
		logs.Error("Update pod of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := fromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) UpdateStatus(pod *model.Pod) (*model.Pod, error) {
	k8sPod := toK8sPod(pod)
	k8sPod, err := p.pod.UpdateStatus(k8sPod)
	if err != nil {
		logs.Error("Create pod status of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := fromK8sPod(k8sPod)
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

	modelPod := fromK8sPod(pod)
	return modelPod, nil
}

func (p *pods) List() (*model.PodList, error) {
	podList, err := p.pod.List(metav1.ListOptions{})
	if err != nil {
		logs.Error("list pods failed. Err:%+v", err)
		return nil, err
	}

	modelPodList := fromK8sPodList(podList)
	return modelPodList, nil
}

var _ PodCliInterface = &pods{}

func NewPods(namespace string) PodCliInterface {
	//TODO: init the clientset.
	var client *kubernetes.Clientset
	return &pods{
		namespace: namespace,
		pod:       client.CoreV1().Pods(namespace),
	}
}

// PodCli has methods to work with Pod resources in k8s-assist.
// How to:  podCli, err := k8sassist.NewPods(nameSpace)
//          _, err := podCli.Update(&pod)
type PodCliInterface interface {
	Create(*model.Pod) (*model.Pod, error)
	Update(*model.Pod) (*model.Pod, error)
	UpdateStatus(*model.Pod) (*model.Pod, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Pod, error)
	List() (*model.PodList, error)
	//List(opts v1.ListOptions) (*v1.PodList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Pod, err error)
}
