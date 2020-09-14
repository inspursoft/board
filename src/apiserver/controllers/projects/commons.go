package projects

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/models/vm"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/adapting"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

// Operations about projects
type CommonController struct {
	c.BaseController
}

// @Title Get project by ID
// @Description Get projects by ID.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	true	"Query project ID for projects"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [get]
func (p *CommonController) Get() {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":project_id"))
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

// @Title List all projects
// @Description List all for projects.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_name	query	string	false	"Query name of projects"
// @Param   project_public	query	int	false	"Query public of projects"
// @Param	member_only	query	int	false	"Query member only of projects"
// @Param	page_index query	int	false	"Page index for pagination."
// @Param	page_size	query	int	false	"Page per size for pagination."
// @Param	order_field	string	false	"Order by field. (Default is creation_time)"
// @Param	order_asc	int	false	"Order option for ascend or descend. (asc 0, desc, 1)"
// @Param	search	query	string	false	"Query item for projects"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (p *CommonController) List() {
	projectName := p.GetString("project_name")
	strPublic := p.GetString("project_public")
	memberOnly, _ := p.GetInt("member_only", 0)

	pageIndex, _ := p.GetInt("page_index", 0)
	pageSize, _ := p.GetInt("page_size", 0)
	orderField := p.GetString("order_field", "creation_time")
	orderAsc, _ := p.GetInt("order_asc", 0)

	query := vm.Project{Name: projectName, OwnerName: p.CurrentUser.Username, Public: 0}

	public, err := strconv.Atoi(strPublic)
	if err == nil {
		query.Public = public
	}

	if pageIndex == 0 && pageSize == 0 {
		var projects []*vm.Project
		var err error
		if memberOnly == 1 {
			projects, err = adapting.GetProjectsByMember(query, p.CurrentUser.ID)
		} else {
			projects, err = adapting.GetProjectsByUser(query, p.CurrentUser.ID)
		}
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(projects)
	} else {
		paginatedProjects, err := adapting.GetPaginatedProjectsByUser(query, p.CurrentUser.ID, pageIndex, pageSize, orderField, orderAsc)
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(paginatedProjects)
	}
}

// @Title Add project
// @Description Add project.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	body	body 	vm.Project	true	"View model for projects."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (p *CommonController) Add() {
	var reqProject vm.Project
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

	isSuccess, err := adapting.CreateProject(reqProject)
	if err != nil {
		p.InternalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Project name: %s is illegal.", reqProject.Name))
		return
	}

	isSuccess, err = adapting.CreateNamespace(&reqProject)
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

// @Title Update project by ID
// @Description Update project by ID.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	true	"ID of projects"
// @Param	body	body	models.projects.vm.Project	true	"View model for projects."
// @Param	action	query	string	true	"Option of update."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [put]
func (c *CommonController) Update() {

}

// @Title Delete project by ID
// @Description Delete project by ID.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	true	"ID of projects"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id [delete]
func (p *CommonController) Delete() {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":project_id"))
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
