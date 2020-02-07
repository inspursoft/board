package jobs

import (
	"github.com/astaxie/beego"
)

// Operations about Jobs supplements
type SupplementController struct {
	beego.Controller
}

// @Title Get specified Job status
// @Description Get specified Job status
// @Param	job_id	path	int	true	"ID of Jobs"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:job_id/status [get]
func (c *SupplementController) Status() {

}

// @Title Get specified Job pods
// @Description Get specified Job status
// @Param	job_id	path	int	true	"ID of Jobs"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:job_id/pods [get]
func (c *SupplementController) Pods() {

}

// @Title Get specified Job pods console logs
// @Description Get specified Job pods console logs
// @Param	job_id	path	int	true	"ID of Jobs"
// @Param	pod_name	path	string	true	"Name of Pods to Job."
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:job_id/logs/:pod_name [get]
func (c *SupplementController) PodLogs() {

}

// @Title Check Job existing
// @Description Check Job existing
// @Param	job_name	query	string	true	"Name of Job."
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /existing [get]
func (c *SupplementController) Existing() {

}

// @Title Check selectable Job
// @Description Check selectable Job
// @Param	job_name	query	string	true	"Name of Job."
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /selectable [get]
func (c *SupplementController) Selectable() {

}
