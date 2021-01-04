package users

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/service/adapting"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

// Operations about users
type CommonController struct {
	c.BaseController
}

func (u *CommonController) Prepare() {
	u.EnableXSRF = false
	u.ResolveSignedInUser()
	u.RecordOperationAudit()
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

// @Title Change signed in user self info
// @Description Change signed in user self info
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	body	body 	"vm.User"	true	"View model for users."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [put]
func (u *CommonController) Update() {
	if u.IsExternalAuth && u.CurrentUser.Username != "boardadmin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusPreconditionFailed, "Current AUTH_MODE is not available to the user.")
		return
	}

	var reqUser vm.User
	var err error
	err = u.ResolveBody(&reqUser)
	if err != nil {
		return
	}

	reqUser.ID = u.CurrentUser.ID
	users, err := service.GetUsers("email", reqUser.Email, "id", "email")
	if err != nil {
		u.InternalError(err)
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbortAudit(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	if len(users) > 0 && users[0].ID != reqUser.ID {
		u.CustomAbortAudit(http.StatusConflict, "Email already exists.")
		return
	}

	if !utils.ValidateWithMaxLength(reqUser.Realname, 40) {
		u.CustomAbortAudit(http.StatusBadRequest, "Realname maximum length is 40 characters.")
		return
	}

	if !utils.ValidateWithMaxLength(reqUser.Comment, 127) {
		u.CustomAbortAudit(http.StatusBadRequest, "Comment maximum length is 127 characters.")
		return
	}

	reqUser.Email = strings.TrimSpace(reqUser.Email)
	reqUser.Realname = strings.TrimSpace(reqUser.Realname)
	reqUser.Comment = strings.TrimSpace(reqUser.Comment)

	isSuccess, err := adapting.UpdateUser(reqUser, "email", "realname", "comment")
	if err != nil {
		u.InternalError(err)
		return
	}

	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to change user account.")
	}
}
