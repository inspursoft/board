package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"path/filepath"
	"strconv"

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
	NA                     = "NA"
	NotFound               = "NotFound"
	serviceNamespace       = "default" //TODO create in project post
)

const (
	preparing = iota
	running
	stopped
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
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}

//Get request massage parameters
func handleReqPara(p *ServiceController) (model.ServiceConfig, int, error) {
	var err error
	var reqServiceConfig model.ServiceConfig

	//get the request data
	reqData, err := p.resolveBody()
	if err != nil {
		return reqServiceConfig, 0, err
	}

	//prase the request data to struct
	err = json.Unmarshal(reqData, &reqServiceConfig)
	if err != nil {
		return reqServiceConfig, 0, err
	}

	//Add registry to container images for deployment
	for index, container := range reqServiceConfig.DeploymentYaml.ContainerList {
		reqServiceConfig.DeploymentYaml.ContainerList[index].BaseImage =
			filepath.Join(registryBaseURI(), container.BaseImage)
	}

	logs.Debug("%+v", reqServiceConfig)

	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		return reqServiceConfig, 0, err
	}
	logs.Debug("To check serviceID existing %d", serviceID) //TODO

	return reqServiceConfig, serviceID, err
}

func handleTestReqPara(reqServiceConfig model.ServiceConfig) model.ServiceConfig {
	//alter parameter for test
	//totest use pointer transfer parameter
	reqServiceConfig.ServiceYaml.Name = test + reqServiceConfig.ServiceYaml.Name
	reqServiceConfig.ServiceYaml.Selectors[0] = test + reqServiceConfig.ServiceYaml.Selectors[0]
	reqServiceConfig.DeploymentYaml.Name = test + reqServiceConfig.DeploymentYaml.Name
	logs.Debug("%+v", reqServiceConfig)

	return reqServiceConfig
}

// API to deploy service
func deployServiceCommonAction(p *ServiceController, depFileName string, serFileName string,
	handleTestRepPara ...func(model.ServiceConfig) model.ServiceConfig) {
	var err error
	var pushobject pushObject

	//Get request massage parameters
	reqServiceConfig, serviceID, err := handleReqPara(p)
	if err != nil {
		p.internalError(err)
		return
	}

	// Check request parameters
	err = service.CheckReqPara(reqServiceConfig)
	if err != nil {
		p.customAbort(http.StatusBadRequest, err.Error())
		return
	}

	//set deployment path
	serviceConfigPath := filepath.Join(repoPath(),
		reqServiceConfig.ProjectName, strconv.Itoa(serviceID))
	logs.Debug("Service config path: %s", serviceConfigPath)
	service.SetDeploymentPath(serviceConfigPath)

	//Handle parameter for test
	if len(handleTestRepPara) > 0 {
		reqServiceConfig = handleTestRepPara[0](reqServiceConfig)
	}

	//Build deployment yaml file
	err = service.BuildDeploymentYaml(reqServiceConfig, depFileName)
	if err != nil {
		p.internalError(err)
		return
	}

	//Build service yaml file
	err = service.BuildServiceYaml(reqServiceConfig, serFileName)
	if err != nil {
		p.internalError(err)
		return
	}

	//serviceNamespace = reqServiceConfig.ProjectName TODO in project

	// Push deployment to jenkins
	pushobject.FileName = depFileName
	pushobject.JobName = serviceProcess
	pushobject.Value = filepath.Join(reqServiceConfig.ProjectName, strconv.Itoa(serviceID))
	pushobject.Message = fmt.Sprintf("Create deployment for project %s service %d",
		reqServiceConfig.ProjectName, reqServiceConfig.ServiceID)
	pushobject.Extras = filepath.Join(kubeMasterURL(), deploymentAPI, serviceNamespace, "deployments")

	// Add deployment file
	pushobject.Items = []string{filepath.Join(pushobject.Value, depFileName)}

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push deployment object: %d %s", ret, msg)

	//TODO: If fail to create deployment, should not continue to create service

	//Push service to jenkins
	pushobject.FileName = serFileName
	pushobject.Message = fmt.Sprintf("Create service for project %s service %d",
		reqServiceConfig.ProjectName, reqServiceConfig.ServiceID)
	pushobject.Extras = filepath.Join(kubeMasterURL(), serviceAPI, serviceNamespace, "services")

	// Add deployment file
	pushobject.Items = []string{filepath.Join(pushobject.Value, serFileName)}

	ret, msg, err = InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Debug("Internal push service object: %d %s", ret, msg)

	// Update service status in database
	updateService := model.ServiceStatus{ID: int64(serviceID), Status: running,
		Name: reqServiceConfig.ServiceYaml.Name, OwnerID: p.currentUser.ID}
	_, err = service.UpdateService(updateService, "name", "status", "owner_id")
	if err != nil {
		p.internalError(err)
		return
	}

	p.serveStatus(ret, msg)
}

// API to deploy service
func (p *ServiceController) DeployServiceAction() {
	//Judge authority
	if !(p.isSysAdmin && p.isProjectAdmin) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to manipulate user.")
		return
	}
	deployServiceCommonAction(p, deploymentFilename, serviceFilename)
}

// API to deploy test service
func (p *ServiceController) DeployServiceTestAction() {
	//Judge authority
	if !(p.isSysAdmin && p.isProjectAdmin) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to manipulate user.")
		return
	}
	deployServiceCommonAction(p, deploymentTestFilename, serviceTestFilename, handleTestReqPara)
}

//get service list
func (p *ServiceController) GetServiceListAction() {
	serviceList, err := service.GetServiceList()
	if err != nil {
		p.internalError(err)
		return
	}
	p.Data["json"] = serviceList
	p.ServeJSON()
}

// API to create service config
func (p *ServiceController) CreateServiceConfigAction() {
	//Judge authority
	if !(p.isSysAdmin && p.isProjectAdmin) {
		p.customAbort(http.StatusForbidden, "Insuffient privileges to manipulate user.")
		return
	}
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
	logs.Debug("%+v", reqServiceProject)
	//Assign and return Service ID with mysql
	var newservice model.ServiceStatus
	newservice.ProjectID = reqServiceProject.ProjectID
	newservice.ProjectName = reqServiceProject.ProjectName
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.currentUser.ID

	serviceID, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.internalError(err)
		return
	}
	p.Data["json"] = strconv.Itoa(int(serviceID))
	p.ServeJSON()
}

func (p *ServiceController) DeleteServiceAction() {
	if !(p.isSysAdmin && p.isProjectAdmin) {
		p.customAbort(http.StatusForbidden, "Insuffient privileges to manipulate user.")
		return
	}

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
	if !(p.isSysAdmin && p.isProjectAdmin) {
		p.customAbort(http.StatusForbidden, "Insuffient privileges to manipulate user.")
		return
	}
	var err error
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
		p.Data["json"] = NA
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
	var err error
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
	serviceConfigPath := filepath.Join(repoPath, s.ProjectName,
		strconv.Itoa(serviceID))
	logs.Debug("Service config path: %s", serviceConfigPath)

	// Get the service config by yaml files
	var serviceConfig model.ServiceConfig
	serviceConfig.ServiceID = s.ID
	serviceConfig.ProjectID = s.ProjectID
	serviceConfig.ProjectName = s.ProjectName
	//TODO
	// err = service.InterpretServiceConfigYaml(&serviceConfig, serviceConfigPath)
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
	reqServiceConfig, serviceID, err := handleReqPara(p)
	if err != nil {
		p.internalError(err)
		return
	}

	// Check request parameters
	err = service.CheckReqPara(reqServiceConfig)
	if err != nil {
		p.customAbort(http.StatusBadRequest, err.Error())
		return
	}

	serviceConfigPath := filepath.Join(repoPath, reqServiceConfig.ProjectName,
		strconv.Itoa(serviceID))
	// TODO update yaml files by service config
	// err = service.UpdateServiceConfigYaml(reqServiceConfig, serviceConfigPath)
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
	serviceConfigPath := filepath.Join(repoPath, s.ProjectName,
		strconv.Itoa(serviceID))
	logs.Debug("Service config path: %s", serviceConfigPath)

	// Delete yaml files
	// TODO
	// err = service.DeleteServiceConfigYaml(serviceConfigPath)
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
