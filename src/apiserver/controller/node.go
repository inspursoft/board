package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
)

type NodeController struct {
	baseController
}

func (p *NodeController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}
func (n *NodeController) GetNode() {
	para := n.GetString("node_name")
	res, err := service.GetNode(para)
	if err != nil {
		n.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.Data["json"] = res
	n.ServeJSON()
}

func (n *NodeController) NodeToggle() {
	var responseStatus bool
	var err error
	paraName := n.GetString("node_name")
	paraStatus, _ := n.GetBool("node_status")
	if !n.isSysAdmin {
		n.CustomAbort(http.StatusForbidden, "user should be admin")
		return

	}
	switch paraStatus {
	case true:
		responseStatus, err = service.ResumeNode(paraName)
	case false:
		responseStatus, err = service.SuspendNode(paraName)
	}
	if err != nil {
		n.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	if responseStatus != true {
		n.CustomAbort(http.StatusPreconditionFailed, fmt.Sprint(err))
	}

}
func (n *NodeController) NodeList() {
	if !n.isSysAdmin {
		n.CustomAbort(http.StatusForbidden, "user should be admin")
		return
	}
	res := service.GetNodeList()
	n.Data["json"] = res
	n.ServeJSON()
}
