package apps

import (
	"git/inspursoft/board/src/common/model"

	"io"

	"k8s.io/client-go/kubernetes/typed/core/v1"
)

type services struct {
	namespace string
	service   v1.ServiceInterface
}

func (d *services) Create(*model.Service) (*model.Service, error) {
	return nil, nil
}

func (d *services) Update(*model.Service) (*model.Service, error) {
	return nil, nil
}

func (d *services) UpdateStatus(*model.Service) (*model.Service, error) {
	return nil, nil
}

func (d *services) Delete(name string) error {
	return nil
}

func (d *services) Get(name string) (*model.Service, error) {
	return nil, nil
}

func (d *services) List() (*model.ServiceList, error) {
	return nil, nil
}

func (d *services) CreateByYaml(r io.Reader) (*model.Service, error) {

	return nil, nil
}

// newNodes returns a Nodes
func NewServices(namespace string, service v1.ServiceInterface) *services {
	return &services{
		namespace: namespace,
		service:   service,
	}
}
