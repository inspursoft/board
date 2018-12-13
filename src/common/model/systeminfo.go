package model

type SystemInfo struct {
	BoardHost         string `json:"board_host"`
	AuthMode          string `json:"auth_mode"`
	SetAdminPassword  string `json:"set_auth_password"`
	InitProjectRepo   string `json:"init_project_repo"`
	SyncK8s           string `json:"sync_k8s"`
	RedirectionURL    string `json:"redirection_url"`
	Version           string `json:"board_version"`
	DNSSuffix         string `json:"dns_suffix"`
	KubernetesVersion string `json:"kubernetes_version"`
}

// Info contains versioning information.
// TODO: Add []string of api versions supported? It's still unclear
// how we'll want to distribute that information.
type KubernetesInfo struct {
	Major        string `json:"major"`
	Minor        string `json:"minor"`
	GitVersion   string `json:"git_version"`
	GitCommit    string `json:"git_commit"`
	GitTreeState string `json:"git_tree_state"`
	BuildDate    string `json:"build_date"`
	GoVersion    string `json:"go_version"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}
