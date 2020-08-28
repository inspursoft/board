package registries

import (
	"github.com/astaxie/beego"
)

// Operations about image registries
type CommonController struct {
	beego.Controller
}

// @Title List all images from registry by project or query by image name.
// @Description List all images from registries by project or query by image name.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	image_name	path	string	false	"Query with name for images by project"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:image_name [get]
func (c *CommonController) List() {

}

// @Title  Delete image from registry by project
// @Description Delete image from registry by project.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	image_name	path	string	true	"Image name to be deleted."
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:image_name [delete]
func (c *CommonController) Delete() {

}
