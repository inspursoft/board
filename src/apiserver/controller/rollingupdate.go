package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/logs"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
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

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceImageConfigAction() {
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

func (p *ServiceRollingUpdateController) getServiceConfig() (*v1beta1.Deployment, string) {
	projectName := p.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		p.internalError(err)
	}
	if isExistence != true {
		p.customAbort(http.StatusBadRequest, "Project doesn't exist.")
	}

	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.internalError(err)
	}
	if serviceStatus == nil {
		p.customAbort(http.StatusBadRequest, "Service name doesn't exist.")
	}

	cli, err := service.K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		p.internalError(err)
	}

	d := apiSet.Deployments(projectName)
	deploymentConfig, err := d.Get(serviceName)
	if err != nil {
		logs.Error("Failed to get service info %+v\n", err)
		p.internalError(err)
	}

	return deploymentConfig, projectName
}

func (p *ServiceRollingUpdateController) PatchRollingUpdateServiceImageAction() {

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

	serviceConfig, _ := p.getServiceConfig()
	if len(serviceConfig.Spec.Template.Spec.Containers) != len(imageList) {
		p.customAbort(http.StatusConflict, "Image's config is invalid.")
	}

	var rollingUpdateConfig v1beta1.Deployment
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

	p.PatchServiceAction(&rollingUpdateConfig)
}

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceNodeGroupConfigAction() {
	serviceConfig, _ := p.getServiceConfig()
	for key, value := range serviceConfig.Spec.Template.Spec.NodeSelector {
		if key == "kubernetes.io/hostname" {
			p.Data["json"] = value
		} else {
			p.Data["json"] = key
		}
	}

	p.ServeJSON()
}

func (p *ServiceRollingUpdateController) PatchRollingUpdateServiceNodeGroupAction() {
	nodeGroup := p.GetString("node_selector")
	if nodeGroup == "" {
		p.customAbort(http.StatusBadRequest, "nodeGroup is empty.")
	}
	rollingUpdateConfig, _ := p.getServiceConfig()
	nodeGroupExists, _ := service.NodeGroupExists(nodeGroup)
	if nodeGroupExists {
		rollingUpdateConfig.Spec.Template.Spec.NodeSelector = map[string]string{nodeGroup: "true"}
	} else {
		rollingUpdateConfig.Spec.Template.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": nodeGroup}
	}

	p.PatchServiceAction(rollingUpdateConfig)
}

func (p *ServiceRollingUpdateController) PatchServiceAction(rollingUpdateConfig *v1beta1.Deployment) {
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
	if serviceStatus.Status == uncompleted {
		logs.Debug("Service is uncompleted, cannot be updated %s\n", serviceName)
		p.customAbort(http.StatusMethodNotAllowed, "Service is in uncompleted")
	}

	serviceRollConfig, err := json.Marshal(rollingUpdateConfig)
	if err != nil {
		logs.Debug("rollingUpdateConfig %+v\n", rollingUpdateConfig)
		p.internalError(err)
	}

	cli, err := service.K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		p.internalError(err)
	}

	d := apiSet.Deployments(projectName)
	patchType := api.StrategicMergePatchType
	deployData, err := d.Patch(serviceName, patchType, serviceRollConfig)
	if err != nil {
		logs.Error("Failed to update service %+v\n", err)
		p.internalError(err)
	}
	logs.Debug("New updated deployment: %+v\n", deployData)
}
