package model

import "time"

type Service struct {
	ID           int64     `json:"service_id" orm:"column(id)"`
	Name         string    `json:"service_name" orm:"column(name)"`
	ProjectID    string    `json:"service_project_id" orm:"column(project_id)"`
	ProjectName  string    `json:"service_project_name"`
	Comment      string    `json:"service_comment" orm:"column(comment)"`
	OwnerID      string    `json:"service_owner_id" orm:"column(owner_id)"`
	Status       int       `json:"service_status" orm:"column(status)"`
	Public       int       `json:"service_public" orm:"column(public)"`
	Deleted      int       `json:"service_deleted" orm:"column(deleted)"`
	CreationTime time.Time `json:"service_creation_time" orm:"column(creation_time)"`
	UpdateTime   time.Time `json:"service_update_time" orm:"column(update_time)"`
}
