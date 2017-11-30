package model

type Container struct {
	Name          string       `json:"name"`
	WorkingDir    string       `json:"working_dir"`
	VolumeMounts  VolumeStruct `json:"volume_mount"`
	Image         ImageIndex   `json:"image"`
	Env           []EnvStruct  `json:"env"`
	ContainerPort []int        `json:"container_port"`
	Command       string       `json:"command"`
}

type VolumeStruct struct {
	TargetStorageService string `json:"target_storage_service"`
	TargetPath           string `json:"target_path"`
	VolumeName           string `json:"volume_name"`
	ContainerPath        string `json:"container_path"`
}
