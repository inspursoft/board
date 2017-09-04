package model

import "git/inspursoft/board/src/common/model/yaml"

type ServiceConfig struct {
	ServiceID      int64           `json:"service_config_service_id" orm:"column(id)"`
	ProjectID      int64           `json:"service_config_project_id" orm:"column(project_id)"`
	Phase          string          `json:"service_config_phase"`
	ImageList      []string        `json:"service_image_list"`
	ServiceYaml    yaml.Service    `json:"service_yaml"`
	DeploymentYaml yaml.Deployment `json:"deployment_yaml"`
}

type ServiceConfigImage struct {
	ServiceID int64 `json:"service_config_id" orm:"column(service_id)"`
	ImageID   int64 `json:"service_image_id" orm:"column(image_id)"`
}
