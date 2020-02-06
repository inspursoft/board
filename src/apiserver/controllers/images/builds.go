package images

import (
	"github.com/astaxie/beego"
)

// Operations about images
type BuildController struct {
	beego.Controller
}

// @Title Create docker image by template
// @Description Create docker image by template.
// @Param	body	body 	models.images.vm.Image	true	"View model for image."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /by_template [post]
func (c *BuildController) ByTemplate() {

}

// @Title Create docker image by uploaded package
// @Description Create docker image by uploaded package.
// @Param	body	body 	models.images.vm.Image	true	"View model for image."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /by_uploaded_package [post]
func (c *BuildController) ByUploadedPackage() {

}

// @Title Create docker image by dockerfile
// @Description Create docker image by dockerfile.
// @Param	body	body 	models.images.vm.Image	true	"View model for image."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /by_dockerfile [post]
func (c *BuildController) ByDockerfile() {

}
