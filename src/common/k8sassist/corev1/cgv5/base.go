package base

import (
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
	"k8s.io/client-go/tools/clientcmd"
)

func NewBaseClient(masterURL, kubeConfigPath string) (*rest.Config, *kubernetes.Clientset, scale.ScalesGetter, error) {
	//get config
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeConfigPath)
	if err != nil {
		return nil, nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	discoveryClient := memory.NewMemCacheClient(clientset.Discovery())
	expander := restmapper.NewShortcutExpander(restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient), discoveryClient)
	scaleGetter, err := scale.NewForConfig(config, expander, dynamic.LegacyAPIPathResolverFunc, scale.NewDiscoveryScaleKindResolver(discoveryClient))
	if err != nil {
		return nil, nil, nil, err
	}

	return config, clientset, scaleGetter, nil
}
