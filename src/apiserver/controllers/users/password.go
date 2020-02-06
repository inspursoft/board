package users

import "github.com/astaxie/beego"

type PasswordController struct {
	beego.Controller
}

// @Title Change signed in user password
// @Description Change signed in user password
// @Param	body	body 	models.users.vm.Password	true	"View model for user password."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /password [put]
func (p *PasswordController) UpdatePassword() {

}
