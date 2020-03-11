package model

type ProjectMember struct {
	ID        int64  `json:"project_member_id" orm:"column(id)"`
	UserID    int64  `json:"project_member_user_id" orm:"column(user_id)"`
	Username  string `json:"project_member_username" orm:"column(username)"`
	ProjectID int64  `json:"project_member_project_id" orm:"column(project_id)"`
	RoleID    int64  `json:"project_member_role_id" orm:"column(role_id)"`
}
