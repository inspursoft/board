package rollingupdates

import (
	"github.com/astaxie/beego"
	_ "git/inspursoft/board/src/apiserver/models"
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
func (c *ImageController) List() {

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
// @router /:project_id/:service_id/images [patch]
func (c *ImageController) Patch() {

}
