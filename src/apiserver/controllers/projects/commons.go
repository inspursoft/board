package projects

import (
	"github.com/astaxie/beego"
)

// Operations about projects
type CommonController struct {
	beego.Controller
}

// @Title List all projects
// @Description List all for projects.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	search	query	string	false	"Query item for projects"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [get]
func (c *CommonController) List() {

}

// @Title Add project
// @Description Add project.
// @Param	body	body 	models.projects.vm.Project	true	"View model for projects."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}

// @Title Update project by ID
// @Description Update project by ID.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	body	body	models.projects.vm.Project	true	"View model for projects."
// @Param	action	query	string	true	"Option of update."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [put]
func (c *CommonController) Update() {

}

// @Title Delete project by ID
// @Description Delete project by ID.
// @Param	project_id	path	int	true	"ID of projects"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [delete]
func (c *CommonController) Delete() {

}
