package corev1

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/apps"
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	"io"
	"net/http"
)

func NewAppV1Client(config *types.Config, clientset *types.Clientset, scaleGetter types.ScaleGetter) AppV1ClientInterface {
	return &AppV1Client{
		Config:      config,
		Clientset:   clientset,
		ScaleGetter: scaleGetter,
	}
}

type AppV1Client struct {
	Config      *types.Config
	Clientset   *types.Clientset
	ScaleGetter types.ScaleGetter
}

func (p *AppV1Client) Discovery() ServerVersionInterface {
	return apps.NewServerVersion(p.Clientset.Discovery())
}

func (p *AppV1Client) Service(namespace string) ServiceClientInterface {
	return apps.NewServices(namespace, p.Clientset.CoreV1().Services(namespace))
}

func (p *AppV1Client) Deployment(namespace string) DeploymentClientInterface {
	return apps.NewDeployments(namespace, p.Clientset.AppsV1().Deployments(namespace))
}

func (p *AppV1Client) Node() NodeClientInterface {
	return apps.NewNodes(p.Clientset.CoreV1().Nodes())
}

func (p *AppV1Client) Namespace() NamespaceClientInterface {
	return apps.NewNamespaces(p.Clientset.CoreV1().Namespaces())
}

func (p *AppV1Client) Scale(namespace string) ScaleClientInterface {
	return apps.NewScales(namespace, p.ScaleGetter.Scales(namespace))
}

func (p *AppV1Client) ReplicaSet(namespace string) ReplicaSetClientInterface {
	return apps.NewReplicaSets(namespace, p.Clientset.AppsV1().ReplicaSets(namespace))
}

func (p *AppV1Client) Pod(namespace string) PodClientInterface {
	return apps.NewPods(p.Clientset, p.Config, namespace, p.Clientset.CoreV1().Pods(namespace))
}

func (p *AppV1Client) AutoScale(namespace string) AutoscaleInterface {
	return apps.NewAutoScales(namespace, p.Clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace))
}

func (p *AppV1Client) PersistentVolume() PersistentVolumeInterface {
	return apps.NewPersistentVolume(p.Clientset.CoreV1().PersistentVolumes())
}

func (p *AppV1Client) PersistentVolumeClaim(namespace string) PersistentVolumeClaimInterface {
	return apps.NewPersistentVolumeClaim(namespace, p.Clientset.CoreV1().PersistentVolumeClaims(namespace))
}

func (p *AppV1Client) ConfigMap(namespace string) ConfigMapInterface {
	return apps.NewConfigMap(namespace, p.Clientset.CoreV1().ConfigMaps(namespace))
}

func (p *AppV1Client) StatefulSet(namespace string) StatefulSetClientInterface {
	return apps.NewStatefulSets(namespace, p.Clientset.AppsV1().StatefulSets(namespace))
}

func (p *AppV1Client) Job(namespace string) JobInterface {
	return apps.NewJob(namespace, p.Clientset.BatchV1().Jobs(namespace))
}

func (p *AppV1Client) DaemonSet(namespace string) DaemonSetInterface {
	return apps.NewDaemonSets(namespace, p.Clientset.AppsV1().DaemonSets(namespace))
}

func (p *AppV1Client) Proxy() ProxyInterface {
	return apps.NewProxy(p.Config)
}

func (p *AppV1Client) Extend() ExtendInterface {
	return apps.NewExtend(p.Clientset)
}

// AppV1ClientInterface level 1 interface to access others
type AppV1ClientInterface interface {
	Discovery() ServerVersionInterface
	Service(namespace string) ServiceClientInterface
	Deployment(namespace string) DeploymentClientInterface
	Node() NodeClientInterface
	Namespace() NamespaceClientInterface
	Scale(namespace string) ScaleClientInterface
	ReplicaSet(namespace string) ReplicaSetClientInterface
	Pod(namespace string) PodClientInterface
	AutoScale(namespace string) AutoscaleInterface
	PersistentVolume() PersistentVolumeInterface
	PersistentVolumeClaim(namespace string) PersistentVolumeClaimInterface
	ConfigMap(namespace string) ConfigMapInterface
	StatefulSet(namespace string) StatefulSetClientInterface
	Job(namespace string) JobInterface
	DaemonSet(namespace string) DaemonSetInterface
	Proxy() ProxyInterface
	Extend() ExtendInterface
}

// ServerVersionInterface has a method for retrieving the server's version.
type ServerVersionInterface interface {
	// ServerVersion retrieves and parses the server's version (git version).
	ServerVersion() (*model.KubernetesInfo, error)
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
	Patch(name string, pt model.PatchType, modelService *model.Service) (*model.Service, []byte, error)
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
	List(args ...string) (*model.NodeList, error)
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
	Get(resource model.GroupResource, name string) (*model.Scale, error)
	Update(resource model.GroupResource, scale *model.Scale) (*model.Scale, error)
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
	List(opts model.ListOptions) (*model.PodList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Pod, err error)
	GetLogs(name string, opts *model.PodLogOptions) (io.ReadCloser, error)
	ShellExec(podName, containerName string, cmd []string, ptyHandler model.PtyHandler) error
	CopyFromPod(podName, container, src, dest string, cmd []string) error
	CopyToPod(podName, container, src, dest string) error
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

// How to:  autoscaleCli, err := k8sassist.NewAutoscale(nameSpace)
//          _, err := autoscaleCli.Update(&autoscale)
type AutoscaleInterface interface {
	Create(*model.AutoScale) (*model.AutoScale, error)
	Update(*model.AutoScale) (*model.AutoScale, error)
	UpdateStatus(*model.AutoScale) (*model.AutoScale, error)
	Delete(name string) error
	//	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string) (*model.AutoScale, error)
	List() (*model.AutoScaleList, error)
}

// How to:  autoscaleCli, err := k8sassist.NewAutoscale(nameSpace)
//          _, err := autoscaleCli.Update(&autoscale)
type PersistentVolumeInterface interface {
	Create(*model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error)
	Update(*model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error)
	UpdateStatus(*model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error)
	Delete(name string) error
	//	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string) (*model.PersistentVolumeK8scli, error)
	List() (*model.PersistentVolumeList, error)
}

type PersistentVolumeClaimInterface interface {
	Create(*model.PersistentVolumeClaimK8scli) (*model.PersistentVolumeClaimK8scli, error)
	Update(*model.PersistentVolumeClaimK8scli) (*model.PersistentVolumeClaimK8scli, error)
	UpdateStatus(*model.PersistentVolumeClaimK8scli) (*model.PersistentVolumeClaimK8scli, error)
	Delete(name string) error
	//	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string) (*model.PersistentVolumeClaimK8scli, error)
	List() (*model.PersistentVolumeClaimList, error)
}

type ConfigMapInterface interface {
	Create(*model.ConfigMap) (*model.ConfigMap, error)
	Update(*model.ConfigMap) (*model.ConfigMap, error)
	//UpdateStatus(*model.ConfigMap) (*model.ConfigMap, error)
	Delete(name string) error
	//	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string) (*model.ConfigMap, error)
	List() (*model.ConfigMapList, error)
}

// StatefulSetClientInterface has methods to work with StatefulSet resources of adapter
type StatefulSetClientInterface interface {
	Create(*model.StatefulSet) (*model.StatefulSet, []byte, error)
	Update(*model.StatefulSet) (*model.StatefulSet, []byte, error)
	UpdateStatus(*model.StatefulSet) (*model.StatefulSet, []byte, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.StatefulSet, []byte, error)
	//List(opts v1.ListOptions) (*DeploymentList, error)
	List() (*model.StatefulSetList, error)
	Patch(name string, pt model.PatchType, data []byte, subresources ...string) (result *model.StatefulSet, err error)
	PatchToK8s(string, model.PatchType, *model.StatefulSet) (*model.StatefulSet, []byte, error)
	CreateByYaml(io.Reader) (*model.StatefulSet, error)
	CheckYaml(io.Reader) (*model.StatefulSet, error)
}
type JobInterface interface {
	Create(*model.Job) (*model.Job, []byte, error)
	Update(*model.Job) (*model.Job, []byte, error)
	UpdateStatus(*model.Job) (*model.Job, []byte, error)
	Delete(name string) error
	Get(name string) (*model.Job, []byte, error)
	List() (*model.JobList, error)
	Patch(name string, pt model.PatchType, data []byte, subresources ...string) (result *model.Job, err error)
	PatchToK8s(string, model.PatchType, *model.Job) (*model.Job, []byte, error)
	CreateByYaml(io.Reader) (*model.Job, error)
	CheckYaml(io.Reader) (*model.Job, error)
}

// DaemonSetInterface has methods to work with DaemonSet resources.
type DaemonSetInterface interface {
	Create(*model.DaemonSet) (*model.DaemonSet, []byte, error)
	Update(*model.DaemonSet) (*model.DaemonSet, []byte, error)
	UpdateStatus(*model.DaemonSet) (*model.DaemonSet, []byte, error)
	Delete(name string) error
	Get(name string) (*model.DaemonSet, []byte, error)
	List(opts model.ListOptions) (*model.DaemonSetList, error)
}

type ProxyInterface interface {
	ProxyAPI(apiProxyPrefix string) (http.Handler, error)
}

type ExtendInterface interface {
	ListSelectRelatePods(infos []*model.K8sInfo) (*model.PodList, error)
}
