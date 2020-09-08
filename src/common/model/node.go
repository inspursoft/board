package model

import (
	"time"
)

type NodeGroup struct {
	ID           int64     `json:"nodegroup_id" orm:"column(id)"`
	GroupName    string    `json:"nodegroup_name" orm:"column(name)"`
	Comment      string    `json:"nodegroup_comment" orm:"column(comment)"`
	OwnerID      int64     `json:"nodegroup_owner_id" orm:"column(owner_id)"`
	CreationTime time.Time `json:"nodegroup_creation_time" orm:"column(creation_time)"`
	UpdateTime   time.Time `json:"nodegroup_update_time" orm:"column(update_time)"`
	Deleted      int       `json:"nodegroup_deleted" orm:"column(deleted)"`
	Project      string    `json:"nodegroup_project" orm:"column(project_name)"`
	ProjectID    int64     `json:"nodegroup_project" orm:"column(project_id)"`
}

type NodeGroupDetail struct {
	NodeGroup NodeGroup `json:"nodegroup"`
	NodeList  []string  `json:"nodelist,omitempty"`
}

type NodeAvailableResources struct {
	NodeID       int    `json:"node_id"`
	NodeName     string `json:"node_name"`
	CPUAvail     string `json:"cpu_available"`
	MemAvail     string `json:"mem_available"`
	StorageAvail string `json:"storage_available"`
}

type NodeCli struct {
	NodeName string            `json:"node_name"`
	NodeType string            `json:"node_type"`
	NodeIP   string            `json:"node_ip"`
	Password string            `json:"node_password"`
	Labels   map[string]string `json:"labels,omitempty"`
	Taints   []Taint           `json:"taints,omitempty"`
}

type NodeControlStatus struct {
	NodeName          string            `json:"node_name"`
	NodeType          string            `json:"node_type"`
	NodeIP            string            `json:"node_ip"`
	NodePhase         string            `json:"node_phase"`
	NodeUnschedule    bool              `json:"node_unschedulable"`
	NodeDeletable     bool              `json:"node_deletable"`
	Service_Instances []ServiceInstance `json:"service_instances"`
}

type ServiceInstance struct {
	ProjectName         string   `json:"project_name"`
	ServiceInstanceName string   `json:"service_instance_name"`
	OwnerReferenceKind  []string `json:"owner_reference_kind"`
}

type EdgeNodeCli struct {
	NodeName     string `json:"name"`
	NodeIP       string `json:"node_ip"`
	CPUType      string `json:"cpu_type"`
	Password     string `json:"node_password"`
	Master       string `json:"master"`
	RegistryMode string `json:"registry_mode"`
}
