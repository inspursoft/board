package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
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
		memoryCache.Put(token.TokenString, token.TokenString, time.Second*1800)
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
	if strings.TrimSpace(reqUser.Username) == "" {
		u.CustomAbort(http.StatusBadRequest, "Username cannot be empty.")
		return
	}
	if strings.TrimSpace(reqUser.Email) == "" {
		u.CustomAbort(http.StatusBadRequest, "Email cannot be empty.")
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
	emailExists, err := service.EmailExists(reqUser.Email)
	if err != nil {
		u.internalError(err)
		return
	}
	if emailExists {
		u.serveStatus(http.StatusConflict, "Email already exists.")
		return
	}
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
	token := u.GetString("token")
	payload, err := verifyToken(token)
	if err != nil || payload == nil {
		u.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	payload["token"] = token
	u.Data["json"] = payload
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
