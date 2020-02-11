package rollingupdates

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about services
type ImageController struct {
	beego.Controller
}

// @Title List rolling updates images for services
// @Description List rolling updates for services.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Param	search	query	string	false	"Query item for services"
// @Success 200 {body} Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/images [get]
func (p *ImageController) List() {

}

// @Title Patch rolling updates images for services
// @Description Patch rolling updates for services.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Param	body	body	models.Image	"View model of rolling update image."
// @Success 200 {object} Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /image [patch]
func (p *ImageController) Patch() {

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

func (p *ImageController) PatchServiceAction(rollingUpdateConfig *model.Deployment) {
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

func (p *ImageController) getServiceConfig() (deploymentConfig *model.Deployment, err error) {
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
