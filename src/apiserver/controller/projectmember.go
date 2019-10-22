package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
)

type ProjectMemberController struct {
	BaseController
}

func (pm *ProjectMemberController) AddOrUpdateProjectMemberAction() {

	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":id"))
	if err != nil {
		pm.internalError(err)
		return
	}
	pm.resolveProjectOwnerByID(int64(projectID))

	var reqProjectMember model.ProjectMember
	err = pm.resolveBody(&reqProjectMember)
	if err != nil {
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
		return
	}
	baseRepoName := pm.project.Name
	service.ForkRepo(user, baseRepoName)
}

func (pm *ProjectMemberController) DeleteProjectMemberAction() {

	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":projectId"))
	if err != nil {
		pm.internalError(err)
		return
	}

	pm.resolveProjectOwnerByID(int64(projectID))

	userID, err := strconv.Atoi(pm.Ctx.Input.Param(":userId"))
	if err != nil {
		pm.internalError(err)
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
	pm.resolveProjectMemberByID(int64(projectID))
	projectMembers, err := service.GetProjectMembers(int64(projectID))
	if err != nil {
		pm.internalError(err)
		return
	}
	pm.renderJSON(projectMembers)
}
