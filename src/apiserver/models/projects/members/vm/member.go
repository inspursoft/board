package vm

type ProjectMember struct {
	ID        int64  `json:"project_member_id"`
	UserID    int64  `json:"project_member_user_id"`
	Username  string `json:"project_member_username"`
	ProjectID int64  `json:"project_member_project_id"`
	RoleID    int64  `json:"project_member_role_id"`
}
