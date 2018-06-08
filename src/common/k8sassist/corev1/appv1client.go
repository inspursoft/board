package corev1

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/apps"
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	"io"
)

func NewAppV1Client(clientset *types.Clientset) AppV1ClientInterface {
	return &AppV1Client{
		Clientset: clientset,
	}
}

type AppV1Client struct {
	Clientset *types.Clientset
}

func (p *AppV1Client) Service(namespace string) ServiceClientInterface {
	return apps.NewServices(namespace, p.Clientset.CoreV1().Services(namespace))
}

func (p *AppV1Client) Deployment(namespace string) DeploymentClientInterface {
	return apps.NewDeployments(namespace, p.Clientset.AppsV1beta2().Deployments(namespace))
}

func (p *AppV1Client) Node() NodeClientInterface {
	return apps.NewNodes(p.Clientset.CoreV1().Nodes())
}

func (p *AppV1Client) Namespace() NamespaceClientInterface {
	return apps.NewNamespaces(p.Clientset.CoreV1().Namespaces())
}

func (p *AppV1Client) Scale(namespace string) ScaleClientInterface {
	return apps.NewScales(namespace, p.Clientset.ExtensionsV1beta1().Scales(namespace))
}

func (p *AppV1Client) ReplicaSet(namespace string) ReplicaSetClientInterface {
	return apps.NewReplicaSets(namespace, p.Clientset.AppsV1beta2().ReplicaSets(namespace))
}

func (p *AppV1Client) Pod(namespace string) PodClientInterface {
	return apps.NewPods(namespace, p.Clientset.CoreV1().Pods(namespace))
}

// AppV1ClientInterface level 1 interface to access others
type AppV1ClientInterface interface {
	Service(namespace string) ServiceClientInterface
	Deployment(namespace string) DeploymentClientInterface
	Node() NodeClientInterface
	Namespace() NamespaceClientInterface
	Scale(namespace string) ScaleClientInterface
	ReplicaSet(namespace string) ReplicaSetClientInterface
	Pod(namespace string) PodClientInterface
}

// ServiceCli interface has methods to work with Service resources in k8s-assist.
// How to:  serviceCli, err := k8sassist.NewServices(nameSpace)
//          service, err := serviceCli.Get(serviceName)
type ServiceClientInterface interface {
	Create(*model.Service) (*model.Service, []byte, error)
	Update(*model.Service) (*model.Service, []byte, error)
	UpdateStatus(*model.Service) (*model.Service, []byte, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Service, []byte, error)
	List() (*model.ServiceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Service, err error)
	CreateByYaml(io.Reader) (*model.Service, error)
	CheckYaml(io.Reader) (*model.Service, error)
}

// NodeCli Interface has methods to work with Node resources in k8s-assist.
// How to:  nodeCli, err := k8sassist.NewNodes()
//          nodeInstance, err := nodeCli.Get(nodename)
type NodeClientInterface interface {
	Create(*model.Node) (*model.Node, error)
	Update(*model.Node) (*model.Node, error)
	UpdateStatus(*model.Node) (*model.Node, error)
	Delete(name string) error
	Get(name string) (*model.Node, error)
	List() (*model.NodeList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Node, err error)
}

// NamespaceClientInterface Interface has methods to work with Namespace resources.
type NamespaceClientInterface interface {
	Create(*model.Namespace) (*model.Namespace, error)
	Update(*model.Namespace) (*model.Namespace, error)
	Delete(name string) error
	Get(name string) (*model.Namespace, error)
	List() (*model.NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

// ScaleClientInterface interface has methods on Scale resources in k8s-assist.
type ScaleClientInterface interface {
	Get(kind string, name string) (*model.Scale, error)
	Update(kind string, scale *model.Scale) (*model.Scale, error)
}

// ReplicaSetInterface has methods to work with ReplicaSet resources.
type ReplicaSetClientInterface interface {
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

// PodClientInterface has methods to work with Pod resources in k8s-assist.
// How to:  podCli, err := k8sassist.NewPods(nameSpace)
//          _, err := podCli.Update(&pod)
type PodClientInterface interface {
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
type DeploymentClientInterface interface {
	Create(*model.Deployment) (*model.Deployment, []byte, error)
	Update(*model.Deployment) (*model.Deployment, []byte, error)
	UpdateStatus(*model.Deployment) (*model.Deployment, []byte, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Deployment, []byte, error)
	//List(opts v1.ListOptions) (*DeploymentList, error)
	List() (*model.DeploymentList, error)
	Patch(name string, pt model.PatchType, data []byte, subresources ...string) (result *model.Deployment, err error)
	PatchToK8s(string, model.PatchType, *model.Deployment) (*model.Deployment, []byte, error)
	CreateByYaml(io.Reader) (*model.Deployment, error)
	CheckYaml(io.Reader) (*model.Deployment, error)
}
