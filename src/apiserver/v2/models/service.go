package models

import (
	"time"
)

type Service struct {
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
	ServiceYaml    string    `json:"service_yaml" orm:"column(service_yaml)"`
	DeploymentYaml string    `json:"deployment_yaml" orm:"column(deployment_yaml)"`
}
