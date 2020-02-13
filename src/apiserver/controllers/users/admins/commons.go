package admins

import (
	"github.com/astaxie/beego"
)

// Operations about users
type CommonController struct {
	beego.Controller
}

// @Title List all users by admin
// @Description List all for users.
// @Param	user_id	path	int	false	"ID of users"
// @Param	search	query	string	false	"Query item for users"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id [get]
func (u *CommonController) List() {

}

// @Title Add user by admin
// @Description Add user by admin.
// @Param	body	body 	models.users.vm.User	true	"View model for users."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (u *CommonController) Add() {

}

// @Title Update user info by admin
// @Description Update user info by admin.
// @Param	user_id	path	int	false	"ID of users"
// @Param	body	body 	models.users.vm.User	true	"View model for users."
// @Param	action	query	string	true	"Option of update."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id [put]
func (u *CommonController) Update() {

}

// @Title Delete user info by admin
// @Description Delete user info by admin.
// @Param	user_id	path	int	false	"ID of users"
// @Param	body	body 	models.users.vm.User	true	"View model for users."
// @Param	action	query	string	true	"Option of update."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id [delete]
func (u *CommonController) Delete() {

}
