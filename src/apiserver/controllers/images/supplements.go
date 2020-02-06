package images

import (
	"github.com/astaxie/beego"
)

// Operations about images supplements
type SupplementController struct {
	beego.Controller
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
// @Param	project_id	path	int	true	"ID of project"
// @Param 	body	body	models.images.vm.Image	true	"Request models for image."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/exists [post]
func (c *SupplementController) Exists() {

}
