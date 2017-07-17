package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
	"strings"
)

type ProjectController struct {
	baseController
}

func (p *ProjectController) Prepare() {
	user, err := p.getCurrentUser()
	if err != nil {
		p.internalError(err)
		return
	}
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}

func (p *ProjectController) CreateProjectAction() {
	if !p.isProjectAdmin {
		p.CustomAbort(http.StatusForbidden, "Insuffient privileges for creating projects.")
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
	if strings.TrimSpace(reqProject.Name) == "" {
		p.serveStatus(http.StatusBadRequest, "Project name cannot be empty.")
		return
	}
	projectExists, err := service.ProjectExists(reqProject.Name)
	if err != nil {
		p.internalError(err)
		return
	}
	if projectExists {
		p.serveStatus(http.StatusConflict, "Project name already exists.")
		return
	}

	reqProject.OwnerID = int(p.currentUser.ID)
	reqProject.OwnerName = p.currentUser.Username

	isSuccess, err := service.CreateProject(reqProject)
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.serveStatus(http.StatusBadRequest, "Project contains invalid characters.")
	}
}

func (p *ProjectController) GetProjectsAction() {
	projectName := p.GetString("project_name")
	strPublic := p.GetString("project_public")

	query := model.Project{Name: projectName, Public: 0}

	var err error
	public, err := strconv.Atoi(strPublic)
	if err == nil {
		query.Public = public
	}

	var projects []*model.Project
	if p.isSysAdmin {
		projects, err = service.GetAllProjects(query)
	} else {
		projects, err = service.GetProjectsByUser(query, p.currentUser.ID)
	}

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
		p.CustomAbort(http.StatusNotFound, "No project was found with provided ID.")
		return
	}
	p.Data["json"] = project
	p.ServeJSON()
}

func (p *ProjectController) DeleteProjectAction() {
	if !p.isProjectAdmin {
		p.CustomAbort(http.StatusForbidden, "Insuffient privileges for creating projects.")
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
		p.CustomAbort(http.StatusNotFound, "Cannot find project by ID")
		return
	}
	isSuccess, err := service.DeleteProject(int64(projectID))
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbort(http.StatusBadRequest, "Failed to delete project.")
	}
}

func (p *ProjectController) ToggleProjectPublicAction() {
	if !p.isProjectAdmin {
		p.CustomAbort(http.StatusForbidden, "Insuffient privileges for creating projects.")
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
		p.CustomAbort(http.StatusNotFound, "Cannot find project by ID")
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
		p.CustomAbort(http.StatusBadRequest, "Failed to update project public.")
	}
}
