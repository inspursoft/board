package model

type Config struct {
	Name    string `json:"name" orm:"column(name);pk"`
	Value   string `json:"value" orm:"column(value)"`
	Comment string `json:"comment" orm:"column(comment)"`
}

type K8SProxyConfig struct {
	Enable bool `json:"enable"`
}
