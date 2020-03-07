package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"net/http"

	c "git/inspursoft/board/src/common/controller"

	"github.com/astaxie/beego"
)

type DashboardNodeController struct {
	c.BaseController
}

type NodeBodyReqPara struct {
	TimeUnit      string `json:"node_time_unit"`
	TimeCount     int    `json:"node_time_count"`
	TimestampBase int    `json:"node_timestamp"`
	DurationTime  int    `json:"node_duration_time"`
}

func (s *DashboardNodeController) GetNodeData() {
	var getNodeDataBodyReq NodeBodyReqPara
	err := s.ResolveBody(&getNodeDataBodyReq)
	if err != nil {
		return
	}
	nodeName := s.GetString("node_name")
	beego.Debug("node_name", nodeName)
	if getNodeDataBodyReq.TimeCount == 0 {
		s.CustomAbortAudit(http.StatusBadRequest, "Time count for node data retrieval cannnot be empty.")
		return
	}
	if getNodeDataBodyReq.TimestampBase == 0 {
		s.CustomAbortAudit(http.StatusBadRequest, "Time stamp for node data retrieval cannot be empty.")
		return
	}
	if getNodeDataBodyReq.TimeUnit == "" {
		s.CustomAbortAudit(http.StatusBadRequest, "Time unit for node data retrieval cannot be empty.")
		return
	}
	var dashboardNodeDataResp service.Dashboard
	dashboardNodeDataResp.SetNodeParaFromBodyReq(getNodeDataBodyReq.TimeUnit, getNodeDataBodyReq.TimeCount,
		getNodeDataBodyReq.TimestampBase, nodeName, getNodeDataBodyReq.DurationTime)
	beego.Debug(getNodeDataBodyReq.TimeUnit, getNodeDataBodyReq.TimeCount,
		getNodeDataBodyReq.TimestampBase, nodeName)
	err = dashboardNodeDataResp.GetNodeDataToObj()
	if err != nil {
		s.InternalError(err)
		return
	}
	_, err = dashboardNodeDataResp.GetNodeListToObj()
	if err != nil {
		s.InternalError(err)
		return
	}
	s.RenderJSON(dashboardNodeDataResp.NodeResp)
}
