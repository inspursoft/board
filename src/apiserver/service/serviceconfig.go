package service

import (
	"errors"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/model/yaml"

	modelK8s "k8s.io/client-go/pkg/api/v1"
)

func InitServiceConfig() (*model.ServiceConfig, error) {
	return &model.ServiceConfig{}, nil
}

func SelectProject(config *model.ServiceConfig, projectID int64) (*model.ServiceConfig, error) {
	config.Phase = "SELECT_PROJECT"
	config.ProjectID = projectID
	return config, nil
}

func ConfigureContainers(config *model.ServiceConfig, containers []yaml.Container) (*model.ServiceConfig, error) {
	config.Phase = "CONFIGURE_CONTAINERS"
	config.DeploymentYaml = yaml.Deployment{}
	config.DeploymentYaml.ContainerList = containers
	return config, nil
}

func ConfigureService(config *model.ServiceConfig, service yaml.Service, deployment yaml.Deployment) (*model.ServiceConfig, error) {
	config.Phase = "CONFIGURE_SERVICE"
	config.ServiceYaml = service
	config.DeploymentYaml = deployment
	return config, nil
}

func ConfigureTest(config *model.ServiceConfig) error {
	config.Phase = "CONFIGURE_TESTING"
	return nil
}

func Deploy(config *model.ServiceConfig) error {
	config.Phase = "CONFIGURE_DEPLOY"
	return nil
}

func CreateServiceConfig(s model.ServiceStatus) (int64, error) {
	serviceID, err := dao.AddService(s)
	if err != nil {
		return 0, err
	}
	return serviceID, err
}

func UpdateService(s model.ServiceStatus, fieldNames ...string) (bool, error) {
	if s.ID == 0 {
		return false, errors.New("no Service ID provided")
	}
	_, err := dao.UpdateService(s, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetServiceList() ([]model.ServiceStatus, error) {

	serviceList, err := dao.GetServiceData()
	if err != nil {
		return nil, err
	}
	return serviceList, err
}

func DeleteService(serviceID int64) (bool, error) {
	s := model.ServiceStatus{ID: serviceID, Deleted: 1}
	_, err := dao.UpdateService(s, "deleted")
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetServiceStatus(serviceUrl string) (modelK8s.Service, error, bool) {
	var service modelK8s.Service

	flag, err := k8sGet(&service, serviceUrl)
	if flag == false {
		return service, err, false
	}
	if err != nil {
		return service, err, true
	}

	return service, err, true
}

func GetEndpointStatus(serviceUrl string) (modelK8s.Endpoints, error, bool) {
	var endpoint modelK8s.Endpoints

	flag, err := k8sGet(&endpoint, serviceUrl)
	if flag == false {
		return endpoint, err, false
	}
	if err != nil {
		return endpoint, err, true
	}

	return endpoint, err, true
}

func GetService(service model.ServiceStatus, selectedFields ...string) (*model.ServiceStatus, error) {
	s, err := dao.GetService(service, selectedFields...)
	if err != nil {
		return nil, err
	}
	return s, nil
}
