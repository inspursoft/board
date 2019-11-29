package model

import (
	"time"
)

type JobStatusMO struct {
	ID           int64     `json:"job_id" orm:"column(id)"`
	Name         string    `json:"job_name" orm:"column(name)"`
	ProjectID    int64     `json:"job_project_id" orm:"column(project_id)"`
	ProjectName  string    `json:"job_project_name" orm:"column(project_name)"`
	Comment      string    `json:"job_comment" orm:"column(comment)"`
	OwnerID      int64     `json:"job_owner_id" orm:"column(owner_id)"`
	OwnerName    string    `json:"job_owner_name" orm:"column(owner_name)"`
	Status       int       `json:"job_status" orm:"column(status)"`
	Deleted      int       `json:"job_deleted" orm:"column(deleted)"`
	CreationTime time.Time `json:"job_creation_time" orm:"column(creation_time)"`
	UpdateTime   time.Time `json:"job_update_time" orm:"column(update_time)"`
	Source       int       `json:"job_source" orm:"column(source)"`
	Yaml         string    `json:"job_yaml" orm:"column(yaml)"`
}

func (j *JobStatusMO) TableName() string {
	return "job_status"
}

type PaginatedJobStatus struct {
	Pagination    *Pagination    `json:"pagination"`
	JobStatusList []*JobStatusMO `json:"job_status_list"`
}

type JobConfig struct {
	ID                    int64         `json:"job_id"`
	Name                  string        `json:"job_name"`
	ProjectID             int64         `json:"project_id"`
	ProjectName           string        `json:"project_name"`
	ContainerList         []Container   `json:"container_list"`
	NodeSelector          string        `json:"node_selector"`
	AffinityList          []JobAffinity `json:"affinity_list"`
	Parallelism           *int32        `json:"parallelism,omitempty"`
	Completions           *int32        `json:"completions,omitempty"`
	ActiveDeadlineSeconds *int64        `json:"active_Deadline_Seconds,omitempty"`
	BackoffLimit          *int32        `json:"backoff_Limit,omitempty"`
}

type JobAffinity struct {
	AntiFlag int      `json:"anti_flag"`
	JobNames []string `json:"job_names"`
}
