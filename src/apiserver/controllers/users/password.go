package users

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/service/adapting"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type PasswordController struct {
	c.BaseController
}

func (u *PasswordController) Prepare() {
	u.EnableXSRF = false
	u.ResolveSignedInUser()
	u.RecordOperationAudit()
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

// @Title Change signed in user password
// @Description Change signed in user password
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	body	body 	"vm.ChangePassword"	true	"View model for user password."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /password [put]
func (u *PasswordController) UpdatePassword() {
	var err error

	if u.IsExternalAuth && u.CurrentUser.Username != "boardadmin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}

	currentUser := u.GetCurrentUser()
	user, err := service.GetUserByID(currentUser.ID)
	if err != nil {
		u.InternalError(err)
		return
	}
	if user == nil {
		u.CustomAbortAudit(http.StatusNotFound, "No found user with provided User ID.")
		return
	}

	if !(u.IsSysAdmin || u.CurrentUser.ID == user.ID) {
		u.CustomAbortAudit(http.StatusForbidden, "Only system admin can change others' password.")
		return
	}

	var changePassword vm.ChangePassword
	err = u.ResolveBody(&changePassword)
	if err != nil {
		return
	}

	changePassword.OldPassword = utils.Encrypt(changePassword.OldPassword, u.CurrentUser.Salt)

	if changePassword.OldPassword != user.Password {
		u.CustomAbortAudit(http.StatusForbidden, "Old password input is incorrect.")
		return
	}
	if !utils.ValidateWithLengthRange(changePassword.NewPassword, 8, 20) {
		u.CustomAbortAudit(http.StatusBadRequest, "Password does not satisfy complexity requirement.")
		return
	}
	updateUser := vm.User{
		ID:       user.ID,
		Password: utils.Encrypt(changePassword.NewPassword, u.CurrentUser.Salt),
	}
	isSuccess, err := adapting.UpdateUser(updateUser, "password")
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to change password.")
	}
}
