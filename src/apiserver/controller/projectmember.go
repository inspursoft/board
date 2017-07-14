package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
)

func addOrUpdateProjectMemberAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("POST") {

		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
		}

		projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}

		isExists, err := service.ProjectExistsByID(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !isExists {
			c.customAbort(http.StatusNotFound, "Cannot find project by ID")
			return
		}

		hasPermission, err := checkUserChangePermission(int64(projectID), currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}

		if !hasPermission {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		reqData := c.resolveBody()
		if reqData != nil {
			var reqProjectMember model.ProjectMember
			err := json.Unmarshal(reqData, &reqProjectMember)
			if err != nil {
				c.internalError(err)
				return
			}
			if reqProjectMember.RoleID == 0 {
				c.customAbort(http.StatusBadRequest, "Invalid role ID for project.")
				return
			}
			if reqProjectMember.UserID == 0 {
				c.customAbort(http.StatusBadRequest, "Invalid user ID.")
				return
			}

			isSuccess, err := service.AddOrUpdateProjectMember(int64(projectID), reqProjectMember.UserID, reqProjectMember.RoleID)
			if err != nil {
				c.internalError(err)
				return
			}
			if !isSuccess {
				c.customAbort(http.StatusBadRequest, "Failed to add or upate project member.")
			}
		}
	}
}

func deleteProjectMemberAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("DELETE") {
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

		projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}

		isExists, err := service.ProjectExistsByID(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !isExists {
			c.customAbort(http.StatusNotFound, "Cannot find project by ID")
			return
		}

		hasPermission, err := checkUserPermission(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !hasPermission {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		var reqProjectMember model.ProjectMember

		reqData := c.resolveBody()
		if reqData != nil {
			err = json.Unmarshal(reqData, &reqProjectMember)
			if err != nil {
				c.internalError(err)
				return
			}
			if reqProjectMember.UserID == 0 {
				c.customAbort(http.StatusBadRequest, "Invalid project member user ID.")
				return
			}

			user, err := service.GetUserByID(reqProjectMember.UserID)
			if err != nil {
				c.internalError(err)
				return
			}
			if user == nil {
				c.customAbort(http.StatusNotFound, "No user was found with provided user ID.")
				return
			}

			if reqProjectMember.UserID == currentUser.ID {
				c.customAbort(http.StatusConflict, "Self privilege to the current project cannot be deleted.")
				return
			}
			isSuccess, err := service.DeleteProjectMember(int64(projectID), reqProjectMember.UserID)
			if err != nil {
				c.internalError(err)
				return
			}
			if !isSuccess {
				c.customAbort(http.StatusBadRequest, "Failed to delete project member.")
			}
		}
	}
}

func getProjectMembersAction(resp http.ResponseWriter, req *http.Request) {
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

		projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}

		isExists, err := service.ProjectExistsByID(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !isExists {
			c.customAbort(http.StatusNotFound, "Cannot find project by ID")
			return
		}

		hasPermission, err := checkUserPermission(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !hasPermission {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		projectMembers, err := service.GetProjectMembers(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		c.serveJSON(projectMembers)
	}
}

func OperateProjectMembersAction(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		getProjectMembersAction(resp, req)
	case "POST":
		addOrUpdateProjectMemberAction(resp, req)
	case "DELETE":
		deleteProjectMemberAction(resp, req)
	}
}
