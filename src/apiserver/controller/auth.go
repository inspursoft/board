package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/auth"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

var reservdUsernames = [...]string{"explore", "create", "assets", "css", "img", "js", "less", "plugins", "debug", "raw", "install", "api", "avatar", "user", "org", "help", "stars", "issues", "pulls", "commits", "repo", "template", "new", ".", ".."}

type AuthController struct {
	BaseController
}

func (u *AuthController) Prepare() {
	u.isExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
	u.recordOperationAudit()
}

func (u *AuthController) processAuth(principal, password string) (string, bool) {
	var currentAuth *auth.Auth
	var err error
	if principal == "admin" {
		currentAuth, err = auth.GetAuth("db_auth")
	} else {
		currentAuth, err = auth.GetAuth(authMode())
	}
	if err != nil {
		u.internalError(err)
		return "", false
	}
	user, err := (*currentAuth).DoAuth(principal, password)
	if err != nil {
		u.internalError(err)
		return "", false
	}

	if user == nil {
		u.serveStatus(http.StatusBadRequest, "Incorrect username or password.")
		return "", false
	}
	payload := make(map[string]interface{})
	payload["id"] = strconv.Itoa(int(user.ID))
	payload["username"] = user.Username
	payload["email"] = user.Email
	payload["realname"] = user.Realname
	payload["is_system_admin"] = user.SystemAdmin
	token, err := signToken(payload)
	if err != nil {
		u.internalError(err)
		return "", false
	}
	memoryCache.Put(user.Username, token.TokenString, time.Second*time.Duration(tokenCacheExpireSeconds))
	memoryCache.Put(token.TokenString, payload, time.Second*time.Duration(tokenCacheExpireSeconds))
	u.auditUser, _ = service.GetUserByName(user.Username)
	return token.TokenString, true
}

func (u *AuthController) SignInAction() {
	var reqUser model.User
	err := u.resolveBody(&reqUser)
	if err != nil {
		return
	}
	token, _ := u.processAuth(reqUser.Username, reqUser.Password)
	u.renderJSON(model.Token{TokenString: token})
}

func (u *AuthController) ExternalAuthAction() {
	externalToken := u.GetString("external_token")
	if externalToken == "" {
		u.customAbort(http.StatusBadRequest, "Missing token for verification.")
		return
	}
	if token, isSuccess := u.processAuth(externalToken, ""); isSuccess {
		u.Redirect(fmt.Sprintf("http://%s/dashboard?token=%s", utils.GetStringValue("BOARD_HOST_IP"), token), http.StatusFound)
		logs.Debug("Successful logged in.")
	}
}

func (u *AuthController) SignUpAction() {
	if u.isExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.customAbort(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}
	var reqUser model.User
	err := u.resolveBody(&reqUser)
	if err != nil {
		return
	}

	if !utils.ValidateWithPattern("username", reqUser.Username) {
		u.customAbort(http.StatusBadRequest, "Username content is illegal.")
		return
	}

	// can't be the reserved name.
	for _, rsdname := range reservdUsernames {
		if rsdname == reqUser.Username {
			u.customAbort(http.StatusBadRequest, fmt.Sprintf("Username %s is reserved.", reqUser.Username))
			return
		}
	}

	usernameExists, err := service.UserExists("username", reqUser.Username, 0)
	if err != nil {
		u.internalError(err)
		return
	}

	if usernameExists {
		u.customAbort(http.StatusConflict, "Username already exists.")
		return
	}

	if !utils.ValidateWithLengthRange(reqUser.Password, 8, 20) {
		u.customAbort(http.StatusBadRequest, "Password length should be between 8 and 20 characters.")
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.customAbort(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	emailExists, err := service.UserExists("username", reqUser.Email, 0)
	if err != nil {
		u.internalError(err)
		return
	}
	if emailExists {
		u.serveStatus(http.StatusConflict, "Email already exists.")
		return
	}

	if !utils.ValidateWithMaxLength(reqUser.Realname, 40) {
		u.customAbort(http.StatusBadRequest, "Realname maximum length is 40 characters.")
		return
	}

	if !utils.ValidateWithMaxLength(reqUser.Comment, 127) {
		u.customAbort(http.StatusBadRequest, "Comment maximum length is 127 characters.")
		return
	}

	reqUser.Username = strings.TrimSpace(reqUser.Username)
	reqUser.Email = strings.TrimSpace(reqUser.Email)
	reqUser.Realname = strings.TrimSpace(reqUser.Realname)
	reqUser.Comment = strings.TrimSpace(reqUser.Comment)

	isSuccess, err := service.SignUp(reqUser)
	if err != nil {
		u.internalError(err)
		return
	}
	if !isSuccess {
		u.serveStatus(http.StatusBadRequest, "Failed to sign up user.")
	}
}

func (u *AuthController) CurrentUserAction() {
	user := u.getCurrentUser()
	if user == nil {
		u.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}

	u.renderJSON(user)
}

func (u *AuthController) GetSystemInfo() {
	systemInfo, err := service.GetSystemInfo()
	if err != nil {
		u.internalError(err)
		return
	}
	u.renderJSON(systemInfo)
}

func (u *AuthController) GetKubernetesInfo() {
	kubernetesInfo, err := service.GetKubernetesInfo()
	if err != nil {
		u.internalError(err)
		return
	}
	u.renderJSON(kubernetesInfo)
}

func (u *AuthController) LogOutAction() {
	err := u.signOff()
	if err != nil {
		u.customAbort(http.StatusBadRequest, "Incorrect username to log out.")
	}
}

func (u *AuthController) UserExists() {
	target := u.GetString("target")
	value := u.GetString("value")
	userID, _ := u.GetInt64("user_id")
	isExists, err := service.UserExists(target, value, userID)
	if err != nil {
		u.internalError(err)
		return
	}
	if isExists {
		u.customAbort(http.StatusConflict, target+" already exists.")
	}
}

func (u *AuthController) ResetPassword() {
	if utils.GetBoolValue("IS_EXTERNAL_AUTH") {
		u.customAbort(http.StatusPreconditionFailed, "Resetting password doesn't support in external auth.")
		return
	}
	resetUUID := u.GetString("reset_uuid")
	user, err := service.GetUserByResetUUID(resetUUID)
	if err != nil {
		logs.Error("Failed to get user by reset UUID: %s, error: %+v", resetUUID, err)
		u.internalError(err)
		return
	}
	if user == nil {
		logs.Error("Invalid reset UUID: %s", resetUUID)
		u.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid reset UUID: %s", resetUUID))
		return
	}
	newPassword := u.GetString("password")
	if strings.TrimSpace(newPassword) == "" {
		logs.Error("No password provided.")
		u.customAbort(http.StatusBadRequest, "No password provided.")
		return
	}
	_, err = service.ResetUserPassword(*user, newPassword)
	if err != nil {
		logs.Error("Failed to reset user password for user ID: %d, error: %+v", user.ID, err)
		u.internalError(err)
	}
}
