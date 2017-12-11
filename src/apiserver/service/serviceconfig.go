package service

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/model/yaml"
	"time"

	"github.com/astaxie/beego/logs"

	"k8s.io/client-go/kubernetes"
	modelK8s "k8s.io/client-go/pkg/api/v1"
	modelK8sExt "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	//"k8s.io/client-go/pkg/api/resource"
	//"k8s.io/client-go/pkg/api/v1"
	//"k8s.io/client-go/rest"
	//apiCli "k8s.io/client-go/tools/clientcmd/api"
)

const (
	defaultProjectName = "library"
	defaultProjectID   = 1
	defaultOwnerID     = 1
	defaultOwnerName   = "anonymous"
	defaultPublic      = 0
	defaultComment     = "init service"
	defaultDeleted     = 0
	defaultStatus      = 1
	serviceNamespace   = "default" //TODO create namespace in project post
	scaleKind          = "Deployment"
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

func DeleteServiceByID(s model.ServiceStatus) (int64, error) {
	if s.ID == 0 {
		return 0, errors.New("no Service ID provided")
	}
	num, err := dao.DeleteService(s)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func GetServiceList(name string, userID int64) ([]*model.ServiceStatus, error) {
	query := model.ServiceStatus{Name: name}
	serviceList, err := dao.GetServiceData(query, userID)
	if err != nil {
		return nil, err
	}
	return serviceList, err
}

func GetPaginatedServiceList(name string, userID int64, pageIndex int, pageSize int) (*model.PaginatedServiceStatus, error) {
	query := model.ServiceStatus{Name: name}
	paginatedServiceStatus, err := dao.GetPaginatedServiceData(query, userID, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	return paginatedServiceStatus, nil
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

func GetNodesStatus(nodesUrl string) (modelK8s.NodeList, error, bool) {
	var nodes modelK8s.NodeList

	flag, err := k8sGet(&nodes, nodesUrl)
	if flag == false {
		return nodes, err, false
	}
	if err != nil {
		return nodes, err, true
	}

	return nodes, err, true
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

func SyncServiceWithK8s() error {

	serviceUrl := fmt.Sprintf("%s/api/v1/services", kubeMasterURL())
	logs.Debug("Get Service Status serviceUrl:%+s", serviceUrl)

	//obtain serviceList data
	var serviceList modelK8s.ServiceList
	_, err := GetK8sData(&serviceList, serviceUrl)
	if err != nil {
		return err
	}

	//handle the serviceList data
	var servicequery model.ServiceStatus
	for _, item := range serviceList.Items {
		servicequery.Name = item.ObjectMeta.Name
		servicequery.OwnerID = defaultOwnerID
		servicequery.OwnerName = defaultOwnerName
		servicequery.ProjectName = defaultProjectName
		servicequery.ProjectID = defaultProjectID
		servicequery.Public = defaultPublic
		servicequery.Comment = defaultComment
		servicequery.Deleted = defaultDeleted
		servicequery.Status = defaultStatus
		servicequery.CreationTime, _ = time.Parse(time.RFC3339, item.CreationTimestamp.Format(time.RFC3339))
		servicequery.UpdateTime, _ = time.Parse(time.RFC3339, item.CreationTimestamp.Format(time.RFC3339))
		_, err = dao.SyncServiceData(servicequery)
		if err != nil {
			logs.Error("Sync Service %s failed.", servicequery.Name)
		}
	}

	return nil
}

func ScaleReplica(serviceInfo model.ServiceStatus, number int32) (bool, error) {

	cli, err := K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return false, err
	}
	s := apiSet.Scales(serviceNamespace)
	scale, err := s.Get(scaleKind, serviceInfo.Name)

	if scale.Spec.Replicas != number {
		scale.Spec.Replicas = number
		_, err = s.Update(scaleKind, scale)
		if err != nil {
			logs.Info("Failed to update service replicas", scale)
			return false, err
		}
	} else {
		logs.Info("Service replicas needn't change %d", scale.Spec.Replicas)
	}
	return true, err
}

func GetSelectableServices(pname string, sName string) ([]string, error) {
	serviceList, err := dao.GetSelectableServices(pname, sName)
	if err != nil {
		return nil, err
	}
	return serviceList, err
}

func GetDeployment(pName string, sName string) (*modelK8sExt.Deployment, error) {
	cli, err := K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return nil, err
	}
	d := apiSet.Deployments(pName)
	deployment, err := d.Get(sName)
	if err != nil {
		logs.Info("Failed to get deployment", pName, sName)
		return nil, err
	}
	return deployment, err
}

func GetK8sService(pName string, sName string) (*modelK8s.Service, error) {
	cli, err := K8sCliFactory("", kubeMasterURL(), "v1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return nil, err
	}
	s := apiSet.Services(pName)
	k8sService, err := s.Get(sName)
	if err != nil {
		logs.Info("Failed to get K8s service", pName, sName)
		return nil, err
	}
	return k8sService, err
}
