package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

type AuthController struct {
	c.BaseController
}

func (u *AuthController) Prepare() {
	initStatus := utils.GetConfig("INIT_STATUS")
	if initStatus() != "READY" {
		u.ServeStatus(http.StatusNotAcceptable, initStatus())
		return
	}
	u.EnableXSRF = false
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
	u.RecordOperationAudit()
}

func (u *AuthController) SignInAction() {
	var reqUser model.User
	err := u.ResolveBody(&reqUser)
	if err != nil {
		return
	}
	token, _ := u.ProcessAuth(reqUser.Username, reqUser.Password)
	if token != "" {
		u.RenderJSON(model.Token{TokenString: token})
	}
}

func (u *AuthController) ExternalAuthAction() {
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

func (u *AuthController) SignUpAction() {
	if u.IsExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}
	var reqUser model.User
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

	reqUser.Password, err = service.DecodeUserPassword(reqUser.Password)
	if err != nil {
		logs.Error("Password encode error %v", err)
		u.CustomAbortAudit(http.StatusBadRequest, "Password encode error.")
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
	reqUser.SystemAdmin = 0

	isSuccess, err := service.SignUp(reqUser)
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.ServeStatus(http.StatusBadRequest, "Failed to sign up user.")
	}
}

func (u *AuthController) CurrentUserAction() {
	user := u.GetCurrentUser()
	if user == nil {
		u.CustomAbortAudit(http.StatusUnauthorized, "Need to login first.")
		return
	}

	u.RenderJSON(user)
}

func (u *AuthController) GetSystemInfo() {
	systemInfo, err := service.GetSystemInfo()
	if err != nil {
		u.InternalError(err)
		return
	}
	u.RenderJSON(systemInfo)
}

func (u *AuthController) GetSystemResources() {
	systemResources, err := service.GetSystemResourcesInfo()
	if err != nil {
		u.InternalError(err)
		return
	}
	u.RenderJSON(systemResources)
}

func (u *AuthController) GetKubernetesInfo() {
	kubernetesInfo, err := service.GetKubernetesInfo()
	if err != nil {
		u.InternalError(err)
		return
	}
	u.RenderJSON(kubernetesInfo)
}

func (u *AuthController) LogOutAction() {
	err := u.SignOff()
	if err != nil {
		u.CustomAbortAudit(http.StatusBadRequest, "Incorrect username to log out.")
	}
}

func (u *AuthController) UserExists() {
	target := u.GetString("target")
	value := u.GetString("value")
	userID, _ := u.GetInt64("user_id")
	isExists, err := service.UserExists(target, value, userID)
	if err != nil {
		u.InternalError(err)
		return
	}
	if isExists {
		u.CustomAbortAudit(http.StatusConflict, target+" already exists.")
	}
}

func (u *AuthController) ResetPassword() {
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

	newPassword, err = service.DecodeUserPassword(newPassword)
	if err != nil {
		logs.Error("Decode error %s %v", newPassword, err)
		u.CustomAbortAudit(http.StatusBadRequest, "No password encoded.")
		return
	}

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
