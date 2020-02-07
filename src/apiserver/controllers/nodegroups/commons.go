package nodegroups

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/models"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
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
// @Param       token   header  string  "Current available token"
// @Param	body	body 	models.nodegroups.NodeGroup	true	"View model for node group."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @Failure 409 Node group already exists.
// @Failure 500 Internal errors.
// @router / [post]
func (c *CommonController) Add() {
	var reqNodeGroup models.nodegroups.NodeGroup
	var err error
	err = n.resolveBody(&reqNodeGroup)
	if err != nil {
		return
	}

	if !utils.ValidateWithLengthRange(reqNodeGroup.GroupName, 1, 63) {
		n.customAbort(http.StatusBadRequest, "NodeGroup Name must be not empty and no more than 63 characters ")
		return
	}

	nodeGroupExists, err := service.NodeGroupExists(reqNodeGroup.GroupName)
	if err != nil {
		n.internalError(err)
		return
	}
	if nodeGroupExists {
		n.customAbort(http.StatusConflict, "Node Group name already exists.")
		return
	}

	reqNodeGroup.GroupName = strings.TrimSpace(reqNodeGroup.GroupName)
	reqNodeGroup.OwnerID = int64(n.currentUser.ID)

	group, err := service.CreateNodeGroup(reqNodeGroup)
	if err != nil {
		logs.Debug("Failed to add node group %s", reqNodeGroup.GroupName)
		n.internalError(err)
		return
	}
	logs.Info("Added node group %s %d", reqNodeGroup.GroupName, group.ID)
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
