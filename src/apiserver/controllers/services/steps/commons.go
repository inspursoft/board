package steps

import "github.com/astaxie/beego"

// Operations about service deployments
type CommonController struct {
	beego.Controller
}

// @Title Get service config with steps
// @Description Create service config with steps.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of service"
// @Param	step	query	string	true	"Step of service config."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [get]
func (c *CommonController) Get() {

}

// @Title Set service config with steps
// @Description Set service config with steps.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of service"
// @Param	body	body	models.services.vm.Step	"View model of service with step."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [put]
func (c *CommonController) Update() {

}

// @Title Delete service with step by ID
// @Description Delete service with step by ID.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [delete]
func (c *CommonController) Delete() {

}
