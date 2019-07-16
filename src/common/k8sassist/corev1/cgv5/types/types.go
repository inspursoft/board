package types

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	betchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	kubernetes "k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	config "k8s.io/client-go/rest"
)

//define Deployment type
type DeploymentList = appsv1beta2.DeploymentList
type Deployment = appsv1beta2.Deployment
type TypeMeta = metav1.TypeMeta
type ObjectMeta = metav1.ObjectMeta
type DeploymentSpec = appsv1beta2.DeploymentSpec
type DeploymentStatus = appsv1beta2.DeploymentStatus
type LabelSelector = metav1.LabelSelector
type PodTemplateSpec = v1.PodTemplateSpec
type PodSpec = v1.PodSpec
type Container = v1.Container
type ContainerPort = v1.ContainerPort

// These are valid conditions of a job.
const (
	JobComplete JobConditionType = "Complete"
	JobFailed   JobConditionType = "Failed"
)

//define Job type
type JobList = betchv1.JobList
type Job = betchv1.Job
type JobSpec = betchv1.JobSpec
type JobStatus = betchv1.JobStatus
type JobCondition = betchv1.JobCondition
type JobConditionType = betchv1.JobConditionType
type ConditionStatus = v1.ConditionStatus

//define Service type
type ServiceList = v1.ServiceList
type Service = v1.Service
type ServiceSpec = v1.ServiceSpec
type ServicePort = v1.ServicePort
type ServiceType = v1.ServiceType
type Protocol = v1.Protocol
type IntOrString = intstr.IntOrString
type Type = intstr.Type
type SessionAffinity = v1.ServiceAffinity
type SessionAffinityConfig = v1.SessionAffinityConfig
type ClientIPConfig = v1.ClientIPConfig

//define namespace type
type Namespace = v1.Namespace
type NamespaceList = v1.NamespaceList

//define Options
type GetOptions = metav1.GetOptions
type ListOptions = metav1.ListOptions
type DeleteOptions = metav1.DeleteOptions

//define config
type Config = config.Config
type Clientset = kubernetes.Clientset
type NamespacePhase = v1.NamespacePhase

//define api
type NamespaceInterface = corev1.NamespaceInterface

//define time
type Time = metav1.Time

//define var
const (
	DeletePropagationForeground = metav1.DeletePropagationForeground
	NamespaceActive             = v1.NamespaceActive
	NamespaceTerminating        = v1.NamespaceTerminating
	SessionAffinityClientIP     = v1.ServiceAffinityClientIP
	ServiceAffinityNone         = v1.ServiceAffinityNone
)

const (
	Int    Type = iota // The IntOrString holds an int.
	String             // The IntOrString holds a string.
)

const (
	serviceAPIVersion    = "v1"
	serviceKind          = "Service"
	nodePort             = "NodePort"
	deploymentAPIVersion = "apps/v1beta2"
	deploymentKind       = "Deployment"
	namespaceKind        = "Namespace"
	namespaceAPIVersion  = "v1"
	podKind              = "Pod"
	podVersion           = "v1"
	maxPort              = 32765
	minPort              = 30000
)
