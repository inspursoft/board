package model

type Role struct {
	ID      int64  `json:"role_id" orm:"column(id)"`
	Name    string `json:"role_name" orm:"column(name)"`
	Comment string `json:"role_comment" orm:"column(comment)"`
}
