package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"strings"

	"github.com/astaxie/beego/logs"
)

type NodeController struct {
	BaseController
}

func (n *NodeController) GetNode() {
	para := n.GetString("node_name")
	res, err := service.GetNode(para)
	if err != nil {
		n.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.renderJSON(res)
}

func (n *NodeController) NodeToggle() {
	if !n.isSysAdmin {
		n.customAbort(http.StatusForbidden, "user should be admin")
		return
	}

	var responseStatus bool
	var err error
	paraName := n.GetString("node_name")
	paraStatus, _ := n.GetBool("node_status")

	switch paraStatus {
	case true:
		responseStatus, err = service.ResumeNode(paraName)
	case false:
		responseStatus, err = service.SuspendNode(paraName)
	}
	if err != nil {
		n.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	if responseStatus != true {
		n.customAbort(http.StatusPreconditionFailed, fmt.Sprint(err))
	}
}

func (n *NodeController) NodeList() {
	ping, _ := n.GetBool("ping")
	nodeList := service.GetNodeList()
	if ping {
		availableNodeList := []service.NodeListResult{}
		for _, node := range nodeList {
			status, err := utils.PingIPAddr(node.NodeIP)
			if err != nil {
				logs.Error("Failed to ping IPAddr: %s, error: %+v", node.NodeIP, err)
			}
			if status {
				availableNodeList = append(availableNodeList, node)
				break
			}
		}
		n.renderJSON(availableNodeList)
		return
	}
	n.renderJSON(nodeList)
}

func (n *NodeController) AddNodeToGroupAction() {
	//TODO node_id is not reay, should implement it
	//nodeID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))

	nodeName := n.GetString("node_name")
	groupName := n.GetString("groupname")
	logs.Debug("Adding %s to %s", nodeName, groupName)

	//TODO check existing
	err := service.AddNodeToGroup(nodeName, groupName)
	if err != nil {
		n.internalError(err)
		return
	}
}

func (n *NodeController) GetGroupsOfNodeAction() {

	//TODO node_id is not reay, should implement it
	//nodeID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))

	nodeName := n.GetString("node_name")

	// Get the nodegroups of this node
	groups, err := service.GetGroupOfNode(nodeName)
	if err != nil {
		logs.Error("Failed to get node %s group", nodeName)
		n.internalError(err)
		return
	}
	n.renderJSON(groups)
}

func (n *NodeController) RemoveNodeFromGroupAction() {
	//TODO node_id is not reay, should implement it
	//nodeID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))

	nodeName := n.GetString("node_name")
	groupName := n.GetString("groupname")
	//logs.Debug("Remove %s from %s", nodeName, groupName)

	err := service.RemoveNodeFromGroup(nodeName, groupName)
	if err != nil {
		n.internalError(err)
		return
	}
	logs.Debug("Removed %s from %s", nodeName, groupName)
}

func (n *NodeController) NodesAvailalbeResources() {
	logs.Debug("GetNodesResources")
	resources, err := service.GetNodesAvailableResources()
	if err != nil {
		n.internalError(err)
		return
	}

	n.renderJSON(resources)
}

func (n *NodeController) AddNodeAction() {
	var reqNode model.NodeCli
	var err error
	err = n.resolveBody(&reqNode)
	if err != nil {
		return
	}

	if !utils.ValidateWithLengthRange(reqNode.NodeName, 1, 63) {
		n.customAbort(http.StatusBadRequest, "NodeName must be not empty and no more than 63 characters ")
		return
	}

	nodeExists, err := service.NodeExists(reqNode.NodeName)
	if err != nil {
		n.internalError(err)
		return
	}
	if nodeExists {
		n.customAbort(http.StatusConflict, "Nodename already exists.")
		return
	}

	reqNode.NodeName = strings.TrimSpace(reqNode.NodeName)

	node, err := service.CreateNode(reqNode)
	if err != nil {
		logs.Debug("Failed to add node %s", reqNode.NodeName)
		n.internalError(err)
		return
	}
	logs.Info("Added node %s", node.ObjectMeta.Name)
}
