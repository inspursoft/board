package model

type SystemInfo struct {
	BoardHost        string `json:"board_host"`
	AuthMode         string `json:"auth_mode"`
	SetAdminPassword string `json:"set_auth_password"`
	InitProjectRepo  string `json:"init_project_repo"`
}
