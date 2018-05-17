package corev1

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/apps"
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
)

type AppV1Client struct {
	Clientset *types.Clientset
}

type AppV1ClientInterface interface {
	Service(namespace string) ServiceCliInterface
	Deployment(namespace string) DeploymentCliInterface
	Node() NodeCliInterface
	Namespace() NamespaceCliInterface
	Scale(namespace string) ScaleCliInterface
	ReplicaSet(namespace string) ReplicaSetCliInterface
	Pod(namespace string) PodCliInterface
}

func NewAppV1Client(clientset *types.Clientset) *AppV1Client {
	return &AppV1Client{
		Clientset: clientset,
	}
}

func (p *AppV1Client) Namespace() NamespaceClientInterface {
	return &apps.NamespaceClient{Namespace: p.Clientset.CoreV1().Namespaces()}
}

func (p *AppV1Client) Node() NodeClientInterface {
	return nil
}

func (p *AppV1Client) Deployment() DeploymentClientInterface {
	return nil
}

func (p *AppV1Client) Service() ServiceClientInterface {
	return nil
}

// ServiceCli interface has methods to work with Service resources in k8s-assist.
// How to:  serviceCli, err := k8sassist.NewServices(nameSpace)
//          service, err := serviceCli.Get(serviceName)
type ServiceCliInterface interface {
	Create(*model.Service) (*model.Service, error)
	Update(*model.Service) (*model.Service, error)
	UpdateStatus(*model.Service) (*model.Service, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Service, error)
	List() (*model.ServiceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Service, err error)
}

// NamespaceCli Interface has methods to work with Namespace resources.
type DeploymentClientInterface interface {
	Create(*model.Namespace) (*model.Namespace, error)
	//Delete(*model.Namespace) error
	//Get(name string) (*model.Namespace, error)
	// List() (*model.NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

// NodeCli Interface has methods to work with Node resources in k8s-assist.
// How to:  nodeCli, err := k8sassist.NewNodes()
//          nodeInstance, err := nodeCli.Get(nodename)
type NodeCliInterface interface {
	Create(*model.Node) (*model.Node, error)
	Update(*model.Node) (*model.Node, error)
	UpdateStatus(*model.Node) (*model.Node, error)
	Delete(name string) error
	Get(name string) (*model.Node, error)
	List() (*model.NodeList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Node, err error)
}

// NamespaceCli Interface has methods to work with Namespace resources.
type NamespaceCliInterface interface {
	Create(*model.Namespace) (*model.Namespace, error)
	Update(*model.Namespace) (*model.Namespace, error)
	Delete(*model.Namespace) error
	Get(name string) (*model.Namespace, error)
	List() (*model.NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

// The ScaleCli interface has methods on Scale resources in k8s-assist.
type ScaleCliInterface interface {
	Get(kind string, name string) (*model.Scale, error)
	Update(kind string, scale *model.Scale) (*model.Scale, error)
}

// ReplicaSetInterface has methods to work with ReplicaSet resources.
type ReplicaSetCliInterface interface {
	Create(*model.ReplicaSet) (*model.ReplicaSet, error)
	Update(*model.ReplicaSet) (*model.ReplicaSet, error)
	UpdateStatus(*model.ReplicaSet) (*model.ReplicaSet, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.ReplicaSet, error)
	List(opts model.ListOptions) (*model.ReplicaSetList, error)
	//Watch(opts v1.ListOptions) (watch.Interface, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1beta1.ReplicaSet, err error)
	//ReplicaSetExpansion
}

// PodCli has methods to work with Pod resources in k8s-assist.
// How to:  podCli, err := k8sassist.NewPods(nameSpace)
//          _, err := podCli.Update(&pod)
type PodCliInterface interface {
	Create(*model.Pod) (*model.Pod, error)
	Update(*model.Pod) (*model.Pod, error)
	UpdateStatus(*model.Pod) (*model.Pod, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Pod, error)
	List() (*model.PodList, error)
	//List(opts v1.ListOptions) (*v1.PodList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Pod, err error)
}

// How to:  deploymentCli, err := k8sassist.NewDeployments(nameSpace)
//          _, err := deploymentCli.Update(&deployment)
type DeploymentCliInterface interface {
	Create(*model.Deployment) (*model.Deployment, error)
	Update(*model.Deployment) (*model.Deployment, error)
	UpdateStatus(*model.Deployment) (*model.Deployment, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Deployment, error)
	//List(opts v1.ListOptions) (*DeploymentList, error)
	List() (*model.DeploymentList, error)
	Patch(name string, pt model.PatchType, data []byte, subresources ...string) (result *model.Deployment, err error)
}
