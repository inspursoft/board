package controller

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type PodController struct {
	c.BaseController
}

func (p *PodController) PodShellAction() {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":projectid"))
	if err != nil {
		p.InternalError(err)
		return
	}
	project, err := service.GetProjectByID(int64(projectID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if project == nil {
		p.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("No project was found with provided ID: %d", projectID))
		return
	}
	pod := p.Ctx.Input.Param(":podname")
	container := p.GetString("container")

	// upgrade the connection to websocket.
	logs.Info("Requested Pod %s/%s web console", project.Name, pod)
	err = service.PodShell(project.Name, pod, container, p.Ctx.ResponseWriter, p.Ctx.Request)
	if _, ok := err.(websocket.HandshakeError); ok {
		p.CustomAbortAudit(http.StatusBadRequest, "Not a websocket handshake.")
		return
	} else if err != nil {
		p.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

}
