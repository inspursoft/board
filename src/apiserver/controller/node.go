package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type NodeController struct {
	baseController
}

func (p *NodeController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
}
func (n *NodeController) GetNode() {
	para := n.GetString("node_name")
	res, err := service.GetNode(para)
	if err != nil {
		n.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.Data["json"] = res
	n.ServeJSON()
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
	res := service.GetNodeList()
	n.Data["json"] = res
	n.ServeJSON()
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
	n.Data["json"] = groups
	n.ServeJSON()
}

func (n *NodeController) RemoveNodeFromGroupAction() {
	//TODO node_id is not reay, should implement it
	//nodeID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))

	//TODO
	return
}
