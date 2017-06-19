package models

import (
	"time"
)

// User holds the details of a user.
type User struct {
	UserID              int       `orm:"pk;auto;column(user_id)" json:"user_id"`
	Username            string    `orm:"column(username)" json:"user_name"`
	Email               string    `orm:"column(email)" json:"user_email"`
	Password            string    `orm:"column(password)" json:"user_password"`
	Realname            string    `orm:"column(realname)" json:"user_realname"`
	Comment             string    `orm:"column(comment)" json:"user_comment"`
	Deleted             int       `orm:"column(deleted)" json:"user_deleted"`
	HasSystemAdminRole  int       `orm:"column(sysadmin_flag)" json:"user_system_admin"`
	HasProjectAdminRole int       `orm:"column(sysadmin_flag)" json:"user_project_admin"`
	ResetUUID           string    `orm:"column(reset_uuid)" json:"user_reset_uuid"`
	Salt                string    `orm:"column(salt)" json:"user_salt"`
	CreationTime        time.Time `orm:"creation_time" json:"user_creation_time"`
	UpdateTime          time.Time `orm:"update_time" json:"user_update_time"`
}
