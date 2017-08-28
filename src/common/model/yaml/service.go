package yaml

type Service struct {
	Name          string   `json:"service_name"`
	NodePorts     []int    `json:"service_nodeports"`
	ExternalPaths []string `json:"service_externalpaths"`
	Selectors     []string `json:"service_selectors"`
}
