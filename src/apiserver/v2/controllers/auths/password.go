package auths

import (
	"github.com/astaxie/beego"
)

// Operations about user password
type PasswordController struct {
	beego.Controller
}

// @Title Reset password for user
// @Description Reset password for users.
// @Param	body	body 	"models.users.vm.ResetPassword"	true	"View model for user resetting password."
// @Success 200 Successful reset.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /reset [post]
func (u *PasswordController) Reset() {

}

// @Title Change password for user
// @Description Change password for users.
// @Param	body	body 	"models.users.vm.ChangePassword"	true	"View model for user changing password."
// @Success 200 Successful changed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /change [put]
func (u *PasswordController) Change() {

}

// @Title Forgot password for user
// @Description Forgot password for users.
// @Param	body	body 	"models.users.vm.ForgotPassword"	true	"View model for user changing password."
// @Success 200 Successful changed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /forgot [post]
func (u *PasswordController) Forgot() {

}
