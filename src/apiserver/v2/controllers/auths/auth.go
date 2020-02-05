package auths

import (
	"github.com/astaxie/beego"
)

// Operations about authorization
type AuthController struct {
	beego.Controller
}

// @Title Authoriazation for user sign in
// @Description Sign in for users.
// @Param	body	body 	"models.users.vm.SignIn"	true	"View model for user sign in."
// @Success 200 Successful signed in.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /sign-in [post]
func (u *AuthController) SignIn() {
	beego.Debug("Signed in")
}

// @Title Operation for user sign up
// @Description Sign up for users.
// @Param	body	body 	"models.users.vm.SignUp"	true	"View model for user sign up."
// @Success 200 Successful signed up.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /sign-up [post]
func (s *AuthController) SignUp() {

}

// @Title Authoriazation for user sign out
// @Description Sign out for users.
// @Param	username	query	string 	true	"Request for user sign out."
// @Success 200 Successful signed out.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /sign-out [get]
func (u *AuthController) SignOut() {
	beego.Debug("Signed out")
}

// @Title Authoriazation for third-party
// @Description Sign out for third-party.
// @Param	token	query	string 	true	"Request for third-party token."
// @Success 200 Successful signed out.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /third-party [post]
func (u *AuthController) ThirdParty() {

}
