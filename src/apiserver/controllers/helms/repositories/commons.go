package repositories

import (
	"github.com/astaxie/beego/logs"

	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
)

// Operations about Helm repositories
type CommonController struct {
	c.BaseController
}

// @Title List all Helm repositories
// @Description List all for Helm repositories
// @Param	search	query	string	false	"Query item for Helm repository"
// @Success 200 {array} []vm.HelmRepository Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (cc *CommonController) List() {
	logs.Info("list all helm repos")

	// list the repos from storage
	repos, err := service.ListVMHelmRepositories()
	if err != nil {
		cc.InternalError(err)
		return
	}

	cc.RenderJSON(repos)
}

// @Title Get Helm repository detail
// @Description Get Helm repository detail
// @Param	repository_id	path	int	true	"ID of Helm repository"
// @Param	chart_name	path	string	true	"Name of Helm Chart"
// @Param	chart_version	path	string	true	"Version of Helm Chart"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:repository_id/charts/:chart_name/:chart_version [get]
func (c *CommonController) Get() {

}

// @Title Upload chart to Helm repository
// @Description Upload chart to Helm repositoryh
// @Param	repository_id	path	int	true	"ID of Helm repository"
// @Param	body	body	models.helms.vm.Chart	"View model of Helm Chart."
// @Success 200 Successful uploaded.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:repository_id/charts/uploads [post]
func (c *CommonController) Post() {

}

// @Title Delete Helm repository
// @Description Delete Helm repository
// @Param	repository_id	path	int	true	"ID of Helm repository"
// @Param	chart_name	path	string	true	"Name of Helm Chart"
// @Param	chart_version	path	string	true	"Version of Helm Chart"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:repository_id/charts/:chart_name/:chart_version [delete]
func (c *CommonController) Delete() {

}
