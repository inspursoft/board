package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"io"
	"io/ioutil"

	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"

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

func (s *services) CreateByYaml(r io.Reader) (*model.Service, error) {
	context, err := ioutil.ReadAll(r)
	if err != nil {
		logs.Error("Read file failed, error: %v", err)
		return nil, err
	}

	var service types.Service
	err = yaml.Unmarshal(context, &service)
	if err != nil {
		logs.Error("Unmarshal service failed, error: %v", err)
		return nil, err
	}

	serviceInfo, err := s.service.Create(&service)
	if err != nil {
		logs.Error("Create service failed, error: %v", err)
		return nil, err
	}

	return types.FromK8sService(serviceInfo), err
}

// newNodes returns a Nodes
func NewServices(namespace string, service v1.ServiceInterface) *services {
	return &services{
		namespace: namespace,
		service:   service,
	}
}
