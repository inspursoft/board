package members

import (
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"net/http"
	"strconv"

	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
)

// Operations about project members
type CommonController struct {
	c.BaseController
}

// @Title List all members of project
// @Description List all for projects.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	true	"ID of projects"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [get]
func (pm *CommonController) List() {
	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":project_id"))
	if err != nil {
		pm.InternalError(err)
		return
	}
	pm.ResolveProjectMemberByID(int64(projectID))
	projectMembers, err := service.GetProjectMembers(int64(projectID))
	if err != nil {
		pm.InternalError(err)
		return
	}
	pm.RenderJSON(projectMembers)
}

// @Title Add member to the project
// @Description Add member to the project.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	false	"ID of projects"
// @Param	body	body	vm.ProjectMember	true	"View model of project member."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [post]
func (pm *CommonController) Add() {
	pm.Update()
}

// @Title Update project by ID
// @Description Update project by ID.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	false	"ID of projects"
// @Param	body	body	vm.ProjectMember	true	"View model of project member."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [put]
func (pm *CommonController) Update() {
	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":project_id"))
	if err != nil {
		pm.InternalError(err)
		return
	}
	pm.ResolveProjectOwnerByID(int64(projectID))

	var reqProjectMember vm.ProjectMember
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

// @Title Delete project member by project and user ID
// @Description Delete project member by project and user ID
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	true	"ID of projects"
// @Param	user_id	path	int	true	"ID of users"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:user_id [delete]
func (pm *CommonController) Delete() {
	projectID, err := strconv.Atoi(pm.Ctx.Input.Param(":project_id"))
	if err != nil {
		pm.InternalError(err)
		return
	}

	pm.ResolveProjectOwnerByID(int64(projectID))

	userID, err := strconv.Atoi(pm.Ctx.Input.Param(":user_id"))
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
