package services

import "github.com/astaxie/beego"

// Operations about services
type CommonController struct {
	beego.Controller
}

// @Title List all services
// @Description List all for services.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Param	search	query	string	false	"Query item for services"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [get]
func (c *CommonController) List() {

}

// @Title Create service
// @Description Create project.
// @Param	body	body 	models.services.vm.Service	true	"View model for service."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}

// @Title Update service by ID
// @Description Update service by ID.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Param	body	body	models.services.vm.Service	true	"View model for service."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [put]
func (c *CommonController) Update() {

}

// @Title Delete service by ID
// @Description Delete service by ID.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [delete]
func (c *CommonController) Delete() {

}
