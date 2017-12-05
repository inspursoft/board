package model

type SystemInfo struct {
	AuthMode         string `json:"auth_mode"`
	SetAdminPassword string `json:"set_auth_password"`
	InitProjectRepo  string `json:"init_project_repo"`
}
