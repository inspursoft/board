package nodegroups

import "github.com/astaxie/beego"

// Operation about supplementary for nodegroups
type SupplementController struct {
	beego.Controller
}

// @Title Check node group existing status by ID
// @Description Check node group existing status by ID.
// @Param	nodegroup_id	path	int	true	"ID of node group"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:nodegroup_id/existing [get]
func (s *SupplementController) Toggle() {

}
