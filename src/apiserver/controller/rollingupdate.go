package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/logs"
)

type ServiceRollingUpdateController struct {
	BaseController
}

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceImageConfigAction() {
	serviceConfig, err := p.getServiceConfig()
	if err != nil {
		return
	}
	if len(serviceConfig.Spec.Template.Spec.Containers) < 1 {
		p.customAbort(http.StatusBadRequest, "Requested service's config is invalid.")
		return
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

func (p *ServiceRollingUpdateController) getServiceConfig() (deploymentConfig *model.Deployment, err error) {
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

	deploymentConfig, _, err = service.GetDeployment(projectName, serviceName)
	if err != nil {
		logs.Error("Failed to get service info %+v\n", err)
		p.parseError(err, parseGetK8sError)
		return
	}
	return
}
func (s *ServiceRollingUpdateController) GetServiceSessionFlagAction() {
	serviceName := s.GetString("service_name")
	projectName := s.GetString("project_name")
	s.resolveProjectMember(projectName)
	svc, err := service.GetK8sService(projectName, serviceName)
	if err != nil {
		s.internalError(err)
		return
	}
	s.renderJSON(svc)
}

func (s *ServiceRollingUpdateController) PatchServiceSessionAction() {
	sessionAffinityFlag, err := s.GetInt("session_affinity_flag", 0)
	if err != nil {
		s.internalError(err)
		return
	}
	serviceName := s.GetString("service_name")
	projectName := s.GetString("project_name")
	s.resolveProjectMember(projectName)
	svc, err := service.GetK8sService(projectName, serviceName)
	if err != nil {
		s.internalError(err)
		return
	}
	svc.SessionAffinityFlag = sessionAffinityFlag
	_, svcFileInfo, err := service.PatchK8sService(projectName, serviceName, svc)
	if err != nil {
		s.internalError(err)
		return
	}

	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		s.internalError(err)
		return
	}
	updateService := model.ServiceStatus{ID: serviceStatus.ID, ServiceYaml: string(svcFileInfo)}
	_, err = service.UpdateService(updateService, "service_yaml")
	if err != nil {
		s.internalError(err)
		return
	}
	logs.Debug("Update service Successful.And the services config:%+v\n", updateService)

	s.resolveRepoServicePath(projectName, serviceName)
	err = utils.GenerateFile(svcFileInfo, s.repoServicePath, serviceFilename)
	if err != nil {
		s.internalError(err)
		return
	}
	s.pushItemsToRepo(filepath.Join(serviceName, serviceFilename))
}

func (p *ServiceRollingUpdateController) PatchRollingUpdateServiceImageAction() {

	var imageList []model.ImageIndex
	err := p.resolveBody(&imageList)
	if err != nil {
		return
	}

	serviceConfig, err := p.getServiceConfig()
	if err != nil {
		return
	}
	if len(serviceConfig.Spec.Template.Spec.Containers) != len(imageList) {
		p.customAbort(http.StatusConflict, "Image's config is invalid.")
	}

	//var rollingUpdateConfig model.Deployment
	var rollingUpdateConfig model.PodSpec
	for index, container := range serviceConfig.Spec.Template.Spec.Containers {
		image := registryBaseURI() + "/" + imageList[index].ImageName + ":" + imageList[index].ImageTag
		if serviceConfig.Spec.Template.Spec.Containers[index].Image != image {
			rollingUpdateConfig.Containers = append(rollingUpdateConfig.Containers, model.K8sContainer{
				Name:  container.Name,
				Image: image,
			})
		}
	}
	if len(rollingUpdateConfig.Containers) == 0 {
		logs.Info("Nothing to be updated")
		return
	}
	serviceConfig.Spec.Template.Spec = rollingUpdateConfig
	p.PatchServiceAction(serviceConfig)
}

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceNodeGroupConfigAction() {
	serviceConfig, err := p.getServiceConfig()
	if err != nil {
		return
	}
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
		return
	}
	rollingUpdateConfig, err := p.getServiceConfig()
	if err != nil {
		return
	}
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
	logs.Debug("Action updating nodeselector: ", rollingUpdateConfig)
	p.UpdateServiceAction(rollingUpdateConfig)
}

func (p *ServiceRollingUpdateController) PatchServiceAction(rollingUpdateConfig *model.Deployment) {
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

	deploymentConfig, deploymentFileInfo, err := service.PatchDeployment(projectName, serviceName, rollingUpdateConfig)
	if err != nil {
		logs.Error("Failed to get service info %+v\n", err)
		p.parseError(err, parsePostK8sError)
		return
	}

	p.resolveRepoServicePath(projectName, serviceName)
	err = utils.GenerateFile(deploymentFileInfo, p.repoServicePath, deploymentFilename)
	if err != nil {
		p.internalError(err)
		return
	}
	p.pushItemsToRepo(filepath.Join(serviceName, deploymentFilename))

	logs.Debug("New updated deployment: %+v\n", deploymentConfig)
}

// a temp fix, need to refactor
func (p *ServiceRollingUpdateController) UpdateServiceAction(rollingUpdateConfig *model.Deployment) {
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)

	serviceName := p.GetString("service_name")
	deploymentConfig, deploymentFileInfo, err := service.UpdateDeployment(projectName, serviceName, rollingUpdateConfig)
	if err != nil {
		logs.Error("Failed to get service info %+v\n", err)
		p.parseError(err, parsePostK8sError)
		return
	}
	logs.Debug("Updated deployment: ", deploymentConfig)
	p.resolveRepoServicePath(projectName, serviceName)
	err = utils.GenerateFile(deploymentFileInfo, p.repoServicePath, deploymentFilename)
	if err != nil {
		p.internalError(err)
		return
	}
	p.pushItemsToRepo(filepath.Join(serviceName, deploymentFilename))

	logs.Debug("New updated deployment: %+v\n", deploymentConfig)
}
