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
