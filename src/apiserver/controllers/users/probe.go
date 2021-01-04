package users

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"net/http"
)

//Operation about user probe
type ProbeController struct {
	c.BaseController
}

func (p *ProbeController) Prepare() {
	p.EnableXSRF = false
}

// @Title Probe for getting current user.
// @Description Get current user information from session.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /current [get]
func (p *ProbeController) Current() {
	user := p.GetCurrentUser()
	if user == nil {
		p.CustomAbortAudit(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.RenderJSON(user)
}
