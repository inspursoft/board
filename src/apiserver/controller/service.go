package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego/logs"
)

const deploymentFilename = "deployment.yaml"
const serviceFilename = "service.yaml"
const serviceProcess = "process_service"
const apiheader = "Content-Type: application/yaml"
const deploymentAPI = "/apis/extensions/v1beta1/namespaces/"
const serviceAPI = "/api/v1/namespaces/"

var KubeMasterStatus bool
var serviceNamespace = "default" //TODO create in project post
var registryprefix = os.Getenv("REGISTRY_HOST") + ":" + os.Getenv("REGISTRY_PORT")
var KubeMasterUrl = os.Getenv("KUBEMASTER_IP") + ":" + os.Getenv("KUBEMASTER_PORT")

type ServiceController struct {
	baseController
}

//  Checking the user priviledge by token
func (p *ServiceController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}

// API to deploy service
func (p *ServiceController) DeployServiceAction() {
	var err error
	var reqServiceConfig model.ServiceConfig
	var pushobject pushObject

	//Judge authority
	if p.isSysAdmin == false && p.isProjectAdmin == false {
		p.CustomAbort(http.StatusForbidden, "Insuffient privileges to manipulate user.")
		return
	}

	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("To check serviceID existing %s", serviceID) //TODO

	//get the request data
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}

	//prase the request data to struct
	err = json.Unmarshal(reqData, &reqServiceConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	// Check deployment parameters
	err = service.CheckDeploymentYamlPara(reqServiceConfig)
	if err != nil {
		p.CustomAbort(http.StatusBadRequest, err.Error())
		return
	}

	// Check service parameters
	err = service.CheckServiceYamlPara(reqServiceConfig)
	if err != nil {
		p.CustomAbort(http.StatusBadRequest, err.Error())
		return
	}

	//set deployment path
	serviceId := int(reqServiceConfig.ServiceID)
	serviceConfigPath := filepath.Join(repoPath,
		reqServiceConfig.ProjectName, strconv.Itoa(serviceId))
	logs.Debug("Service config path: %s", serviceConfigPath)
	service.SetDeploymentPath(serviceConfigPath)

	//Add registry to container images for deployment
	for index, container := range reqServiceConfig.DeploymentYaml.ContainerList {
		reqServiceConfig.DeploymentYaml.ContainerList[index].BaseImage =
			filepath.Join(registryprefix, container.BaseImage)
	}
	logs.Info(reqServiceConfig)

	//Build deployment yaml file
	err = service.BuildDeploymentYaml(reqServiceConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	//Build service yaml file
	err = service.BuildServiceYaml(reqServiceConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	//serviceNamespace = reqServiceConfig.ProjectName TODO in project

	// Push deployment to jenkins
	pushobject.FileName = deploymentFilename
	pushobject.JobName = serviceProcess
	pushobject.Value = filepath.Join(reqServiceConfig.ProjectName, strconv.Itoa(serviceId))
	pushobject.Message = fmt.Sprintf("Create deployment for project %s service %d",
		reqServiceConfig.ProjectName, reqServiceConfig.ServiceID)
	pushobject.Extras = filepath.Join(KubeMasterUrl, deploymentAPI, serviceNamespace, "/deployments")

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
		reqServiceConfig.ProjectName, reqServiceConfig.ServiceID)
	pushobject.Extras = filepath.Join(KubeMasterUrl, serviceAPI, serviceNamespace, "/services")

	// Add deployment file
	pushobject.Items = []string{filepath.Join(pushobject.Value, serviceFilename)}

	ret, msg, err = InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push service object: %d %s", ret, msg)

	// Update service status in database
	updateService := model.ServiceStatus{ID: int64(serviceID), Status: 1,
		Name: reqServiceConfig.ServiceYaml.Name}
	_, err = service.UpdateService(updateService, "name", "status")
	if err != nil {
		p.internalError(err)
		return
	}

	p.CustomAbort(ret, msg)
}

// API to create service config
func (p *ServiceController) CreateServiceConfigAction() {
	var err error
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
	logs.Info(reqServiceProject)
	//Assign and return Service ID with mysql
	var newservice model.ServiceStatus
	newservice.ProjectID = reqServiceProject.ProjectID
	newservice.ProjectName = reqServiceProject.ProjectName
	newservice.Status = 0 // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.currentUser.ID

	serviceID, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.internalError(err)
		return
	}
	p.Data["json"] = strconv.FormatInt(serviceID, 10)
	p.ServeJSON()
}
