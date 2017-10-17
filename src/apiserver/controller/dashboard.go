package controller

import (
	"encoding/json"

	"io/ioutil"

	"git/inspursoft/board/src/apiserver/service"
	"net/http"

	"fmt"

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
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
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
		s.CustomAbort(http.StatusBadRequest, "should provide time count")
		return
	}
	if req.Node.TimestampBase == 0 && req.Service.TimestampBase == 0 {
		s.CustomAbort(http.StatusBadRequest, "should provide timestamp")

		return
	}
	if req.Node.TimeUnit == "" && req.Service.TimeUnit == "" {
		s.CustomAbort(http.StatusBadRequest, "should provide time unit")
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
		s.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	_, err = para.GetNodeListToObj()
	if err != nil {
		s.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
		beego.Error(err)
		return
	}
	resp.Node = para.NodeResp
	para.SetServicePara(req.Service.TimeUnit,
		req.Service.TimeCount, req.Service.TimestampBase, serviceName,
		req.Service.DurationTime)
	err = para.GetServiceDataToObj()
	if err != nil {
		beego.Error(err)
	}
	_, err = para.GetServiceListToObj()
	if err != nil {
		s.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	resp.Service = para.ServiceResp
	if err != nil {
		beego.Error(err)
	}
	s.Data["json"] = resp
	s.ServeJSON()

}
