package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

type NodeGroupController struct {
	baseController
}

func (n *NodeGroupController) Prepare() {
	user := n.getCurrentUser()
	if user == nil {
		n.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	n.currentUser = user
	n.isSysAdmin = (user.SystemAdmin == 1)
}

func (n *NodeGroupController) GetNodeGroupsAction() {
	res, err := service.GetNodeGroupList()
	if err != nil {
		logs.Debug("Failed to get Node Group List")
		n.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.Data["json"] = res
	n.ServeJSON()
}

func (n *NodeGroupController) AddNodeGroupAction() {
	reqData, err := n.resolveBody()
	if err != nil {
		n.internalError(err)
		return
	}
	var reqNodeGroup model.NodeGroup
	err = json.Unmarshal(reqData, &reqNodeGroup)
	if err != nil {
		n.internalError(err)
		return
	}

	if reqNodeGroup.GroupName == "" {
		n.customAbort(http.StatusBadRequest, "NodeGroup Name should not null")
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
