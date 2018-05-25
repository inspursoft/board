package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"net/http"

	"github.com/astaxie/beego"
)

type DashboardNodeController struct {
	BaseController
}

type NodeBodyReqPara struct {
	TimeUnit      string `json:"node_time_unit"`
	TimeCount     int    `json:"node_time_count"`
	TimestampBase int    `json:"node_timestamp"`
	DurationTime  int    `json:"node_duration_time"`
}

func (s *DashboardNodeController) GetNodeData() {
	var getNodeDataBodyReq NodeBodyReqPara
	s.resolveBody(&getNodeDataBodyReq)
	nodeName := s.GetString("node_name")
	beego.Debug("node_name", nodeName)
	if getNodeDataBodyReq.TimeCount == 0 {
		s.customAbort(http.StatusBadRequest, "Time count for node data retrieval cannnot be empty.")
		return
	}
	if getNodeDataBodyReq.TimestampBase == 0 {
		s.customAbort(http.StatusBadRequest, "Time stamp for node data retrieval cannot be empty.")
		return
	}
	if getNodeDataBodyReq.TimeUnit == "" {
		s.customAbort(http.StatusBadRequest, "Time unit for node data retrieval cannot be empty.")
		return
	}
	var dashboardNodeDataResp service.Dashboard
	dashboardNodeDataResp.SetNodeParaFromBodyReq(getNodeDataBodyReq.TimeUnit, getNodeDataBodyReq.TimeCount,
		getNodeDataBodyReq.TimestampBase, nodeName, getNodeDataBodyReq.DurationTime)
	beego.Debug(getNodeDataBodyReq.TimeUnit, getNodeDataBodyReq.TimeCount,
		getNodeDataBodyReq.TimestampBase, nodeName)
	err := dashboardNodeDataResp.GetNodeDataToObj()
	if err != nil {
		s.internalError(err)
		return
	}
	_, err = dashboardNodeDataResp.GetNodeListToObj()
	if err != nil {
		s.internalError(err)
		return
	}
	s.renderJSON(dashboardNodeDataResp.NodeResp)
}
