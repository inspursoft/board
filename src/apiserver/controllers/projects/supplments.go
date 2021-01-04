package projects

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/service/adapting"
	"net/http"
	"strconv"
)

//Operation about Project's supplement actions
type SupplementController struct {
	c.BaseController
}

// @Title Supplement for toggling project publicity option.
// @Description Supplement for toggling project publicity option.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_id	path	int	true	"Request for project ID."
// @Success 200 Successful toggled.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/toggle [head]
func (p *SupplementController) Toggle() {
	projectID, err := strconv.Atoi(p.Ctx.Input.Param(":project_id"))
	if err != nil {
		p.InternalError(err)
		return
	}

	p.ResolveProjectOwnerByID(int64(projectID))

	var reqProject vm.Project
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

// @Title Supplement for checking project existing by name.
// @Description Supplement for checking project existing by name.
// @Param	token	query	string	true	"Token acquired when signed in."
// @Param	project_name	query	string	true	"Request for project name."
// @Success 200 Successful toggled.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /existing [head]
func (p *SupplementController) Existing() {
	projectName := p.GetString("project_name")
	query := vm.Project{Name: projectName}
	project, err := adapting.GetProject(query, "name")
	if err != nil {
		p.InternalError(err)
		return
	}
	if project != nil {
		p.CustomAbortAudit(http.StatusConflict, fmt.Sprintf("Project name: %s already exists.", projectName))
	}
}
