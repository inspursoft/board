package rollingupdates

import "github.com/astaxie/beego"

// Operations about services with node group
type NodeGroupController struct {
	beego.Controller
}

// @Title List rolling updates with node group for services
// @Description List rolling updates with node group for services.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Param	search	query	string	false	"Query item for services"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/nodegroups [get]
func (c *NodeGroupController) List() {

}

// @Title Patch rolling updates with node group for services
// @Description Patch rolling updates with node group for services.
// @Param	project_id	path	int	false	"ID of projects"
// @Param	service_id	path	int	false	"ID of services"
// @Param	body	body	models.services.rollingupdates.vm.NodeGroup	"View model of rolling update image."
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id/nodegroups [patch]
func (c *NodeGroupController) Patch() {

}
