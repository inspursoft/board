package members

import (
	"github.com/astaxie/beego"
)

// Operations about project members
type CommonController struct {
	beego.Controller
}

// @Title List all members of project
// @Description List all for projects.
// @Param	project_id	path	int	true	"ID of projects"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [get]
func (c *CommonController) List() {

}

// @Title Add member to the project
// @Description Add member to the project.
// @Param	project_id	path	int	true	"ID of projects"
// @Param	body	body 	models.projects.members.vm.Member	true	"View model for project member."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [post]
func (c *CommonController) Add() {

}

// @Title Update project by ID
// @Description Update project by ID.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	body	body	models.projects.members.vm.Member	true	"View model for projects."
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
// @Param	project_id	path	int	false	"ID of projects"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [delete]
func (c *CommonController) Delete() {

}
