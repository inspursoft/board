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

type ServiceDeployController struct {
	baseController
}

//  Checking the user priviledge by token
func (p *ServiceDeployController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
}

func (p *ServiceDeployController) DeployServiceAction() {
	key := p.token
	configService := NewConfigServiceStep(key)

	var newservice model.ServiceStatus
	newservice.Name = configService.ServiceName
	newservice.ProjectID = configService.ProjectID
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.currentUser.ID
	newservice.OwnerName = p.currentUser.Username

	project, err := service.GetProject(model.Project{ID: configService.ProjectID}, "id")
	if err != nil {
		p.internalError(err)
		return
	}
	if project == nil {
		p.customAbort(http.StatusBadRequest, projectIDInvalidErr.Error())
		return
	}
	newservice.ProjectName = project.Name

	serviceID, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.internalError(err)
		return
	}

	loadPath := filepath.Join(repoPath(), project.Name, strconv.Itoa(int(serviceID)))
	err = service.CheckDeploymentPath(loadPath)
	if err != nil {
		p.internalError(err)
		return
	}

	err = service.AssembleDeploymentYaml((*model.ConfigServiceStep)(configService), loadPath)
	if err != nil {
		p.internalError(err)
		return
	}

	deployPushobject := assemblePushObject(deploymentFilename, serviceID, project.Name, "deployments")
	ret, msg, err := InternalPushObjects(&deployPushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push deployment object: %d %s", ret, msg)
	err = service.AssembleServiceYaml((*model.ConfigServiceStep)(configService), loadPath)
	if err != nil {
		p.internalError(err)
		return
	}

	servicePushobject := assemblePushObject(serviceFilename, serviceID, project.Name, "services")
	ret, msg, err = InternalPushObjects(&servicePushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push deployment object: %d %s", ret, msg)

	serviceConfig, err := json.Marshal(&configService)
	if err != nil {
		p.internalError(err)
		return
	}

	updateService := model.ServiceStatus{ID: serviceID, Status: running, ServiceConfig: string(serviceConfig)}
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
	logs.Info("Service with ID:%d has been deleted in cache.", serviceID)

	configService.ServiceID = serviceID
	p.Data["json"] = configService
	p.ServeJSON()
}

func assemblePushObject(fileName string, serviceID int64, projectName string, extras string) pushObject {
	var pushobject pushObject
	pushobject.FileName = fileName
	pushobject.JobName = serviceProcess
	pushobject.Value = filepath.Join(projectName, strconv.Itoa(int(serviceID)))
	pushobject.Message = fmt.Sprintf("Create %s for project %s service %d", extras,
		projectName, serviceID)
	if extras == "deployments" {
		pushobject.Extras = filepath.Join(kubeMasterURL(), deploymentAPI, projectName, extras)
	} else {
		pushobject.Extras = filepath.Join(kubeMasterURL(), serviceAPI, projectName, extras)
	}
	pushobject.Items = []string{filepath.Join(pushobject.Value, fileName)}
	logs.Info("pushobject.FileName:%+v\n", pushobject.FileName)
	logs.Info("pushobject.Value:%+v\n", pushobject.Value)
	logs.Info("pushobject.Extras:%+v\n", pushobject.Extras)
	return pushobject
}

func (p *ServiceDeployController) DeployServiceTestAction() {
	key := p.token
	configService := NewConfigServiceStep(key)
	configService.ServiceName = test + configService.ServiceName
	SetConfigServiceStep(key, configService)
	p.DeployServiceAction()
}
