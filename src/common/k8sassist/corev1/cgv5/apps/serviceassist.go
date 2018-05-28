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

func (s *services) Create(modelService *model.Service) (*model.Service, []byte, error) {
	typeService := types.ToK8sService(modelService)
	svc, err := s.service.Create(typeService)
	if err != nil {
		logs.Error("Create service failed, error: %v", err)
		return nil, nil, err
	}

	svcfileInfo, err := yaml.Marshal(svc)
	if err != nil {
		logs.Error("Marshal service failed, error: %v", err)
		return types.FromK8sService(svc), nil, err
	}
	return types.FromK8sService(svc), svcfileInfo, nil
}

func (s *services) Update(modelService *model.Service) (*model.Service, []byte, error) {
	typeService := types.ToK8sService(modelService)
	svc, err := s.service.Update(typeService)
	if err != nil {
		logs.Error("Update service failed, error: %v", err)
		return nil, nil, err
	}

	svcfileInfo, err := yaml.Marshal(svc)
	if err != nil {
		logs.Error("Marshal service info failed, error: %v", err)
		return types.FromK8sService(svc), nil, err
	}
	return types.FromK8sService(svc), svcfileInfo, nil
}

func (s *services) UpdateStatus(modelService *model.Service) (*model.Service, []byte, error) {
	typeService := types.ToK8sService(modelService)
	svc, err := s.service.UpdateStatus(typeService)
	if err != nil {
		logs.Error("Updatestatus service failed, error: %v", err)
		return nil, nil, err
	}

	svcfileInfo, err := yaml.Marshal(svc)
	if err != nil {
		logs.Error("Marshal service info failed, error: %v", err)
		return types.FromK8sService(svc), nil, err
	}
	return types.FromK8sService(svc), svcfileInfo, nil
}

func (s *services) Delete(name string) error {
	err := s.service.Delete(name, &types.DeleteOptions{})
	if err != nil {
		logs.Error("Delete service failed, error: %v", err)
		return err
	}
	return nil
}

func (s *services) Get(name string) (*model.Service, []byte, error) {
	svc, err := s.service.Get(name, types.GetOptions{})
	if err != nil {
		logs.Error("Get service failed, error: %v", err)
		return nil, nil, err
	}
	svcfileInfo, err := yaml.Marshal(svc)
	if err != nil {
		logs.Error("Marshal service info failed, error: %v", err)
		return types.FromK8sService(svc), nil, err
	}

	return types.FromK8sService(svc), svcfileInfo, nil
}

func (s *services) List() (*model.ServiceList, error) {
	svcList, err := s.service.List(types.ListOptions{})
	if err != nil {
		logs.Error("Get service list failed, error: %v", err)
		return nil, err
	}
	return types.FromK8sServiceList(svcList), nil
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

	err = types.CheckServiceConfig(s.namespace, service)
	if err != nil {
		logs.Error("Service config is error, error: %v", err)
		return nil, err
	}

	serviceInfo, err := s.service.Create(&service)
	if err != nil {
		logs.Error("Create service failed, error: %v", err)
		return nil, err
	}

	return types.FromK8sService(serviceInfo), err
}

func (s *services) CheckYaml(r io.Reader) (*model.Service, error) {
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

	err = types.CheckServiceConfig(s.namespace, service)
	if err != nil {
		logs.Error("Service config is error, error: %v", err)
		return nil, err
	}

	return types.FromK8sService(&service), nil
}

// newNodes returns a Nodes
func NewServices(namespace string, service v1.ServiceInterface) *services {
	return &services{
		namespace: namespace,
		service:   service,
	}
}
