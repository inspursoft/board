package collect

import (
	"github.com/astaxie/beego/orm"
	"git/inspursoft/board/src/collector/model/collect/dashboard"
)

func init() {
	orm.RegisterModel(
		new(Node), new(Pod), new(Service),
		new(ServiceKvMap), new(PodKvMap), new(dashboard.ServiceDashboardSecond),
		new(dashboard.ServiceDashboardMinute), new(dashboard.ServiceDashboardHour),
		new(dashboard.ServiceDashboardDay), new(dashboard.TimeListLog),
		new(dashboard.NodeDashboardMinute),new(dashboard.NodeDashboardHour),
		new(dashboard.NodeDashboardDay))
}

type Node struct {
	Id             int64    `json:"id" orm:"pk;auto"`
	NodeName       string `json:"pod_name" orm:"column(node_name)"`
	NumbersCpuCore string  `json:"pod_name" orm:"column(numbers_cpu_core)"`
	NumbersGpuCore string  `json:"pod_name" orm:"column(numbers_gpu_core)"`
	MemorySize     string  `json:"pod_name" orm:"column(memory_size)"`
	PodLimit       string  `json:"pod_name" orm:"column(pod_limit)"`
	CreateTime     string `json:"Creat_time" orm:"column(create_time)"`
	InternalIp     string `json:"ip" orm:"column(ip)"`
	CpuUsage       float32 `json:"cpu_usage" orm:"column(cpu_usage)"`
	MemUsage       float32 `json:"mem_usage" orm:"column(mem_usage)"`
	TimeListId     int64 `json:"pod_name" orm:"column(time_list_id)"`
	StorageTotal   int64 `json:"pod_name" orm:"column(storage_total)"`
	StorageUse     int64 `json:"pod_name" orm:"column(storage_use)"`
}

type Pod struct {
	Id         int64    `json:"id" orm:"pk;auto"`
	PodName    string `json:"pod_name" orm:"column(pod_name)"`
	PodHostIP  string`json:"pod_name" orm:"column(pod_hostIP)"`
	CreateTime string `json:"Creat_time" orm:"column(create_time)"`
	TimeListId int64 `json:"pod_name" orm:"column(time_list_id)"`
}

type Service struct {
	Id          int64    `json:"id" orm:"pk;auto"`
	ServiceName string `json:"pod_name" orm:"column(service_name)"`
	CreateTime  string `json:"Creat_time" orm:"column(create_time)"`
	TimeListId  int64 `json:"pod_name" orm:"column(time_list_id)"`
}
type ServiceKvMap struct {
	Id         int64    `json:"id" orm:"pk;auto"`
	Name       string    `json:"pod_name" orm:"column(name)"`
	Value      string    `json:"pod_name" orm:"column(value)"`
	Belong     string    `json:"pod_name" orm:"column(belong)"`
	TimeListId int64 `json:"pod_name" orm:"column(time_list_id)"`
}
type PodKvMap struct {
	Id         int64    `json:"id" orm:"pk;auto"`
	Name       string    `json:"pod_name" orm:"column(name)"`
	Value      string    `json:"pod_name" orm:"column(value)"`
	Belong     string    `json:"pod_name" orm:"column(belong)"`
	TimeListId int64 `json:"pod_name" orm:"column(time_list_id)"`
}
