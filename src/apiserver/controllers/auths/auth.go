package auths

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/service/adapting"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

// Operations about authorization
type AuthController struct {
	c.BaseController
}

func (u *AuthController) Prepare() {
	u.EnableXSRF = false
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
	u.RecordOperationAudit()
}

// @Title Authoriazation for user sign in
// @Description Sign in for users.
// @Param	body	body 	"vm.SignIn"	true	"View model for user sign in."
// @Success 200 Successful signed in.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /sign-in [post]
func (u *AuthController) SignIn() {
	var reqUser vm.SignIn
	err := u.ResolveBody(&reqUser)
	if err != nil {
		return
	}
	logs.Debug("Resolved body: %+v", reqUser)
	token, _ := u.ProcessAuth(reqUser.Username, reqUser.Password)
	u.RenderJSON(vm.Token{TokenString: token})
}

// @Title Operation for user sign up
// @Description Sign up for users.
// @Param	body	body 	"vm.SignUp"	true	"View model for user sign up."
// @Success 200 Successful signed up.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /sign-up [post]
func (u *AuthController) SignUp() {
	if u.IsExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}
	var reqUser vm.User
	err := u.ResolveBody(&reqUser)
	if err != nil {
		return
	}

	if !utils.ValidateWithPattern("username", reqUser.Username) {
		u.CustomAbortAudit(http.StatusBadRequest, "Username content is illegal.")
		return
	}

	// can't be the reserved name.
	for _, rsdname := range c.ReservedUsernames {
		if rsdname == reqUser.Username {
			u.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Username %s is reserved.", reqUser.Username))
			return
		}
	}

	usernameExists, err := service.UserExists("username", reqUser.Username, 0)
	if err != nil {
		u.InternalError(err)
		return
	}

	if usernameExists {
		u.CustomAbortAudit(http.StatusConflict, "Username already exists.")
		return
	}

	if !utils.ValidateWithLengthRange(reqUser.Password, 8, 20) {
		u.CustomAbortAudit(http.StatusBadRequest, "Password length should be between 8 and 20 characters.")
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbortAudit(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	emailExists, err := service.UserExists("username", reqUser.Email, 0)
	if err != nil {
		u.InternalError(err)
		return
	}
	if emailExists {
		u.ServeStatus(http.StatusConflict, "Email already exists.")
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

	reqUser.Username = strings.TrimSpace(reqUser.Username)
	reqUser.Email = strings.TrimSpace(reqUser.Email)
	reqUser.Realname = strings.TrimSpace(reqUser.Realname)
	reqUser.Comment = strings.TrimSpace(reqUser.Comment)

	isSuccess, err := adapting.SignUp(reqUser)
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.ServeStatus(http.StatusBadRequest, "Failed to sign up user.")
	}
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
	err := u.SignOff()
	if err != nil {
		u.CustomAbortAudit(http.StatusBadRequest, "Incorrect username to log out.")
	}
}

// @Title Authoriazation for third-party
// @Description Sign out for third-party.
// @Param	external_token	query	string 	true	"Request for third-party token."
// @Success 200 Successful signed out.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /third-party [post]
func (u *AuthController) ThirdParty() {
	externalToken := u.GetString("external_token")
	if externalToken == "" {
		u.CustomAbortAudit(http.StatusBadRequest, "Missing token for verification.")
		return
	}
	if token, isSuccess := u.ProcessAuth(externalToken, ""); isSuccess {
		u.Redirect(fmt.Sprintf("http://%s/dashboard?token=%s", utils.GetStringValue("BOARD_HOST_IP"), token), http.StatusFound)
		logs.Debug("Successful logged in.")
	}
}
