package controller

import (
	"encoding/json"
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

type AuthController struct {
	baseController
}

func (u *AuthController) Prepare() {
	u.isExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

func (u *AuthController) processAuth(principal, password string) (string, bool) {
	currentAuth, err := auth.GetAuth(authMode())
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
	payload["is_project_admin"] = user.ProjectAdmin
	payload["is_system_admin"] = user.SystemAdmin
	token, err := signToken(payload)
	if err != nil {
		u.internalError(err)
		return "", false
	}
	memoryCache.Put(user.Username, token.TokenString, time.Second*time.Duration(tokenCacheExpireSeconds))
	memoryCache.Put(token.TokenString, payload, time.Second*time.Duration(tokenCacheExpireSeconds))
	return token.TokenString, true
}

func (u *AuthController) SignInAction() {
	var err error
	reqData, err := u.resolveBody()
	if err != nil {
		u.internalError(err)
		return
	}
	if reqData != nil {
		var reqUser model.User
		err = json.Unmarshal(reqData, &reqUser)
		if err != nil {
			u.internalError(err)
			return
		}
		token, _ := u.processAuth(reqUser.Username, reqUser.Password)
		u.Data["json"] = model.Token{TokenString: token}
		u.ServeJSON()
	}
}

func (u *AuthController) ExternalAuthAction() {
	externalToken := u.GetString("external_token")
	if externalToken == "" {
		u.customAbort(http.StatusBadRequest, "Missing token for verification.")
		return
	}
	if token, isSuccess := u.processAuth(externalToken, ""); isSuccess {
		u.Redirect(fmt.Sprintf("http://%s/dashboard?token=%s", utils.GetStringValue("BOARD_HOST"), token), http.StatusFound)
		logs.Debug("Successful logged in.")
	}

}

func (u *AuthController) SignUpAction() {
	var err error
	if u.isExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.customAbort(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}

	reqData, err := u.resolveBody()
	if err != nil {
		u.internalError(err)
		return
	}
	var reqUser model.User
	err = json.Unmarshal(reqData, &reqUser)
	if err != nil {
		u.internalError(err)
		return
	}

	if !utils.ValidateWithPattern("username", reqUser.Username) {
		u.customAbort(http.StatusBadRequest, "Username content is illegal.")
		return
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
	u.Data["json"] = user
	u.ServeJSON()
}

func (u *AuthController) GetSystemInfo() {
	systemInfo, err := service.GetSystemInfo()
	if err != nil {
		u.internalError(err)
		return
	}
	u.Data["json"] = systemInfo
	u.ServeJSON()
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
