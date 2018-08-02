package service

import (
	"errors"
	"strconv"
	"strings"
	//"fmt"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/model/yaml"
	"time"

	"github.com/astaxie/beego/logs"
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

const (
	board = iota
	k8s
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

func UpdateServiceStatus(serviceID int64, status int) (bool, error) {
	return UpdateService(model.ServiceStatus{ID: serviceID, Status: status, Deleted: 0}, "status", "deleted")
}

func UpdateServicePublic(serviceID int64, public int) (bool, error) {
	return UpdateService(model.ServiceStatus{ID: serviceID, Public: public, Deleted: 0}, "public", "deleted")
}

func DeleteServiceByID(serviceID int64) (int64, error) {
	if serviceID == 0 {
		return 0, errors.New("no Service ID provided")
	}
	num, err := dao.DeleteService(model.ServiceStatus{ID: serviceID})
	if err != nil {
		return 0, err
	}
	return num, nil
}

func GetServiceList(name string, userID int64) ([]*model.ServiceStatusMO, error) {
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
	s := model.ServiceStatus{ID: serviceID}
	_, err := dao.DeleteService(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetServiceStatus(serviceURL string) (*model.Service, error) {
	var service model.Service
	logs.Debug("Get Service info serviceURL(service): %+s", serviceURL)
	err := k8sGet(&service, serviceURL)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func GetServiceByK8sassist(pName string, sName string) (*model.Service, error) {
	logs.Debug("Get Service info %s/%s", pName, sName)

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: kubeMasterURL(),
	})
	service, _, err := k8sclient.AppV1().Service(pName).Get(sName)

	if err != nil {
		return nil, err
	}
	return service, nil
}

func GetNodesStatus(nodesURL string) (*model.NodeList, error) {
	logs.Debug("Get Node info nodeURL (endpoint): %+s", nodesURL)

	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nodes, err := k8sclient.AppV1().Node().List()

	if err != nil {
		return nil, err
	}
	return nodes, nil
}

/*
func GetEndpointStatus(serviceUrl string) (*modelK8s.Endpoints, error) {
	var endpoint modelK8s.Endpoints
	err := k8sGet(&endpoint, serviceUrl)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}
*/

func GetService(service model.ServiceStatus, selectedFields ...string) (*model.ServiceStatus, error) {
	s, err := dao.GetService(service, selectedFields...)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetServiceByID(serviceID int64) (*model.ServiceStatus, error) {
	return GetService(model.ServiceStatus{ID: serviceID, Deleted: 0}, "id", "deleted")
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

func GetDeployConfig(deployConfigURL string) (model.Deployment, error) {
	var deployConfig model.Deployment
	err := k8sGet(&deployConfig, deployConfigURL)
	return deployConfig, err
}

func SyncServiceWithK8s(pName string) error {
	logs.Debug("Sync Service Status of namespace %s", pName)

	//obtain serviceList data of
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: kubeMasterURL(),
	})

	serviceList, err := k8sclient.AppV1().Service(pName).List()
	if err != nil {
		logs.Error("Failed to get service list %s", pName)
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
		servicequery.Source = k8s
		servicequery.CreationTime, _ = time.Parse(time.RFC3339, item.CreationTimestamp.Format(time.RFC3339))
		servicequery.UpdateTime, _ = time.Parse(time.RFC3339, item.CreationTimestamp.Format(time.RFC3339))
		_, err = dao.SyncServiceData(servicequery)
		if err != nil {
			logs.Error("Sync Service %s failed.", servicequery.Name)
		}
	}

	return nil
}

func ScaleReplica(serviceInfo *model.ServiceStatus, number int32) (bool, error) {

	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	s := k8sclient.AppV1().Scale(serviceInfo.ProjectName)

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

func GetDeployment(pName string, sName string) (*model.Deployment, error) {
	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	d := k8sclient.AppV1().Deployment(pName)

	deployment, _, err := d.Get(sName)
	if err != nil {
		logs.Info("Failed to get deployment", pName, sName)
		return nil, err
	}
	return deployment, err
}

func PatchDeployment(pName string, sName string, deploymentConfig *model.Deployment) (*model.Deployment, []byte, error) {
	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	d := k8sclient.AppV1().Deployment(pName)

	//deployment, deploymentFileInfo, err := d.Update(deploymentConfig)
	deployment, deploymentFileInfo, err := d.PatchToK8s(sName, model.StrategicMergePatchType, deploymentConfig)
	if err != nil {
		logs.Info("Failed to patch deployment", pName, deploymentConfig.Name)
		return nil, nil, err
	}
	return deployment, deploymentFileInfo, err
}

func GetK8sService(pName string, sName string) (*model.Service, error) {
	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	s := k8sclient.AppV1().Service(pName)

	k8sService, _, err := s.Get(sName)
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

func StopServiceK8s(s *model.ServiceStatus) error {
	logs.Info("stop service in cluster %s", s.Name)
	// Stop deployment
	config := k8sassist.K8sAssistConfig{}
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	d := k8sclient.AppV1().Deployment(s.ProjectName)
	err := d.Delete(s.Name)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		logs.Error("Failed to delete deployment in cluster, error:%v", err)
		return err
	}
	svc := k8sclient.AppV1().Service(s.ProjectName)
	err = svc.Delete(s.Name)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		logs.Error("Failed to delete service in cluster, error:%v", err)
		return err
	}
	return nil
}

func MarshalService(serviceConfig *model.ConfigServiceStep) *model.Service {
	if serviceConfig == nil {
		return nil
	}
	ports := make([]model.ServicePort, 0)
	for index, port := range serviceConfig.ExternalServiceList {
		ports = append(ports, model.ServicePort{
			Name:     "port" + strconv.Itoa(index),
			Port:     int32(port.NodeConfig.TargetPort),
			NodePort: int32(port.NodeConfig.NodePort),
		})
	}

	return &model.Service{
		ObjectMeta: model.ObjectMeta{Name: serviceConfig.ServiceName},
		Ports:      ports,
		Selector:   map[string]string{"app": serviceConfig.ServiceName},
		Type:       "NodePort",
	}
}

func setDeploymentNodeSelector(nodeOrNodeGroupName string) map[string]string {
	if nodeOrNodeGroupName == "" {
		return nil
	}
	nodegroup, _ := dao.GetNodeGroup(model.NodeGroup{GroupName: nodeOrNodeGroupName}, "name")
	if nodegroup != nil && nodegroup.ID != 0 {
		return map[string]string{nodeOrNodeGroupName: "true"}
	} else {
		return map[string]string{"kubernetes.io/hostname": nodeOrNodeGroupName}
	}
}

func setDeploymentContainers(containerList []model.Container, registryURI string) []model.K8sContainer {
	if containerList == nil {
		return nil
	}
	k8sContainerList := make([]model.K8sContainer, 0)
	for _, cont := range containerList {
		container := model.K8sContainer{}
		container.Name = cont.Name

		if cont.WorkingDir != "" {
			container.WorkingDir = cont.WorkingDir
		}

		if len(cont.Command) > 0 {
			container.Command = append(container.Command, "/bin/sh")
			container.Args = append(container.Args, "-c", cont.Command)
		}

		if cont.VolumeMounts.VolumeName != "" {
			container.VolumeMounts = append(container.VolumeMounts, model.VolumeMount{
				Name:      cont.VolumeMounts.VolumeName,
				MountPath: cont.VolumeMounts.ContainerPath,
			})
		}

		if len(cont.Env) > 0 {
			for _, enviroment := range cont.Env {
				if enviroment.EnvName != "" {
					container.Env = append(container.Env, model.EnvVar{
						Name:  enviroment.EnvName,
						Value: enviroment.EnvValue,
					})
				}
			}
		}

		for _, port := range cont.ContainerPort {
			container.Ports = append(container.Ports, model.ContainerPort{
				ContainerPort: int32(port),
			})
		}

		container.Image = registryURI + "/" + cont.Image.ImageName + ":" + cont.Image.ImageTag

		k8sContainerList = append(k8sContainerList, container)
	}
	return k8sContainerList
}

func setDeploymentVolumes(containerList []model.Container) []model.Volume {
	if containerList == nil {
		return nil
	}
	volumes := make([]model.Volume, 0)
	for _, cont := range containerList {
		if strings.ToLower(cont.VolumeMounts.TargetStorageService) == "hostpath" {
			volumes = append(volumes, model.Volume{
				Name: cont.VolumeMounts.VolumeName,
				VolumeSource: model.VolumeSource{
					HostPath: &model.HostPathVolumeSource{
						Path: cont.VolumeMounts.TargetPath,
					},
				},
			})
		} else if strings.ToLower(cont.VolumeMounts.TargetStorageService) == "nfs" {
			index := strings.IndexByte(cont.VolumeMounts.TargetPath, '/')
			volumes = append(volumes, model.Volume{
				Name: cont.VolumeMounts.VolumeName,
				VolumeSource: model.VolumeSource{
					NFS: &model.NFSVolumeSource{
						Server: cont.VolumeMounts.TargetPath[:index],
						Path:   cont.VolumeMounts.TargetPath[index:],
					},
				},
			})
		}
	}
	return volumes
}

func MarshalDeployment(serviceConfig *model.ConfigServiceStep, registryURI string) *model.Deployment {
	if serviceConfig == nil {
		return nil
	}
	podTemplate := model.PodTemplateSpec{
		ObjectMeta: model.ObjectMeta{
			Name:   serviceConfig.ServiceName,
			Labels: map[string]string{"app": serviceConfig.ServiceName},
		},
		Spec: model.PodSpec{
			Volumes:      setDeploymentVolumes(serviceConfig.ContainerList),
			Containers:   setDeploymentContainers(serviceConfig.ContainerList, registryURI),
			NodeSelector: setDeploymentNodeSelector(serviceConfig.NodeSelector),
		},
	}

	return &model.Deployment{
		ObjectMeta: model.ObjectMeta{
			Name:      serviceConfig.ServiceName,
			Namespace: serviceConfig.ProjectName,
		},
		Spec: model.DeploymentSpec{
			Replicas: int32(serviceConfig.Instance),
			Selector: map[string]string{"app": serviceConfig.ServiceName},
			Template: podTemplate,
		},
	}
}

func MarshalNode(nodeName, labelKey string, schedFlag bool) *model.Node {
	label := make(map[string]string)
	if labelKey != "" {
		label[labelKey] = "true"
	}
	return &model.Node{
		ObjectMeta: model.ObjectMeta{
			Name:   nodeName,
			Labels: label,
		},
		Unschedulable: schedFlag,
	}
}

func MarshalNamespace(namespace string) *model.Namespace {
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{Name: namespace},
	}
}

func GetPods() (*model.PodList, error) {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: kubeMasterURL(),
	})
	l, err := k8sclient.AppV1().Pod("").List()
	if err != nil {
		return nil, err
	}
	return l, nil
}

func UpdateDeployment(pName string, sName string, deploymentConfig *model.Deployment) (*model.Deployment, []byte, error) {
	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	d := k8sclient.AppV1().Deployment(pName)

	deployment, deploymentFileInfo, err := d.Update(deploymentConfig)
	if err != nil {
		logs.Info("Failed to update deployment", pName, deploymentConfig.Name)
		return nil, nil, err
	}
	return deployment, deploymentFileInfo, err
}
