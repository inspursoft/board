package releases

import (
	"github.com/astaxie/beego"
)

// Operations about Helm releases
type CommonController struct {
	beego.Controller
}

// @Title List all Helm releases
// @Description List all for Helm releases
// @Param	search	query	string	false	"Query item for Helm repository"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (c *CommonController) List() {

}

// @Title Install chart to Helm releases
// @Description Install chart to Helm releases
// @Param	body	body	models.helm.vm.Release	true	"View model of Helm release"
// @Success 200 Successful installed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Post() {

}

// @Title Get Helm release detail
// @Description Get Helm repository detail
// @Param	release_id	path	int	true	"ID of Helm release"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:release_id [get]
func (c *CommonController) Get() {

}

// @Title Delete Helm release
// @Description Delete Helm release
// @Param	release_id	path	int	true	"ID of Helm release"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:release_id [delete]
func (c *CommonController) Delete() {

}
