package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

type NodeGroupController struct {
	BaseController
}

func (n *NodeGroupController) GetNodeGroupsAction() {
	res, err := service.GetNodeGroupList()
	if err != nil {
		logs.Debug("Failed to get Node Group List")
		n.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.renderJSON(res)
}

func (n *NodeGroupController) AddNodeGroupAction() {
	var reqNodeGroup model.NodeGroup
	var err error
	err = n.resolveBody(&reqNodeGroup)
	if err != nil {
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

func (n *NodeGroupController) CheckNodeGroupNameExistingAction() {
	nodeGroupName := n.GetString("nodegroup_name")
	isExists, err := service.NodeOrNodeGroupExists(nodeGroupName)
	if err != nil {
		n.internalError(err)
		return
	}
	if isExists {
		n.customAbort(http.StatusConflict, "This nodegroup name is already existing.")
		return
	}

	logs.Info("Group name of %s is available", nodeGroupName)
}

func (n *NodeGroupController) DeleteNodeGroupAction() {
	groupName := n.GetString("groupname")
	logs.Debug("Removing nodegroup %s", groupName)

	if groupName == "" {
		n.customAbort(http.StatusBadRequest, "NodeGroup Name should not null")
		return
	}

	nodeGroupExists, err := service.NodeGroupExists(groupName)
	if err != nil {
		n.internalError(err)
		return
	}
	if !nodeGroupExists {
		n.customAbort(http.StatusBadRequest, "Node Group name not exists.")
		return
	}

	err = service.RemoveNodeGroup(groupName)
	if err != nil {
		logs.Debug("Failed to remove nodegroup %s", groupName)
		n.internalError(err)
		return
	}
	logs.Info("Removed nodegroup %s", groupName)
}
