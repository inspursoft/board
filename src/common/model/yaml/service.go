package yaml

type ExternalStruct struct {
	ContainerName string `json:"service_containername"`
	ContainerPort int    `json:"service_containerport"`
	NodePort      int    `json:"service_nodeport"`
	ExternalPath  string `json:"service_externalpath"`
}

type Service struct {
	Name      string           `json:"service_name"`
	External  []ExternalStruct `json:"service_external"`
	Selectors []string         `json:"service_selectors"`
}
