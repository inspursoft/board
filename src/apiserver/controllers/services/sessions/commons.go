package sessions

import "github.com/astaxie/beego"

// Operations about service sessions
type CommonController struct {
	beego.Controller
}

// @Title Get service session
// @Description Get service session.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of service"
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [get]
func (c *CommonController) List() {

}

// @Title Update service session by ID
// @Description Update service session by ID.
// @Param	project_id	path	int	true	"ID of project"
// @Param	service_id	path	int	true	"ID of services"
// @Param	body	body 	models.services.vm.Session	true	"View model for service session."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/:service_id [patch]
func (c *CommonController) Update() {

}
