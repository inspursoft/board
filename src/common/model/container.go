package model

type Container struct {
	Name          string
	WorkingDir    string
	VolumeMounts  string
	Env           EnvStruct
	ContainerPort []int
	Command       string
}
