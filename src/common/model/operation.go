package model

import (
	"time"
)

type Operation struct {
	ID int64 `json:"operation_id" orm:"column(id)"`
	//Comment string `json:"operation_comment" orm:"column(comment)"`
	//Tag          string    `json:"operation_tag" orm:"column(tag)"`
	CreationTime time.Time `json:"operation_creation_time" orm:"column(creation_time)"`
	UpdateTime   time.Time `json:"operation_update_time" orm:"column(update_time)"`
	Deleted      int       `json:"operation_deleted" orm:"column(deleted)"`
	UserID       int64     `json:"operation_user_id" orm:"column(user_id)"`
	UserName     string    `json:"operation_user_name" orm:"column(user_name)"`
	ProjectID    int64     `json:"operation_project_id" orm:"column(project_id)"`
	ProjectName  string    `json:"operation_project_name" orm:"column(project_name)"`
	ObjectType   string    `json:"operation_object_type" orm:"column(object_type)"`
	ObjectName   string    `json:"operation_object_name" orm:"column(object_name)"`
	Action       string    `json:"operation_action" orm:"column(action)"`
	Status       string    `json:"operation_status" orm:"column(status)"`
	Path         string    `json:"operation_path" orm:"column(path)"`
}

type PaginatedOperations struct {
	Pagination    *Pagination  `json:"pagination"`
	OperationList []*Operation `json:"operation_list"`
}

type OperationParam struct {
	Action   string
	User     string
	Object   string
	Status   string
	Fromdate int64
	Todate   int64
}
