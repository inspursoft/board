package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"

	"strings"

	"github.com/astaxie/beego/logs"
)

type NodeController struct {
	c.BaseController
}

func (n *NodeController) Prepare() {
	n.EnableXSRF = false
	n.ResolveSignedInUser()
	n.RecordOperationAudit()
}

func (n *NodeController) GetNode() {
	para := n.GetString("node_name")
	res, err := service.GetNode(para)
	if err != nil {
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.RenderJSON(res)
}

func (n *NodeController) NodeToggle() {
	if !n.IsSysAdmin {
		n.CustomAbortAudit(http.StatusForbidden, "user should be admin")
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
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	if responseStatus != true {
		n.CustomAbortAudit(http.StatusPreconditionFailed, fmt.Sprint(err))
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
		n.RenderJSON(availableNodeList)
		return
	}
	n.RenderJSON(nodeList)
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
		n.InternalError(err)
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
		n.InternalError(err)
		return
	}
	n.RenderJSON(groups)
}

func (n *NodeController) RemoveNodeFromGroupAction() {
	//TODO node_id is not reay, should implement it
	//nodeID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))

	nodeName := n.GetString("node_name")
	groupName := n.GetString("groupname")
	//logs.Debug("Remove %s from %s", nodeName, groupName)

	err := service.RemoveNodeFromGroup(nodeName, groupName)
	if err != nil {
		n.InternalError(err)
		return
	}
	logs.Debug("Removed %s from %s", nodeName, groupName)
}

func (n *NodeController) NodesAvailalbeResources() {
	logs.Debug("GetNodesResources")
	resources, err := service.GetNodesAvailableResources()
	if err != nil {
		n.InternalError(err)
		return
	}

	n.RenderJSON(resources)
}

func (n *NodeController) AddNodeAction() {
	var reqNode model.NodeCli
	var err error
	err = n.ResolveBody(&reqNode)
	if err != nil {
		return
	}

	if !utils.ValidateWithLengthRange(reqNode.NodeName, 1, 63) {
		n.CustomAbortAudit(http.StatusBadRequest, "NodeName must be not empty and no more than 63 characters ")
		return
	}

	nodeExists, err := service.NodeExists(reqNode.NodeName)
	if err != nil {
		n.InternalError(err)
		return
	}
	if nodeExists {
		n.CustomAbortAudit(http.StatusConflict, "Nodename already exists.")
		return
	}

	reqNode.NodeName = strings.TrimSpace(reqNode.NodeName)

	node, err := service.CreateNode(reqNode)
	if err != nil {
		logs.Debug("Failed to add node %s", reqNode.NodeName)
		n.InternalError(err)
		return
	}
	logs.Info("Added node %s", node.ObjectMeta.Name)
}

//Remove a node from cluster by node name
func (n *NodeController) RemoveNodeAction() {
	if n.IsSysAdmin == false {
		n.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to control node.")
		return
	}
	nodeName := strings.TrimSpace(n.Ctx.Input.Param(":nodename"))

	// Workaround delete node by IP, deprecated
	nodeIP := n.GetString("node_ip")
	if nodeIP != "" && nodeIP != nodeName {
		nodeitem, err := service.GetNodebyIP(nodeIP)
		if err != nil {
			n.InternalError(err)
			return
		}
		if nodeitem == nil {
			logs.Debug("Not found node by IP %s", nodeIP)
			n.CustomAbortAudit(http.StatusNotFound, "Not found this node IP.")
			return
		}
		nodeName = nodeitem.Name
	}

	logs.Debug("To delete node %s", nodeName)
	res, err := service.DeleteNode(nodeName)
	if err != nil {
		n.InternalError(err)
		return
	}

	if res == false {
		n.CustomAbortAudit(http.StatusNotFound, "Nodename Not Found.")
		return
	}
	logs.Debug("Removed %s", nodeName)
}

// Get the running status of a node
func (n *NodeController) GetNodeStatusAction() {
	if n.IsSysAdmin == false {
		n.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to control node.")
		return
	}
	nodeName := strings.TrimSpace(n.Ctx.Input.Param(":nodename"))
	logs.Debug("Get node status %s", nodeName)

	nExists, err := service.NodeExists(nodeName)
	if err != nil {
		logs.Debug("Failed to list nodes for %s", nodeName)
		n.InternalError(err)
		return
	}
	if !nExists {
		logs.Info("Node name %s not existing in cluster.", nodeName)
		n.CustomAbortAudit(http.StatusNotFound, "Node name not found.")
		return
	}

	nodestatus, err := service.GetNodeControlStatus(nodeName)
	if err != nil {
		logs.Debug("Failed to get node status %s", nodeName)
		n.InternalError(err)
		return
	}
	n.RenderJSON(*nodestatus)
}

// Drain the service instances from the node
func (n *NodeController) NodeDrainAction() {
	if n.IsSysAdmin == false {
		n.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to control node.")
		return
	}
	nodeName := strings.TrimSpace(n.Ctx.Input.Param(":nodename"))
	logs.Debug("Drain the node %s", nodeName)

	nodeDel, err := service.GetNodebyName(nodeName)
	if err != nil {
		logs.Debug("Failed to get node %s", nodeName)
		n.InternalError(err)
		return
	}

	if !nodeDel.Unschedulable {
		logs.Debug("Cannot drain a schedulable node %s", nodeName)
		n.CustomAbortAudit(http.StatusPreconditionRequired, "Cannot drain a schedulable node.")
		return
	}

	//TODO drain services by adminserver

	//Drain services by k8s api server
	err = service.DrainNodeServiceInstance(nodeName)
	if err != nil {
		logs.Debug("Failed to drain node %s", nodeName)
		n.InternalError(err)
		return
	}
	logs.Debug("Drained node %s", nodeName)
}

// Get edge node list
func (n *NodeController) EdgeNodeList() {
	var edgelist []string
	nodeList := service.GetNodeList()
	for _, v := range nodeList {
		if v.NodeType == service.NodeTypeEdge {
			edgelist = append(edgelist, v.NodeName)
		}
	}
	n.RenderJSON(edgelist)
}

// Create a new edge node
func (n *NodeController) AddEdgeNodeAction() {
	var reqNode model.EdgeNodeCli
	var err error
	err = n.ResolveBody(&reqNode)
	if err != nil {
		return
	}

	if !utils.ValidateWithLengthRange(reqNode.NodeName, 1, 63) {
		n.CustomAbortAudit(http.StatusBadRequest, "NodeName must be not empty and no more than 63 characters ")
		return
	}

	nodeExists, err := service.NodeExists(reqNode.NodeName)
	if err != nil {
		n.InternalError(err)
		return
	}
	if nodeExists {
		n.CustomAbortAudit(http.StatusConflict, "Nodename already exists.")
		return
	}

	//reqNode.NodeName = strings.TrimSpace(reqNode.NodeName)

	//TODO Check the hostname config
	res, err := service.CheckEdgeHostname(reqNode)
	if res != true || err != nil {
		n.CustomAbortAudit(http.StatusBadRequest, "edgenode config error.")
		return
	}

	//TODO create edge node, label yaml and run script
	node, err := service.CreateEdgeNode(reqNode)
	if err != nil {
		logs.Debug("Failed to add edge node %s", reqNode.NodeName)
		n.InternalError(err)
		return
	}
	logs.Info("Added edge node %s", node.ObjectMeta.Name)
}

// Get the edge node
func (n *NodeController) GetEdgeNodeAction() {
	nodeName := strings.TrimSpace(n.Ctx.Input.Param(":nodename"))
	logs.Debug("Get the edge node %s", nodeName)

	nodeDel, err := service.GetNodebyName(nodeName)
	if err != nil {
		logs.Debug("Failed to get node %s", nodeName)
		n.InternalError(err)
		return
	}
	n.RenderJSON(*nodeDel)
}

// Delete the edge node
func (n *NodeController) RemoveEdgeNodeAction() {
	if n.IsSysAdmin == false {
		n.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to control node.")
		return
	}
	nodeName := strings.TrimSpace(n.Ctx.Input.Param(":nodename"))
	logs.Debug("Get the edge node %s", nodeName)

	//TODO remove an edge node
	//TODO Check the edge status, autonomous offline
	logs.Debug("To delete node %s", nodeName)
	res, err := service.DeleteNode(nodeName)
	if err != nil {
		n.InternalError(err)
		return
	}

	if res == false {
		n.CustomAbortAudit(http.StatusNotFound, "Nodename Not Found.")
		return
	}
	logs.Debug("Removed Edge %s", nodeName)
}

// Get edge node hostname and check exsiting
func (n *NodeController) CheckEdgeName() {
	edgeIP := n.GetString("edge_ip")
	edgePassword := n.GetString("edge_password")

	if edgeIP == "" || edgePassword == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "IP or password invalid")
		return
	}

	edgeHostname, err := service.GetEdgeHostname(edgeIP, edgePassword)
	if err != nil {
		logs.Debug("Failed to get edge hostname %v", err)
		n.CustomAbortAudit(http.StatusBadRequest, "IP or password invalid")
		return
	}

	n.RenderJSON(edgeHostname)
}
