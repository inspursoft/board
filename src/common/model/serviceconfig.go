package model

import "github.com/inspursoft/board/src/common/model/yaml"

type ServiceConfig struct {
	ServiceID      int64           `json:"service_id" orm:"column(id)"`
	ProjectID      int64           `json:"project_id" orm:"column(project_id)"`
	ProjectName    string          `json:"project_name"`
	Phase          string          `json:"config_phase"`
	DeploymentYaml yaml.Deployment `json:"deployment_yaml"`
	ServiceYaml    yaml.Service    `json:"service_yaml"`
}

type ServiceConfigImage struct {
	ServiceID int64 `json:"service_config_id" orm:"column(service_id)"`
	ImageID   int64 `json:"service_image_id" orm:"column(image_id)"`
}

type PortsServiceYaml struct {
	Port       int
	TargetPort int `yaml:"targetPort,flow"`
	NodePort   int `yaml:"nodePort,flow"`
}

type SelectorServiceYaml struct {
	App string
}

//type ServiceStructYaml struct {
//	ApiVersion string `yaml:"apiVersion,flow"`
//	Kind       string
//	Metadata   struct {
//		Name   string
//		Labels struct {
//			App string
//		}
//	}
//	Spec struct {
//		Tpe      string              `yaml:"type,flow,omitempty"`
//		Ports    []PortsServiceYaml  `yaml:",omitempty"`
//		Selector SelectorServiceYaml `yaml:",omitempty"`
//	} `yaml:",omitempty"`
//}

type PortsDeploymentYaml struct {
	ContainerPort int `yaml:"containerPort,flow"`
}

type VolumeMountDeploymentYaml struct {
	Name      string
	MountPath string `yaml:"mountPath,flow"`
}

type EnvDeploymentYaml struct {
	Name  string
	Value string
}

type VolumesDeploymentYaml struct {
	Name string `yaml:",omitempty"`
	Nfs  struct {
		Server string
		Path   string
	} `yaml:",omitempty"`
	HostPath struct {
		Path string
	} `yaml:"hostPath,flow,omitempty"`
}

type ContainersDeploymentYaml struct {
	Name       string
	Image      string
	Workingdir string   `yaml:",omitempty"`
	Command    []string `yaml:",omitempty"`
	Resources  struct {
		Requests struct {
			Cpu    string `yaml:",omitempty"`
			Memory string `yaml:",omitempty"`
		}
		Limits struct {
			Cpu    string `yaml:",omitempty"`
			Memory string `yaml:",omitempty"`
		}
	} `yaml:",omitempty"`
	Ports       []PortsDeploymentYaml       `yaml:",omitempty"`
	VolumeMount []VolumeMountDeploymentYaml `yaml:"VolumeMount,omitempty,flow"`
	Env         []EnvDeploymentYaml         `yaml:",omitempty"`
}

type DeploymentStructYaml struct {
	ApiVersion string `yaml:"apiVersion,flow"`
	Kind       string
	Metadata   struct {
		Name string
	} `yaml:",omitempty"`
	Spec struct {
		Replicas int `yaml:",omitempty"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string
				} `yaml:",omitempty"`
			} `yaml:",omitempty"`
			Spec struct {
				Volumes    []VolumesDeploymentYaml    `yaml:",omitempty"`
				Containers []ContainersDeploymentYaml `yaml:",omitempty"`
			} `yaml:",omitempty"`
		} `yaml:",omitempty"`
	} `yaml:",omitempty"`
}

type ServiceProject struct {
	ProjectID   int64  `json:"project_id"`
	ProjectName string `json:"project_name"`
}

type ConfigServiceStep struct {
	ProjectID           int64             `json:"project_id"`
	ProjectName         string            `json:"project_name"`
	ServiceID           int64             `json:"service_id"`
	ServiceName         string            `json:"service_name"`
	ServiceType         int               `json:"service_type"`
	Public              int               `json:"service_public"`
	NodeSelector        string            `json:"node_selector"`
	Instance            int               `json:"instance"`
	ClusterIP           string            `json:"cluster_ip"`
	ContainerList       []Container       `json:"container_list"`
	InitContainerList   []Container       `json:"initcontainer_list"`
	ExternalServiceList []ExternalService `json:"external_service_list"`
	AffinityList        []Affinity        `json:"affinity_list"`
	SessionAffinityFlag int               `json:"session_affinity_flag"`
	SessionAffinityTime int               `json:"session_affinity_time"`
}

type Affinity struct {
	AntiFlag     int      `json:"anti_flag"`
	ServiceNames []string `json:"service_names"`
}
