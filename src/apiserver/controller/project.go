package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

const (
	defaultProject = "library"
)

type ProjectController struct {
	baseController
}

func (p *ProjectController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}

func (p *ProjectController) CreateProjectAction() {
	if !p.isProjectAdmin {
		p.customAbort(http.StatusForbidden, "Insufficient privileges for creating projects.")
		return
	}
	var err error
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}
	var reqProject model.Project
	err = json.Unmarshal(reqData, &reqProject)
	if err != nil {
		p.internalError(err)
		return
	}
	if !utils.ValidateWithLengthRange(reqProject.Name, 2, 30) {
		p.customAbort(http.StatusBadRequest, "Project name length should be between 2 and 30 characters.")
		return
	}
	if !utils.ValidateWithPattern("project", reqProject.Name) {
		p.customAbort(http.StatusBadRequest, "Project name is invalid.")
		return
	}

	projectExists, err := service.ProjectExists(reqProject.Name)
	if err != nil {
		p.internalError(err)
		return
	}
	if projectExists {
		p.customAbort(http.StatusConflict, "Project name already exists.")
		return
	}

	reqProject.Name = strings.TrimSpace(reqProject.Name)

	reqProject.OwnerID = int(p.currentUser.ID)
	reqProject.OwnerName = p.currentUser.Username

	isSuccess, err := service.CreateProject(reqProject)
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Project name: %s is illegal.", reqProject.Name))
	}

}

func (p *ProjectController) ProjectExists() {
	projectName := p.GetString("project_name")
	query := model.Project{Name: projectName}
	project, err := service.GetProject(query, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if project != nil {
		p.customAbort(http.StatusConflict, fmt.Sprintf("Project name: %s already exists.", projectName))
	}
}

func (p *ProjectController) GetProjectsAction() {
	projectName := p.GetString("project_name")
	strPublic := p.GetString("project_public")

	query := model.Project{Name: projectName, OwnerName: p.currentUser.Username, Public: 0}

	var err error
	public, err := strconv.Atoi(strPublic)
	if err == nil {
		query.Public = public
	}

	projects, err := service.GetProjectsByUser(query, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	p.Data["json"] = projects
	p.ServeJSON()
}

func (p *ProjectController) GetProjectAction() {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}
	projectQuery := model.Project{ID: int64(projectID), Deleted: 0}
	project, err := service.GetProject(projectQuery, "id", "deleted")
	if err != nil {
		p.internalError(err)
		return
	}
	if project == nil {
		p.customAbort(http.StatusNotFound, fmt.Sprintf("No project was found with provided ID: %d", projectID))
		return
	}
	p.Data["json"] = project
	p.ServeJSON()
}

func (p *ProjectController) DeleteProjectAction() {
	if !p.isProjectAdmin {
		p.customAbort(http.StatusForbidden, "Insuffient privileges for creating projects.")
		return
	}

	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}
	isExists, err := service.ProjectExistsByID(int64(projectID))
	if err != nil {
		p.internalError(err)
		return
	}
	if !isExists {
		p.customAbort(http.StatusNotFound, fmt.Sprintf("Cannot find project with ID: %d", projectID))
		return
	}
	isSuccess, err := service.DeleteProject(int64(projectID))
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, "Failed to delete project.")
	}
}

func (p *ProjectController) ToggleProjectPublicAction() {
	if !p.isProjectAdmin {
		p.customAbort(http.StatusForbidden, "Insufficient privileges for creating projects.")
		return
	}

	var err error
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}
	isExists, err := service.ProjectExistsByID(int64(projectID))
	if err != nil {
		p.internalError(err)
		return
	}
	if !isExists {
		p.customAbort(http.StatusNotFound, fmt.Sprintf("Cannot find project by ID: %d", projectID))
		return
	}

	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}

	var reqProject model.Project
	err = json.Unmarshal(reqData, &reqProject)
	if err != nil {
		p.internalError(err)
		return
	}
	reqProject.ID = int64(projectID)
	isSuccess, err := service.UpdateProject(reqProject, "public")
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, "Failed to update project public.")
	}
}

// TODO
func init() {
	logs.Info("Init git repo for default project %s", defaultProject)
	_, err := service.InitRepo(repoServeURL(), repoPath())
	if err != nil {
		logs.Error("Failed to initialize default user's repo: %+v\n", err)
		return
	}

	subPath := defaultProject
	if subPath != "" {
		os.MkdirAll(filepath.Join(repoPath(), subPath), 0755)
		if err != nil {
			logs.Error("Failed to make default user's repo: %+v\n", err)
		}
	}
}
