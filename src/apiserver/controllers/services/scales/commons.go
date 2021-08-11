package scales

import "github.com/astaxie/beego"

// Operations about service scales
type CommonController struct {
	beego.Controller
}

// @Title Get service scales
// @Description Get service scales.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of service"
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [get]
func (c *CommonController) Get() {

}

// @Title Update service scale by ID
// @Description Update service scale by ID.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Param	body	body 	models.services.vm.Scale	true	"View model for service scale."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [put]
func (c *CommonController) Update() {

}
