package auths

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

// Operations about user password
type PasswordController struct {
	c.BaseController
}

// @Title Reset password for user
// @Description Reset password for users.
// @Param	reset_uuid	query 	string	true	"UUID for requesting to reset password"
// @Param	password	query	string	true	"New password to be reset."
// @Success 200 Successful reset.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /reset [post]
func (u *PasswordController) Reset() {
	if utils.GetBoolValue("IS_EXTERNAL_AUTH") {
		u.CustomAbortAudit(http.StatusPreconditionFailed, "Resetting password doesn't support in external auth.")
		return
	}
	resetUUID := u.GetString("reset_uuid")
	user, err := service.GetUserByResetUUID(resetUUID)
	if err != nil {
		logs.Error("Failed to get user by reset UUID: %s, error: %+v", resetUUID, err)
		u.InternalError(err)
		return
	}
	if user == nil {
		logs.Error("Invalid reset UUID: %s", resetUUID)
		u.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Invalid reset UUID: %s", resetUUID))
		return
	}
	newPassword := u.GetString("password")
	if strings.TrimSpace(newPassword) == "" {
		logs.Error("No password provided.")
		u.CustomAbortAudit(http.StatusBadRequest, "No password provided.")
		return
	}
	_, err = service.ResetUserPassword(*user, newPassword)
	if err != nil {
		logs.Error("Failed to reset user password for user ID: %d, error: %+v", user.ID, err)
		u.InternalError(err)
	}
}
