package nodes

import (
	"github.com/astaxie/beego"
)

// Operations about nodes
type CommonController struct {
	beego.Controller
}

// @Title List all nodes
// @Description List all for nodes.
// @Param	nodegroup_id	path	int	false	"ID of persistence"
// @Param	node_id	path	int	false	"ID of persistence"
// @Param	search	query	string	false	"Query item for nodes"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:nodegroup_id/:node_id [get]
func (c *CommonController) List() {

}

// @Title Add node
// @Description Add node.
// @Param	nodegroup_id	path	int	true	"ID of node group"
// @Param	body	body 	models.nodes.vm.Node	true	"View model for node."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}

// @Title Delete node by ID
// @Description Delete node by ID.
// @Param	nodegroup_id	path	int	true	"ID of node group"
// @Param	node_id	path	int	true	"ID of node"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:nodegroup_id/:node_id [delete]
func (c *CommonController) Delete() {

}
