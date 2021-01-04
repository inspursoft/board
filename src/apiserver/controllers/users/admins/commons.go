package admins

import (
	"fmt"
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/service/adapting"
	"github.com/inspursoft/board/src/apiserver/service/devops/gogs"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"

	c "github.com/inspursoft/board/src/apiserver/controllers/commons"

	"github.com/astaxie/beego/logs"
)

// Operations about admins
type CommonController struct {
	c.BaseController
}

func (u *CommonController) Prepare() {
	u.EnableXSRF = false
	u.ResolveSignedInUser()
	u.RecordOperationAudit()
	if !u.IsSysAdmin {
		u.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to manipulate user.")
		return
	}
	u.IsExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

// @Title List all users by admin
// @Description List all for users.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	username	query	string	false	"Query username for users."
// @Param	email	query	string	false	"Query email for users."
// @Param	page_index query	int	false	"Page index for pagination."
// @Param	page_size	query	int	false	"Page per size for pagination."
// @Param	order_field	string	false	"Order by field. (Default is creation_time)"
// @Param	order_asc	int	false	"Order option for ascend or descend. (asc 0, desc, 1)"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (u *CommonController) List() {
	username := u.GetString("username")
	email := u.GetString("email")
	pageIndex, _ := u.GetInt("page_index", 0)
	pageSize, _ := u.GetInt("page_size", 0)
	isPaginated := !(pageIndex == 0 && pageSize == 0)
	orderField := u.GetString("order_field", "creation_time")
	orderAsc, _ := u.GetInt("order_asc", 0)

	fieldName := "deleted"
	var fieldValue interface{}
	if strings.TrimSpace(username) != "" {
		fieldName = "username"
		fieldValue = username
	} else if strings.TrimSpace(email) != "" {
		fieldName = "email"
		fieldValue = email
	}
	if isPaginated {
		paginatedUsers, err := service.GetPaginatedUsers(fieldName, fieldValue, pageIndex, pageSize, orderField, orderAsc)
		if err != nil {
			u.InternalError(err)
			return
		}
		u.Data["json"] = paginatedUsers
	} else {
		users, err := service.GetUsers(fieldName, fieldValue, "id", "username", "email", "realname")
		if err != nil {
			u.InternalError(err)
			return
		}
		u.Data["json"] = users
	}

	u.ServeJSON()
}

// @Title Get user by ID for admin
// @Description Get user by ID for admin.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	user_id	path	int	true	"Query user ID for users."
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id [get]
func (u *CommonController) Get() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":user_id"))
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

// @Title Add user by admin
// @Description Add user by admin.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	body	body 	"vm.User"	true	"View model for users."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (u *CommonController) Add() {
	if u.IsExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}
	var reqUser vm.User
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

	isSuccess, err := adapting.SignUp(reqUser)
	if err != nil {
		u.InternalError(err)
		return
	}
	if !isSuccess {
		u.CustomAbortAudit(http.StatusBadRequest, "Failed to sign up user.")
	}
}

// @Title Update user info by admin
// @Description Update user info by admin.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	user_id	path	int	false	"ID of users"
// @Param	body	body 	"vm.User"	true	"View model for users."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id [put]
func (u *CommonController) Update() {
	if u.IsExternalAuth {
		logs.Debug("Current AUTH_MODE is external auth.")
		u.CustomAbortAudit(http.StatusMethodNotAllowed, "Current AUTH_MODE is external auth.")
		return
	}
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":user_id"))
	if err != nil {
		u.InternalError(err)
		return
	}
	var reqUser model.User
	err = u.ResolveBody(&reqUser)
	if err != nil {
		return
	}
	reqUser.ID = int64(userID)

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

// @Title Delete user info by admin
// @Description Delete user info by admin.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	user_id	path	int	false	"ID of users"
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:user_id [delete]
func (u *CommonController) Delete() {
	userID, err := strconv.Atoi(u.Ctx.Input.Param(":user_id"))
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
	adminUser, err := service.GetUserByID(1)
	if err != nil {
		u.InternalError(err)
		return
	}
	err = gogs.NewGogsHandler(adminUser.Username, adminUser.RepoToken).DeleteUser(user.Username)
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
