package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

type UserController struct {
	BaseController
}

func (u *UserController) Prepare() {
	u.resolveSignedInUser()
	u.isExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

func (u *UserController) GetUsersAction() {
	username := u.GetString("username")
	email := u.GetString("email")
	pageIndex, _ := u.GetInt("page_index", 0)
	pageSize, _ := u.GetInt("page_size", 0)
	isPaginated := !(pageIndex == 0 && pageSize == 0)
	orderField := u.GetString("order_field", "CREATE_TIME")
	orderAsc, _ := u.GetInt("order_asc", 0)

	var paginatedUsers *model.PaginatedUsers
	var users []*model.User
	var err error
	var fieldName string
	var fieldValue interface{}
	if strings.TrimSpace(username) != "" {
		fieldName = "username"
		fieldValue = username
	} else if strings.TrimSpace(email) != "" {
		fieldName = "email"
		fieldValue = email
	}
	if isPaginated {
		paginatedUsers, err = service.GetPaginatedUsers(fieldName, fieldValue, pageIndex, pageSize, orderField, orderAsc)
		u.Data["json"] = paginatedUsers
	} else {
		users, err = service.GetUsers(fieldName, fieldValue)
		u.Data["json"] = users
	}
	if err != nil {
		u.internalError(err)
		return
	}
	u.ServeJSON()
}

func (u *UserController) ChangeUserAccount() {
	if u.isExternalAuth && u.currentUser.Username != "admin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.customAbort(http.StatusPreconditionFailed, "Current AUTH_MODE is not available to the user.")
		return
	}

	var reqUser model.User
	var err error
	u.resolveBody(&reqUser)

	reqUser.ID = u.currentUser.ID
	users, err := service.GetUsers("email", reqUser.Email)
	if err != nil {
		u.internalError(err)
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.customAbort(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	if len(users) > 0 && users[0].ID != reqUser.ID {
		u.customAbort(http.StatusConflict, "Email already exists.")
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

	reqUser.Email = strings.TrimSpace(reqUser.Email)
	reqUser.Realname = strings.TrimSpace(reqUser.Realname)
	reqUser.Comment = strings.TrimSpace(reqUser.Comment)

	isSuccess, err := service.UpdateUser(reqUser, "email", "realname", "comment")
	if err != nil {
		u.internalError(err)
		return
	}

	if !isSuccess {
		u.customAbort(http.StatusBadRequest, "Failed to change user account.")
	}
}

func (u *UserController) ChangePasswordAction() {
	var err error

	if u.isExternalAuth && u.currentUser.Username != "admin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.customAbort(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}

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
		u.customAbort(http.StatusNotFound, "No found user with provided User ID.")
		return
	}

	if !(u.isSysAdmin || u.currentUser.ID == user.ID) {
		u.customAbort(http.StatusForbidden, "Only system admin can change others' password.")
		return
	}

	var changePassword model.ChangePassword
	u.resolveBody(&changePassword)

	changePassword.OldPassword = utils.Encrypt(changePassword.OldPassword, u.currentUser.Salt)

	if changePassword.OldPassword != user.Password {
		u.customAbort(http.StatusForbidden, "Old password input is incorrect.")
		return
	}
	if !utils.ValidateWithLengthRange(changePassword.NewPassword, 8, 20) {
		u.customAbort(http.StatusBadRequest, "Password does not satisfy complexity requirement.")
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
		u.customAbort(http.StatusBadRequest, "Failed to change password.")
	}
}

type SystemAdminController struct {
	BaseController
}

func (u *SystemAdminController) Prepare() {
	u.resolveSignedInUser()
	if !u.isSysAdmin {
		u.customAbort(http.StatusForbidden, "Insufficient privileges to manipulate user.")
		return
	}
	u.isExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

func (u *SystemAdminController) AddUserAction() {

	if u.isExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.customAbort(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}
	var reqUser model.User
	var err error
	u.resolveBody(&reqUser)

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
	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.customAbort(http.StatusBadRequest, "Email content is illegal.")
		return
	}
	emailExists, err := service.UserExists("email", reqUser.Email, 0)
	if err != nil {
		u.internalError(err)
		return
	}
	if emailExists {
		u.customAbort(http.StatusConflict, "Email already exists.")
		return
	}

	if !utils.ValidateWithLengthRange(reqUser.Password, 8, 20) {
		u.customAbort(http.StatusBadRequest, "Password does not satisfy complexity requirement.")
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
		u.customAbort(http.StatusBadRequest, "Failed to sign up user.")
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
		u.customAbort(http.StatusNotFound, "No user found with provided User ID.")
		return
	}
	user.Password = ""
	u.Data["json"] = user
	u.ServeJSON()
}

func (u *SystemAdminController) DeleteUserAction() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid user ID: %d", userID))
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if user == nil {
		u.customAbort(http.StatusNotFound, "No user was found with provided ID.")
		return
	}
	if userID == 1 || int64(userID) == u.currentUser.ID {
		u.customAbort(http.StatusBadRequest, "System admin user or current user cannot be deleted.")
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

	if u.isExternalAuth && u.currentUser.Username != "admin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.customAbort(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}

	var err error
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.serveStatus(http.StatusBadRequest, "Invalid user ID.")
		return
	}
	if u.currentUser.ID == int64(userID) {
		u.customAbort(http.StatusForbidden, "Insufficient privileges.")
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.internalError(err)
		return
	}
	if user == nil {
		u.customAbort(http.StatusNotFound, "No user was found with provided ID.")
		return
	}

	var reqUser model.User
	u.resolveBody(&reqUser)

	reqUser.ID = user.ID
	users, err := service.GetUsers("email", reqUser.Email)
	if err != nil {
		u.internalError(err)
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.customAbort(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	if len(users) > 0 && users[0].ID != reqUser.ID {
		u.customAbort(http.StatusConflict, "Email already exists.")
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

func (u *SystemAdminController) ToggleSystemAdminAction() {

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
		u.customAbort(http.StatusNotFound, "No found user with provided user ID.")
		return
	}
	if u.currentUser.ID == user.ID {
		u.customAbort(http.StatusBadRequest, "Self system admin cannot be changed.")
		return
	}

	var reqUser model.User
	u.resolveBody(&reqUser)

	user.SystemAdmin = reqUser.SystemAdmin
	isSuccess, err := service.UpdateUser(*user, "system_admin")
	if err != nil {
		u.internalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbort(http.StatusBadRequest, "Failed to toggle user system admin.")
	}
}
