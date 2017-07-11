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
			var sessionUserID int64 = 1
			projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
			if err != nil {
				c.internalError(err)
				return
			}
			isSuccess, err := service.AddOrUpdateProjectMember(int64(projectID), sessionUserID, reqProjectMember.RoleID)
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
		var sessionUserID int64 = 1
		projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}
		isSuccess, err := service.DeleteProjectMember(int64(projectID), sessionUserID)
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSuccess {
			c.customAbort(http.StatusBadRequest, "Failed to delete project member.")
		}
	}
}

func getProjectMembersAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("GET") {
		var sessionUserID int64 = 1
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

		projectMembers, err := service.GetProjectMembers(int64(projectID), sessionUserID)
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
