package rollingupdates

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	c "git/inspursoft/board/src/common/controller"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego/logs"
)

// Operations about services
type ImageController struct {
	c.BaseController
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
// @Param	body	body	models.ImageIndex	"View model of rolling update image."
// @Success 200 {object} Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/images [patch]
func (p *ImageController) PatchRollingUpdateServiceImageAction() {
	var imageList []model.ImageIndex
	err := p.ResolveBody(&imageList)
	if err != nil {
		return
	}

	serviceConfig, err := p.getServiceConfig()
	if err != nil {
		return
	}
	if len(serviceConfig.Spec.Template.Spec.Containers) != len(imageList) {
		p.CustomAbortAudit(http.StatusConflict, "Image's config is invalid.")
	}

	//var rollingUpdateConfig model.Deployment
	var rollingUpdateConfig model.PodSpec
	for index, container := range serviceConfig.Spec.Template.Spec.Containers {
		image := c.RegistryBaseURI() + "/" + imageList[index].ImageName + ":" + imageList[index].ImageTag
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
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":project_id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	project, err := service.GetProjectByID(int64(projectID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if project == nil {
		p.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("No project was found with provided ID: %d", projectID))
		return
	}
	projectName := project.Name
	//	projectName := p.GetString("project_name")
	p.ResolveProjectMember(projectName)

	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":service_id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	s, err := service.GetServiceByID(int64(serviceID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if s == nil {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	serviceName := s.Name
	//	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.InternalError(err)
		return
	}
	if serviceStatus.Status == 3 {
		logs.Debug("Service is uncompleted, cannot be updated %s\n", serviceName)
		p.CustomAbortAudit(http.StatusMethodNotAllowed, "Service is in uncompleted")
		return
	}

	deploymentConfig, deploymentFileInfo, err := service.PatchDeployment(projectName, serviceName, rollingUpdateConfig)
	if err != nil {
		logs.Error("Failed to get service info %+v\n", err)
		p.ParseError(err, c.ParsePostK8sError)
		return
	}

	p.ResolveRepoServicePath(projectName, serviceName)
	err = utils.GenerateFile(deploymentFileInfo, p.RepoServicePath, deploymentFilename)
	if err != nil {
		p.InternalError(err)
		return
	}
	p.PushItemsToRepo(filepath.Join(serviceName, deploymentFilename))

	logs.Debug("New updated deployment: %+v\n", deploymentConfig)
}

func (p *ImageController) getServiceConfig() (deploymentConfig *model.Deployment, err error) {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":project_id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	project, err := service.GetProjectByID(int64(projectID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if project == nil {
		p.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("No project was found with provided ID: %d", projectID))
		return
	}
	projectName := project.Name
	//	projectName := p.GetString("project_name")
	p.ResolveProjectMember(projectName)

	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":service_id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	s, err := service.GetServiceByID(int64(serviceID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if s == nil {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	serviceName := s.Name
	//	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.InternalError(err)
		return
	}
	if serviceStatus == nil {
		p.CustomAbortAudit(http.StatusBadRequest, "Service name doesn't exist.")
		return
	}

	deploymentConfig, _, err = service.GetDeployment(projectName, serviceName)
	if err != nil {
		logs.Error("Failed to get service info %+v\n", err)
		p.ParseError(err, c.ParseGetK8sError)
		return
	}
	return
}
