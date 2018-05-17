// a temp file for building and guiding
package k8sassist

import (
	//api "k8s.io/client-go/pkg/api"
	//v1 "k8s.io/client-go/pkg/api/v1"
	//watch "k8s.io/client-go/pkg/watch"
	//rest "k8s.io/client-go/rest"
	"git/inspursoft/board/src/common/model"
)

type services struct {
	ns string
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

var _ ServiceCliInterface = &services{}

// newNodes returns a Nodes
func NewServices(namespace string) (*services, error) {
	return &services{ns: namespace}, nil
}
