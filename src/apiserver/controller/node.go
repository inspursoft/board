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
	/*user := p.getCurrentUser()
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)*/
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
	return

}
