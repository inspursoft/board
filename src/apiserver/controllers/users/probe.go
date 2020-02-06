package users

import "github.com/astaxie/beego"

//Operation about user probe
type ProbeController struct {
	beego.Controller
}

// @Title Probe for getting current user.
// @Description Get current user information from session.
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /current [get]
func (p *ProbeController) Current() {

}
