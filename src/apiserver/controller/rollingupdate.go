package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	yaml "yaml-2"

	"k8s.io/client-go/pkg/api/v1"

	"github.com/astaxie/beego/logs"
)

type ServiceRollingUpdateController struct {
	baseController
}

func (p *ServiceRollingUpdateController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
}

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceConfigAction() {
	projectName, serviceID := p.isServiceValid()
	serviceConfig := p.getServiceConfig(projectName, serviceID)

	if serviceConfig.Spec.Template == nil || len(serviceConfig.Spec.Template.Spec.Containers) < 1 {
		p.customAbort(http.StatusBadRequest, "Requested service's config is invalid.")
	}

	var imageList []model.ImageIndex
	for _, container := range serviceConfig.Spec.Template.Spec.Containers {
		indexProject := strings.IndexByte(container.Image, '/')
		indexImage := strings.LastIndexByte(container.Image, '/')
		indexTag := strings.LastIndexByte(container.Image, ':')
		imageList = append(imageList, model.ImageIndex{ImageName: container.Image[indexProject+1 : indexTag],
			ImageTag:    container.Image[indexTag+1:],
			ProjectName: container.Image[indexProject+1 : indexImage]})
	}

	p.Data["json"] = imageList
	p.ServeJSON()
}
func (p *ServiceRollingUpdateController) PostRollingUpdateServiceConfigAction() {
	projectName, serviceID := p.isServiceValid()

	var imageList []model.ImageIndex
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
	}
	err = json.Unmarshal(reqData, &imageList)
	if err != nil {
		p.internalError(err)
	}
	logs.Debug("Image list info: %+v\n", imageList)

	serviceConfig := p.getServiceConfig(projectName, serviceID)
	if len(serviceConfig.Spec.Template.Spec.Containers) != len(imageList) {
		p.customAbort(http.StatusConflict, "Image's config is invalid.")
	}

	var rollingUpdateConfig v1.ReplicationController
	for index, container := range serviceConfig.Spec.Template.Spec.Containers {
		image := registryBaseURI() + "/" + imageList[index].ImageName + ":" + imageList[index].ImageTag
		if serviceConfig.Spec.Template.Spec.Containers[index].Image != image {
			rollingUpdateConfig.Spec.Template.Spec.Containers = append(rollingUpdateConfig.Spec.Template.Spec.Containers, v1.Container{
				Name:  container.Name,
				Image: image,
			})
		}
	}

	deploymentAbsName := filepath.Join(repoPath(), projectName, serviceID, rollingUpdateFilename)
	err = service.GenerateYamlFile(deploymentAbsName, rollingUpdateConfig)
	if err != nil {
		p.internalError(err)
	}

	extras := filepath.Join("deployments", serviceConfig.Name)
	deployPushobject := p.rollingUpdateService(rollingUpdateFilename, serviceID, projectName, extras, rollingUpdate, deploymentAPI)
	ret, msg, err := InternalPushObjects(&deployPushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push deployment object: %d %s", ret, msg)

}

func (p *ServiceRollingUpdateController) isServiceValid() (string, string) {
	projectName := p.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		p.internalError(err)
	}
	if isExistence != true {
		p.customAbort(http.StatusBadRequest, "Project don't exist.")
	}

	serviceName := p.GetString("service_name")
	serviceID, err := getServiceID(serviceName, projectName)
	if err != nil {
		p.internalError(err)
	}
	if serviceID == "" {
		p.customAbort(http.StatusBadRequest, "Service name don't exist.")
	}
	return projectName, serviceID
}

func (p *ServiceRollingUpdateController) getServiceConfig(projectName string, serviceID string) *v1.ReplicationController {
	absFileName := filepath.Join(repoPath(), projectName, serviceID, deploymentFilename)
	logs.Info("User: %s get deployment.yaml images info from %s.", p.currentUser.Username, absFileName)

	yamlFile, err := ioutil.ReadFile(absFileName)
	if err != nil {
		p.internalError(err)
	}

	var serviceConfig v1.ReplicationController
	err = yaml.Unmarshal(yamlFile, &serviceConfig)
	if err != nil {
		p.internalError(err)
	}

	return &serviceConfig
}
func (p *ServiceRollingUpdateController) rollingUpdateService(fileName string, serviceID string, projectName string,
	extras string, jobName string, apiVersion string) pushObject {
	var pushobject pushObject
	pushobject.FileName = fileName
	pushobject.JobName = jobName
	pushobject.Value = filepath.Join(projectName, serviceID)
	pushobject.Message = fmt.Sprintf("Create %s for project %s service %s", extras,
		projectName, serviceID)
	pushobject.Extras = filepath.Join(kubeMasterURL(), apiVersion, projectName, extras)

	pushobject.Items = []string{filepath.Join(pushobject.Value, fileName)}
	logs.Info("pushobject.FileName:%+v\n", pushobject.FileName)
	logs.Info("pushobject.Value:%+v\n", pushobject.Value)
	logs.Info("pushobject.Extras:%+v\n", pushobject.Extras)
	return pushobject
}
