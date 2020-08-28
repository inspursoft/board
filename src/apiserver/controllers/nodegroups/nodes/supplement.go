package nodes

import "github.com/astaxie/beego"

// Operation about supplementary for node
type SupplementController struct {
	beego.Controller
}

// @Title Toggle node status by ID
// @Description Toggle node status by ID.
// @Param	nodegroup_id	path	int	true	"ID of node group"
// @Param	node_id	path	int	true	"ID of node"
// @Param	status	query	string	true	"Status of node"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:nodegroup_id/:node_id/toggle [get]
func (s *SupplementController) Toggle() {

}
