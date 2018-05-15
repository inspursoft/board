package k8sassist

import (
	"git/inspursoft/board/src/common/model"
)

type pods struct {
}

func (p *pods) Create(pod *model.Pod) (*model.Pod, error) {
	return nil, nil
}

func (p *pods) Update(pod *model.Pod) (*model.Pod, error) {
	return nil, nil
}

func (p *pods) UpdateStatus(*model.Pod) (*model.Pod, error) {
	return nil, nil
}

func (p *pods) Delete(name string) error {
	return nil
}

func (p *pods) Get(name string) (*model.Pod, error) {
	return nil, nil
}

func (p *pods) List() (*model.PodList, error) {
	return nil, nil
}

var _ PodCli = &pods{}

// PodCli has methods to work with Pod resources in k8s-assist.
// How to:  podCli, err := k8sassist.NewPods(nameSpace)
//          _, err := podCli.Update(&pod)
type PodCli interface {
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
