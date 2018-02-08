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

func (p *ServiceDeployController) generateRepoPathByProject(project *model.Project) string {
	if project == nil {
		p.customAbort(http.StatusBadRequest, "Failed to generate repo path since project is nil.")
	}
	return filepath.Join(baseRepoPath(), p.currentUser.Username, project.Name)
}

func (p *ServiceDeployController) getKey() string {
	return strconv.Itoa(int(p.currentUser.ID))
}

func (p *ServiceDeployController) DeployServiceAction() {
	key := p.getKey()
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

	serviceInfo, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.internalError(err)
		return
	}

	repoPath := p.generateRepoPathByProject(project)
	loadPath := filepath.Join(repoPath, serviceProcess, strconv.Itoa(int(serviceInfo.ID)))
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

	err = service.AssembleServiceYaml((*model.ConfigServiceStep)(configService), loadPath)
	if err != nil {
		p.internalError(err)
		return
	}

	var pushObject pushObject
	pushObject.UserID = p.currentUser.ID
	pushObject.FileName = fmt.Sprintf("%s,%s", deploymentFilename, serviceFilename)
	pushObject.JobName = serviceProcess
	pushObject.ProjectName = project.Name
	pushObject.Extras = fmt.Sprintf("%s,%s", fmt.Sprintf("%s%s%s/%s", kubeMasterURL(), deploymentAPI, project.Name, "deployments"),
		fmt.Sprintf("%s%s%s/%s", kubeMasterURL(), serviceAPI, project.Name, "services"))
	pushObject.Value = filepath.Join(serviceProcess, strconv.Itoa(int(serviceInfo.ID)))
	pushObject.Message = fmt.Sprintf("Create service for project %s with service %d", project.Name, serviceInfo.ID)

	relPath := filepath.Join(serviceProcess, strconv.Itoa(int(serviceInfo.ID)))

	generateMetaConfiguration(&pushObject, repoPath)
	pushObject.Items = []string{"META.cfg", filepath.Join(relPath, deploymentFilename), filepath.Join(relPath, serviceFilename)}

	ret, msg, err := InternalPushObjects(&pushObject, &(p.baseController))
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
	p.Data["json"] = configService
	p.ServeJSON()
}

func (p *ServiceDeployController) DeployServiceTestAction() {
	key := p.getKey()
	configService := NewConfigServiceStep(key)
	configService.ServiceName = test + configService.ServiceName
	SetConfigServiceStep(key, configService)
	p.DeployServiceAction()
}
