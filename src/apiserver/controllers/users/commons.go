package users

import (
	"github.com/astaxie/beego"
)

// Operations about users
type CommonController struct {
	beego.Controller
}

// @Title Change signed in user self info
// @Description Change signed in user self info
// @Param	body	body 	models.users.vm.User	true	"View model for users."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [put]
func (u *CommonController) Update() {

}
