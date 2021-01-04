package controller

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
	"strings"
)

type ProjectMemberController struct {
	c.BaseController
}

func (pm *ProjectMemberController) AddOrUpdateProjectMemberAction() {

	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":id"))
	if err != nil {
		pm.InternalError(err)
		return
	}
	pm.ResolveProjectOwnerByID(int64(projectID))

	var reqProjectMember model.ProjectMember
	err = pm.ResolveBody(&reqProjectMember)
	if err != nil {
		return
	}

	role, err := service.GetRoleByID(reqProjectMember.RoleID)
	if err != nil {
		pm.InternalError(err)
		return
	}
	if role == nil {
		pm.CustomAbortAudit(http.StatusNotFound, "No role found with provided role ID.")
		return
	}

	user, err := service.GetUserByID(reqProjectMember.UserID)
	if err != nil {
		pm.InternalError(err)
		return
	}
	if user == nil {
		pm.CustomAbortAudit(http.StatusNotFound, "No user found with provided user ID.")
		return
	}

	isSuccess, err := service.AddOrUpdateProjectMember(int64(projectID), reqProjectMember.UserID, reqProjectMember.RoleID)
	if err != nil {
		pm.InternalError(err)
		return
	}
	if !isSuccess {
		pm.CustomAbortAudit(http.StatusBadRequest, "Failed to add or upate project member.")
		return
	}
	baseRepoName := pm.Project.Name
	service.CurrentDevOps().ForkRepo(*user, baseRepoName)
}

func (pm *ProjectMemberController) DeleteProjectMemberAction() {

	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":projectId"))
	if err != nil {
		pm.InternalError(err)
		return
	}

	pm.ResolveProjectOwnerByID(int64(projectID))

	userID, err := strconv.Atoi(pm.Ctx.Input.Param(":userId"))
	if err != nil {
		pm.InternalError(err)
		return
	}

	isSuccess, err := service.DeleteProjectMember(int64(projectID), int64(userID))
	if err != nil {
		pm.InternalError(err)
		return
	}
	if !isSuccess {
		pm.CustomAbortAudit(http.StatusBadRequest, "Failed to delete project member.")
	}
}

func (pm *ProjectMemberController) GetProjectMembersAction() {
	membersType := strings.ToLower(pm.GetString("type", "current"))
	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":id"))
	if err != nil {
		pm.InternalError(err)
		return
	}

	pm.ResolveProjectMemberByID(int64(projectID))
	if membersType == "current" {
		projectMembers, err := service.GetProjectMembers(int64(projectID))
		if err != nil {
			pm.InternalError(err)
			return
		}
		pm.RenderJSON(projectMembers)
	} else if membersType == "available" {
		availableMembers, err := service.GetProjectAvailableMembers(int64(projectID))
		if err != nil {
			pm.InternalError(err)
			return
		}
		pm.RenderJSON(availableMembers)
	} else {
		pm.CustomAbortAudit(http.StatusBadRequest, "Invalid value of the query parameter of type.")
	}
}
