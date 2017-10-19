package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
)

type DashboardNodeController struct {
	baseController
}

type NodeBodyReqPara struct {
	TimeUnit      string `json:"node_time_unit"`
	TimeCount     int    `json:"node_time_count"`
	TimestampBase int    `json:"node_timestamp"`
	DurationTime  int    `json:"node_duration_time"`
}

func (p *DashboardNodeController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
}

func (b *DashboardNodeController) resolveBody() (in NodeBodyReqPara, err error) {
	data, err := ioutil.ReadAll(b.Ctx.Request.Body)
	json.Unmarshal(data, &in)
	if err != nil {
		return in, err
	}
	return in, nil
}

func (s *DashboardNodeController) GetNodeData() {
	getNodeDataBodyReq, _ := s.resolveBody()
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
	beego.Error(err)
	s.Data["json"] = dashboardNodeDataResp.NodeResp
	s.ServeJSON()

}
