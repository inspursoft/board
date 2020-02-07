package dashboards

import "github.com/astaxie/beego"

// Operation about data of dashboard.
type DataController struct {
	beego.Controller
}

// @Title Fetch service data for dashboard
// @Description Fetch service data for dashboard
// @Param	search	query	string	false	"Query item for service data"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /services [post]
func (d *DataController) Service() {

}

// @Title Fetch data for dashboard
// @Description Fetch data for dashboard
// @Param	search	query	string	false	"Query item for data"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /data [post]
func (d *DataController) Data() {

}

// @Title Fetch node info for dashboard
// @Description Fetch node info for dashboard
// @Param	search	query	string	false	"Query item for data"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /nodes [post]
func (d *DataController) Node() {

}

// @Title Fetch server time for dashboard
// @Description Fetch server time for dashboard
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /server_time [get]
func (d *DataController) ServerTime() {

}
