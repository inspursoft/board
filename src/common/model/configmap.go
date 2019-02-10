package model

type ConfigMapStruct struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	DataList  map[string]string `json:"datalist,omitempty"`
}
