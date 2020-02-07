package jobs

import (
	"github.com/astaxie/beego"
)

// Operations about Jobs
type CommonController struct {
	beego.Controller
}

// @Title List all Jobs
// @Description List all Jobs
// @Param	search	query	string	false	"Query item for Jobs"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (c *CommonController) List() {

}

// @Title Deploy Job
// @Description Deploy Job
// @Param	body	body	models.jobs.vm.Job	"View model of Job."
// @Success 200 Successful deployed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /deploy [post]
func (c *CommonController) Post() {

}

// @Title Delete Job
// @Description Delete Job
// @Param	job_id	path	int	true	"ID of Job"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:job_id [delete]
func (c *CommonController) Delete() {

}
