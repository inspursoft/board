package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
	"strings"
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
		}

		payload := make(map[string]interface{})
		payload["id"] = strconv.Itoa(int(user.ID))

		token, err := signToken(payload)
		if err != nil {
			u.internalError(err)
			return
		}
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
