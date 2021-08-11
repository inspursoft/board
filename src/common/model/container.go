package model

type Container struct {
	Name       string `json:"name"`
	WorkingDir string `json:"working_dir"`
	//VolumeMounts  VolumeStruct `json:"volume_mount"`
	VolumeMounts  []VolumeMountStruct `json:"volume_mounts"`
	Image         ImageIndex          `json:"image"`
	Env           []EnvStructCont     `json:"env"`
	ContainerPort []int               `json:"container_port"`
	Command       string              `json:"command"`
	CPURequest    string              `json:"cpu_request"`
	MemRequest    string              `json:"mem_request"`
	CPULimit      string              `json:"cpu_limit"`
	MemLimit      string              `json:"mem_limit"`
	GPULimit      string              `json:"gpu_limit"`
}

type EnvStructCont struct {
	EnvName          string `json:"dockerfile_envname"`
	EnvValue         string `json:"dockerfile_envvalue"`
	EnvConfigMapName string `json:"configmap_name"`
	EnvConfigMapKey  string `json:"configmap_key"`
}

type VolumeStruct struct {
	TargetStorageService string `json:"target_storage_service"`
	TargetPath           string `json:"target_path"`
	TargetFile           string `json:"target_file"`
	VolumeName           string `json:"volume_name"`
	ContainerPath        string `json:"container_path"`
	ContainerFile        string `json:"container_file"`
	//mount type: 0, folder; 1, file
	MountTypeFlag int `json:"mount_type_flag"`
}

type VolumeMountStruct struct {
	VolumeType    string `json:"volume_type"`
	VolumeName    string `json:"volume_name"`
	ContainerPath string `json:"container_path"`
	ContainerFile string `json:"container_file"`
	//mount type: 0, folder; 1, file
	ContainerPathFlag    int    `json:"container_path_flag"`
	TargetStorageService string `json:"target_storage_service"`
	TargetPath           string `json:"target_path"`
	TargetFile           string `json:"target_file"`
	TargetPVC            string `json:"target_pvc"`
	TargetConfigMap      string `json:"target_configmap"`
}
