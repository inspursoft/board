package yaml

type NortPort struct {
	ContainerPort int `json:"container_port"`
	ExternalPort  int `json:"node_port"`
}

type ExternalPath struct {
	ContainerPort int    `json:"container_port"`
	Path          string `json:"external_path"`
}

type Service struct {
	Name          string         `json:"service_name"`
	NodePorts     []NortPort     `json:"service_nodeports"`
	ExternalPaths []ExternalPath `json:"service_externalpaths"`
	Selectors     []string       `json:"service_selectors"`
}
