package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions"

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

func (p *ServiceRollingUpdateController) generateRepoPathByProjectName(projectName string) string {
	return filepath.Join(baseRepoPath(), p.currentUser.Username, projectName)
}

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceConfigAction() {
	serviceConfig, _ := p.getServiceConfig()
	if len(serviceConfig.Spec.Template.Spec.Containers) < 1 {
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

func (p *ServiceRollingUpdateController) getServiceConfig() (*extensions.Deployment, string) {
	projectName := p.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		p.internalError(err)
	}
	if isExistence != true {
		p.customAbort(http.StatusBadRequest, "Project don't exist.")
	}

	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.internalError(err)
	}
	if serviceStatus == nil {
		p.customAbort(http.StatusBadRequest, "Service name don't exist.")
	}

	repoPath := p.generateRepoPathByProjectName(projectName)
	absFileName := filepath.Join(repoPath, serviceProcess, strconv.Itoa(int(serviceStatus.ID)), deploymentFilename)
	logs.Info("User: %s get deployment.yaml images info from %s.", p.currentUser.Username, absFileName)

	yamlFile, err := ioutil.ReadFile(absFileName)
	if err != nil {
		p.internalError(err)
	}

	var deploymentConfig extensions.Deployment
	err = yaml.Unmarshal(yamlFile, &deploymentConfig)
	if err != nil {
		p.internalError(err)
	}

	return &deploymentConfig, projectName
}

func (p *ServiceRollingUpdateController) PatchRollingUpdateServiceAction() {

	var imageList []model.ImageIndex
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
	}
	//logs.Debug("reqData %+v\n", string(reqData))
	err = json.Unmarshal(reqData, &imageList)
	if err != nil {
		p.internalError(err)
	}
	logs.Debug("Image list info: %+v\n", imageList)

	serviceConfig, projectName := p.getServiceConfig()
	if len(serviceConfig.Spec.Template.Spec.Containers) != len(imageList) {
		p.customAbort(http.StatusConflict, "Image's config is invalid.")
	}

	//Can not rollingupdate an uncompleted service
	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.internalError(err)
	}
	if serviceStatus.Status == uncompleted {
		logs.Debug("Service is uncompleted, cannot be updated %s\n", serviceName)
		p.customAbort(http.StatusMethodNotAllowed, "Service is in uncompleted")
	}

	var rollingUpdateConfig v1.ReplicationController
	rollingUpdateConfig.Spec.Template = &v1.PodTemplateSpec{}
	for index, container := range serviceConfig.Spec.Template.Spec.Containers {
		image := registryBaseURI() + "/" + imageList[index].ImageName + ":" + imageList[index].ImageTag
		if serviceConfig.Spec.Template.Spec.Containers[index].Image != image {
			rollingUpdateConfig.Spec.Template.Spec.Containers = append(rollingUpdateConfig.Spec.Template.Spec.Containers, v1.Container{
				Name:  container.Name,
				Image: image,
			})
		}
	}

	if len(rollingUpdateConfig.Spec.Template.Spec.Containers) == 0 {
		logs.Info("Nothing to be updated")
		return
	}

	serviceRollConfig, err := json.Marshal(rollingUpdateConfig)
	if err != nil {
		logs.Debug("rollingUpdateConfig %+v\n", rollingUpdateConfig)
		p.internalError(err)
	}
	//logs.Debug("Marshal serviceRollConfig %+v\n", string(serviceRollConfig))

	cli, err := service.K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		p.internalError(err)
	}

	d := apiSet.Deployments(projectName)
	patchType := api.StrategicMergePatchType
	deployData, err := d.Patch(serviceConfig.Name, patchType, serviceRollConfig)
	if err != nil {
		logs.Error("Failed to update service %+v\n", err)
		p.internalError(err)
	}
	logs.Debug("New updated deployment: %+v\n", deployData)

	//update deployment yaml file
	repoPath := p.generateRepoPathByProjectName(projectName)
	err = service.RollingUpdateDeploymentYaml(repoPath, deployData)
	if err != nil {
		logs.Error("Failed to update deployment yaml file:%+v\n", err)
		p.internalError(err)
	}

	serviceInfo, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.internalError(err)
	}
	if serviceInfo == nil {
		logs.Error("Failed to get service info by service name: %s and project name: %s", serviceName, projectName)
		p.internalError(err)
	}
	var pushObject pushObject
	pushObject.UserID = p.currentUser.ID
	pushObject.FileName = deploymentFilename
	pushObject.JobName = rollingUpdate
	pushObject.ProjectName = projectName

	pushObject.Message = fmt.Sprintf("Rolling update service for project %s with service ID %d", projectName, serviceInfo.ID)

	generateMetaConfiguration(&pushObject, repoPath)
	pushObject.Items = []string{"META.cfg", filepath.Join(serviceProcess, strconv.Itoa(int(serviceInfo.ID)), deploymentFilename)}

	statusCode, message, err := InternalPushObjects(&pushObject, &(p.baseController))
	if err != nil {
		p.internalError(err)
	}
	logs.Info("Internal pushed for rolling update with status: %d and message: %s.", statusCode, message)
}
