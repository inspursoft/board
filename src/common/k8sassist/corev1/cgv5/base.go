package base

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewBaseClient(masterURL string) (*kubernetes.Clientset, error) {
	//get config
	config, err := clientcmd.BuildConfigFromFlags(masterURL, "")
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
