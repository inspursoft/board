package nodegroups

import (
	"github.com/astaxie/beego"
)

// Operations about node groups
type CommonController struct {
	beego.Controller
}

// @Title List all node groups
// @Description List all for node groups.
// @Param	nodegroup_id	path	int	false	"ID of node group"
// @Param	search	query	string	false	"Query item for node groups"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:nodegroup_id [get]
func (c *CommonController) List() {

}

// @Title Add node group
// @Description Add node group.
// @Param	body	body 	models.nodegroups.vm.NodeGroup	true	"View model for node group."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}

// @Title Delete node group by ID
// @Description Delete node group by ID.
// @Param	nodegroup_id	path	int	true	"ID of node group"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:nodegroup_id [delete]
func (c *CommonController) Delete() {

}
