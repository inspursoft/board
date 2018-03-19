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
	defaultOwnerName   = "admin"
	defaultPublic      = 0
	defaultComment     = "init service"
	defaultDeleted     = 0
	defaultStatus      = 1
	serviceNamespace   = "default" //TODO create namespace in project post
	scaleKind          = "Deployment"
	k8sService         = "kubernetes"
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

func CreateServiceConfig(serviceConfig model.ServiceStatus) (*model.ServiceStatus, error) {
	query := model.Project{Name: serviceConfig.ProjectName}
	project, err := GetProject(query, "name")
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project is invalid")
	}

	serviceConfig.ProjectID = project.ID
	serviceID, err := dao.AddService(serviceConfig)
	if err != nil {
		return nil, err
	}
	serviceConfig.ID = serviceID
	return &serviceConfig, err
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

func GetPaginatedServiceList(name string, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedServiceStatus, error) {
	query := model.ServiceStatus{Name: name}
	paginatedServiceStatus, err := dao.GetPaginatedServiceData(query, userID, pageIndex, pageSize, orderField, orderAsc)
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

func GetServiceStatus(serviceURL string) (*modelK8s.Service, error) {
	var service modelK8s.Service
	logs.Debug("Get Service info serviceURL(service): %+s", serviceURL)
	err := k8sGet(&service, serviceURL)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func GetNodesStatus(nodesURL string) (*modelK8s.NodeList, error) {
	var nodes modelK8s.NodeList
	logs.Debug("Get Node info nodeURL (endpoint): %+s", nodesURL)
	err := k8sGet(&nodes, nodesURL)
	if err != nil {
		return nil, err
	}
	return &nodes, nil
}

func GetEndpointStatus(serviceUrl string) (*modelK8s.Endpoints, error) {
	var endpoint modelK8s.Endpoints
	err := k8sGet(&endpoint, serviceUrl)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func GetService(service model.ServiceStatus, selectedFields ...string) (*model.ServiceStatus, error) {
	s, err := dao.GetService(service, selectedFields...)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetServiceByProject(serviceName string, projectName string) (*model.ServiceStatus, error) {
	var servicequery model.ServiceStatus
	servicequery.Name = serviceName
	servicequery.ProjectName = projectName
	service, err := GetService(servicequery, "name", "project_name")
	if err != nil {
		return nil, err
	}
	return service, nil
}

func GetDeployConfig(deployConfigURL string) (modelK8sExt.Deployment, error) {
	var deployConfig modelK8sExt.Deployment
	err := k8sGet(&deployConfig, deployConfigURL)
	return deployConfig, err
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
		queryProject := model.Project{Name: item.Namespace}
		project, err := GetProject(queryProject, "name")
		if err != nil {
			logs.Error("Failed to check project in DB %s", item.Namespace)
			return err
		}
		if project == nil {
			logs.Error("not found project in DB: %s", item.Namespace)
			continue
		}
		if item.ObjectMeta.Name == k8sService {
			continue
		}
		servicequery.Name = item.ObjectMeta.Name
		servicequery.OwnerID = int64(project.OwnerID) //owner or admin TBD
		servicequery.OwnerName = project.OwnerName
		servicequery.ProjectName = project.Name
		servicequery.ProjectID = project.ID
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
	s := apiSet.Scales(serviceInfo.ProjectName)
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

func GetScaleStatus(serviceInfo *model.ServiceStatus) (model.ScaleStatus, error) {
	var scaleStatus model.ScaleStatus
	deployment, err := GetDeployment(serviceInfo.ProjectName, serviceInfo.Name)
	if err != nil {
		logs.Debug("Failed to get deployment %s", serviceInfo.Name)
		return scaleStatus, err
	}
	scaleStatus.DesiredInstance = deployment.Status.Replicas
	scaleStatus.AvailableInstance = deployment.Status.AvailableReplicas
	return scaleStatus, nil
}
