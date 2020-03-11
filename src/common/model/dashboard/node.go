package dashboard

var NodeStatusDashboardList []NodeStatusDashboard

type NodeStatusDashboard struct {
	NodeName       string `json:"node_name"`
	NodeTimeUnit   string `json:"node_time_unit"`
	NodeCount      string `json:"node_count"`
	NodeStatusLogs []NodeStatus `json:"node_status_log"`
}
type NodeStatus struct {
	CpuCoreInfo    string `json:"cpu_core"`
	MemorySizeInfo string `json:"memory_size"`
	PodLimitInfo   string `json:"pod_limit"`
}
type NodeDashboardMinute struct {
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
type NodeDashboardHour struct {
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
type NodeDashboardDay struct {
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
