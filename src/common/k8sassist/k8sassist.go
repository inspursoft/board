package k8sassist

import (
	v1 "github.com/inspursoft/board/src/common/k8sassist/corev1"
	base "github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5"
)

type K8sAssistConfig struct {
	K8sMasterURL   string
	KubeConfigPath string
}

type K8sAssistClient struct {
	config *K8sAssistConfig
	appV1  *v1.AppV1Client
}

func NewK8sAssistClient(c *K8sAssistConfig) *K8sAssistClient {
	return &K8sAssistClient{
		config: c,
	}
}

func (c *K8sAssistClient) AppV1() v1.AppV1ClientInterface {
	config, clientset, scaleGetter, err := base.NewBaseClient(c.config.K8sMasterURL, c.config.KubeConfigPath)
	if err != nil {
		panic(err)
	}
	return v1.NewAppV1Client(config, clientset, scaleGetter)
}
