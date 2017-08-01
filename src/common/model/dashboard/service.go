package dashboard

type TimeListLog struct {
	Id         int64    `json:"id" orm:"pk;auto"`
	RecordTime int64    `json:"pod_name" orm:"column(record_time)"`
}
type ServiceDashboardSecond struct {
	Id              int64    `json:"id" orm:"pk;auto"`
	ServiceName     string    `json:"ServiceName" orm:"column(service_name)"`
	PodNumber       int64    `json:"PodNumber" orm:"column(pod_number)"`
	ContainerNumber int64    `json:"ContainerNumber" orm:"column(container_number)"`
	TimeListId      int64    `json:"TimeListId" orm:"column(time_list_id)"`
}
type ServiceDashboardMinute struct {
	Id              int64    `json:"id" orm:"pk;auto"`
	ServiceName     string    `json:"pod_name" orm:"column(service_name)"`
	PodNumber       int64    `json:"pod_name" orm:"column(pod_number)"`
	ContainerNumber int64    `json:"pod_name" orm:"column(container_number)"`
	TimeListId      int64    `json:"TimeListId" orm:"column(time_list_id)"`
}
type ServiceDashboardHour struct {
	Id              int64    `json:"id" orm:"pk;auto"`
	ServiceName     string    `json:"pod_name" orm:"column(service_name)"`
	PodNumber       int64    `json:"pod_name" orm:"column(pod_number)"`
	ContainerNumber int64    `json:"pod_name" orm:"column(container_number)"`
	TimeListId      int64    `json:"TimeListId" orm:"column(time_list_id)"`
}
type ServiceDashboardDay struct {
	Id              int64    `json:"id" orm:"pk;auto"`
	ServiceName     string    `json:"pod_name" orm:"column(service_name)"`
	PodNumber       int64    `json:"pod_name" orm:"column(pod_number)"`
	ContainerNumber int64    `json:"pod_name" orm:"column(container_number)"`
	TimeListId      int64    `json:"TimeListId" orm:"column(time_list_id)"`
}
