package model

import (
	"time"
)

type PodPhase string

const (
	PodPending   PodPhase = "Pending"
	PodRunning   PodPhase = "Running"
	PodSucceeded PodPhase = "Succeeded"
	PodFailed    PodPhase = "Failed"
	PodUnknown   PodPhase = "Unknown"
)

type ResourceName string

const (
	ResourceCPU              ResourceName = "cpu"
	ResourceMemory           ResourceName = "memory"
	ResourceStorage          ResourceName = "storage"
	ResourceEphemeralStorage ResourceName = "ephemeral-storage"
	ResourceNvidiaGPU        ResourceName = "alpha.kubernetes.io/nvidia-gpu"
)

type NodePhase string

const (
	NodePending    NodePhase = "Pending"
	NodeRunning    NodePhase = "Running"
	NodeTerminated NodePhase = "Terminated"
)

type NodeConditionType string

const (
	NodeReady              NodeConditionType = "Ready"
	NodeOutOfDisk          NodeConditionType = "OutOfDisk"
	NodeMemoryPressure     NodeConditionType = "MemoryPressure"
	NodeDiskPressure       NodeConditionType = "DiskPressure"
	NodeNetworkUnavailable NodeConditionType = "NetworkUnavailable"
	NodeConfigOK           NodeConditionType = "ConfigOK"
)

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

type NodeAddressType string

const (
	NodeHostName    NodeAddressType = "Hostname"
	NodeExternalIP  NodeAddressType = "ExternalIP"
	NodeInternalIP  NodeAddressType = "InternalIP"
	NodeExternalDNS NodeAddressType = "ExternalDNS"
	NodeInternalDNS NodeAddressType = "InternalDNS"
)

// should call kubernetes Quantity String() func.
type Quantity string

type ResourceList map[ResourceName]Quantity

type ObjectMeta struct {
	Name              string
	Namespace         string
	CreationTimestamp time.Time
	DeletionTimestamp *time.Time
	Labels            map[string]string
}

type Node struct {
	ObjectMeta
	NodeIP        string
	Unschedulable bool
	Groups        map[string]string
	Status        NodeStatus
}

type NodeStatus struct {
	Capacity    ResourceList
	Allocatable ResourceList
	Phase       NodePhase
	Conditions  []NodeCondition
	Addresses   []NodeAddress
}

type NodeAddress struct {
	Type    NodeAddressType
	Address string
}

type NodeCondition struct {
	Type    NodeConditionType
	Status  ConditionStatus
	Reason  string
	Message string
}

type NodeList struct {
	Items []Node
}

type Namespace struct {
	ObjectMeta
	NamespacePhase string
}

type NamespaceList struct {
	Items []Namespace
}

type ServicePort struct {
	Name       string
	Protocol   string
	Port       int32
	TargetPort int32
	NodePort   int32
}

type Service struct {
	ObjectMeta
	Ports       []ServicePort
	Selector    map[string]string
	ClusterIP   string
	Type        string
	ExternalIPs []string
	//SessionAffinity ServiceAffinity
	ExternalName string
}

type ServiceList struct {
	Items []Service
}

type ScaleState struct {
	Replicas       int32
	Selector       map[string]string
	TargetSelector string
}

// represents a scaling request for a resource.
type Scale struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Replicas  int32
	Status    ScaleState
}

// DeploymentStatus is the most recently observed status of the Deployment.
type DeploymentStatus struct {
	Replicas            int32
	UpdatedReplicas     int32
	AvailableReplicas   int32
	UnavailableReplicas int32
}

// DeploymentSpec is the specification of the desired behavior of the Deployment.
type DeploymentSpec struct {
	Replicas int32
	Selector map[string]string
	Template PodTemplateSpec
	Paused   bool //TODO for pause
	//RollbackTo *RollbackConfig //TODO
}

type Deployment struct {
	ObjectMeta
	Spec   DeploymentSpec
	Status DeploymentStatus
}

type DeploymentList struct {
	Items []Deployment
}

type PodList struct {
	Items []Pod
}

type Pod struct {
	ObjectMeta
	Spec   PodSpec
	Status PodStatus
}

type PodSpec struct {
	Volumes        []Volume
	InitContainers []Container
	Containers     []Container
	NodeSelector   map[string]string
	NodeName       string
	HostNetwork    bool
}

type PodStatus struct {
	Phase     PodPhase
	Reason    string
	HostIP    string
	PodIP     string
	StartTime *time.Time
}

type PodTemplateSpec struct {
	ObjectMeta
	Spec PodSpec
}

type Volume struct {
	Name string
	VolumeSource
}

type VolumeSource struct {
	HostPath *HostPathVolumeSource
	NFS      *NFSVolumeSource
}

type HostPathVolumeSource struct {
	Path string
}

type NFSVolumeSource struct {
	Server string
	Path   string
}

// DeploymentCli has methods to work with Deployment resources in k8s-assist.
// How to:  deploymentCli, err := k8sassist.NewDeployments(nameSpace)
//          _, err := deploymentCli.Update(&deployment)
type DeploymentCli interface {
	Create(*Deployment) (*Deployment, error)
	Update(*Deployment) (*Deployment, error)
	UpdateStatus(*Deployment) (*Deployment, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*Deployment, error)
	//List(opts v1.ListOptions) (*DeploymentList, error)
	List() (*DeploymentList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1beta1.Deployment, err error)
}

// NamespaceCli Interface has methods to work with Namespace resources.
// How to:  namespaceCli, err := k8sassist.NewNamespaces()
//          nl, err := namespaceCli.List()
type NamespaceCli interface {
	Create(*Namespace) (*Namespace, error)
	Update(*Namespace) (*Namespace, error)
	UpdateStatus(*Namespace) (*Namespace, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*Namespace, error)
	List() (*NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

// ServiceCli interface has methods to work with Service resources in k8s-assist.
// How to:  serviceCli, err := k8sassist.NewServices(nameSpace)
//          service, err := serviceCli.Get(serviceName)
type ServiceCli interface {
	Create(*Service) (*Service, error)
	Update(*Service) (*Service, error)
	UpdateStatus(*Service) (*Service, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*Service, error)
	List() (*ServiceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Service, err error)
}

// The ScaleCli interface has methods on Scale resources in k8s-assist.
type ScaleCli interface {
	Get(kind string, name string) (*Scale, error)
	Update(kind string, scale *Scale) (*Scale, error)
}
