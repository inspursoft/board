package admins

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
)

//Operation about Admin's supplement actions
type SupplementController struct {
	c.BaseController
}

func (u *SupplementController) Prepare() {
	u.EnableXSRF = false
	u.ResolveSignedInUser()
	u.RecordOperationAudit()
	if !u.IsSysAdmin {
		u.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to manipulate user.")
		return
	}
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

// @Title Supplement for toggling user's admin option.
// @Description Supplement for toggling user's admin option.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	user_id	path	int	true	"Request for user ID."
// @Success 200 Successful toggled.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id/toggle [head]
func (u *SupplementController) Toggle() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":user_id"))
	if err != nil {
		u.InternalError(err)
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.InternalError(err)
		return
	}
	if user == nil {
		u.CustomAbortAudit(http.StatusNotFound, "No found user with provided user ID.")
		return
	}
	if userID == 1 || u.CurrentUser.ID == user.ID {
		u.CustomAbortAudit(http.StatusBadRequest, "Self or system admin cannot be changed.")
		return
	}
	logs.Debug("User: %+v", user)
	if user.SystemAdmin == 0 {
		user.SystemAdmin = 1
	} else {
		user.SystemAdmin = 0
	}
	isSuccess, err := service.UpdateUser(*user, "id", "system_admin")
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to toggle user system admin.")
	}
}
