package controller

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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

// BoardInfo from adminserver
type BoardModuleInfo struct {
	ID        string `json:"id"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Ports     string `json:"ports"`
	Name      string `json:"name"`
	CPUPerc   string `json:"cpu_perc"`
	MemUsage  string `json:"mem_usage"`
	NetIO     string `json:"net_io"`
	BlockIO   string `json:"block_io"`
	MemPerc   string `json:"mem_perc"`
	PIDs      string `json:"pids"`
}

type InitStatus int

const (
	InitStatusFirst  InitStatus = 1
	InitStatusSecond InitStatus = 2
	InitStatusThird  InitStatus = 3
)

type InitSysStatus struct {
	Status InitStatus `json:"status"`
}

var AdminServerURL = utils.GetConfig("ADMINSERVER_URL")

var ErrInvalidToken = errors.New("error for invalid token")
var ErrServerAccessFailed = errors.New("error for access server failed")

type Dashboard struct {
	c.BaseController
}

func (s *Dashboard) GetData() {
	var req DsBodyPara
	err := s.ResolveBody(&req)
	if err != nil {
		return
	}
	nodeName := s.GetString("node_name")
	serviceName := s.GetString("service_name")

	if req.Node.TimeCount == 0 && req.Service.TimeCount == 0 {
		s.CustomAbortAudit(http.StatusBadRequest, "Time count for dashboard data retrieval cannot be empty.")
		return
	}
	if req.Node.TimestampBase == 0 && req.Service.TimestampBase == 0 {
		s.CustomAbortAudit(http.StatusBadRequest, "Timestamp for dashboard data retrieval cannot be empty.")
		return
	}
	if req.Node.TimeUnit == "" && req.Service.TimeUnit == "" {
		s.CustomAbortAudit(http.StatusBadRequest, "Time unit for dashboard data retrieval cannot be empty.")
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
	err = para.GetNodeDataToObj()
	if err != nil {
		s.InternalError(err)
		return
	}
	_, err = para.GetNodeListToObj()
	if err != nil {
		s.InternalError(err)
		return
	}
	resp.Node = para.NodeResp
	para.SetServicePara(req.Service.TimeUnit,
		req.Service.TimeCount, req.Service.TimestampBase, serviceName,
		req.Service.DurationTime)
	err = para.GetServiceDataToObj()
	if err != nil {
		s.InternalError(err)
		return
	}
	_, err = para.GetServiceListToObj()
	if err != nil {
		s.InternalError(err)
		return
	}
	resp.Service = para.ServiceResp
	if err != nil {
		s.InternalError(err)
		return
	}
	s.RenderJSON(resp)
}

//Check the adminserver monitor modules
func (s *Dashboard) AdminserverCheck() {
	if s.IsSysAdmin == false {
		s.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to control node.")
		return
	}
	moduleName := s.GetString("module_name")

	token := s.Ctx.Request.Header.Get("token")
	if token == "" {
		token = s.GetString("token")
	}
	// Just use the token, skip verify and refresh it

	adminServerURL := AdminServerURL()

	logs.Debug("%s/monitor?token=%s", adminServerURL, token)

	var boardinfo []BoardModuleInfo
	err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/monitor?token=%s", adminServerURL, token), nil, nil, func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode >= http.StatusInternalServerError {
			logs.Error("Access adminserver failed %s.", req.URL)
			return ErrServerAccessFailed
		}

		if resp.StatusCode == http.StatusUnauthorized {
			logs.Error("Invalid token due to session timeout.")
			return ErrInvalidToken
		}
		return utils.UnmarshalToJSON(resp.Body, &boardinfo)
	})

	//TODO filter the module name
	logs.Debug("Check the module %s", moduleName)

	if err != nil {
		if err.Error() == ErrServerAccessFailed.Error() {
			logs.Debug("Adminserver internal failed %v", err)
			s.CustomAbortAudit(http.StatusNotFound, "Cannot access adminserver.")
			return
		}
		if err.Error() == ErrInvalidToken.Error() {
			logs.Debug("Token failed %v", err)
			s.CustomAbortAudit(http.StatusUnauthorized, "Invalid token to access adminserver.")
			return
		}
		logs.Error("Access adminserver err %v", err)
		s.CustomAbortAudit(http.StatusBadRequest, "Access adminserver failed.")
		return
	}
	s.RenderJSON(boardinfo)
}

//Check the sys status by adminserver
func (s *Dashboard) CheckSysByAdminserver() {
	if s.IsSysAdmin == false {
		s.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to control node.")
		return
	}
	adminServerURL := AdminServerURL()

	logs.Debug("%s/boot/checksysstatus", adminServerURL)

	var sysstatus InitSysStatus
	err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/boot/checksysstatus", adminServerURL), nil, nil, func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode >= http.StatusInternalServerError {
			logs.Error("Access adminserver failed %s.", req.URL)
			return ErrServerAccessFailed
		}

		if resp.StatusCode == http.StatusUnauthorized {
			logs.Error("Invalid token due to session timeout.")
			return ErrInvalidToken
		}
		return utils.UnmarshalToJSON(resp.Body, &sysstatus)
	})

	if err != nil {
		if err.Error() == ErrServerAccessFailed.Error() {
			logs.Debug("Adminserver internal failed %v", err)
			s.CustomAbortAudit(http.StatusNotFound, "Cannot access adminserver.")
			return
		}
		if err.Error() == ErrInvalidToken.Error() {
			logs.Debug("Token failed %v", err)
			s.CustomAbortAudit(http.StatusUnauthorized, "Invalid token to access adminserver.")
			return
		}
		logs.Error("Access adminserver err %v", err)
		s.CustomAbortAudit(http.StatusBadRequest, "Access adminserver failed.")
		return
	}
	s.RenderJSON(sysstatus)
}
