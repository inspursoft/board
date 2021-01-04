package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
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

// Get nodegroup detail
func (n *NodeGroupController) GetNodeGroupAction() {
	var nodegroupdetail model.NodeGroupDetail
	var nodegroup *model.NodeGroup
	var err error
	groupName := n.GetString("groupname")
	logs.Debug("Get nodegroup %s", groupName)

	if groupName == "" {
		nodegroupID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
		if err != nil {
			n.InternalError(err)
			return
		}
		nodegroup, err = service.GetNodeGroup(model.NodeGroup{ID: int64(nodegroupID)}, "id")
		groupName = nodegroup.GroupName
	} else {
		nodegroup, err = service.GetNodeGroup(model.NodeGroup{GroupName: groupName}, "name")
	}
	if err != nil {
		n.InternalError(err)
		return
	}
	if nodegroup == nil && nodegroup.ID == 0 {
		n.CustomAbortAudit(http.StatusBadRequest, "Node Group name not exists.")
		return
	}
	nodegroupdetail.NodeGroup = *nodegroup

	// Get the node list of the nodegroup
	nodelist, err := service.GetNodeListbyLabel(groupName)
	if err != nil {
		logs.Error("Failed to get nodelist: %s, error: %+v", groupName, err)
		n.InternalError(err)
		return
	}
	for _, nodeitem := range nodelist.Items {
		nodegroupdetail.NodeList = append(nodegroupdetail.NodeList, nodeitem.Name)
	}
	logs.Debug("Get nodegroup %v", nodegroupdetail)
	n.RenderJSON(nodegroupdetail)
}

// Update node group comments
func (n *NodeGroupController) UpdateNodeGroupAction() {
	nodegroupID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
	if err != nil {
		n.InternalError(err)
		return
	}
	var reqNodegroup model.NodeGroup
	err = n.ResolveBody(&reqNodegroup)
	if err != nil {
		logs.Error("Failed to get request nodegroup: %d, error: %+v", nodegroupID, err)
		n.CustomAbortAudit(http.StatusBadRequest, "Node group update failed.")
		return
	}
	res, err := service.UpdateNodeGroup(reqNodegroup, "comment")
	if err != nil {
		logs.Error("Failed to update nodegroup: %v, error: %+v", reqNodegroup, err)
		n.InternalError(err)
		return
	}
	logs.Debug("Update nodegroup %v %b", reqNodegroup, res)
}
