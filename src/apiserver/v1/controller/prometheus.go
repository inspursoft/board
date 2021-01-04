package controller

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type PrometheusController struct {
	c.BaseController
}

// @Title GetData
// @Description Get DashBoard Data
// @Param	body	body 	service.RequestPayload	true	"request payload"
// @Param	node	query 	string	true	"node"
// @Param	service	query 	string	true	"service"
// @Success 200 {object} service.DashboardInfo	success
// @Failure 400 Bad Request
// @Failure 500 Internal Server Error
// @router /prometheus [post]
func (p *PrometheusController) GetData() {
	var request service.RequestPayload
	var err error

	nodename := p.GetString("node")
	servicename := p.GetString("service")
	err = utils.UnmarshalToJSON(p.Ctx.Request.Body, &request)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		p.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	if request.TimeCount == 0 {
		p.CustomAbort(http.StatusBadRequest, "Time count for dashboard data retrieval cannot be empty.")
		return
	}
	if request.TimeStamp == 0 {
		p.CustomAbort(http.StatusBadRequest, "Timestamp for dashboard data retrieval cannot be empty.")
		return
	}
	if request.TimeUnit == "" {
		p.CustomAbort(http.StatusBadRequest, "Time unit for dashboard data retrieval cannot be empty.")
		return
	}

	data, err := service.GetDashBoardData(request, nodename, servicename)
	if err != nil {
		logs.Error(err)
		p.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	p.Data["json"] = data
	p.ServeJSON()
}
