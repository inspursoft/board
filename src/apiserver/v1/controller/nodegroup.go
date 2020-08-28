package controller

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

type NodeGroupController struct {
	c.BaseController
}

func (n *NodeGroupController) GetNodeGroupsAction() {
	res, err := service.GetNodeGroupList()
	if err != nil {
		logs.Debug("Failed to get Node Group List")
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.RenderJSON(res)
}

func (n *NodeGroupController) AddNodeGroupAction() {
	var reqNodeGroup model.NodeGroup
	var err error
	err = n.ResolveBody(&reqNodeGroup)
	if err != nil {
		return
	}

	if !utils.ValidateWithLengthRange(reqNodeGroup.GroupName, 1, 63) {
		n.CustomAbortAudit(http.StatusBadRequest, "NodeGroup Name must be not empty and no more than 63 characters ")
		return
	}

	nodeGroupExists, err := service.NodeGroupExists(reqNodeGroup.GroupName)
	if err != nil {
		n.InternalError(err)
		return
	}
	if nodeGroupExists {
		n.CustomAbortAudit(http.StatusConflict, "Node Group name already exists.")
		return
	}

	reqNodeGroup.GroupName = strings.TrimSpace(reqNodeGroup.GroupName)
	reqNodeGroup.OwnerID = int64(n.CurrentUser.ID)

	group, err := service.CreateNodeGroup(reqNodeGroup)
	if err != nil {
		logs.Debug("Failed to add node group %s", reqNodeGroup.GroupName)
		n.InternalError(err)
		return
	}
	logs.Info("Added node group %s %d", reqNodeGroup.GroupName, group.ID)
}

func (n *NodeGroupController) CheckNodeGroupNameExistingAction() {
	nodeGroupName := n.GetString("nodegroup_name")
	isExists, err := service.NodeOrNodeGroupExists(nodeGroupName)
	if err != nil {
		n.InternalError(err)
		return
	}
	if isExists {
		n.CustomAbortAudit(http.StatusConflict, "This nodegroup name is already existing.")
		return
	}

	logs.Info("Group name of %s is available", nodeGroupName)
}

func (n *NodeGroupController) DeleteNodeGroupAction() {
	groupName := n.GetString("groupname")
	logs.Debug("Removing nodegroup %s", groupName)

	if groupName == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "NodeGroup Name should not null")
		return
	}

	nodeGroupExists, err := service.NodeGroupExists(groupName)
	if err != nil {
		n.InternalError(err)
		return
	}
	if !nodeGroupExists {
		n.CustomAbortAudit(http.StatusBadRequest, "Node Group name not exists.")
		return
	}

	err = service.RemoveNodeGroup(groupName)
	if err != nil {
		logs.Debug("Failed to remove nodegroup %s", groupName)
		n.InternalError(err)
		return
	}
	logs.Info("Removed nodegroup %s", groupName)
}
