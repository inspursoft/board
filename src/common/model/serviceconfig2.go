package model

import (
	"k8s.io/client-go/pkg/api/v1"
)

type ServiceConfig2 struct {
	Project    ProjectInfo              `json:"projectinfo"`
	Deployment v1.ReplicationController `json:"deployment_yaml"`
	Service    v1.Service               `json:"service_yaml"`
}

type ProjectInfo struct {
	ServiceID           int64    `json:"service_id" orm:"column(id)"`
	ProjectID           int64    `json:"project_id" orm:"column(project_id)"`
	ServiceName         string   `json:"service_name"`
	ProjectName         string   `json:"project_name"`
	Namespace           string   `json:"namespace"`
	Comment             string   `json:"comment"`
	Phase               string   `json:"config_phase"`
	ServiceExternalPath []string `json:"service_externalpath"`
}

type ConfigServiceStep struct {
	ProjectID           int64             `json:"project_id"`
	ProjectName         string            `json:"project_name"`
	ServiceID           int64             `json:"service_id"`
	ServiceName         string            `json:"service_name"`
	ImageList           []ImageIndex      `json:"image_list"`
	Instance            int               `json:"instance"`
	ContainerList       []Container       `json:"container_list"`
	ExternalServiceList []ExternalService `json:"external_service_list"`
}
