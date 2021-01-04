package images

import (
	"fmt"
	"net/http"
	"strings"

	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
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
// @router /:project_name/existing [post]
func (c *SupplementController) Exists() {
	projectName := strings.TrimSpace(c.Ctx.Input.Param(":project_name"))
	c.ResolveUserPrivilege(projectName)

	// Check this image:tag in system
	var image vm.Image
	err := c.ResolveBody(&image)
	if err != nil {
		return
	}
	existing, err := service.ExistRegistry(projectName, image.ImageName, image.ImageTag)
	if err != nil {
		c.InternalError(err)
		return
	}

	if existing {
		c.CustomAbortAudit(http.StatusConflict, fmt.Sprintf("This image %s:%s already existing.", image.ImageName, image.ImageTag))
		return
	}
}
