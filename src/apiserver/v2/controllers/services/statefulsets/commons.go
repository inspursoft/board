package statefulsets

import "github.com/astaxie/beego"

// Operations about service deployments
type CommonController struct {
	beego.Controller
}

// @Title Create service with statefulsets
// @Description Create service with statefulsets.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of service"
// @Param	body	body	models.services.vm.Statefulset	"View model of service with statefulset."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [post]
func (c *CommonController) Add() {

}

// @Title Delete service with statefulset by ID
// @Description Delete service with statefulset by ID.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [delete]
func (c *CommonController) Delete() {

}
