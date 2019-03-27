package controller

import (
	"fmt"
	"io"
	"os"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/travis"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	dockerfileName         = "Dockerfile"
	deploymentFilename     = "deployment.yaml"
	serviceFilename        = "service.yaml"
	rollingUpdateFilename  = "rollingUpdateDeployment.yaml"
	deploymentTestFilename = "testdeployment.yaml"
	serviceTestFilename    = "testservice.yaml"

	apiheader        = "Content-Type: application/yaml"
	deploymentAPI    = "/apis/extensions/v1beta1/namespaces/"
	serviceAPI       = "/api/v1/namespaces/"
	test             = "test"
	serviceNamespace = "default" //TODO create in project post
	k8sServices      = "kubernetes"
	deploymentType   = "deployment"
	serviceType      = "service"
	startingDuration = 300 * time.Second //300 seconds
)

const (
	preparing = iota
	running
	stopped
	uncompleted
	warning
	deploying
)

type ServiceController struct {
	BaseController
}

func (p *ServiceController) generateDeploymentTravis(serviceName, deploymentURL, serviceURL string) error {
	userID := p.currentUser.ID
	var travisCommand travis.TravisCommand
	travisCommand.Script.Commands = []string{}
	items := []string{
		fmt.Sprintf("curl \"%s/jenkins-job/%d/$BUILD_NUMBER\"", boardAPIBaseURL(), userID),
	}
	if deploymentURL != "" {
		items = append(items, fmt.Sprintf("#curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/deployment.yaml %s", serviceName, deploymentURL))
	}
	if serviceURL != "" {
		items = append(items, fmt.Sprintf("#curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/service.yaml %s", serviceName, serviceURL))
	}
	travisCommand.Script.Commands = items
	return travisCommand.GenerateCustomTravis(p.repoPath)
}

func (p *ServiceController) getKey() string {
	return strconv.Itoa(int(p.currentUser.ID))
}

func (p *ServiceController) resolveServiceInfo() (s *model.ServiceStatus, err error) {
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}
	// Get the project info of this service
	s, err = service.GetServiceByID(int64(serviceID))
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	return
}

func (p *ServiceController) DeployServiceAction() {
	key := p.getKey()
	configService := NewConfigServiceStep(key)

	//Judge authority
	project := p.resolveUserPrivilegeByID(configService.ProjectID)

	var newservice model.ServiceStatus
	newservice.Name = configService.ServiceName
	newservice.ProjectID = configService.ProjectID
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.currentUser.ID
	newservice.OwnerName = p.currentUser.Username
	newservice.Public = configService.Public
	newservice.ProjectName = project.Name

	serviceInfo, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.internalError(err)
		return
	}

	deployInfo, err := service.DeployService((*model.ConfigServiceStep)(configService), registryBaseURI())
	if err != nil {
		p.parseError(err, parsePostK8sError)
		return
	}

	p.resolveRepoServicePath(project.Name, newservice.Name)
	err = service.GenerateDeployYamlFiles(deployInfo, p.repoServicePath)
	if err != nil {
		p.internalError(err)
		return
	}

	deploymentFile := filepath.Join(newservice.Name, deploymentFilename)
	serviceFile := filepath.Join(newservice.Name, serviceFilename)
	items := []string{deploymentFile, serviceFile}
	p.pushItemsToRepo(items...)

	var serviceType int
	if deployInfo.Service.Type == "NodePort" {
		serviceType = model.ServiceTypeNormalNodePort
	} else if deployInfo.Service.Type == "ClusterIP" {
		serviceType = model.ServiceTypeClusterIP
	}
	// TODO support deployment only

	updateService := model.ServiceStatus{ID: serviceInfo.ID, Status: uncompleted, Type: serviceType, ServiceYaml: string(deployInfo.ServiceFileInfo),
		DeploymentYaml: string(deployInfo.DeploymentFileInfo)}
	_, err = service.UpdateService(updateService, "status", "type", "service_yaml", "deployment_yaml")
	if err != nil {
		p.internalError(err)
		return
	}

	err = DeleteConfigServiceStep(key)
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Service with ID:%d has been deleted in cache.", serviceInfo.ID)

	configService.ServiceID = serviceInfo.ID
	p.renderJSON(configService)
}

func (p *ServiceController) DeployServiceTestAction() {
	key := p.getKey()
	configService := NewConfigServiceStep(key)
	configService.ServiceName = test + configService.ServiceName
	SetConfigServiceStep(key, configService)
	p.DeployServiceAction()
}

//
func syncK8sStatus(serviceList []*model.ServiceStatusMO) error {
	var err error
	// synchronize service status with the cluster system
	for _, serviceStatusMO := range serviceList {
		// Get serviceStatus from serviceStatusMO to adapt for updating services
		serviceStatus := &serviceStatusMO.ServiceStatus
		if (*serviceStatus).Status == stopped {
			continue
		}
		// Check the deployment status
		deployment, _, err := service.GetDeployment((*serviceStatus).ProjectName, (*serviceStatus).Name)
		if deployment == nil && serviceStatus.Name != k8sServices {
			logs.Info("Failed to get deployment", err)
			var reason = "The deployment is not established in cluster system"
			(*serviceStatus).Status = uncompleted
			// TODO create a new field in serviceStatus for reason
			(*serviceStatus).Comment = "Reason: " + reason
			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
			if err != nil {
				logs.Error("Failed to update deployment.")
				break
			}
			continue
		} else {
			if deployment.Status.Replicas > deployment.Status.AvailableReplicas {
				logs.Debug("The desired replicas number is not available",
					deployment.Status.Replicas, deployment.Status.AvailableReplicas)
				(*serviceStatus).Status = uncompleted
				reason := "The desired replicas number is not available"
				(*serviceStatus).Comment = "Reason: " + reason
				_, err = service.UpdateService(*serviceStatus, "status", "Comment")
				if err != nil {
					logs.Error("Failed to update deployment replicas.")
					break
				}
				continue
			}
		}

		// Check the service in k8s cluster status
		serviceK8s, err := service.GetK8sService((*serviceStatus).ProjectName, (*serviceStatus).Name)
		if serviceK8s == nil {
			logs.Info("Failed to get service in cluster", err)
			var reason = "The service is not established in cluster system"
			(*serviceStatus).Status = uncompleted
			(*serviceStatus).Comment = "Reason: " + reason
			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
			if err != nil {
				logs.Error("Failed to update service in cluster.")
				break
			}
			continue
		}

		if serviceStatus.Status == uncompleted {
			logs.Info("The service is restored to running")
			(*serviceStatus).Status = running
			(*serviceStatus).Comment = ""
			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
			if err != nil {
				logs.Error("Failed to update service status.")
				break
			}
			continue
		}
	}
	return err
}

//get service list
func (p *ServiceController) GetServiceListAction() {
	serviceName := p.GetString("service_name")
	pageIndex, _ := p.GetInt("page_index", 0)
	pageSize, _ := p.GetInt("page_size", 0)
	orderField := p.GetString("order_field", "creation_time")
	orderAsc, _ := p.GetInt("order_asc", 0)
	if pageIndex == 0 && pageSize == 0 {
		serviceStatus, err := service.GetServiceList(serviceName, p.currentUser.ID)
		if err != nil {
			p.internalError(err)
			return
		}
		err = syncK8sStatus(serviceStatus)
		if err != nil {
			p.internalError(err)
			return
		}
		p.renderJSON(serviceStatus)
	} else {
		paginatedServiceStatus, err := service.GetPaginatedServiceList(serviceName, p.currentUser.ID, pageIndex, pageSize, orderField, orderAsc)
		if err != nil {
			p.internalError(err)
			return
		}
		err = syncK8sStatus(paginatedServiceStatus.ServiceStatusList)
		if err != nil {
			p.internalError(err)
			return
		}
		p.renderJSON(paginatedServiceStatus)
	}
}

// API to create service config
func (p *ServiceController) CreateServiceConfigAction() {
	var reqServiceProject model.ServiceProject
	var err error

	err = p.resolveBody(&reqServiceProject)
	if err != nil {
		return
	}

	//Judge authority
	p.resolveUserPrivilegeByID(reqServiceProject.ProjectID)

	//Assign and return Service ID with mysql
	var newservice model.ServiceStatus
	newservice.ProjectID = reqServiceProject.ProjectID
	newservice.ProjectName = reqServiceProject.ProjectName
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.currentUser.ID

	serviceInfo, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.internalError(err)
		return
	}
	p.renderJSON(serviceInfo.ID)
}

func (p *ServiceController) DeleteServiceAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)

	if err = service.CheckServiceDeletable(s); err != nil {
		p.internalError(err)
		return
	}
	// Call stop service if running
	if s.Status != stopped {
		err = service.StopServiceK8s(s)
		if err != nil {
			p.internalError(err)
			return
		}
	}

	// 	Delete the service's autoscale rule
	hpas, err := service.ListAutoScales(s)
	if err != nil {
		p.internalError(err)
		return
	}
	for _, hpa := range hpas {
		err := service.DeleteAutoScale(s, hpa.ID)
		if err != nil {
			p.internalError(err)
			return
		}
		logs.Debug("Deleted Hpa %d %s", hpa.ID, hpa.HPAName)
	}

	isSuccess, err := service.DeleteService(s.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", s.ID))
		return
	}

	//delete repo files of the service
	p.resolveRepoServicePath(s.ProjectName, s.Name)
	p.removeItemsToRepo(filepath.Join(s.Name, serviceFilename), filepath.Join(s.Name, deploymentFilename))

}

// API to deploy service
func (p *ServiceController) ToggleServiceAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	var reqServiceToggle model.ServiceToggle
	err = p.resolveBody(&reqServiceToggle)
	if err != nil {
		return
	}

	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)

	if s.Status == stopped && reqServiceToggle.Toggle == 0 {
		p.customAbort(http.StatusBadRequest, "Service already stopped.")
		return
	}

	if s.Status == running && reqServiceToggle.Toggle == 1 {
		p.customAbort(http.StatusBadRequest, "Service already running.")
		return
	}

	p.resolveRepoServicePath(s.ProjectName, s.Name)
	if _, err := os.Stat(p.repoServicePath); os.IsNotExist(err) {
		p.customAbort(http.StatusPreconditionFailed, "Service restored from initialization, cannot be switched.")
		return
	}
	if reqServiceToggle.Toggle == 0 {
		// stop service
		err = service.StopServiceK8s(s)
		if err != nil {
			p.internalError(err)
			return
		}
		// Update service status DB
		_, err = service.UpdateServiceStatus(s.ID, stopped)
		if err != nil {
			p.internalError(err)
			return
		}
	} else {
		// start service
		err := service.DeployServiceByYaml(s.ProjectName, p.repoServicePath)
		if err != nil {
			p.parseError(err, parsePostK8sError)
			return
		}
		// Push deployment to Git repo
		items := []string{filepath.Join(s.Name, deploymentFilename), filepath.Join(s.Name, serviceFilename)}
		p.pushItemsToRepo(items...)
		p.collaborateWithPullRequest("master", "master", items...)

		// Update service status DB
		_, err = service.UpdateServiceStatus(s.ID, running)
		if err != nil {
			p.internalError(err)
			return
		}
	}
}

func (p *ServiceController) GetServiceInfoAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	if s.Public != 1 {
		p.resolveUserPrivilegeByID(s.ProjectID)
	}
	serviceStatus, err := service.GetServiceByK8sassist(s.ProjectName, s.Name)
	if err != nil {
		p.parseError(err, parseGetK8sError)
		return
	}
	//Get NodeIP
	//endpointUrl format /api/v1/namespaces/default/endpoints/
	nodesStatus, err := service.GetNodesStatus()
	if err != nil {
		p.parseError(err, parseGetK8sError)
		return
	}
	if len(serviceStatus.Ports) == 0 || len(nodesStatus.Items) == 0 {
		p.renderJSON("NA")
		return
	}

	var serviceInfo model.ServiceInfoStruct
	for _, ports := range serviceStatus.Ports {
		serviceInfo.NodePort = append(serviceInfo.NodePort, ports.NodePort)
	}
	for _, items := range nodesStatus.Items {
		serviceInfo.NodeName = append(serviceInfo.NodeName, items.Status.Addresses...)
	}
	p.renderJSON(serviceInfo)
}

func (p *ServiceController) GetServiceStatusAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)
	serviceStatus, err := service.GetServiceByK8sassist(s.ProjectName, s.Name)
	if err != nil {
		p.parseError(err, parseGetK8sError)
		return
	}
	p.renderJSON(serviceStatus)
}

func (p *ServiceController) ServicePublicityAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)
	var reqServiceUpdate model.ServicePublicityUpdate
	err = p.resolveBody(&reqServiceUpdate)
	if err != nil {
		return
	}
	if s.Public != reqServiceUpdate.Public {
		_, err := service.UpdateServicePublic(s.ID, reqServiceUpdate.Public)
		if err != nil {
			p.internalError(err)
			return
		}
	}
}

func (p *ServiceController) DeleteServiceConfigAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	// Get the path of the service config files
	p.resolveUserPrivilege(s.ProjectName)
	p.resolveRepoServicePath(s.ProjectName, s.Name)
	logs.Debug("Service config path: %s", p.repoServicePath)

	// Delete yaml files
	// TODO
	err = service.DeleteServiceConfigYaml(p.repoServicePath)
	if err != nil {
		logs.Info("failed to delete service yaml", p.repoServicePath)
		p.internalError(err)
		return
	}

	// For terminated service config, actually delete it in DB
	_, err = service.DeleteServiceByID(s.ID)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ServiceController) DeleteDeploymentAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	// Get the path of the service config files
	p.resolveUserPrivilege(s.ProjectName)
	p.resolveRepoServicePath(s.ProjectName, s.Name)
	logs.Debug("Service config path: %s", p.repoServicePath)

	// TODO clear kube-master, even if the service is not deployed successfully
	p.removeItemsToRepo(filepath.Join(s.Name, deploymentFilename))

	// Delete yaml files
	err = service.DeleteServiceConfigYaml(p.repoServicePath)
	if err != nil {
		logs.Info("Failed to delete service yaml under path: %s", p.repoServicePath)
		p.internalError(err)
		return
	}

	// For terminated service config, actually delete it in DB
	_, err = service.DeleteServiceByID(s.ID)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ServiceController) StoreServiceRoute() {
	serviceIdentity := p.GetString("service_identity")
	serviceURL := p.GetString("service_url")
	memoryCache.Put(strings.ToLower(serviceIdentity), serviceURL, time.Second*time.Duration(tokenCacheExpireSeconds))
	logs.Debug("Service identity: %s, URL: %s", serviceIdentity, serviceURL)
}

func (p *ServiceController) ServiceExists() {
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)
	serviceName := p.GetString("service_name")
	isServiceExists, err := service.ServiceExists(serviceName, projectName)
	if err != nil {
		p.internalError(err)
		logs.Error("Check service name failed, error: %+v", err.Error())
		return
	}
	if isServiceExists == true {
		p.customAbort(http.StatusConflict, serverNameDuplicateErr.Error())
		return
	}
}

func (p *ServiceController) ScaleServiceAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)

	var reqServiceScale model.ServiceScale
	err = p.resolveBody(&reqServiceScale)
	if err != nil {
		return
	}
	// change the replica number of service
	res, err := service.ScaleReplica(s, reqServiceScale.Replica)
	if res != true {
		logs.Info("Failed to scale service replica", s, reqServiceScale.Replica)
		p.internalError(err)
		return
	}
	logs.Info("Scale service replica successfully")

	_, deploymentFileInfo, err := service.GetDeployment(s.ProjectName, s.Name)
	if err != nil {
		logs.Error("Failed to get deployment %s", s.Name)
		return
	}
	p.resolveRepoServicePath(s.ProjectName, s.Name)
	err = utils.GenerateFile(deploymentFileInfo, p.repoServicePath, deploymentFilename)
	if err != nil {
		logs.Error("Failed to update file of deployment %s", s.Name)
		return
	}
	p.pushItemsToRepo(filepath.Join(s.Name, deploymentFilename))
}

//get selectable service list
func (p *ServiceController) GetSelectableServicesAction() {
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)
	logs.Info("Get selectable service list for", projectName)
	serviceList, err := service.GetServicesByProjectName(projectName)
	if err != nil {
		logs.Error("Failed to get selectable services.")
		p.internalError(err)
		return
	}
	p.renderJSON(serviceList)
}

func (f *ServiceController) resolveUploadedYamlFile(uploadedFileName string) (func(fileName string, serviceInfo *model.ServiceStatus) error, io.Reader, error) {
	uploadedFile, _, err := f.GetFile(uploadedFileName)
	if err != nil {
		if err.Error() == "http: no such file" {
			f.customAbort(http.StatusBadRequest, "Missing file: "+uploadedFileName)
			return nil, nil, err
		}
		f.internalError(err)
		return nil, nil, err
	}

	return func(fileName string, serviceInfo *model.ServiceStatus) error {
		f.resolveRepoServicePath(serviceInfo.ProjectName, serviceInfo.Name)
		err = utils.CheckFilePath(f.repoServicePath)
		if err != nil {
			f.internalError(err)
			return nil
		}
		return f.SaveToFile(uploadedFileName, filepath.Join(f.repoServicePath, fileName))
	}, uploadedFile, nil
}

func (f *ServiceController) UploadYamlFileAction() {
	projectName := f.GetString("project_name")
	f.resolveProjectMember(projectName)

	fhDeployment, deploymentFile, err := f.resolveUploadedYamlFile("deployment_file")
	if err != nil {
		return
	}
	fhService, serviceFile, err := f.resolveUploadedYamlFile("service_file")
	if err != nil {
		return
	}
	deployInfo, err := service.CheckDeployYamlConfig(serviceFile, deploymentFile, projectName)
	if err != nil {
		f.customAbort(http.StatusBadRequest, err.Error())
		return
	}

	serviceName := deployInfo.Service.ObjectMeta.Name
	serviceInfo, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if serviceInfo != nil {
		f.customAbort(http.StatusBadRequest, "Service name has been used.")
		return
	}
	serviceInfo, err = service.CreateServiceConfig(model.ServiceStatus{
		Name:        serviceName,
		ProjectName: projectName,
		Status:      preparing, // 0: preparing 1: running 2: suspending
		OwnerID:     f.currentUser.ID,
		OwnerName:   f.currentUser.Username,
	})
	if err != nil {
		f.internalError(err)
		return
	}
	err = fhDeployment(deploymentFilename, serviceInfo)
	if err != nil {
		f.internalError(err)
		return
	}
	err = fhService(serviceFilename, serviceInfo)
	if err != nil {
		f.internalError(err)
		return
	}
	f.renderJSON(serviceInfo)
}

func (f *ServiceController) DownloadDeploymentYamlFileAction() {
	projectName := f.GetString("project_name")
	serviceName := f.GetString("service_name")
	serviceInfo, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if serviceInfo == nil {
		f.customAbort(http.StatusBadRequest, "Service name is invalid.")
		return
	}
	f.resolveRepoServicePath(projectName, serviceName)
	yamlType := f.GetString("yaml_type")
	if yamlType == "" {
		f.customAbort(http.StatusBadRequest, "No YAML type found.")
		return
	}
	if yamlType == deploymentType {
		f.resolveDownloadYaml(serviceInfo, deploymentFilename, service.GenerateDeploymentYamlFileFromK8s)
	} else if yamlType == serviceType {
		f.resolveDownloadYaml(serviceInfo, serviceFilename, service.GenerateServiceYamlFileFromK8s)
	}
}

func (f *ServiceController) resolveDownloadYaml(serviceConfig *model.ServiceStatus, fileName string, generator func(*model.ServiceStatus, string) error) {
	logs.Debug("Current download yaml file: %s", fileName)
	//checkout the path of download
	err := utils.CheckFilePath(f.repoServicePath)
	if err != nil {
		f.internalError(err)
		return
	}
	absFileName := filepath.Join(f.repoServicePath, fileName)
	err = generator(serviceConfig, f.repoServicePath)
	if err != nil {
		f.parseError(err, parseGetK8sError)
		return
	}
	logs.Info("User: %s downloaded %s YAML file.", f.currentUser.Username, fileName)
	f.Ctx.Output.Download(absFileName, fileName)
}

func (p *ServiceController) GetScaleStatusAction() {
	// Get the current service status
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)
	scaleStatus, err := service.GetScaleStatus(s)
	if err != nil {
		logs.Debug("Get scale deployment status failed %s", s.Name)
		p.internalError(err)
		return
	}
	p.renderJSON(scaleStatus)
}

func (p *ServiceController) DeleteDeployAction() {
	var err error

	key := p.getKey()
	configService := NewConfigServiceStep(key)

	//Judge authority
	p.resolveUserPrivilegeByID(configService.ProjectID)

	// Clean deployment and service

	s := model.ServiceStatus{Name: configService.ServiceName,
		ProjectName: configService.ProjectName,
	}

	err = service.StopServiceK8s(&s)
	if err != nil {
		logs.Error("Failed to clean service %s", s.Name)
		p.internalError(err)
		return
	}

	//Clean data DB if existing
	serviceData, err := service.GetService(s, "name", "project_name")
	if serviceData != nil {
		isSuccess, err := service.DeleteService(serviceData.ID)
		if err != nil {
			p.internalError(err)
			return
		}
		if !isSuccess {
			p.customAbort(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", s.ID))
			return
		}
	}

	//delete repo files of the service
	p.resolveRepoServicePath(s.ProjectName, s.Name)
	p.removeItemsToRepo(filepath.Join(s.Name, serviceFilename), filepath.Join(s.Name, deploymentFilename))

	//clean the config step
	err = DeleteConfigServiceStep(key)
	if err != nil {
		logs.Debug("Failed to clean the config steps")
		p.internalError(err)
		return
	}
}

//get service nodeport list
func (p *ServiceController) GetServiceNodePorts() {
	projectName := p.GetString("project_name")

	if projectName != "" {
		p.resolveProjectMember(projectName)
		logs.Info("Get service nodeport list for", projectName)
	}

	//Check nodeports in cluster
	nodeportList, err := service.GetNodePortsK8s(projectName)
	if err != nil {
		logs.Error("Failed to get selectable services.")
		p.internalError(err)
		return
	}

	//Check nodeports in DB
	p.renderJSON(nodeportList)
}

//import cluster services
func (p *ServiceController) ImportServicesAction() {

	if p.isSysAdmin == false {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to import services.")
		return
	}

	projectList, err := service.GetProjectsByUser(model.Project{}, p.currentUser.ID)
	if err != nil {
		logs.Error("Failed to get projects.")
		p.internalError(err)
		return
	}

	for _, project := range projectList {
		err := service.SyncServiceWithK8s(project.Name)
		if err != nil {
			logs.Error("Failed to sync service for project %s.", project.Name)
			p.internalError(err)
			return
		}
	}
	logs.Debug("imported services from cluster successfully")
}
