package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
)

type ProjectMemberController struct {
	baseController
}

func (pm *ProjectMemberController) Prepare() {
	user := pm.getCurrentUser()
	if user == nil {
		pm.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	pm.currentUser = user
	pm.isSysAdmin = (user.SystemAdmin == 1)
	pm.isProjectAdmin = (user.ProjectAdmin == 1)
	if !pm.isProjectAdmin {
		pm.customAbort(http.StatusForbidden, "Insufficient privileges to for manipulating projects.")
		return
	}
}

func (pm *ProjectMemberController) AddOrUpdateProjectMemberAction() {
	var err error
	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":id"))
	if err != nil {
		pm.internalError(err)
		return
	}
	isExists, err := service.ProjectExistsByID(int64(projectID))
	if err != nil {
		pm.internalError(err)
		return
	}
	if !isExists {
		pm.customAbort(http.StatusNotFound, "Cannot find project by ID")
		return
	}
	reqData, err := pm.resolveBody()
	if err != nil {
		pm.internalError(err)
		return
	}

	var reqProjectMember model.ProjectMember
	err = json.Unmarshal(reqData, &reqProjectMember)
	if err != nil {
		pm.internalError(err)
		return
	}
	role, err := service.GetRoleByID(reqProjectMember.RoleID)
	if err != nil {
		pm.internalError(err)
		return
	}
	if role == nil {
		pm.customAbort(http.StatusNotFound, "No role found with provided role ID.")
		return
	}

	user, err := service.GetUserByID(reqProjectMember.UserID)
	if err != nil {
		pm.internalError(err)
		return
	}
	if user == nil {
		pm.customAbort(http.StatusNotFound, "No user found with provided user ID.")
		return
	}

	isSuccess, err := service.AddOrUpdateProjectMember(int64(projectID), reqProjectMember.UserID, reqProjectMember.RoleID)
	if err != nil {
		pm.internalError(err)
		return
	}
	if !isSuccess {
		pm.customAbort(http.StatusBadRequest, "Failed to add or upate project member.")
	}
}

func (pm *ProjectMemberController) DeleteProjectMemberAction() {
	var err error
	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":projectId"))
	if err != nil {
		pm.internalError(err)
		return
	}
	isExists, err := service.ProjectExistsByID(int64(projectID))
	if err != nil {
		pm.internalError(err)
		return
	}
	if !isExists {
		pm.customAbort(http.StatusNotFound, "Cannot find project by ID")
		return
	}

	userID, err := strconv.Atoi(pm.Ctx.Input.Param(":userId"))
	if err != nil {
		pm.internalError(err)
		return
	}
	if int64(userID) == pm.currentUser.ID {
		pm.customAbort(http.StatusConflict, "Self privilege to the current project cannot be deleted.")
		return
	}
	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		pm.internalError(err)
		return
	}
	if user == nil {
		pm.customAbort(http.StatusNotFound, "No user was found with provided user ID.")
		return
	}
	isSuccess, err := service.DeleteProjectMember(int64(projectID), int64(userID))
	if err != nil {
		pm.internalError(err)
		return
	}
	if !isSuccess {
		pm.customAbort(http.StatusBadRequest, "Failed to delete project member.")
	}
}

func (pm *ProjectMemberController) GetProjectMembersAction() {
	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":id"))
	if err != nil {
		pm.internalError(err)
		return
	}
	isExists, err := service.ProjectExistsByID(int64(projectID))
	if err != nil {
		pm.internalError(err)
		return
	}
	if !isExists {
		pm.customAbort(http.StatusNotFound, "Cannot find project by ID")
		return
	}
	projectMembers, err := service.GetProjectMembers(int64(projectID))
	if err != nil {
		pm.internalError(err)
		return
	}
	pm.Data["json"] = projectMembers
	pm.ServeJSON()
}
