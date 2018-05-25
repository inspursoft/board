package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type ServiceRollingUpdateController struct {
	BaseController
}

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceImageConfigAction() {
	serviceConfig := p.getServiceConfig()
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
	p.renderJSON(imageList)
}

func (p *ServiceRollingUpdateController) getServiceConfig() (deploymentConfig *v1beta1.Deployment) {
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)

	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.internalError(err)
		return
	}
	if serviceStatus == nil {
		p.customAbort(http.StatusBadRequest, "Service name doesn't exist.")
		return
	}

	cli, err := service.K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		p.internalError(err)
		return
	}

	d := apiSet.Deployments(projectName)
	deploymentConfig, err = d.Get(serviceName)
	if err != nil {
		logs.Error("Failed to get service info %+v\n", err)
		p.internalError(err)
		return
	}
	return
}

func (p *ServiceRollingUpdateController) PatchRollingUpdateServiceImageAction() {

	var imageList []model.ImageIndex
	p.resolveBody(&imageList)

	serviceConfig := p.getServiceConfig()
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
	serviceConfig := p.getServiceConfig()
	for key, value := range serviceConfig.Spec.Template.Spec.NodeSelector {
		if key == "kubernetes.io/hostname" {
			p.renderJSON(value)
		} else {
			p.renderJSON(key)
		}
	}
}

func (p *ServiceRollingUpdateController) PatchRollingUpdateServiceNodeGroupAction() {
	nodeGroup := p.GetString("node_selector")
	if nodeGroup == "" {
		p.customAbort(http.StatusBadRequest, "nodeGroup is empty.")
	}
	rollingUpdateConfig := p.getServiceConfig()
	nodeGroupExists, err := service.NodeGroupExists(nodeGroup)
	if err != nil {
		p.internalError(err)
		return
	}
	if nodeGroupExists {
		rollingUpdateConfig.Spec.Template.Spec.NodeSelector = map[string]string{nodeGroup: "true"}
	} else {
		rollingUpdateConfig.Spec.Template.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": nodeGroup}
	}
	p.PatchServiceAction(rollingUpdateConfig)
}

func (p *ServiceRollingUpdateController) PatchServiceAction(rollingUpdateConfig *v1beta1.Deployment) {
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)

	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.internalError(err)
		return
	}
	if serviceStatus.Status == uncompleted {
		logs.Debug("Service is uncompleted, cannot be updated %s\n", serviceName)
		p.customAbort(http.StatusMethodNotAllowed, "Service is in uncompleted")
		return
	}

	serviceRollConfig, err := json.Marshal(rollingUpdateConfig)
	if err != nil {
		logs.Debug("rollingUpdateConfig %+v\n", rollingUpdateConfig)
		p.internalError(err)
		return
	}

	cli, err := service.K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		p.internalError(err)
		return
	}

	d := apiSet.Deployments(projectName)
	patchType := api.StrategicMergePatchType
	deployData, err := d.Patch(serviceName, patchType, serviceRollConfig)
	if err != nil {
		logs.Error("Failed to update service %+v\n", err)
		p.internalError(err)
		return
	}
	logs.Debug("New updated deployment: %+v\n", deployData)
}
