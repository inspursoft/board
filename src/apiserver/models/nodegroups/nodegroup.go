package nodegroups

import (
	"time"
)

type NodeGroup struct {
	ID           int64     `json:"nodegroup_id"`
	GroupName    string    `json:"nodegroup_name"`
	Comment      string    `json:"nodegroup_comment"`
	OwnerID      int64     `json:"nodegroup_owner_id"`
	CreationTime time.Time `json:"nodegroup_creation_time"`
	UpdateTime   time.Time `json:"nodegroup_update_time"`
	Deleted      int       `json:"nodegroup_deleted"`
	Project      string    `json:"nodegroup_project"`
	ProjectID    int64     `json:"nodegroup_project"`
}
