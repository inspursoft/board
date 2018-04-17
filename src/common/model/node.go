package model

import (
	"time"
)

type NodeGroup struct {
	ID           int64     `json:"nodegroup_id" orm:"column(id)"`
	GroupName    string    `json:"nodegroup_name" orm:"column(name)"`
	Comment      string    `json:"nodegroup_comment" orm:"column(comment)"`
	CreationTime time.Time `json:"nodegroup_creation_time" orm:"column(creation_time)"`
	UpdateTime   time.Time `json:"nodegroup_update_time" orm:"column(update_time)"`
	Deleted      int       `json:"nodegroup_deleted" orm:"column(deleted)"`
	Project      string    `json:"nodegroup_project" orm:"column(group_project)"`
}
