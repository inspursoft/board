package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthController struct {
	baseController
}

func (u *AuthController) Prepare() {}

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
		user, err := service.SignIn(reqUser.Username, reqUser.Password)
		if err != nil {
			u.internalError(err)
			return
		}
		if user == nil {
			u.serveStatus(http.StatusBadRequest, "Incorrect username or password.")
			return
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
			return
		}
		memoryCache.Put(token.TokenString, token.TokenString, time.Second*time.Duration(tokenCacheExpireSeconds))
		u.Data["json"] = token
		u.ServeJSON()
	}
}

func (u *AuthController) SignUpAction() {
	var err error
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
		u.CustomAbort(http.StatusBadRequest, "Username content is illegal.")
		return
	}

	usernameExists, err := service.UsernameExists(reqUser.Username)
	if err != nil {
		u.internalError(err)
		return
	}

	if usernameExists {
		u.CustomAbort(http.StatusConflict, "Username already exists.")
		return
	}

	if !utils.ValidateWithLengthRange(reqUser.Password, 8, 20) {
		u.CustomAbort(http.StatusBadRequest, "Password length should be between 8 and 20 characters.")
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbort(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	emailExists, err := service.EmailExists(reqUser.Email)
	if err != nil {
		u.internalError(err)
		return
	}
	if emailExists {
		u.serveStatus(http.StatusConflict, "Email already exists.")
		return
	}

	if !utils.ValidateWithMaxLength(reqUser.Realname, 40) {
		u.CustomAbort(http.StatusBadRequest, "Realname maximum length is 40 characters.")
		return
	}

	if !utils.ValidateWithMaxLength(reqUser.Comment, 127) {
		u.CustomAbort(http.StatusBadRequest, "Comment maximum length is 127 characters.")
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
		u.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	u.Data["json"] = user
	u.ServeJSON()
}

func (u *AuthController) LogOutAction() {
	u.signOff()
}

func (u *AuthController) UserExists() {

	target := u.GetString("target")
	value := u.GetString("value")

	var isExists bool
	var err error
	switch target {
	case "username":
		isExists, err = service.UsernameExists(value)
	case "email":
		isExists, err = service.EmailExists(value)
	default:
		u.CustomAbort(http.StatusBadRequest, "unsupported check target.")
		return
	}
	if err != nil {
		u.internalError(err)
		return
	}
	if isExists {
		u.CustomAbort(http.StatusConflict, target+" already exists.")
	}
}
