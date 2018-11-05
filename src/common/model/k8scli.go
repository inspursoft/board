package model

import (
	"time"
)

//Path Type
type PatchType string

const (
	JSONPatchType           PatchType = "application/json-patch+json"
	MergePatchType          PatchType = "application/merge-patch+json"
	StrategicMergePatchType PatchType = "application/strategic-merge-patch+json"
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
type Quantity int64
type QuantityStr string

type ResourceList map[ResourceName]QuantityStr

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
	Type               NodeConditionType
	Status             ConditionStatus
	LastHeartbeatTime  time.Time
	LastTransitionTime time.Time
	Reason             string
	Message            string
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
	Ports               []ServicePort
	Selector            map[string]string
	ClusterIP           string
	Type                string
	ExternalIPs         []string
	SessionAffinityFlag int
	SessionAffinityTime int
	ExternalName        string
}

type ServiceList struct {
	Items []Service
}

//type ScaleState struct {
//	Replicas       int32
//	Selector       map[string]string
//	TargetSelector string
//}

// represents a scaling request for a resource.
//type Scale struct {
//	Name      string
//	Namespace string
//	Labels    map[string]string
//	Replicas  int32
//	Status    ScaleState
//}

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
	InitContainers []K8sContainer
	Containers     []K8sContainer
	NodeSelector   map[string]string
	NodeName       string
	HostNetwork    bool
	Affinity       K8sAffinity
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

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +optional
	Limits ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +optional
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`
}

// A single application container that you want to run within a pod.
type K8sContainer struct {
	Name       string
	Image      string
	Command    []string
	Args       []string
	WorkingDir string
	Ports      []ContainerPort
	Env        []EnvVar
	// Compute Resources required by this container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	// Pod volumes to mount into the container's filesystem.
	VolumeMounts []VolumeMount
}

// Protocol defines network protocols supported for things like container ports.
type Protocol string

const (
	ProtocolTCP Protocol = "TCP"
	ProtocolUDP Protocol = "UDP"
)

// ContainerPort represents a network port in a single container.
type ContainerPort struct {
	Name          string
	HostPort      int32
	ContainerPort int32
	Protocol      Protocol
	HostIP        string
}

type EnvVar struct {
	Name  string
	Value string
}

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	Name      string
	MountPath string
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

// A new trying, keep ReplicaSet the same with k8s client
// ReplicaSet represents the configuration of a ReplicaSet.
type ReplicaSet struct {
	//unversioned.TypeMeta `json:",inline"`

	// If the Labels of a ReplicaSet are empty, they are defaulted to
	// be the same as the Pod(s) that the ReplicaSet manages.
	// Standard object's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	// +optional
	ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the specification of the desired behavior of the ReplicaSet.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status
	// +optional
	Spec ReplicaSetSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status is the most recently observed status of the ReplicaSet.
	// This data may be out of date by some window of time.
	// Populated by the system.
	// Read-only.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status
	// +optional
	Status ReplicaSetStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ReplicaSetList is a collection of ReplicaSets.
type ReplicaSetList struct {
	//unversioned.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#types-kinds
	// +optional
	//unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of ReplicaSets.
	// More info: http://kubernetes.io/docs/user-guide/replication-controller
	Items []ReplicaSet `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ReplicaSetSpec is the specification of a ReplicaSet.
type ReplicaSetSpec struct {
	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	// More info: http://kubernetes.io/docs/user-guide/replication-controller#what-is-a-replication-controller
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	// +optional
	MinReadySeconds int32 `json:"minReadySeconds,omitempty" protobuf:"varint,4,opt,name=minReadySeconds"`

	// Selector is a label query over pods that should match the replica count.
	// If the selector is empty, it is defaulted to the labels present on the pod template.
	// Label keys and values that must match in order to be controlled by this replica set.
	// More info: http://kubernetes.io/docs/user-guide/labels#label-selectors
	// +optional
	Selector *LabelSelector `json:"selector,omitempty" protobuf:"bytes,2,opt,name=selector"`

	// Template is the object that describes the pod that will be created if
	// insufficient replicas are detected.
	// More info: http://kubernetes.io/docs/user-guide/replication-controller#pod-template
	// +optional
	Template PodTemplateSpec `json:"template,omitempty" protobuf:"bytes,3,opt,name=template"`
}

// ReplicaSetStatus represents the current status of a ReplicaSet.
type ReplicaSetStatus struct {
	// Replicas is the most recently oberved number of replicas.
	// More info: http://kubernetes.io/docs/user-guide/replication-controller#what-is-a-replication-controller
	Replicas int32 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`

	// The number of pods that have labels matching the labels of the pod template of the replicaset.
	// +optional
	FullyLabeledReplicas int32 `json:"fullyLabeledReplicas,omitempty" protobuf:"varint,2,opt,name=fullyLabeledReplicas"`

	// The number of ready replicas for this replica set.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,4,opt,name=readyReplicas"`

	// The number of available replicas (ready for at least minReadySeconds) for this replica set.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty" protobuf:"varint,5,opt,name=availableReplicas"`

	// ObservedGeneration reflects the generation of the most recently observed ReplicaSet.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`

	// Represents the latest available observations of a replica set's current state.
	// +optional
	Conditions []ReplicaSetCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`
}

type ReplicaSetConditionType string

// These are valid conditions of a replica set.
const (
	// ReplicaSetReplicaFailure is added in a replica set when one of its pods fails to be created
	// due to insufficient quota, limit ranges, pod security policy, node selectors, etc. or deleted
	// due to kubelet being down or finalizers are failing.
	ReplicaSetReplicaFailure ReplicaSetConditionType = "ReplicaFailure"
)

// ReplicaSetCondition describes the state of a replica set at a certain point.
type ReplicaSetCondition struct {
	// Type of replica set condition.
	Type ReplicaSetConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ReplicaSetConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/kubernetes/pkg/api/v1.ConditionStatus"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// label selector matches no objects.
type LabelSelector struct {
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	// +optional
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,2,rep,name=matchExpressions"`
}

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	// A selector to restrict the list of returned objects by their labels.
	// Defaults to everything.
	// +optional
	LabelSelector string `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// A selector to restrict the list of returned objects by their fields.
	// Defaults to everything.
	// +optional
	FieldSelector string `json:"fieldSelector,omitempty" protobuf:"bytes,2,opt,name=fieldSelector"`
	// Watch for changes to the described resources and return them as a stream of
	// add, update, and remove notifications. Specify resourceVersion.
	// +optional
	Watch bool `json:"watch,omitempty" protobuf:"varint,3,opt,name=watch"`
	// When specified with a watch call, shows changes that occur after that particular version of a resource.
	// Defaults to changes from the beginning of history.
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,4,opt,name=resourceVersion"`
	// Timeout for the list/watch call.
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty" protobuf:"varint,5,opt,name=timeoutSeconds"`
}

// describes the attributes of a scale subresource
type ScaleSpec struct {
	// desired number of instances for the scaled object.
	// +optional
	Replicas int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
}

// represents the current status of a scale subresource.
type ScaleStatusK8s struct {
	// actual number of observed instances of the scaled object.
	Replicas int32 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`

	// label query over pods that should match the replicas count. More info: http://kubernetes.io/docs/user-guide/labels#label-selectors
	// +optional
	Selector map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`

	// label selector for pods that should match the replicas count. This is a serializated
	// version of both map-based and more expressive set-based selectors. This is done to
	// avoid introspection in the clients. The string will be in the same format as the
	// query-param syntax. If the target type only supports map-based selectors, both this
	// field and map-based selector field are populated.
	// More info: http://kubernetes.io/docs/user-guide/labels#label-selectors
	// +optional
	TargetSelector string `json:"targetSelector,omitempty" protobuf:"bytes,3,opt,name=targetSelector"`
}

// +genclient=true
// +noMethods=true

// represents a scaling request for a resource.
type Scale struct {
	// +optional
	ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// defines the behavior of the scale. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status.
	// +optional
	Spec ScaleSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// current status of the scale. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status. Read-only.
	// +optional
	Status ScaleStatusK8s `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// CrossVersionObjectReference contains enough information to let you identify the referred resource.
type CrossVersionObjectReference struct {
	// Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds"
	Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
	// API version of the referent
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,3,opt,name=apiVersion"`
}

type HorizontalPodAutoscalerSpec struct {
	// reference to scaled resource; horizontal pod autoscaler will learn the current resource consumption
	// and will set the desired number of pods by using its Scale subresource.
	ScaleTargetRef CrossVersionObjectReference `json:"scaleTargetRef" protobuf:"bytes,1,opt,name=scaleTargetRef"`
	// lower limit for the number of pods that can be set by the autoscaler, default 1.
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty" protobuf:"varint,2,opt,name=minReplicas"`
	// upper limit for the number of pods that can be set by the autoscaler; cannot be smaller than MinReplicas.
	MaxReplicas int32 `json:"maxReplicas" protobuf:"varint,3,opt,name=maxReplicas"`
	// target average CPU utilization (represented as a percentage of requested CPU) over all the pods;
	// if not specified the default autoscaling policy will be used.
	// +optional
	TargetCPUUtilizationPercentage *int32 `json:"targetCPUUtilizationPercentage,omitempty" protobuf:"varint,4,opt,name=targetCPUUtilizationPercentage"`
}

type HorizontalPodAutoscalerStatus struct {
	// most recent generation observed by this autoscaler.
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// last time the HorizontalPodAutoscaler scaled the number of pods;
	// used by the autoscaler to control how often the number of pods is changed.
	// +optional
	LastScaleTime *time.Time `json:"lastScaleTime,omitempty" protobuf:"bytes,2,opt,name=lastScaleTime"`

	// current number of replicas of pods managed by this autoscaler.
	CurrentReplicas int32 `json:"currentReplicas" protobuf:"varint,3,opt,name=currentReplicas"`

	// desired number of replicas of pods managed by this autoscaler.
	DesiredReplicas int32 `json:"desiredReplicas" protobuf:"varint,4,opt,name=desiredReplicas"`

	// current average CPU utilization over all pods, represented as a percentage of requested CPU,
	// e.g. 70 means that an average pod is using now 70% of its requested CPU.
	// +optional
	CurrentCPUUtilizationPercentage *int32 `json:"currentCPUUtilizationPercentage,omitempty" protobuf:"varint,5,opt,name=currentCPUUtilizationPercentage"`
}

type AutoScale struct {
	//metav1.TypeMeta `json:",inline"`
	// +optional
	ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// behaviour of autoscaler. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status.
	// +optional
	Spec HorizontalPodAutoscalerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// current information about the autoscaler.
	// +optional
	Status HorizontalPodAutoscalerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// list of horizontal pod autoscaler objects.
type AutoScaleList struct {
	// list of horizontal pod autoscaler objects.
	Items []AutoScale `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// Affinity is a group of affinity scheduling rules.
type K8sAffinity struct {
	// Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
	// +optional
	PodAffinity []PodAffinityTerm
	// Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)).
	// +optional
	PodAntiAffinity []PodAffinityTerm
}

type PodAffinityTerm struct {
	// A label query over a set of resources, in this case pods.
	// +optional
	LabelSelector LabelSelector `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// namespaces specifies which namespaces the labelSelector applies to (matches against);
	// null or empty list means "this pod's namespace"
	Namespaces []string `json:"namespaces,omitempty" protobuf:"bytes,2,rep,name=namespaces"`
	// This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching
	// the labelSelector in the specified namespaces, where co-located is defined as running on a node
	// whose value of the label with key topologyKey matches that of any node on which any of the
	// selected pods is running.
	// For PreferredDuringScheduling pod anti-affinity, empty topologyKey is interpreted as "all topologies"
	// ("all topologies" here means all the topologyKeys indicated by scheduler command-line argument --failure-domains);
	// for affinity and for RequiredDuringScheduling pod anti-affinity, empty topologyKey is not allowed.
	// +optional
	TopologyKey string `json:"topologyKey,omitempty" protobuf:"bytes,3,opt,name=topologyKey"`
}

type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	// +patchMergeKey=key
	// +patchStrategy=merge
	Key string `json:"key" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,1,opt,name=key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator string `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}
