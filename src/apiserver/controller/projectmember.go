package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/gogs"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego/logs"
)

type ProjectMemberController struct {
	baseController
}

func (pm *ProjectMemberController) AddOrUpdateProjectMemberAction() {

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

	var reqProjectMember model.ProjectMember
	pm.resolveBody(&reqProjectMember)

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

	isMember, err := service.IsProjectMember(int64(projectID), pm.currentUser.ID)
	if err != nil {
		pm.internalError(err)
		return
	}

	if !(isMember || pm.isSysAdmin) {
		pm.customAbort(http.StatusForbidden, "User neither has no member to this project nor isn't a system admin.")
		return
	}

	queryProject := model.Project{ID: int64(projectID)}
	project, err := service.GetProject(queryProject, "id")
	if err != nil {
		pm.internalError(err)
		return
	}
	if !(pm.isSysAdmin || int64(project.OwnerID) == pm.currentUser.ID) {
		pm.customAbort(http.StatusForbidden, "User is not the owner of the project.")
		return
	}

	err = gogs.NewGogsHandler(user.Username, user.RepoToken).ForkRepo(project.OwnerName, project.Name, project.Name, "Forked repo.")
	if err != nil {
		pm.internalError(err)
	}

	projectRepoURL := fmt.Sprintf("%s/%s/%s.git", gogitsSSHURL(), user.Username, project.Name)
	projectRepoPath := filepath.Join(baseRepoPath(), user.Username, project.Name)
	repoHandler, err := service.InitRepo(projectRepoURL, user.Username, user.Email, projectRepoPath)
	if err != nil {
		logs.Error("Failed to initialize project repo: %+v", err)
		pm.internalError(err)
	}
	repoHandler.Pull()

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

	user, err := service.GetUserByID(int64(userID))
	if err != nil {
		pm.internalError(err)
		return
	}
	if user == nil {
		pm.customAbort(http.StatusNotFound, "No user was found with provided user ID.")
		return
	}

	query := model.Project{ID: int64(projectID)}
	project, err := service.GetProject(query, "id")
	if err != nil {
		pm.internalError(err)
		return
	}

	if project.OwnerID == int(user.ID) {
		pm.customAbort(http.StatusForbidden, "Project owner cannnot be removed.")
		return
	}

	isMember, err := service.IsProjectMember(int64(projectID), pm.currentUser.ID)
	if err != nil {
		pm.internalError(err)
		return
	}

	if !(isMember || pm.isSysAdmin) {
		pm.customAbort(http.StatusForbidden, "User neither has no member to this project nor isn't a system admin.")
		return
	}

	if !(pm.isSysAdmin || int64(project.OwnerID) == pm.currentUser.ID) {
		pm.customAbort(http.StatusForbidden, "User is not the owner of the project.")
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
