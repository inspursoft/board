package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
	"strings"
)

func SignInAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodPost) {
		reqData := c.resolveBody()
		if reqData != nil {
			var reqUser model.User
			err := json.Unmarshal(reqData, &reqUser)
			if err != nil {
				c.internalError(err)
				return
			}
			isSuccess, err := service.SignIn(reqUser.Username, reqUser.Password)
			if err != nil {
				c.internalError(err)
				return
			}
			if !isSuccess {
				c.serveStatus(http.StatusBadRequest, "Incorrect username or password.")
			}
		}
	}
}

func addUser(c *controller) {
	reqData := c.resolveBody()
	if reqData != nil {
		var reqUser model.User
		err := json.Unmarshal(reqData, &reqUser)
		if err != nil {
			c.internalError(err)
			return
		}
		if strings.TrimSpace(reqUser.Username) == "" {
			c.serveStatus(400, "Username cannot be empty.")
			return
		}
		if strings.TrimSpace(reqUser.Email) == "" {
			c.serveStatus(400, "Email cannot be empty.")
			return
		}
		usernameExists, err := service.UsernameExists(reqUser.Username)
		if err != nil {
			c.internalError(err)
			return
		}
		if usernameExists {
			c.serveStatus(409, "Username already exists.")
			return
		}
		emailExists, err := service.EmailExists(reqUser.Email)
		if err != nil {
			c.internalError(err)
			return
		}
		if emailExists {
			c.serveStatus(409, "Email already exists.")
			return
		}
		isSuccess, err := service.SignUp(reqUser)
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSuccess {
			c.serveStatus(http.StatusBadRequest, "Failed to sign up user.")
		}
	}
}

func SignUpAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodPost) {
		addUser(c)
	}
}

func AddUserAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodPost) {
		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}

		isSysAdmin, err := service.IsSysAdmin(currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSysAdmin {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}
		addUser(c)
	}
}

func GetUsersAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodGet) {
		var err error

		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}

		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}

		isSysAdmin, err := service.IsSysAdmin(currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSysAdmin {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		var users []*model.User

		username := req.FormValue("username")
		email := req.FormValue("email")
		if strings.TrimSpace(username) != "" {
			users, err = service.GetUsers("username", username)
		} else if strings.TrimSpace(email) != "" {
			users, err = service.GetUsers("email", email)
		} else {
			users, err = service.GetUsers("", nil)
		}
		if err != nil {
			c.internalError(err)
			return
		}
		for _, u := range users {
			u.Password = ""
		}
		c.serveJSON(users)
	}
}

func OperateUserAction(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getUserAction(resp, req)
	case http.MethodPut:
		updateUserAction(resp, req)
	case http.MethodDelete:
		deleteUserAction(resp, req)
	}
}

func getUserAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("GET") {
		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}
		userID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}
		user, err := service.GetUserByID(int64(userID))
		if err != nil {
			c.internalError(err)
			return
		}
		if user == nil {
			c.customAbort(http.StatusNotFound, "No user found with provided User ID.")
			return
		}
		user.Password = ""
		c.serveJSON(user)
	}
}

func deleteUserAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodDelete) {
		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}

		isSysAdmin, err := service.IsSysAdmin(currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSysAdmin {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		userID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.customAbort(http.StatusBadRequest, "Invalid user ID.")
			return
		}

		user, err := service.GetUserByID(int64(userID))
		if err != nil {
			c.internalError(err)
			return
		}

		if user == nil {
			c.customAbort(http.StatusNotFound, "No user was found with provided ID.")
			return
		}

		if userID == 1 || int64(userID) == currentUser.ID {
			c.customAbort(http.StatusBadRequest, "System admin user or current user cannot be deleted.")
			return
		}

		isSuccess, err := service.DeleteUser(int64(userID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSuccess {
			c.serveStatus(http.StatusBadRequest, "Failed to delete user.")
		}
	}
}

func updateUserAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodPut) {
		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}

		isSystemAdmin, err := service.IsSysAdmin(currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}

		userID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.serveStatus(http.StatusBadRequest, "Invalid user ID.")
			return
		}

		if !(isSystemAdmin || currentUser.ID == int64(userID)) {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		user, err := service.GetUserByID(int64(userID))
		if err != nil {
			c.internalError(err)
			return
		}

		if user == nil {
			c.customAbort(http.StatusNotFound, "No user was found with provided ID.")
			return
		}

		reqData := c.resolveBody()
		if reqData != nil {
			var reqUser model.User
			err := json.Unmarshal(reqData, &reqUser)
			if err != nil {
				c.internalError(err)
				return
			}

			reqUser.ID = user.ID

			if strings.TrimSpace(reqUser.Email) != "" {
				user.Email = reqUser.Email
			}

			user.Realname = reqUser.Realname
			user.Comment = reqUser.Comment

			isSuccess, err := service.UpdateUser(reqUser, "email", "realname", "comment")

			if err != nil {
				c.internalError(err)
				return
			}

			if !isSuccess {
				c.serveStatus(http.StatusBadRequest, "Failed to update user.")
			}
		}
	}
}

func ChangePasswordAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("PUT") {
		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}

		userID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}

		user, err := service.GetUserByID(int64(userID))
		if err != nil {
			c.internalError(err)
			return
		}

		if user == nil {
			c.customAbort(http.StatusNotFound, "No found user with provided User ID.")
			return
		}

		isSystemAdmin, err := service.IsSysAdmin(currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}

		if !(isSystemAdmin || currentUser.ID == user.ID) {
			c.customAbort(http.StatusForbidden, "Only system admin can change others' password.")
			return
		}

		reqData := c.resolveBody()
		if reqData != nil {
			var changePassword model.ChangePassword
			err := json.Unmarshal(reqData, &changePassword)
			if err != nil {
				c.internalError(err)
				return
			}
			if changePassword.OldPassword != user.Password {
				c.customAbort(http.StatusForbidden, "Old password input is incorrect.")
				return
			}
			if changePassword.NewPassword == "" {
				c.customAbort(http.StatusBadRequest, "New password cannot be empty.")
				return
			}
			updateUser := model.User{
				ID:       user.ID,
				Password: changePassword.NewPassword,
			}
			isSuccess, err := service.UpdateUser(updateUser, "password")
			if err != nil {
				c.internalError(err)
				return
			}
			if !isSuccess {
				c.customAbort(http.StatusBadRequest, "Failed to change password")
			}
		}
	}
}

func ToggleSystemAdminAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("PUT") {
		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}
		isSysAdmin, err := service.IsSysAdmin(currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSysAdmin {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		userID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}

		user, err := service.GetUserByID(int64(userID))
		if err != nil {
			c.internalError(err)
			return
		}
		if user == nil {
			c.customAbort(http.StatusNotFound, "No found user with provided user ID.")
			return
		}

		if currentUser.ID == user.ID {
			c.customAbort(http.StatusBadRequest, "Self system admin cannot be changed.")
			return
		}

		reqData := c.resolveBody()
		if reqData != nil {
			var reqUser model.User
			err := json.Unmarshal(reqData, &reqUser)
			if err != nil {
				c.internalError(err)
				return
			}
			user.SystemAdmin = reqUser.SystemAdmin
			isSuccess, err := service.UpdateUser(*user, "system_admin")
			if err != nil {
				c.internalError(err)
				return
			}
			if !isSuccess {
				c.customAbort(http.StatusBadRequest, "Failed to toggle user system admin.")
			}
		}
	}
}
