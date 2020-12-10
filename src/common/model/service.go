package model

import (
	"time"
)

const (
	ServiceTypeUnknown = iota
	ServiceTypeNormalNodePort
	ServiceTypeHelm
	ServiceTypeDeloymentOnly
	ServiceTypeClusterIP
	ServiceTypeStatefulSet
	ServiceTypeJob
	ServiceTypeEdgeComputing
)

type ServiceStatus struct {
	ID             int64     `json:"service_id" orm:"column(id)"`
	Name           string    `json:"service_name" orm:"column(name)"`
	ProjectID      int64     `json:"service_project_id" orm:"column(project_id)"`
	ProjectName    string    `json:"service_project_name" orm:"column(project_name)"`
	Comment        string    `json:"service_comment" orm:"column(comment)"`
	OwnerID        int64     `json:"service_owner_id" orm:"column(owner_id)"`
	OwnerName      string    `json:"service_owner_name" orm:"column(owner_name)"`
	Status         int       `json:"service_status" orm:"column(status)"`
	Type           int       `json:"service_type" orm:"column(type)"`
	Public         int       `json:"service_public" orm:"column(public)"`
	Deleted        int       `json:"service_deleted" orm:"column(deleted)"`
	CreationTime   time.Time `json:"service_creation_time" orm:"column(creation_time)"`
	UpdateTime     time.Time `json:"service_update_time" orm:"column(update_time)"`
	Source         int       `json:"service_source" orm:"column(source)"`
	SourceID       int64     `json:"service_sourceid" orm:"column(source_id)"`
	ServiceYaml    string    `json:"service_yaml" orm:"column(service_yaml)"`
	DeploymentYaml string    `json:"deployment_yaml" orm:"column(deployment_yaml)"`
}

type ServiceStatusFilter struct {
	Name      string
	ProjectID int64
	Source    *int
	SourceID  *int64
}

type ServiceStatusMO struct {
	ServiceStatus
	IsMember int `json:"service_is_member" orm:"column(is_member)"`
}

type PaginatedServiceStatus struct {
	Pagination        *Pagination        `json:"pagination"`
	ServiceStatusList []*ServiceStatusMO `json:"service_status_list"`
}

type ServiceInfoStruct struct {
	NodePort          []int32            `json:"node_Port,omitempty"`
	NodeName          []NodeAddress      `json:"node_Name,omitempty"`
	ServiceContainers []ServiceContainer `json:"service_Containers,omitempty"`
}

type ServiceToggle struct {
	Toggle int `json:"service_toggle"`
}

type ServicePublicityUpdate struct {
	Public int `json:"service_public"`
}

type ServiceScale struct {
	Replica int32 `json:"service_scale"`
}

type ScaleStatus struct {
	DesiredInstance   int32 `json:"desired_instance"`
	AvailableInstance int32 `json:"available_instance"`
}

type ExternalService struct {
	ContainerName      string       `json:"container_name"`
	NodeConfig         NodeType     `json:"node_config"`
	LoadBalancerConfig LoadBalancer `json:"load_balancer_config"`
}

type NodeType struct {
	TargetPort int `json:"target_port"`
	NodePort   int `json:"node_port"`
	Port       int `json:"port"`
}

type LoadBalancer struct {
	ExternalAccess string `json:"external_access"`
}

type ServiceAutoScale struct {
	ID         int64  `json:"hpa_id" orm:"column(id)"`
	HPAName    string `json:"hpa_name" orm:"column(name)"`
	ServiceID  int64  `json:"service_id" orm:"column(service_id)"`
	MinPod     int    `json:"min_pod" orm:"column(min_pod)"`
	MaxPod     int    `json:"max_pod" orm:"column(max_pod)"`
	CPUPercent int    `json:"cpu_percent" orm:"column(cpu_percent)"`
	HPAStatus  int    `json:"hpa_status" orm:"column(status)"`
}

type PodMO struct {
	Name        string    `json:"name" `
	ProjectName string    `json:"project_name"`
	Spec        PodSpecMO `json:"spec"`
}

type PodSpecMO struct {
	Containers []ContainerMO `json:"containers"`
}

type ContainerMO struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
