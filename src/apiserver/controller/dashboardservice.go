package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"git/inspursoft/board/src/apiserver/service"
	"net/http"

	"github.com/astaxie/beego"
)

type ServiceBodyPara struct {
	TimeUnit      string `json:"service_time_unit"`
	TimeCount     int    `json:"service_time_count"`
	TimestampBase int    `json:"service_timestamp"`
	DurationTime  int    `json:"service_duration_time"`
}

func (p *DashboardServiceController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}
func (b *DashboardServiceController) resolveBody() (in ServiceBodyPara, err error) {
	data, err := ioutil.ReadAll(b.Ctx.Request.Body)
	json.Unmarshal(data, &in)
	if err != nil {
		return in, err
	}
	return in, nil
}

type DashboardServiceController struct {
	baseController
}

func (s *DashboardServiceController) GetServiceData() {
	getServiceDataBodyReq, _ := s.resolveBody()
	serviceName := s.GetString("service_name")
	beego.Debug("servicename", serviceName, getServiceDataBodyReq.DurationTime)
	if getServiceDataBodyReq.TimeCount == 0 {
		s.CustomAbort(http.StatusBadRequest, "")
		return
	}
	if getServiceDataBodyReq.TimestampBase == 0 {
		s.CustomAbort(http.StatusBadRequest, "")
		return
	}
	if getServiceDataBodyReq.TimeUnit == "" {
		s.CustomAbort(http.StatusBadRequest, "")
		return
	}
	var dashboardServiceDataResp service.Dashboard
	dashboardServiceDataResp.SetServicePara(getServiceDataBodyReq.TimeUnit,
		getServiceDataBodyReq.TimeCount, getServiceDataBodyReq.TimestampBase, serviceName,
		getServiceDataBodyReq.DurationTime)
	fmt.Println(getServiceDataBodyReq.TimeUnit, getServiceDataBodyReq.TimeCount,
		getServiceDataBodyReq.TimestampBase, serviceName)
	err := dashboardServiceDataResp.GetServiceDataToObj()
	beego.Error(err)
	_, err = dashboardServiceDataResp.GetServiceListToObj()
	if err != nil {
		s.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	beego.Error(err)
	s.Data["json"] = dashboardServiceDataResp.ServiceResp
	s.ServeJSON()
}
