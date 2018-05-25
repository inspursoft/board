package controller

import (
	"encoding/json"
	"errors"
	"fmt"

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
)

type ServiceController struct {
	BaseController
}

func (p *ServiceController) generateDeploymentTravis(deploymentURL string, serviceURL string) error {
	userID := p.currentUser.ID
	var travisCommand travis.TravisCommand
	travisCommand.Script.Commands = []string{}
	items := []string{
		fmt.Sprintf("curl \"%s/jenkins-job/%d/$BUILD_NUMBER\"", boardAPIBaseURL(), userID),
	}
	if deploymentURL != "" {
		items = append(items, fmt.Sprintf("curl -X POST -H 'Content-Type: application/yaml' --data-binary @deployment.yaml %s", deploymentURL))
	}
	if serviceURL != "" {
		items = append(items, fmt.Sprintf("curl -X POST -H 'Content-Type: application/yaml' --data-binary @service.yaml %s", serviceURL))
	}
	travisCommand.Script.Commands = items
	return travisCommand.GenerateCustomTravis(p.repoPath)
}

func (p *ServiceController) getKey() string {
	return strconv.Itoa(int(p.currentUser.ID))
}

func (p *ServiceController) resolveServiceInfo() (s *model.ServiceStatus) {
	var err error
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

	p.resolveRepoPath(project.Name)
	err = service.CheckDeploymentPath(p.repoPath)
	if err != nil {
		p.internalError(err)
		return
	}

	err = service.AssembleDeploymentYaml((*model.ConfigServiceStep)(configService), p.repoPath)
	if err != nil {
		p.internalError(err)
		return
	}

	err = service.AssembleServiceYaml((*model.ConfigServiceStep)(configService), p.repoPath)
	if err != nil {
		p.internalError(err)
		return
	}

	deploymentURL := fmt.Sprintf("%s%s%s/%s", kubeMasterURL(), deploymentAPI, project.Name, "deployments")
	serviceURL := fmt.Sprintf("%s%s%s/%s", kubeMasterURL(), serviceAPI, project.Name, "services")
	err = p.generateDeploymentTravis(deploymentURL, serviceURL)
	if err != nil {
		logs.Error("Failed to generate deployement travis.yml: %+v", err)
		p.internalError(err)
		return
	}

	items := []string{".travis.yml", deploymentFilename, serviceFilename}
	p.pushItemsToRepo(items...)

	serviceConfig, err := json.Marshal(&configService)
	if err != nil {
		p.internalError(err)
		return
	}

	updateService := model.ServiceStatus{ID: serviceInfo.ID, Status: running, ServiceConfig: string(serviceConfig)}
	_, err = service.UpdateService(updateService, "id", "status", "service_config")
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
		deployment, err := service.GetDeployment((*serviceStatus).ProjectName, (*serviceStatus).Name)
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
	orderField := p.GetString("order_field", "CREATE_TIME")
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

	p.resolveBody(&reqServiceProject)

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
	s := p.resolveServiceInfo()
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)

	var err error
	// Call stop service if running
	switch s.Status {
	case running:
		//err = stopService(s)
		err = service.StopServiceK8s(s)
		if err != nil {
			p.internalError(err)
			return
		}
	case uncompleted:
		timeInt := time.Now().Sub(s.UpdateTime)
		logs.Debug("uncompleted status in %+v", timeInt)
		if timeInt < startingDuration {
			p.customAbort(http.StatusBadRequest,
				fmt.Sprintf("Invalid request %d in starting status", s.ID))
			return
		}
		err = service.CleanDeploymentK8s(s)
		if err != nil {
			logs.Error("Failed to clean deployment %s", s.Name)
			p.internalError(err)
			return
		}
		err = service.CleanServiceK8s(s)
		if err != nil {
			logs.Error("Failed to clean service %s", s.Name)
			p.internalError(err)
			return
		}
	}

	isSuccess, err := service.DeleteService(s.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", s.ID))
	}

	//delete repo files of the service
	p.resolveRepoPath(s.ProjectName)
	p.removeItemsToRepo(serviceFilename, deploymentFilename)

}

// API to deploy service
func (p *ServiceController) ToggleServiceAction() {
	s := p.resolveServiceInfo()

	var reqServiceToggle model.ServiceToggle
	p.resolveBody(&reqServiceToggle)

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

	var err error
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
		// Push deployment to Git repo
		deploymentURL := fmt.Sprintf("%s%s%s/%s", kubeMasterURL(), deploymentAPI, s.ProjectName, "deployments")
		serviceURL := fmt.Sprintf("%s%s%s/%s", kubeMasterURL(), serviceAPI, s.ProjectName, "services")

		p.resolveRepoPath(s.ProjectName)
		err = p.generateDeploymentTravis(deploymentURL, serviceURL)
		if err != nil {
			logs.Error("Failed to generate deployment travis: %+v", err)
			p.internalError(err)
			return
		}
		// Add deployment file to repo
		items := []string{".travis.yml", deploymentFilename, serviceFilename}
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

func stopService(s *model.ServiceStatus) error {
	// Stop service
	header := http.Header{
		"Content-Type": []string{"application/yaml"},
	}
	deleteServiceURL := kubeMasterURL() + serviceAPI + s.ProjectName + "/services/" + s.Name
	err := utils.SimpleDeleteRequestHandle(deleteServiceURL, header)
	if err != nil {
		logs.Error("Failed to request %s to stop service.", deleteServiceURL)
	}
	logs.Info("Stop service successfully, id: %d, name: %s", s.ID, s.Name)
	// Stop deployment
	deleteDeploymentURL := kubeMasterURL() + deploymentAPI + s.ProjectName + "/deployments/" + s.Name
	err = utils.SimpleDeleteRequestHandle(deleteDeploymentURL, header)
	if err != nil {
		logs.Error("Failed to request %s to stop deployment.", deleteDeploymentURL)
	}
	logs.Info("Stop deployment successfully, id: %d, name: %s", s.ID, s.Name)
	return nil
}

func (p *ServiceController) resolveErrOutput(err error) {
	if err != nil {
		if strings.Index(err.Error(), "StatusNotFound:") == 0 {
			var output interface{}
			json.Unmarshal([]byte(err.Error()[len("StatusNotFound:"):]), &output)
			p.Data["json"] = output
			p.ServeJSON()
			return
		}
		p.internalError(err)
	}
}

func (p *ServiceController) GetServiceInfoAction() {

	s := p.resolveServiceInfo()
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)

	serviceStatus, err := service.GetServiceStatus(kubeMasterURL() + serviceAPI + s.ProjectName + "/services/" + s.Name)
	if err != nil {
		p.resolveErrOutput(err)
		return
	}
	//Get NodeIP
	//endpointUrl format /api/v1/namespaces/default/endpoints/
	nodesStatus, err := service.GetNodesStatus(fmt.Sprintf("%s/api/v1/nodes", kubeMasterURL()))
	if err != nil {
		p.resolveErrOutput(err)
		return
	}
	if len(serviceStatus.Spec.Ports) == 0 || len(nodesStatus.Items) == 0 {
		p.renderJSON("NA")
		return
	}

	var serviceInfo model.ServiceInfoStruct
	for _, ports := range serviceStatus.Spec.Ports {
		serviceInfo.NodePort = append(serviceInfo.NodePort, ports.NodePort)
	}
	for _, items := range nodesStatus.Items {
		serviceInfo.NodeName = append(serviceInfo.NodeName, items.Status.Addresses...)
	}
	p.renderJSON(serviceInfo)
}

func (p *ServiceController) GetServiceStatusAction() {
	s := p.resolveServiceInfo()
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)
	serviceStatus, err := service.GetServiceStatus(kubeMasterURL() + serviceAPI + s.ProjectName + "/services/" + s.Name)
	if err != nil {
		p.resolveErrOutput(err)
		return
	}
	p.renderJSON(serviceStatus)
}

func (p *ServiceController) ServicePublicityAction() {
	s := p.resolveServiceInfo()
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)
	var reqServiceUpdate model.ServicePublicityUpdate
	p.resolveBody(&reqServiceUpdate)
	if s.Public != reqServiceUpdate.Public {
		_, err := service.UpdateServicePublic(s.ID, reqServiceUpdate.Public)
		if err != nil {
			p.internalError(err)
			return
		}
	}
}

func (p *ServiceController) DeleteServiceConfigAction() {
	s := p.resolveServiceInfo()
	// Get the path of the service config files
	p.resolveUserPrivilege(s.ProjectName)
	p.resolveRepoPath(s.ProjectName)
	logs.Debug("Service config path: %s", p.repoPath)

	// Delete yaml files
	// TODO
	err := service.DeleteServiceConfigYaml(p.repoPath)
	if err != nil {
		logs.Info("failed to delete service yaml", p.repoPath)
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
	s := p.resolveServiceInfo()
	// Get the path of the service config files
	p.resolveUserPrivilege(s.ProjectName)
	p.resolveRepoPath(s.ProjectName)
	logs.Debug("Service config path: %s", p.repoPath)

	// TODO clear kube-master, even if the service is not deployed successfully
	deploymentURL := filepath.Join(kubeMasterURL(), deploymentAPI, s.ProjectName, "deployments")
	err := p.generateDeploymentTravis(deploymentURL, "")
	if err != nil {
		logs.Error("Failed to generate deployment travis: %+v", err)
		p.internalError(err)
		return
	}
	// Update git repo
	p.removeItemsToRepo(".travis.yml", deploymentFilename)

	// Delete yaml files
	err = service.DeleteServiceConfigYaml(p.repoPath)
	if err != nil {
		logs.Info("failed to delete service yaml", p.repoPath)
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
	s := p.resolveServiceInfo()
	//Judge authority
	p.resolveUserPrivilegeByID(s.ProjectID)

	var reqServiceScale model.ServiceScale
	p.resolveBody(&reqServiceScale)
	// change the replica number of service
	res, err := service.ScaleReplica(s, reqServiceScale.Replica)

	if res != true {
		logs.Info("Failed to scale service replica", s, reqServiceScale.Replica)
		p.internalError(err)
		return
	}
	logs.Info("Scale service replica successfully")
}

//get selectable service list
func (p *ServiceController) GetSelectableServicesAction() {
	serviceName := p.GetString("service_name")
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)
	logs.Info("Get selectable service list for", projectName, serviceName)
	serviceList, err := service.GetSelectableServices(projectName, serviceName)
	if err != nil {
		logs.Error("Failed to get selectable services.")
		p.internalError(err)
		return
	}
	p.renderJSON(serviceList)
}

func (f *ServiceController) resolveUploadedYamlFile(uploadedFileName string, target interface{}, customError error) func(fileName string, serviceInfo *model.ServiceStatus) error {
	uploadedFile, _, err := f.GetFile(uploadedFileName)
	if err != nil {
		if err.Error() == "http: no such file" {
			f.customAbort(http.StatusBadRequest, "Missing file: "+uploadedFileName)
			return nil
		}
		f.internalError(err)
		return nil
	}
	err = utils.UnmarshalYamlFile(uploadedFile, target)
	if err != nil {
		if strings.Index(err.Error(), "InternalError:") == 0 {
			f.internalError(errors.New(err.Error()[len("InternalError:"):]))
			return nil
		}
		f.customAbort(http.StatusBadRequest, customError.Error())
		return nil
	}

	return func(fileName string, serviceInfo *model.ServiceStatus) error {
		f.resolveRepoPath(serviceInfo.ProjectName)
		err = service.CheckDeploymentPath(f.repoPath)
		if err != nil {
			f.internalError(err)
			return nil
		}
		return f.SaveToFile(uploadedFileName, filepath.Join(f.repoPath, fileName))
	}
}

func (f *ServiceController) UploadYamlFileAction() {
	projectName := f.GetString("project_name")
	f.resolveProjectMember(projectName)

	var deploymentConfig service.Deployment
	fhDeployment := f.resolveUploadedYamlFile("deployment_file", &deploymentConfig, service.DeploymentYamlFileUnmarshalErr)

	var serviceConfig service.Service
	fhService := f.resolveUploadedYamlFile("service_file", &serviceConfig, service.ServiceYamlFileUnmarshalErr)

	var err error
	err = service.CheckDeploymentConfig(projectName, deploymentConfig)
	if err != nil {
		f.customAbort(http.StatusBadRequest, err.Error())
	}
	err = service.CheckServiceConfig(projectName, serviceConfig)
	if err != nil {
		f.customAbort(http.StatusBadRequest, err.Error())
	}
	//check label selector
	serviceInfo, err := service.GetServiceByProject(serviceConfig.Name, projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if serviceInfo != nil {
		f.customAbort(http.StatusBadRequest, "Service name has been used.")
		return
	}
	serviceInfo, err = service.CreateServiceConfig(model.ServiceStatus{
		Name:        serviceConfig.Name,
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
	f.resolveRepoPath(projectName)
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
	yamlType := f.GetString("yaml_type")
	if yamlType == "" {
		f.customAbort(http.StatusBadRequest, "No YAML type found.")
	}
	if yamlType == deploymentType {
		deploymentConfigURL := kubeMasterURL() + filepath.Join(deploymentAPI, projectName, "deployments", serviceName)
		f.resolveDownloadYaml(deploymentConfigURL, deploymentFilename, service.GenerateDeploymentYamlFileFromK8S)
	} else if yamlType == serviceType {
		serviceConfigURL := kubeMasterURL() + filepath.Join(serviceAPI, projectName, "services", serviceName)
		f.resolveDownloadYaml(serviceConfigURL, serviceFilename, service.GenerateServiceYamlFileFromK8S)
	}
}

func (f *ServiceController) resolveDownloadYaml(configURL, fileName string, generator func(targetURL, path string) error) {
	logs.Debug("Current download config URL: %s", configURL)
	absFileName := filepath.Join(f.repoPath, fileName)
	err := generator(configURL, absFileName)
	if err != nil {
		if strings.Index(err.Error(), "StatusNotFound:") == 0 {
			f.customAbort(http.StatusNotFound, service.ServiceNotFoundErr.Error())
			return
		}
		f.internalError(err)
		return
	}
	logs.Info("User: %s downloaded %s YAML file.", f.currentUser.Username, fileName)
	f.Ctx.Output.Download(absFileName, fileName)
}

func (p *ServiceController) GetScaleStatusAction() {
	// Get the current service status
	s := p.resolveServiceInfo()
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
