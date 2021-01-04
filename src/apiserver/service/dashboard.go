package service

import (
	"github.com/inspursoft/board/src/common/dao"

	"github.com/astaxie/beego"
)

type Dashboard struct {
	NodeReqPara
	ServiceReqPara
	NodeResp
	NodeListResp
	ServiceResp
	ServiceListResp
}

type NodeListResp struct {
	NodeListData []dao.NodeListDataLogs `json:"node_list_data"`
}
type ServiceListResp struct {
	ServiceListData []dao.ServiceListDataLogs `json:"service_list_data"`
}

type NodeResp struct {
	NodeName       string             `json:"node_name"`
	IsOverMinLimit bool               `json:"is_over_min_limit"`
	IsOverMaxLimit bool               `json:"is_over_max_limit"`
	NodeTimestamp  int                `json:"node_timestamp"`
	NodeCount      int                `json:"node_count"`
	TimeUnit       string             `json:"time_unit"`
	NodeLogsData   []dao.NodeDataLogs `json:"node_logs_data"`
	NodeListResp
}
type ServiceResp struct {
	ServiceName     string               `json:"service_name"`
	IsOverMinLimit  bool                 `json:"is_over_min_limit"`
	IsOverMaxLimit  bool                 `json:"is_over_max_limit"`
	ServiceTimeUnit string               `json:"service_time_unit"`
	ServiceCount    int                  `json:"service_count"`
	ServiceLogsData []dao.ServiceDataLog `json:"service_logs_data"`
	ServiceListResp
}

type NodeReqPara struct {
	TimeUnit  string
	TimeCount int
	NodeName  string
	TimeStamp int
	DuraTime  int
}

type ServiceReqPara struct {
	TimeUnit    string
	TimeCount   int
	ServiceName string
	TimeStamp   int
	DuraTime    int
}

func (d *Dashboard) SetNodeParaFromBodyReq(timeUnit string, timeCount int, timestamp int,
	nodeName string, daraTime int) (err error) {
	d.NodeReqPara = NodeReqPara{
		TimeUnit:  timeUnit,
		TimeCount: timeCount,
		TimeStamp: timestamp,
		NodeName:  nodeName,
		DuraTime:  daraTime,
	}
	return nil
}

func (d *Dashboard) SetServicePara(timeUnit string, timeCount int,
	timestamp int, serviceName string, daraTime int) (err error) {
	d.ServiceReqPara = ServiceReqPara{
		TimeUnit:    timeUnit,
		TimeCount:   timeCount,
		TimeStamp:   timestamp,
		ServiceName: serviceName,
		DuraTime:    daraTime,
	}
	return nil
}
func (d *Dashboard) GetServiceDataToObj() (err error) {
	var tMin int
	switch d.ServiceReqPara.TimeUnit {
	case "second":
		tMin = d.ServiceReqPara.TimeStamp - d.ServiceReqPara.TimeCount*5
	case "minute":
		tMin = d.ServiceReqPara.TimeStamp - d.ServiceReqPara.TimeCount*60
	case "hour":
		tMin = d.ServiceReqPara.TimeStamp - d.ServiceReqPara.TimeCount*60*60
	case "day":
		tMin = d.ServiceReqPara.TimeStamp - d.ServiceReqPara.TimeCount*60*60*24

	}
	s := dao.DashboardServiceDao{}
	s.QueryPara = dao.QueryPara{
		Name:      d.ServiceReqPara.ServiceName,
		TimeStamp: d.ServiceReqPara.TimeStamp,
		TimeCount: d.ServiceReqPara.TimeCount,
		TimeUnit:  d.ServiceReqPara.TimeUnit,
		DuraTime:  d.ServiceReqPara.DuraTime,
	}
	beego.Debug(s.QueryPara)
	if d.ServiceReqPara.ServiceName == "" {
		d.ServiceResp = ServiceResp{
			ServiceName:     "total",
			ServiceTimeUnit: d.ServiceReqPara.TimeUnit,
		}
		d.ServiceResp.ServiceCount, d.ServiceLogsData, err = s.GetTotalServiceData()
		if err != nil {
			return err
		}
	} else {
		d.ServiceResp = ServiceResp{
			ServiceName:     d.ServiceReqPara.ServiceName,
			ServiceTimeUnit: d.ServiceReqPara.TimeUnit,
		}
		d.ServiceResp.ServiceCount, d.ServiceLogsData, err = s.GetServiceData()
	}
	lt, err := s.GetLimitTime()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	if tMin > lt.MaxTime {
		d.ServiceResp.IsOverMaxLimit = true
	}
	if d.ServiceReqPara.TimeStamp < lt.MinTime {
		d.ServiceResp.IsOverMinLimit = true
	}

	return nil
}
func (d *Dashboard) GetNodeDataToObj() (err error) {
	var tMin int
	switch d.NodeReqPara.TimeUnit {
	case "second":
		tMin = d.NodeReqPara.TimeStamp - d.ServiceReqPara.TimeCount*5
	case "minute":
		tMin = d.NodeReqPara.TimeStamp - d.ServiceReqPara.TimeCount*60
	case "hour":
		tMin = d.NodeReqPara.TimeStamp - d.ServiceReqPara.TimeCount*60*60
	case "day":
		tMin = d.NodeReqPara.TimeStamp - d.ServiceReqPara.TimeCount*60*60*24

	}
	s := dao.DashboardNodeDao{}
	s.QueryPara = dao.QueryPara{
		Name:      d.NodeReqPara.NodeName,
		TimeStamp: d.NodeReqPara.TimeStamp,
		TimeCount: d.NodeReqPara.TimeCount,
		TimeUnit:  d.NodeReqPara.TimeUnit,
		DuraTime:  d.NodeReqPara.DuraTime,
	}
	if d.NodeReqPara.NodeName == "" {
		d.NodeResp = NodeResp{
			NodeName:      "average",
			TimeUnit:      d.ServiceReqPara.TimeUnit,
			NodeTimestamp: d.NodeReqPara.TimeStamp,
		}
		d.NodeResp.NodeCount, d.NodeLogsData, err = s.GetTotalNodeData()
	} else {
		d.NodeResp = NodeResp{
			NodeName:      d.NodeReqPara.NodeName,
			TimeUnit:      d.NodeReqPara.TimeUnit,
			NodeTimestamp: d.NodeReqPara.TimeStamp,
		}
		d.NodeResp.NodeCount, d.NodeLogsData, err = s.GetNodeData()
	}
	lt, err := s.GetLimitTime()
	if err != nil {
		return err
	}

	if tMin > lt.MaxTime {
		d.NodeResp.IsOverMaxLimit = true
	}

	if d.NodeReqPara.TimeStamp < lt.MinTime {
		d.NodeResp.IsOverMinLimit = true
	}
	return nil
}
func (d *Dashboard) GetNodeListToObj() (count int, err error) {
	s := dao.DashboardNodeDao{}
	s.QueryPara = dao.QueryPara{
		Name:      d.NodeReqPara.NodeName,
		TimeStamp: d.NodeReqPara.TimeStamp,
		TimeCount: d.NodeReqPara.TimeCount,
		TimeUnit:  d.NodeReqPara.TimeUnit,
		DuraTime:  d.NodeReqPara.DuraTime,
	}
	count, d.NodeResp.NodeListResp.NodeListData, err = s.GetNodeListData()
	return
}
func (d *Dashboard) GetServiceListToObj() (count int, err error) {
	s := dao.DashboardServiceDao{}
	s.QueryPara = dao.QueryPara{
		Name:      d.ServiceReqPara.ServiceName,
		TimeStamp: d.ServiceReqPara.TimeStamp,
		TimeCount: d.ServiceReqPara.TimeCount,
		TimeUnit:  d.ServiceReqPara.TimeUnit,
		DuraTime:  d.ServiceReqPara.DuraTime,
	}
	count, d.ServiceResp.ServiceListResp.ServiceListData, err = s.GetServiceListData()
	return
}
