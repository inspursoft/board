package controller

import (
	"git/inspursoft/board/src/apiserver/service"
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
	baseController
}

func (s *DashboardServiceController) GetServiceData() {

	var getServiceDataBodyReq ServiceBodyPara
	s.resolveBody(&getServiceDataBodyReq)
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
	err := dashboardServiceDataResp.GetServiceDataToObj()
	_, err = dashboardServiceDataResp.GetServiceListToObj()
	if err != nil {
		s.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	query := model.Project{}
	projectList, err := service.GetProjectsByMember(query, s.currentUser.ID)
	if err != nil {
		s.internalError(err)
		return
	}
	serviceList := make([]dao.ServiceListDataLogs, 0)
	for _, svc := range dashboardServiceDataResp.ServiceResp.ServiceListData {
		svcQuery, err := service.GetService(model.ServiceStatus{Name: svc.NodeName}, "name")
		if err != nil {
			s.internalError(err)
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

	s.Data["json"] = dashboardServiceDataResp.ServiceResp
	s.ServeJSON()
}

func (s *DashboardServiceController) GetServerTime() {
	time := service.GetServerTime()
	s.Data["json"] = time
	s.ServeJSON()

}
