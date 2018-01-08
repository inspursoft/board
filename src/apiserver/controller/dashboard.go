package controller

import (
	"encoding/json"

	"io/ioutil"

	"git/inspursoft/board/src/apiserver/service"
	"net/http"

	"github.com/astaxie/beego"
)

type ServicePara struct {
	TimeUnit      string `json:"service_time_unit"`
	TimeCount     int    `json:"service_time_count"`
	TimestampBase int    `json:"service_timestamp"`
	DurationTime  int    `json:"service_duration_time"`
}
type NodePara struct {
	TimeUnit      string `json:"node_time_unit"`
	TimeCount     int    `json:"node_time_count"`
	TimestampBase int    `json:"node_timestamp"`
	DurationTime  int    `json:"node_duration_time"`
}
type DsBodyPara struct {
	Service ServicePara `json:"service"`
	Node    NodePara    `json:"node"`
}
type DsResp struct {
	Node    service.NodeResp    `json:"node"`
	Service service.ServiceResp `json:"service"`
}

func (p *Dashboard) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
}
func (b *Dashboard) resolveBody() (in DsBodyPara, err error) {
	data, err := ioutil.ReadAll(b.Ctx.Request.Body)
	json.Unmarshal(data, &in)
	if err != nil {
		return in, err
	}
	return in, nil
}

type Dashboard struct {
	baseController
}

func (s *Dashboard) GetData() {
	req, _ := s.resolveBody()
	nodeName := s.GetString("node_name")
	serviceName := s.GetString("service_name")
	beego.Debug("node_name", nodeName)
	if req.Node.TimeCount == 0 && req.Service.TimeCount == 0 {
		s.customAbort(http.StatusBadRequest, "Time count for dashboard data retrieval cannot be empty.")
		return
	}
	if req.Node.TimestampBase == 0 && req.Service.TimestampBase == 0 {
		s.customAbort(http.StatusBadRequest, "Timestamp for dashboard data retrieval cannot be empty.")
		return
	}
	if req.Node.TimeUnit == "" && req.Service.TimeUnit == "" {
		s.customAbort(http.StatusBadRequest, "Time unit for dashboard data retrieval cannot be empty.")
		return
	}
	var (
		para service.Dashboard
		resp DsResp
	)
	para.SetNodeParaFromBodyReq(req.Node.TimeUnit, req.Node.TimeCount,
		req.Node.TimestampBase, nodeName, req.Node.DurationTime)
	beego.Debug(req.Node.TimeUnit, req.Node.TimeCount,
		req.Node.TimestampBase, nodeName)
	err := para.GetNodeDataToObj()
	if err != nil {
		s.internalError(err)
		return
	}
	_, err = para.GetNodeListToObj()
	if err != nil {
		s.internalError(err)
		return
	}
	resp.Node = para.NodeResp
	para.SetServicePara(req.Service.TimeUnit,
		req.Service.TimeCount, req.Service.TimestampBase, serviceName,
		req.Service.DurationTime)
	err = para.GetServiceDataToObj()
	if err != nil {
		s.internalError(err)
		return
	}
	_, err = para.GetServiceListToObj()
	if err != nil {
		s.internalError(err)
		return
	}
	resp.Service = para.ServiceResp
	if err != nil {
		s.internalError(err)
		return
	}
	s.Data["json"] = resp
	s.ServeJSON()

}
