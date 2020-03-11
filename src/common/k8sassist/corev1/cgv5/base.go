package base

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewBaseClient(masterURL string, kubeConfigPath string) (*kubernetes.Clientset, error) {
	
	// config 获取支持 url 和 path 方式，通过 BuildConfigFromFlags() 函数获取 restclient.Config 对象，用来下边根据该 config 对象创建 client 集合
	//get config
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeConfigPath)
	if err != nil {
		return nil, err
	}

	// 根据获取的 config 来创建一个 clientset 对象。通过调用 NewForConfig 函数创建 clientset 对象。
	// NewForConfig 函数具体实现就是初始化 clientset 中的每个 client，基本涵盖了 k8s 内各种类型
	return kubernetes.NewForConfig(config)
}
