package yaml

type Env struct {
	Name  string `json:"env_name"`
	Value string `json:"env_value"`
}

type Volume struct {
	Dir               string `json:"container_dir"`
	TargetStorageName string `json:"target_storagename"`
	TargetDir         string `json:"target_dir"`
}

type Container struct {
	Name      string   `json:"container_name"`
	BaseImage string   `json:"container_baseimage"`
	WorkDir   string   `json:"container_workdir"`
	Ports     []int    `json:"container_ports"`
	Volumes   []Volume `json:"container_volumes"`
	Envs      []Env    `json:"container_envs"`
	Command   []string `json:"container_command"`
	Memory    string   `json:"container_memory"`
	CPU       string   `json:"container_cpu"`
}
