package autoscales

import "github.com/astaxie/beego"

// Operations about service auto-scale
type CommonController struct {
	beego.Controller
}

// @Title Get service auto-scale
// @Description Get service auto-scale.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of service"
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [get]
func (c *CommonController) List() {

}

// @Title Create service auto-scale
// @Description Create service auto-scale.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Param	body	body 	models.services.vm.AutoScale	true	"View model for service auto-scale."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [post]
func (c *CommonController) Add() {

}

// @Title Update service auto-scale by ID
// @Description Update service auto-scale by ID.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Param	body	body 	models.services.vm.AutoScale	true	"View model for service auto-scale."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [put]
func (c *CommonController) Update() {

}

// @Title Delete service auto-scale by ID
// @Description Delete service auto-scale by ID.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [delete]
func (c *CommonController) Delete() {

}
