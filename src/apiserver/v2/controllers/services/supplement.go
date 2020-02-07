package services

import ( 
  "github.com/astaxie/beego"
  _ "git/inspursoft/board/src/apiserver/v2/models"
)
// Operations about supplement of services
type SupplementController struct {
	beego.Controller
}

// @Title Update the
// @Description Get status of service by ID.
// @Param       project_name      path    string     true   "Name of project"
// @Param       service_name      path    string     true   "Name of service"
// @Param       phase             qurey   string     false  "Switch the service status, only for running or stopped"
// @Success 200 {object} models.Service
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @Failure 404 Not found.
// @Failure 500 Unexpected internal errors.
// @router /:project_name/:service_name/toggle [put]
func (c *SupplementController) SwitchSerivceStatus() {

}


// @Title Get status of service by ID
// @Description Get status of service by ID.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/status [get]
func (c *SupplementController) Status() {

}

// @Title Get selectable service
// @Description Get selectable service.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/selectable [get]
func (c *SupplementController) Seletable() {

}

// @Title Get info of service
// @Description Get info of service.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/info [get]
func (c *SupplementController) Info() {

}

// @Title Get URL route of service
// @Description Get URL route of service.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Param	body	body	models.services.vm.Route	true	"Route of service deployment."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/route [post]
func (c *SupplementController) Route() {

}

// @Title Test deployment of service
// @Description Test deployment of service.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	service_id	path	int	true	"ID of services"
// @Param	body	body	models.services.vm.Test	true	"Test options of services"
// @Success 200 Successful tested.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/test [post]
func (c *SupplementController) Test() {

}

// @Title Toggle status of service
// @Description Toggle status of service.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	service_id	path	int	true	"ID of services"
// @Param	status	query	string	true	"Status of service"
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/toggle [put]
func (c *SupplementController) Toggle() {

}

// @Title Change publicity of service
// @Description Change publicity of service.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	service_id	path	int	true	"ID of services"
// @Param	publicity	query	string	true	"Publicity of service"
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/publicity [put]
func (c *SupplementController) Publicity() {

}
