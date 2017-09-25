package model

import "time"

type ServiceStatus struct {
	ID           int64     `json:"service_id" orm:"column(id)"`
	Name         string    `json:"service_name" orm:"column(name)"`
	ProjectID    int64     `json:"service_project_id" orm:"column(project_id)"`
	ProjectName  string    `json:"service_project_name" orm:"column(project_name)"`
	Comment      string    `json:"service_comment" orm:"column(comment)"`
	OwnerID      int64     `json:"service_owner_id" orm:"column(owner_id)"`
	Status       int       `json:"service_status" orm:"column(status)"`
	Public       int       `json:"service_public" orm:"column(public)"`
	Deleted      int       `json:"service_deleted" orm:"column(deleted)"`
	CreationTime time.Time `json:"service_creation_time" orm:"column(creation_time)"`
	UpdateTime   time.Time `json:"service_update_time" orm:"column(update_time)"`
}

//func (s *Service) TableName() string {
//	return "service_status"
//}

type ServiceToggle struct {
	Toggle bool `json:"service_toggle"`
}
