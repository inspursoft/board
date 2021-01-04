package nodegroups

import (
	//"encoding/json"

	//"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/models/nodegroups"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	//"github.com/astaxie/beego"

	//"io/ioutil"

	"github.com/astaxie/beego/logs"
)

// Operations about node groups
type CommonController struct {
	c.BaseController
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
func (n *CommonController) List() {

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
func (n *CommonController) Add() {
	var reqNodeGroup nodegroups.NodeGroup
	var moNodeGroup model.NodeGroup
	var err error

	err = n.ResolveBody(&reqNodeGroup)
	if err != nil {
		return
	}

	// //  TODO Use base.controller later
	// 	data, err := ioutil.ReadAll(c.Ctx.Request.Body)
	// 	if err != nil {
	// 		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	err = json.Unmarshal(data, &reqNodeGroup)
	// 	if err != nil {
	// 		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}

	if !utils.ValidateWithLengthRange(reqNodeGroup.GroupName, 1, 63) {
		n.CustomAbort(http.StatusBadRequest, "NodeGroup Name must be not empty and no more than 63 characters ")
		return
	}

	nodeGroupExists, err := service.NodeGroupExists(reqNodeGroup.GroupName)
	if err != nil {
		n.InternalError(err)
		return
	}
	if nodeGroupExists {
		n.CustomAbort(http.StatusConflict, "Node Group name already exists.")
		return
	}

	moNodeGroup.GroupName = strings.TrimSpace(reqNodeGroup.GroupName)
	moNodeGroup.OwnerID = int64(n.CurrentUser.ID)

	group, err := service.CreateNodeGroup(moNodeGroup)
	if err != nil {
		logs.Debug("Failed to add node group %s", moNodeGroup.GroupName)
		n.InternalError(err)
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
func (n *CommonController) Delete() {

}
