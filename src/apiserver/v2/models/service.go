package models

import (
	"time"
)


type Service struct {
	ID             int64     `json:"service_id" orm:"column(id)"`
	Name           string    `json:"service_name" orm:"column(name)"`
	ProjectName    string    `json:"service_project_name" orm:"column(project_name)"`
	OwnerName      string    `json:"service_owner_name" orm:"column(owner_name)"`
	Public         int       `json:"service_public" orm:"column(public)"`
	Phase          int       `json:"service_phase" orm:"column(phase)"`
        Type           int       `json:"service_type" orm:"column(type)"`
	Yaml           string    `json:"service_yaml" orm:"column(yaml)"`
	Comment        string    `json:"service_comment" orm:"column(comment)"`
	CreationTime   time.Time `json:"service_creation_time" orm:"column(creation_time)"`
	UpdateTime     time.Time `json:"service_update_time" orm:"column(update_time)"`
}
