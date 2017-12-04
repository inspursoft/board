package controller

import (
	"encoding/json"
	"fmt"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	deploymentFilename     = "deployment.yaml"
	serviceFilename        = "service.yaml"
	deploymentTestFilename = "testdeployment.yaml"
	serviceTestFilename    = "testservice.yaml"
	serviceProcess         = "process_service"
	apiheader              = "Content-Type: application/yaml"
	deploymentAPI          = "/apis/extensions/v1beta1/namespaces/"
	serviceAPI             = "/api/v1/namespaces/"
	test                   = "test"
	nacomment              = "NA"
	serviceNamespace       = "default" //TODO create in project post
)

const (
	preparing = iota
	running
	stopped
	warning
)

type ServiceController struct {
	baseController
}

//  Checking the user priviledge by token
func (p *ServiceController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
}

func (p *ServiceController) handleReqData() (model.ServiceConfig2, int, error) {
	var reqServiceConfig model.ServiceConfig2
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		return reqServiceConfig, 0, err
	}

	//get request message
	reqData, err := p.resolveBody()
	if err != nil {
		return reqServiceConfig, serviceID, err
	}

	//Get request data
	err = json.Unmarshal(reqData, &reqServiceConfig)
	if err != nil {
		return reqServiceConfig, serviceID, err
	}
	logs.Debug("Deploy service info: %+v\n", reqServiceConfig)

	return reqServiceConfig, serviceID, err
}

func (p *ServiceController) commonDeployService(reqServiceConfig model.ServiceConfig2, serviceID int, testFlag bool) {
	currentProject, err := service.GetProject(model.Project{ID: reqServiceConfig.Project.ProjectID}, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if currentProject == nil {
		p.customAbort(http.StatusBadRequest, "Invalid project ID.")
		return
	}
	isMember, err := service.IsProjectMember(reqServiceConfig.Project.ProjectID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to deploy service.")
		return
	}

	//get deployment type
	var deploymentFile string
	var serviceFile string
	if testFlag == true {
		deploymentFile = deploymentTestFilename
		serviceFile = serviceTestFilename
	} else {
		deploymentFile = deploymentFilename
		serviceFile = serviceFilename
	}

	var pushobject pushObject
	deploymentAbsName := filepath.Join(repoPath(), reqServiceConfig.Project.ProjectName, strconv.Itoa(serviceID), deploymentFile)
	pushobject.FileName = deploymentFile
	pushobject.JobName = serviceProcess
	pushobject.Value = filepath.Join(reqServiceConfig.Project.ProjectName, strconv.Itoa(serviceID))
	pushobject.Message = fmt.Sprintf("Create deployment for project %s service %d",
		reqServiceConfig.Project.ProjectName, reqServiceConfig.Project.ServiceID)
	pushobject.Extras = filepath.Join(kubeMasterURL(), deploymentAPI, serviceNamespace, "deployments")
	pushobject.Items = []string{filepath.Join(pushobject.Value, deploymentFile)}
	loadPath := filepath.Join(repoPath(), reqServiceConfig.Project.ProjectName, strconv.Itoa(serviceID))
	logs.Info("deployment pushobject.FileName:%+v\n", pushobject.FileName)
	logs.Info("deployment pushobject.Value:%+v\n", pushobject.Value)
	logs.Info("deployment pushobject.Extras:%+v\n", pushobject.Extras)

	err = service.CheckDeploymentPath(loadPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return
	}

	err = service.GenerateYamlFile(deploymentAbsName, reqServiceConfig.Deployment)
	if err != nil {
		p.internalError(err)
		return
	}

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push deployment object: %d %s", ret, msg)

	serviceAbsName := filepath.Join(repoPath(), reqServiceConfig.Project.ProjectName, strconv.Itoa(serviceID), serviceFile)
	pushobject.FileName = serviceFile
	pushobject.Message = fmt.Sprintf("Create service for project %s service %d",
		reqServiceConfig.Project.ProjectName, reqServiceConfig.Project.ServiceID)
	pushobject.Extras = filepath.Join(kubeMasterURL(), serviceAPI, serviceNamespace, "services")
	pushobject.Items = []string{filepath.Join(pushobject.Value, serviceFile)}
	logs.Info("service pushobject.FileName:%+v\n", pushobject.FileName)
	logs.Info("service pushobject.Value:%+v\n", pushobject.Value)
	logs.Info("service pushobject.Extras:%+v\n", pushobject.Extras)

	err = service.GenerateYamlFile(serviceAbsName, reqServiceConfig.Service)
	if err != nil {
		p.internalError(err)
		return
	}

	ret, msg, err = InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}

	ret, msg, err = InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push deployment object: %d %s", ret, msg)

	//update db
	updateService := model.ServiceStatus{ID: int64(serviceID), Status: running,
		Name: reqServiceConfig.Service.ObjectMeta.Name, OwnerID: p.currentUser.ID,
		OwnerName: p.currentUser.Username, Comment: reqServiceConfig.Project.Comment}
	_, err = service.UpdateService(updateService, "name", "status", "owner_id", "OwnerName", "Comment")
	if err != nil {
		p.internalError(err)
		return
	}

}

// API to deploy service
func (p *ServiceController) DeployServiceAction() {
	reqServiceConfig, serviceID, err := p.handleReqData()
	if err != nil {
		p.internalError(err)
		return
	}

	//check data
	err = service.CheckReqPara(reqServiceConfig)
	if err != nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
		logs.Error("Request parameters error when deploy service, error: %+v", err.Error())
		return
	}

	isServiceDuplicated, err := service.ServiceExists(reqServiceConfig.Service.ObjectMeta.Name, reqServiceConfig.Project.ProjectName)
	if err != nil {
		p.internalError(err)
		return
	}
	if isServiceDuplicated == true {
		p.serveStatus(http.StatusBadRequest, serverNameDuplicateErr.Error())
		logs.Error("Request parameters error when deploy service, error: %+v", serverNameDuplicateErr.Error())
		return
	}

	//Add registry to container images for deployment
	for index, container := range reqServiceConfig.Deployment.Spec.Template.Spec.Containers {
		reqServiceConfig.Deployment.Spec.Template.Spec.Containers[index].Image =
			filepath.Join(registryBaseURI(), container.Image)
	}
	logs.Debug("%+v", reqServiceConfig)

	p.commonDeployService(reqServiceConfig, serviceID, false)
}

// API to deploy test service
func (p *ServiceController) DeployServiceTestAction() {
	//Judge authority
	reqServiceConfig, serviceID, err := p.handleReqData()
	if err != nil {
		p.internalError(err)
		return
	}

	//check data
	err = service.CheckReqPara(reqServiceConfig)
	if err != nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
		logs.Error("Request parameters error when deploy service, error: %+v", err.Error())
		return
	}

	isServiceDuplicated, err := service.ServiceExists(reqServiceConfig.Service.ObjectMeta.Name, reqServiceConfig.Project.ProjectName)
	if err != nil {
		p.internalError(err)
		logs.Error("Request parameters error when deploy service, error: %+v", serverNameDuplicateErr.Error())
		return
	}
	if isServiceDuplicated == true {
		p.customAbort(http.StatusBadRequest, serverNameDuplicateErr.Error())
		return
	}

	//Add registry to container images for deployment
	for index, container := range reqServiceConfig.Deployment.Spec.Template.Spec.Containers {
		reqServiceConfig.Deployment.Spec.Template.Spec.Containers[index].Image =
			filepath.Join(registryBaseURI(), container.Image)
	}

	reqServiceConfig.Service.ObjectMeta.Name = test + reqServiceConfig.Service.ObjectMeta.Name
	reqServiceConfig.Service.ObjectMeta.Labels["app"] = test + reqServiceConfig.Service.ObjectMeta.Labels["app"]
	reqServiceConfig.Deployment.ObjectMeta.Name = test + reqServiceConfig.Deployment.ObjectMeta.Name
	reqServiceConfig.Deployment.Spec.Template.ObjectMeta.Labels["app"] = test + reqServiceConfig.Deployment.Spec.Template.ObjectMeta.Labels["app"]

	logs.Debug("%+v", reqServiceConfig)

	p.commonDeployService(reqServiceConfig, serviceID, true)
}

//
func syncK8sStatus(serviceList []*model.ServiceStatus) error {
	var err error
	// synchronize service status with the cluster system
	for _, serviceStatus := range serviceList {
		if (*serviceStatus).Status == stopped {
			continue
		}
		// Check the deployment status
		deployment, err := service.GetDeployment(serviceNamespace, (*serviceStatus).Name)
		if deployment == nil {
			logs.Info("Failed to get deployment", err)
			var reason = "The deployment is not established in cluster system"
			(*serviceStatus).Status = warning
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
				(*serviceStatus).Status = warning
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
		serviceK8s, err := service.GetK8sService(serviceNamespace, (*serviceStatus).Name)
		if serviceK8s == nil {
			logs.Info("Failed to get service in cluster", err)
			var reason = "The service is not established in cluster system"
			(*serviceStatus).Status = warning
			(*serviceStatus).Comment = "Reason: " + reason
			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
			if err != nil {
				logs.Error("Failed to update service in cluster.")
				break
			}
			continue
		}

		if serviceStatus.Status == warning {
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
	serviceName := p.GetString("service_name", "")
	pageIndex, _ := p.GetInt("page_index", 0)
	pageSize, _ := p.GetInt("page_size", 0)
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
		p.Data["json"] = serviceStatus
	} else {
		paginatedServiceStatus, err := service.GetPaginatedServiceList(serviceName, p.currentUser.ID, pageIndex, pageSize)
		if err != nil {
			p.internalError(err)
			return
		}
		err = syncK8sStatus(paginatedServiceStatus.ServiceStatusList)
		if err != nil {
			p.internalError(err)
			return
		}
		p.Data["json"] = paginatedServiceStatus
	}
	p.ServeJSON()
}

// API to create service config
func (p *ServiceController) CreateServiceConfigAction() {
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}
	var reqServiceProject model.ServiceProject
	err = json.Unmarshal(reqData, &reqServiceProject)
	if err != nil {
		p.internalError(err)
		return
	}
	//Assign and return Service ID with mysql
	var newservice model.ServiceStatus
	newservice.ProjectID = reqServiceProject.ProjectID
	newservice.ProjectName = reqServiceProject.ProjectName
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.currentUser.ID

	isMember, err := service.IsProjectMember(newservice.ProjectID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	//Judge authority
	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to create service.")
		return
	}

	serviceID, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.internalError(err)
		return
	}
	p.Data["json"] = strconv.Itoa(int(serviceID))
	p.ServeJSON()
}

func (p *ServiceController) DeleteServiceAction() {
	serviceID, err := strconv.ParseInt(p.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		p.internalError(err)
		return
	}
	// Check service id exist
	var servicequery model.ServiceStatus
	servicequery.ID = int64(serviceID)
	s, err := service.GetService(servicequery, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}

	isMember, err := service.IsProjectMember(s.ProjectID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	//Judge authority
	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to delete service.")
		return
	}

	// Call stop service if running
	if s.Status == running {
		err = stopService(s)
		if err != nil {
			p.internalError(err)
			return
		}
	}

	isSuccess, err := service.DeleteService(serviceID)
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", serviceID))
	}
}

// API to deploy service
func (p *ServiceController) ToggleServiceAction() {
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}

	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}
	var reqServiceToggle model.ServiceToggle
	err = json.Unmarshal(reqData, &reqServiceToggle)
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info(reqServiceToggle)

	// Check the current service status
	var servicequery model.ServiceStatus
	servicequery.ID = int64(serviceID)
	s, err := service.GetService(servicequery, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}

	isMember, err := service.IsProjectMember(s.ProjectID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	//Judge authority
	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to toggle service status.")
		return
	}

	if s.Status == stopped && reqServiceToggle.Toggle == 0 {
		p.customAbort(http.StatusBadRequest, "Service already stopped.")
		return
	}

	if s.Status == running && reqServiceToggle.Toggle == 1 {
		p.customAbort(http.StatusBadRequest, "Service already running.")
		return
	}

	if reqServiceToggle.Toggle == 0 {
		// stop service
		err = stopService(s)
		if err != nil {
			p.internalError(err)
			return
		}
		//logs.Info("Stop service successful")
		// Update service status DB
		servicequery.Status = stopped
		_, err = service.UpdateService(servicequery, "status")
		if err != nil {
			p.internalError(err)
			return
		}
	} else {
		// start service
		//serviceNamespace = reqServiceConfig.ProjectName TODO in project
		// Push deployment to jenkins
		var pushobject pushObject
		pushobject.FileName = deploymentFilename
		pushobject.JobName = serviceProcess
		pushobject.Value = filepath.Join(s.ProjectName, strconv.Itoa(serviceID))
		pushobject.Message = fmt.Sprintf("Create deployment for project %s service %d",
			s.ProjectName, s.ID)
		pushobject.Extras = filepath.Join(kubeMasterURL(), deploymentAPI,
			serviceNamespace, "deployments")

		// Add deployment file
		pushobject.Items = []string{filepath.Join(pushobject.Value, deploymentFilename)}

		ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
		if err != nil {
			p.internalError(err)
			return
		}
		logs.Info("Internal push deployment object: %d %s", ret, msg)

		//TODO: If fail to create deployment, should not continue to create service

		//Push service to jenkins
		pushobject.FileName = serviceFilename
		pushobject.Message = fmt.Sprintf("Create service for project %s service %d",
			s.ProjectName, s.ID)
		pushobject.Extras = filepath.Join(kubeMasterURL(), serviceAPI, serviceNamespace, "services")
		// Add deployment file
		pushobject.Items = []string{filepath.Join(pushobject.Value, serviceFilename)}

		ret, msg, err = InternalPushObjects(&pushobject, &(p.baseController))
		if err != nil {
			p.internalError(err)
			return
		}
		logs.Debug("Internal push service object: %d %s", ret, msg)

		// Update service status DB
		servicequery.Status = running
		_, err = service.UpdateService(servicequery, "status")
		if err != nil {
			p.internalError(err)
			return
		}
		//logs.Info("Start service successful")
	}
}

func stopService(s *model.ServiceStatus) error {
	var err error
	var client = &http.Client{}
	// Stop service
	//deleteServiceURL := filepath.Join(kubeMasterURL(), serviceAPI,
	//	serviceNamespace, "services", s.Name)
	deleteServiceURL := kubeMasterURL() + serviceAPI + serviceNamespace + "/services/" + s.Name
	req, err := http.NewRequest("DELETE", deleteServiceURL, nil)
	if err != nil {
		logs.Error("Failed to new request for delete service: %s", deleteServiceURL)
		return err
	}
	req.Header.Set("Content-Type", "application/yaml")
	resp, err := client.Do(req)
	if err != nil {
		logs.Info(req)
		return err
	}
	defer resp.Body.Close()
	logs.Info("Stop service successfully", s.ID, s.Name, resp)

	// Stop deployment
	//deleteDeploymentURL := filepath.Join(kubeMasterURL(), deploymentAPI,
	//	serviceNamespace, "deployments", s.Name)
	deleteDeploymentURL := kubeMasterURL() + deploymentAPI + serviceNamespace + "/deployments/" + s.Name
	req, err = http.NewRequest("DELETE", deleteDeploymentURL, nil)
	if err != nil {
		logs.Error("Failed to new request for delete deployment: %s", deleteDeploymentURL)
		return err
	}
	req.Header.Set("Content-Type", "application/yaml")
	resp, err = client.Do(req)
	if err != nil {
		logs.Error(req)
		return err
	}
	defer resp.Body.Close()

	logs.Info("Stop deployment successfully, id: %d, name: %s, resp: %+v", s.ID, s.Name, resp)
	return nil
}

func (p *ServiceController) GetServiceInfoAction() {
	var serviceInfo model.ServiceInfoStruct

	//Get Nodeport
	serviceName := p.Ctx.Input.Param(":service_name")
	serviceURL := fmt.Sprintf("%s/api/v1/namespaces/default/services/%s", kubeMasterURL(), serviceName)
	logs.Debug("Get Service info serviceURL(service): %+s", serviceURL)
	serviceStatus, err, flag := service.GetServiceStatus(serviceURL)
	var errOutput interface{}
	if flag == false {
		json.Unmarshal([]byte(err.Error()), &errOutput)
		p.Data["json"] = errOutput
		p.ServeJSON()
		return
	}

	if err != nil {
		p.internalError(err)
		return
	}

	//Get NodeIP
	//endpointUrl format /api/v1/namespaces/default/endpoints/
	nodesURL := fmt.Sprintf("%s/api/v1/nodes", kubeMasterURL())
	logs.Debug("Get Node info nodeURL (endpoint): %+s", nodesURL)
	nodesStatus, err, flag := service.GetNodesStatus(nodesURL)
	if flag == false {
		json.Unmarshal([]byte(err.Error()), &errOutput)
		p.Data["json"] = errOutput
		p.ServeJSON()
		return
	}

	if err != nil {
		p.internalError(err)
		return
	}

	if len(serviceStatus.Spec.Ports) == 0 || len(nodesStatus.Items) == 0 {
		p.Data["json"] = nacomment
		p.ServeJSON()
		return
	}

	for _, ports := range serviceStatus.Spec.Ports {
		serviceInfo.NodePort = append(serviceInfo.NodePort, ports.NodePort)
	}
	for _, items := range nodesStatus.Items {
		serviceInfo.NodeName = append(serviceInfo.NodeName, items.Status.Addresses...)
	}

	p.Data["json"] = serviceInfo
	p.ServeJSON()
}

func (p *ServiceController) GetServiceStatusAction() {
	serviceName := p.Ctx.Input.Param(":service_name")
	serviceURL := fmt.Sprintf("%s/api/v1/namespaces/default/services/%s", kubeMasterURL(), serviceName)
	logs.Debug("Get Service Status serviceURL: %+s", serviceURL)
	serviceStatus, err, flag := service.GetServiceStatus(serviceURL)

	var errOutput interface{}
	if flag == false {
		json.Unmarshal([]byte(err.Error()), &errOutput)
		p.Data["json"] = errOutput
		p.ServeJSON()
		return
	}
	if err != nil {
		p.internalError(err)
		return
	}
	p.Data["json"] = serviceStatus
	p.ServeJSON()
}

func (p *ServiceController) ServicePublicityAction() {
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}

	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}
	var reqServiceUpdate model.ServicePublicityUpdate
	err = json.Unmarshal(reqData, &reqServiceUpdate)
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info(reqServiceUpdate)

	// Check the current service status
	var servicequery model.ServiceStatus
	servicequery.ID = int64(serviceID)
	s, err := service.GetService(servicequery, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}

	isMember, err := service.IsProjectMember(s.ProjectID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	//Judge authority
	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to get publicity of service.")
		return
	}

	if s.Public != reqServiceUpdate.Public {
		servicequery.Public = reqServiceUpdate.Public
		_, err = service.UpdateService(servicequery, "public")
		if err != nil {
			p.internalError(err)
			return
		}
	} else {
		logs.Info("Already in target publicity status")
	}
}

func (p *ServiceController) GetServiceConfigAction() {
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}

	// Get the project info of this service
	var servicequery model.ServiceStatus
	servicequery.ID = int64(serviceID)
	s, err := service.GetService(servicequery, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest,
			fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	logs.Info("service status: ", s)

	// Get the path of the service config files
	serviceConfigPath := filepath.Join(repoPath(), s.ProjectName,
		strconv.Itoa(serviceID))
	logs.Debug("Service config path: %s", serviceConfigPath)

	// Get the service config by yaml files
	var serviceConfig model.ServiceConfig2
	serviceConfig.Project.ServiceID = s.ID
	serviceConfig.Project.ProjectID = s.ProjectID
	serviceConfig.Project.ProjectName = s.ProjectName

	err = service.UnmarshalServiceConfigYaml(&serviceConfig, serviceConfigPath)
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Debug("Service config: ", serviceConfig)

	p.Data["json"] = serviceConfig
	p.ServeJSON()
}

func (p *ServiceController) UpdateServiceConfigAction() {
	var err error
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}

	//Get request massage of service config
	reqServiceConfig, serviceID, err := p.handleReqData()
	if err != nil {
		p.internalError(err)
		return
	}

	serviceConfigPath := filepath.Join(repoPath(), reqServiceConfig.Project.ProjectName,
		strconv.Itoa(serviceID))
	//update yaml files by service config
	err = service.UpdateServiceConfigYaml(reqServiceConfig, serviceConfigPath)
	if err != nil {
		logs.Info("failed to update service yaml", serviceConfigPath)
		p.customAbort(http.StatusBadRequest, err.Error())
		return
	}
}

func (p *ServiceController) DeleteServiceConfigAction() {
	var err error
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}

	// Get the project info of this service
	var servicequery model.ServiceStatus
	servicequery.ID = int64(serviceID)
	s, err := service.GetService(servicequery, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest,
			fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	logs.Info("service status: ", s)

	// Get the path of the service config files
	serviceConfigPath := filepath.Join(repoPath(), s.ProjectName,
		strconv.Itoa(serviceID))
	logs.Debug("Service config path: %s", serviceConfigPath)

	// Delete yaml files
	// TODO
	err = service.DeleteServiceConfigYaml(serviceConfigPath)
	if err != nil {
		logs.Info("failed to delete service yaml", serviceConfigPath)
		p.internalError(err)
		return
	}

	// For terminated service config, actually delete it in DB
	_, err = service.DeleteServiceByID(servicequery)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ServiceController) DeleteDeploymentAction() {
	var err error
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}

	// Get the project info of this service
	var servicequery model.ServiceStatus
	servicequery.ID = int64(serviceID)
	s, err := service.GetService(servicequery, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest,
			fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	logs.Info("service status: ", s)

	// Get the path of the service config files
	serviceConfigPath := filepath.Join(repoPath(), s.ProjectName,
		strconv.Itoa(serviceID))
	logs.Debug("Service config path: %s", serviceConfigPath)

	// TODO clear kube-master, even if the service is not deployed successfully

	// Update git repo
	var pushobject pushObject
	pushobject.FileName = deploymentFilename
	pushobject.JobName = serviceProcess
	pushobject.Value = filepath.Join(s.ProjectName, strconv.Itoa(serviceID))
	pushobject.Message = fmt.Sprintf("Delete yaml files for project %s service %d",
		s.ProjectName, s.ID)
	pushobject.Extras = filepath.Join(kubeMasterURL(), deploymentAPI,
		serviceNamespace, "deployments")

	//Get file list for Jenkis git repo
	uploads, err := service.ListUploadFiles(serviceConfigPath)
	if err != nil {
		p.internalError(err)
		return
	}
	// Add yaml files
	for _, finfo := range uploads {
		filefullname := filepath.Join(pushobject.Value, finfo.FileName)
		pushobject.Items = append(pushobject.Items, filefullname)
	}

	ret, msg, err := InternalCleanObjects(&pushobject, &(p.baseController))
	if err != nil {
		logs.Info("Failed to push object for git repo clean", msg, ret, pushobject)
		p.internalError(err)
		return
	}
	logs.Info("Internal push clean deployment object: %d %s", ret, msg)

	// Delete yaml files
	err = service.DeleteServiceConfigYaml(serviceConfigPath)
	if err != nil {
		logs.Info("failed to delete service yaml", serviceConfigPath)
		p.internalError(err)
		return
	}

	// For terminated service config, actually delete it in DB
	_, err = service.DeleteServiceByID(servicequery)
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
	serviceName := p.GetString("service_name")
	isServiceExists, err := service.ServiceExists(serviceName, projectName)
	if err != nil {
		p.internalError(err)
		logs.Error("Check service name failed, error: %+v", err.Error())
		return
	}
	p.Data["json"] = isServiceExists
	p.ServeJSON()
}

func (p *ServiceController) ScaleServiceAction() {
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}

	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}
	var reqServiceScale model.ServiceScale
	err = json.Unmarshal(reqData, &reqServiceScale)
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info(reqServiceScale)

	// Get the current service status
	var servicequery model.ServiceStatus
	servicequery.ID = int64(serviceID)
	s, err := service.GetService(servicequery, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if s == nil {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}

	isMember, err := service.IsProjectMember(s.ProjectID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	//Judge authority
	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to get publicity of service.")
		return
	}

	// change the replica number of service

	res, err := service.ScaleReplica(*s, reqServiceScale.Replica)

	if res != true {
		logs.Info("Failed to scale service replica", s, reqServiceScale.Replica)
		p.internalError(err)
		return
	}
	logs.Info("Scale service replica successfully")
}

//get selectable service list
func (p *ServiceController) GetSelectableServicesAction() {
	serviceName := p.GetString("service_name", "")
	projectName := p.GetString("project_name", "")
	logs.Info("Get selectable service list for", projectName, serviceName)
	serviceList, err := service.GetSelectableServices(projectName, serviceName)
	if err != nil {
		logs.Error("Failed to get selectable services.")
		p.internalError(err)
		return
	}
	p.Data["json"] = serviceList
	p.ServeJSON()
}
