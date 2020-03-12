package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	c "git/inspursoft/board/src/common/controller"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"net/http"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ServiceBodyPara struct {
	TimeUnit      string `json:"service_time_unit"`
	TimeCount     int    `json:"service_time_count"`
	TimestampBase int    `json:"service_timestamp"`
	DurationTime  int    `json:"service_duration_time"`
}

type DashboardServiceController struct {
	c.BaseController
}

func (s *DashboardServiceController) GetServiceData() {

	var getServiceDataBodyReq ServiceBodyPara
	err := s.ResolveBody(&getServiceDataBodyReq)
	if err != nil {
		return
	}
	serviceName := s.GetString("service_name")

	beego.Debug("servicename", serviceName, getServiceDataBodyReq.DurationTime)
	if getServiceDataBodyReq.TimeCount == 0 {
		s.CustomAbortAudit(http.StatusBadRequest, "")
		return
	}
	if getServiceDataBodyReq.TimestampBase == 0 {
		s.CustomAbortAudit(http.StatusBadRequest, "")
		return
	}
	if getServiceDataBodyReq.TimeUnit == "" {
		s.CustomAbortAudit(http.StatusBadRequest, "")
		return
	}

	var dashboardServiceDataResp service.Dashboard
	dashboardServiceDataResp.SetServicePara(getServiceDataBodyReq.TimeUnit,
		getServiceDataBodyReq.TimeCount, getServiceDataBodyReq.TimestampBase, serviceName,
		getServiceDataBodyReq.DurationTime)
	err = dashboardServiceDataResp.GetServiceDataToObj()
	_, err = dashboardServiceDataResp.GetServiceListToObj()
	if err != nil {
		s.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	query := model.Project{}
	projectList, err := service.GetProjectsByMember(query, s.CurrentUser.ID)
	if err != nil {
		s.InternalError(err)
		return
	}
	serviceList := make([]dao.ServiceListDataLogs, 0)
	for _, svc := range dashboardServiceDataResp.ServiceResp.ServiceListData {
		svcQuery, err := service.GetService(model.ServiceStatus{Name: svc.NodeName}, "name")
		if err != nil {
			s.InternalError(err)
			return
		}
		if svcQuery == nil {
			continue
		}
		if svcQuery.Public == 1 {
			serviceList = append(serviceList, svc)
			continue
		}
		for _, project := range projectList {
			if svcQuery.ProjectName == project.Name {
				serviceList = append(serviceList, svc)
				break
			}
		}
	}
	dashboardServiceDataResp.ServiceResp.ServiceListData = serviceList
	logs.Info("serivcelist:%+v\n", dashboardServiceDataResp.ServiceResp.ServiceListData)

	s.RenderJSON(dashboardServiceDataResp.ServiceResp)
}

func (s *DashboardServiceController) GetServerTime() {
	time := service.GetServerTime()
	s.RenderJSON(time)
}
