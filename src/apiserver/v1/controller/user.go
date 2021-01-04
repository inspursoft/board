package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

type UserController struct {
	c.BaseController
}

func (u *UserController) Prepare() {
	u.EnableXSRF = false
	u.ResolveSignedInUser()
	u.RecordOperationAudit()
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

func (u *UserController) GetUsersAction() {
	if !u.IsSysAdmin {
		u.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to get users info.")
		return
	}
	username := u.GetString("username")
	email := u.GetString("email")
	pageIndex, _ := u.GetInt("page_index", 0)
	pageSize, _ := u.GetInt("page_size", 0)
	isPaginated := !(pageIndex == 0 && pageSize == 0)
	orderField := u.GetString("order_field", "creation_time")
	orderAsc, _ := u.GetInt("order_asc", 0)

	var paginatedUsers *model.PaginatedUsers
	var users []*model.User
	var err error
	fieldName := "deleted"
	var fieldValue interface{} = 0
	if strings.TrimSpace(username) != "" {
		fieldName = "username"
		fieldValue = username
	} else if strings.TrimSpace(email) != "" {
		fieldName = "email"
		fieldValue = email
	}
	if isPaginated {
		selectedFields := []string{"id", "username", "email", "deleted", "realname", "comment", "creation_time", "update_time", "reset_uuid", "system_admin"}
		paginatedUsers, err = service.GetPaginatedUsers(fieldName, fieldValue, pageIndex, pageSize, orderField, orderAsc, selectedFields...)
		u.Data["json"] = paginatedUsers
	} else {
		users, err = service.GetUsers(fieldName, fieldValue, "id", "username", "email", "realname")
		u.Data["json"] = users
	}
	if err != nil {
		u.InternalError(err)
		return
	}
	u.ServeJSON()
}

func (u *UserController) ChangeUserAccount() {
	if u.IsExternalAuth && u.CurrentUser.Username != "boardadmin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusPreconditionFailed, "Current AUTH_MODE is not available to the user.")
		return
	}

	var reqUser model.User
	var err error
	err = u.ResolveBody(&reqUser)
	if err != nil {
		return
	}

	reqUser.ID = u.CurrentUser.ID
	users, err := service.GetUsers("email", reqUser.Email, "id", "email")
	if err != nil {
		u.InternalError(err)
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbortAudit(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	if len(users) > 0 && users[0].ID != reqUser.ID {
		u.CustomAbortAudit(http.StatusConflict, "Email already exists.")
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

	reqUser.Email = strings.TrimSpace(reqUser.Email)
	reqUser.Realname = strings.TrimSpace(reqUser.Realname)
	reqUser.Comment = strings.TrimSpace(reqUser.Comment)

	isSuccess, err := service.UpdateUser(reqUser, "email", "realname", "comment")
	if err != nil {
		u.InternalError(err)
		return
	}

	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to change user account.")
	}
}

func (u *UserController) ChangePasswordAction() {
	var err error

	if u.IsExternalAuth && u.CurrentUser.Username != "boardadmin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}

	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
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
		u.CustomAbortAudit(http.StatusNotFound, "No found user with provided User ID.")
		return
	}

	if !(u.IsSysAdmin || u.CurrentUser.ID == user.ID) {
		u.CustomAbortAudit(http.StatusForbidden, "Only system admin can change others' password.")
		return
	}

	var changePassword model.ChangePassword
	err = u.ResolveBody(&changePassword)
	if err != nil {
		return
	}

	changePassword.OldPassword, err = service.DecodeUserPassword(changePassword.OldPassword)
	if err != nil {
		logs.Error("Password encode error %v", err)
		u.CustomAbortAudit(http.StatusBadRequest, "Password encode error.")
		return
	}

	changePassword.OldPassword = utils.Encrypt(changePassword.OldPassword, u.CurrentUser.Salt)

	if changePassword.OldPassword != user.Password {
		u.CustomAbortAudit(http.StatusForbidden, "Old password input is incorrect.")
		return
	}

	changePassword.NewPassword, err = service.DecodeUserPassword(changePassword.NewPassword)
	if err != nil {
		logs.Error("Password encode error %v", err)
		u.CustomAbortAudit(http.StatusBadRequest, "Password encode error.")
		return
	}

	if !utils.ValidateWithLengthRange(changePassword.NewPassword, 8, 20) {
		u.CustomAbortAudit(http.StatusBadRequest, "Password does not satisfy complexity requirement.")
		return
	}
	updateUser := model.User{
		ID:       user.ID,
		Password: utils.Encrypt(changePassword.NewPassword, u.CurrentUser.Salt),
	}
	isSuccess, err := service.UpdateUser(updateUser, "password")
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to change password.")
	}
}

type SystemAdminController struct {
	c.BaseController
}

func (u *SystemAdminController) Prepare() {
	u.EnableXSRF = true
	u.ResolveSignedInUser()
	u.RecordOperationAudit()
	if !u.IsSysAdmin {
		u.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to manipulate user.")
		return
	}
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

func (u *SystemAdminController) AddUserAction() {

	if u.IsExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}
	var reqUser model.User
	var err error
	err = u.ResolveBody(&reqUser)
	if err != nil {
		return
	}

	if !utils.ValidateWithPattern("username", reqUser.Username) {
		u.CustomAbortAudit(http.StatusBadRequest, "Username content is illegal.")
		return
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
	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbortAudit(http.StatusBadRequest, "Email content is illegal.")
		return
	}
	emailExists, err := service.UserExists("email", reqUser.Email, 0)
	if err != nil {
		u.InternalError(err)
		return
	}
	if emailExists {
		u.CustomAbortAudit(http.StatusConflict, "Email already exists.")
		return
	}

	reqUser.Password, err = service.DecodeUserPassword(reqUser.Password)
	if err != nil {
		logs.Error("Password encode error %v", err)
		u.CustomAbortAudit(http.StatusBadRequest, "Password encode error.")
		return
	}

	if !utils.ValidateWithLengthRange(reqUser.Password, 8, 20) {
		u.CustomAbortAudit(http.StatusBadRequest, "Password does not satisfy complexity requirement.")
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

	isSuccess, err := service.SignUp(reqUser)
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to sign up user.")
	}
}

func (u *SystemAdminController) GetUserAction() {

	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
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
		u.CustomAbortAudit(http.StatusNotFound, "No user found with provided User ID.")
		return
	}
	user.Password = ""
	u.Data["json"] = user
	u.ServeJSON()
}

func (u *SystemAdminController) DeleteUserAction() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Invalid user ID: %d", userID))
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.InternalError(err)
		return
	}
	if user == nil {
		u.CustomAbortAudit(http.StatusNotFound, "No user was found with provided ID.")
		return
	}
	if userID == 1 || int64(userID) == u.CurrentUser.ID {
		u.CustomAbortAudit(http.StatusBadRequest, "System admin user or current user cannot be deleted.")
		return
	}
	err = service.CurrentDevOps().DeleteUser(user.Username)
	if err != nil {
		if err == utils.ErrUnprocessableEntity {
			u.CustomAbortAudit(http.StatusUnprocessableEntity, "User has own project or repo.")
		} else {
			u.InternalError(err)
		}
		return
	}

	isSuccess, err := service.DeleteUser(int64(userID))
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.ServeStatus(http.StatusBadRequest, "Failed to delete user.")
	}
}

func (u *SystemAdminController) UpdateUserAction() {

	if u.IsExternalAuth && u.CurrentUser.Username != "boardadmin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}

	var err error
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil {
		u.ServeStatus(http.StatusBadRequest, "Invalid user ID.")
		return
	}
	if userID == 1 || u.CurrentUser.ID == int64(userID) {
		u.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges.")
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		u.InternalError(err)
		return
	}
	if user == nil {
		u.CustomAbortAudit(http.StatusNotFound, "No user was found with provided ID.")
		return
	}

	var reqUser model.User
	err = u.ResolveBody(&reqUser)
	if err != nil {
		return
	}

	reqUser.ID = user.ID
	users, err := service.GetUsers("email", reqUser.Email, "id", "email")
	if err != nil {
		u.InternalError(err)
		return
	}

	if !utils.ValidateWithPattern("email", reqUser.Email) {
		u.CustomAbortAudit(http.StatusBadRequest, "Email content is illegal.")
		return
	}

	if len(users) > 0 && users[0].ID != reqUser.ID {
		u.CustomAbortAudit(http.StatusConflict, "Email already exists.")
		return
	}

	user.Email = reqUser.Email
	user.Realname = reqUser.Realname
	user.Comment = reqUser.Comment

	isSuccess, err := service.UpdateUser(reqUser, "email", "realname", "comment")
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.ServeStatus(http.StatusBadRequest, "Failed to update user.")
	}
}

func (u *SystemAdminController) ToggleSystemAdminAction() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
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

	var reqUser model.User
	err = u.ResolveBody(&reqUser)
	if err != nil {
		return
	}

	user.SystemAdmin = reqUser.SystemAdmin
	isSuccess, err := service.UpdateUser(*user, "system_admin")
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to toggle user system admin.")
	}
}
