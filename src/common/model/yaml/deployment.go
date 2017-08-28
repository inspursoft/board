package yaml

type NFSVolume struct {
	Name       string `json:"volume_name"`
	ServerName string `json:"server_name"`
	Path       string `json:"volume_path"`
}

type Deployment struct {
	Name          string      `json:"deployment_name"`
	Replicas      int         `json:"deployment_replicas"`
	VolumeList    []NFSVolume `json:"volume_list"`
	ContainerList []Container `json:"container_list"`
}
