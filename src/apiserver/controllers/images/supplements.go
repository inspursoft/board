package images

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"git/inspursoft/board/src/apiserver/models/images/vm"
	"git/inspursoft/board/src/apiserver/service"
	c "git/inspursoft/board/src/common/controller"

	"github.com/astaxie/beego/logs"
)

// Operations about images supplements
type SupplementController struct {
	c.BaseController
}

// @Title Preview dockerfile
// @Description Preview dockerfile.
// @Param	project_id	path	int	true	"ID of project"
// @Param 	body	body	models.images.vm.Image	true	"Request models for image."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/preview [post]
func (c *SupplementController) Preview() {

}

// @Title Clean config about building image
// @Description Clean config about building image
// @Param	project_id	path	int	true	"ID of project"
// @Param 	body	body	models.images.vm.Image	true	"Request models for image."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/clean_config [post]
func (c *SupplementController) CleanConfig() {

}

// @Title Check existing about building image
// @Description Check existing about building image
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_name	path	string	true	"Name of project"
// @Param 	body	body	vm.Image	true	"Request models for image."
// @Success 200 Image not existing.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @Failure 409 Image existing.
// @router /:project_name/exists [post]
func (c *SupplementController) Exists() {
	projectName := strings.TrimSpace(c.Ctx.Input.Param(":project_name"))
	c.ResolveUserPrivilege(projectName)

	// Check this image:tag in system
	var image vm.Image
	err := c.ResolveBody(&image)
	if err != nil {
		return
	}
	existing, err := existRegistry(projectName, image.Name, image.Tag)
	if err != nil {
		c.InternalError(err)
		return
	}

	if existing {
		c.CustomAbortAudit(http.StatusConflict, "This image:tag already existing.")
		return
	}
	logs.Debug("checking image:tag result %t", existing)
}

func existRegistry(projectName string, imageName string, imageTag string) (bool, error) {
	currentName := filepath.Join(projectName, imageName)
	fmt.Println(currentName)
	//check image
	repoList, err := service.GetRegistryCatalog()
	if err != nil {
		logs.Error("Failed to unmarshal repoList body %+v", err)
		return false, err
	}
	for _, imageRegistry := range repoList.Names {
		if imageRegistry == currentName {
			//check tag
			tagList, err := service.GetRegistryImageTags(currentName)
			if err != nil {
				logs.Error("Failed to unmarshal body %+v", err)
				return false, err
			}
			for _, tagID := range tagList.Tags {
				if imageTag == tagID {
					logs.Info("Image tag existing %s:%s", currentName, tagID)
					return true, nil
				}
			}
		}
	}
	return false, err
}
