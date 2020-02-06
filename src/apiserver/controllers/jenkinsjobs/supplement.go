package jenkinsjobs

import (
	"github.com/astaxie/beego"
)

// Operations about Jenkins job supplements
type SupplementController struct {
	beego.Controller
}

// @Title Get build number of Jenkins job
// @Description Get build number of Jenkins job
// @Param	user_id	path	int	true	"ID of user"
// @Param 	build_number	path	int	true	"Build number of Jenkins job."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id/:build_number [get]
func (c *SupplementController) JobNumber() {

}

// @Title Get console logs of Jenkins job
// @Description Get console logs of Jenkins job
// @Param 	build_number	query	int	true	"Build number of Jenkins job."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /console [get]
func (c *SupplementController) Console() {

}

// @Title Stop Jenkins job
// @Description Stop Jenkins job
// @Param 	build_number	query	int	true	"Build number of Jenkins job."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /stop [get]
func (c *SupplementController) Stop() {

}
