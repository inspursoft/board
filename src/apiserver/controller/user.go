package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"
)

type UserController struct {
	baseController
}

func (u *UserController) Prepare() {
	user := u.getCurrentUser()
	if user == nil {
		u.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	u.currentUser = user
	u.isSysAdmin = (u.currentUser.SystemAdmin == 1)
	u.isProjectAdmin = (u.currentUser.ProjectAdmin == 1)
}

func (u *UserController) GetUsersAction() {
	username := u.GetString("username")
	email := u.GetString("email")
	var users []*model.User
	var err error
	if strings.TrimSpace(username) != "" {
		users, err = service.GetUsers("username", username)
	} else if strings.TrimSpace(email) != "" {
		users, err = service.GetUsers("email", email)
	} else {
		users, err = service.GetUsers("", nil)
	}
	if err != nil {
		u.internalError(err)
		return
	}
	for _, u0 := range users {
		u0.Password = ""
	}
	u.Data["json"] = users
	u.ServeJSON()
}

func (u *UserController) ChangeUserAccount() {
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

	reqUser.ID = u.currentUser.ID

	users, err := service.GetUsers("email", reqUser.Email)
	if err != nil {
		u.internalError(err)
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbort(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	if len(users) > 0 && users[0].ID != reqUser.ID {
		u.CustomAbort(http.StatusConflict, "Email already exists.")
		return
	}

	isSuccess, err := service.UpdateUser(reqUser, "email", "realname", "comment")
	if err != nil {
		u.internalError(err)
		return
	}

	if !isSuccess {
		u.CustomAbort(http.StatusBadRequest, "Failed to change user account.")
	}
}

func (u *UserController) ChangePasswordAction() {
	var err error
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.internalError(err)
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if user == nil {
		u.CustomAbort(http.StatusNotFound, "No found user with provided User ID.")
		return
	}

	if !(u.isSysAdmin || u.currentUser.ID == user.ID) {
		u.CustomAbort(http.StatusForbidden, "Only system admin can change others' password.")
		return
	}

	reqData, err := u.resolveBody()
	if err != nil {
		u.internalError(err)
		return
	}

	var changePassword model.ChangePassword
	err = json.Unmarshal(reqData, &changePassword)
	if err != nil {
		u.internalError(err)
		return
	}

	changePassword.OldPassword = utils.Encrypt(changePassword.OldPassword, u.currentUser.Salt)

	if changePassword.OldPassword != user.Password {
		u.CustomAbort(http.StatusForbidden, "Old password input is incorrect.")
		return
	}
	if !utils.ValidateWithLengthRange(changePassword.NewPassword, 8, 20) {
		u.CustomAbort(http.StatusBadRequest, "Password does not satisfy complexity requirement.")
		return
	}
	updateUser := model.User{
		ID:       user.ID,
		Password: utils.Encrypt(changePassword.NewPassword, u.currentUser.Salt),
	}
	isSuccess, err := service.UpdateUser(updateUser, "password")
	if err != nil {
		u.internalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbort(http.StatusBadRequest, "Failed to change password")
	}
}

type SystemAdminController struct {
	baseController
}

func (u *SystemAdminController) Prepare() {
	user := u.getCurrentUser()
	if user == nil {
		u.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	u.currentUser = user
	u.isSysAdmin = (user.SystemAdmin == 1)
	if !u.isSysAdmin {
		u.CustomAbort(http.StatusForbidden, "Insuffient privileges to manipulate user.")
		return
	}
}

func (u *SystemAdminController) AddUserAction() {
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
	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbort(http.StatusBadRequest, "Email content is illegal.")
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

func (u *SystemAdminController) GetUserAction() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.internalError(err)
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if user == nil {
		u.CustomAbort(http.StatusNotFound, "No user found with provided User ID.")
		return
	}
	user.Password = ""
	u.Data["json"] = user
	u.ServeJSON()
}

func (u *SystemAdminController) DeleteUserAction() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.CustomAbort(http.StatusBadRequest, "Invalid user ID.")
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if user == nil {
		u.CustomAbort(http.StatusNotFound, "No user was found with provided ID.")
		return
	}
	if userID == 1 || int64(userID) == u.currentUser.ID {
		u.CustomAbort(http.StatusBadRequest, "System admin user or current user cannot be deleted.")
		return
	}
	isSuccess, err := service.DeleteUser(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if !isSuccess {
		u.serveStatus(http.StatusBadRequest, "Failed to delete user.")
	}
}

func (u *SystemAdminController) UpdateUserAction() {
	var err error
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.serveStatus(http.StatusBadRequest, "Invalid user ID.")
		return
	}
	if u.currentUser.ID == int64(userID) {
		u.CustomAbort(http.StatusForbidden, "Insuffient privileges.")
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if user == nil {
		u.CustomAbort(http.StatusNotFound, "No user was found with provided ID.")
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
	reqUser.ID = user.ID

	users, err := service.GetUsers("email", reqUser.Email)
	if err != nil {
		u.internalError(err)
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbort(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	if len(users) > 0 && users[0].ID != reqUser.ID {
		u.CustomAbort(http.StatusConflict, "Email already exists.")
		return
	}

	user.Email = reqUser.Email
	user.Realname = reqUser.Realname
	user.Comment = reqUser.Comment

	isSuccess, err := service.UpdateUser(reqUser, "email", "realname", "comment")
	if err != nil {
		u.internalError(err)
		return
	}
	if !isSuccess {
		u.serveStatus(http.StatusBadRequest, "Failed to update user.")
	}
}

func toggleUserAction(u *SystemAdminController, actionName string) {
	var err error
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.internalError(err)
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if user == nil {
		u.CustomAbort(http.StatusNotFound, "No found user with provided user ID.")
		return
	}
	if u.currentUser.ID == user.ID {
		u.CustomAbort(http.StatusBadRequest, "Self system admin cannot be changed.")
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
	switch actionName {
	case "system_admin":
		user.SystemAdmin = reqUser.SystemAdmin
	case "project_admin":
		user.ProjectAdmin = reqUser.ProjectAdmin
	}
	isSuccess, err := service.UpdateUser(*user, actionName)
	if err != nil {
		u.internalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbort(http.StatusBadRequest, "Failed to toggle user system admin.")
	}
}

func (u *SystemAdminController) ToggleSystemAdminAction() {
	toggleUserAction(u, "system_admin")
}

func (u *SystemAdminController) ToggleProjectAdminAction() {
	toggleUserAction(u, "project_admin")
}
