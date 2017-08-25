package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
)

type ServerTimeController struct {
	baseController
}

func (p *ServerTimeController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}

func (s *ServerTimeController) GetServerTime() {
	time := service.GetServerTime()
	s.Data["json"] = time
	s.ServeJSON()

}
