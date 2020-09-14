package controller

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

type ProjectController struct {
	c.BaseController
}

func (p *ProjectController) CreateProjectAction() {
	var reqProject model.Project
	var err error
	err = p.ResolveBody(&reqProject)
	if err != nil {
		return
	}

	if !utils.ValidateWithLengthRange(reqProject.Name, 2, 63) {
		p.CustomAbortAudit(http.StatusBadRequest, "Project name length should be between 2 and 63 characters.")
		return
	}
	if !utils.ValidateWithPattern("project", reqProject.Name) {
		p.CustomAbortAudit(http.StatusBadRequest, "Project name is invalid.")
		return
	}

	projectExists, err := service.ProjectExists(reqProject.Name)
	if err != nil {
		p.InternalError(err)
		return
	}
	if projectExists {
		p.CustomAbortAudit(http.StatusConflict, "Project name already exists.")
		return
	}

	// Check namespace in k8s cluster
	projectExists, err = service.NamespaceExists(reqProject.Name)
	if err != nil {
		p.InternalError(err)
		return
	}
	if projectExists {
		p.CustomAbortAudit(http.StatusConflict, fmt.Sprintf("Namespace %s already exists in cluster.", reqProject.Name))
		return
	}

	reqProject.Name = strings.TrimSpace(reqProject.Name)
	reqProject.OwnerID = int(p.CurrentUser.ID)
	reqProject.OwnerName = p.CurrentUser.Username
	reqProject.CreationTime = time.Now()
	reqProject.UpdateTime = reqProject.CreationTime

	isSuccess, err := service.CreateProject(reqProject)
	if err != nil {
		p.InternalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Project name: %s is illegal.", reqProject.Name))
		return
	}

	isSuccess, err = service.CreateNamespace(&reqProject)
	if err != nil {
		p.InternalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Namespace name: %s is illegal.", reqProject.Name))
	}

	err = service.CurrentDevOps().CreateRepoAndJob(p.CurrentUser.ID, reqProject.Name)
	if err != nil {
		logs.Error("Failed to create repo and job for project: %s", reqProject.Name)
		p.InternalError(err)
	}
}

func (p *ProjectController) ProjectExists() {
	projectName := p.GetString("project_name")
	query := model.Project{Name: projectName}
	project, err := service.GetProject(query, "name")
	if err != nil {
		p.InternalError(err)
		return
	}
	if project != nil {
		p.CustomAbortAudit(http.StatusConflict, fmt.Sprintf("Project name: %s already exists.", projectName))
	}
}

func (p *ProjectController) GetProjectsAction() {
	projectName := p.GetString("project_name")
	strPublic := p.GetString("project_public")
	memberOnly, _ := p.GetInt("member_only", 0)

	pageIndex, _ := p.GetInt("page_index", 0)
	pageSize, _ := p.GetInt("page_size", 0)
	orderField := p.GetString("order_field", "creation_time")
	orderAsc, _ := p.GetInt("order_asc", 0)

	orderFieldValue, err := service.ParseOrderField("project", orderField)
	if err != nil {
		p.CustomAbortAudit(http.StatusBadRequest, err.Error())
		return
	}

	query := model.Project{Name: projectName, OwnerName: p.CurrentUser.Username, Public: 0}

	public, err := strconv.Atoi(strPublic)
	if err == nil {
		query.Public = public
	}

	if pageIndex == 0 && pageSize == 0 {
		var projects []*model.Project
		var err error
		if memberOnly == 1 {
			projects, err = service.GetProjectsByMember(query, p.CurrentUser.ID)
		} else {
			projects, err = service.GetProjectsByUser(query, p.CurrentUser.ID)
		}
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(projects)
	} else {
		paginatedProjects, err := service.GetPaginatedProjectsByUser(query, p.CurrentUser.ID, pageIndex, pageSize, orderFieldValue, orderAsc)
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(paginatedProjects)
	}
}

func (p *ProjectController) GetProjectAction() {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	project, err := service.GetProjectByID(int64(projectID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if project == nil {
		p.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("No project was found with provided ID: %d", projectID))
		return
	}
	p.RenderJSON(project)
}

func (p *ProjectController) DeleteProjectAction() {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	project := p.ResolveProjectByID(int64(projectID))
	if !(p.IsSysAdmin || int64(project.OwnerID) == p.CurrentUser.ID) {
		p.CustomAbortAudit(http.StatusForbidden, "User is not the owner of the project.")
		return
	}
	user, err := service.GetUserByName(project.OwnerName)
	if err != nil {
		p.InternalError(err)
		return
	}
	isSuccess, err := service.DeleteProject(user.ID, int64(projectID))
	if err != nil {
		if err == utils.ErrUnprocessableEntity {
			p.CustomAbortAudit(http.StatusUnprocessableEntity, fmt.Sprintf("Project %s has own member, repo or service.", project.Name))
		} else {
			p.InternalError(err)
		}
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, "Failed to delete project.")
		return
	}
}

func (p *ProjectController) ToggleProjectPublicAction() {

	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.InternalError(err)
		return
	}

	p.ResolveProjectOwnerByID(int64(projectID))

	var reqProject model.Project
	err = p.ResolveBody(&reqProject)
	if err != nil {
		return
	}

	isSuccess, err := service.ToggleProjectPublic(int64(projectID), reqProject.Public)
	if err != nil {
		p.InternalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, "Failed to update project public.")
	}
}
