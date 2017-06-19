package controller

import (
	"apiserver/model"
	"apiserver/service"
	"encoding/json"
	"fmt"
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

func SignUpAction(resp http.ResponseWriter, req *http.Request) {
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
}

func GetUsersAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodGet) {
		var users []*model.User
		var err error
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
		c.serveJSON(users)
	}
}

func OperateUserAction(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		updateUserAction(resp, req)
	case http.MethodDelete:
		deleteUserAction(resp, req)
	}
}

func deleteUserAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod(http.MethodDelete) {
		userID, err := strconv.Atoi(c.GetStringFromPath(apiprefix + "/users"))
		if err != nil {
			c.serveStatus(http.StatusBadRequest, "Invalid user ID.")
			return
		}
		fmt.Printf("UserID for deletion: %d\n", userID)
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
		reqData := c.resolveBody()
		if reqData != nil {
			var reqUser model.User
			err := json.Unmarshal(reqData, &reqUser)
			if err != nil {
				c.internalError(err)
				return
			}
			userID, err := strconv.Atoi(c.GetStringFromPath(apiprefix + "/users"))
			if err != nil {
				c.serveStatus(http.StatusBadRequest, "Invalid user ID.")
				return
			}
			fmt.Printf("UserID for update: %d\n", userID)
			reqUser.ID = int64(userID)
			isSuccess, err := service.UpdateUser(reqUser)
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
