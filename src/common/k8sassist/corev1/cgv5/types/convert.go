package types

import (
	"time"

	"git/inspursoft/board/src/common/model"

	"strconv"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	autoscalev1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// generate k8s objectmeta from model objectmeta
func ToK8sObjectMeta(meta model.ObjectMeta) metav1.ObjectMeta {
	var deleteTime *metav1.Time
	if meta.DeletionTimestamp != nil {
		t := metav1.NewTime(*meta.DeletionTimestamp)
		deleteTime = &t
	}
	return metav1.ObjectMeta{
		Name:              meta.Name,
		Namespace:         meta.Namespace,
		CreationTimestamp: metav1.NewTime(meta.CreationTimestamp),
		DeletionTimestamp: deleteTime,
		Labels:            meta.Labels,
	}
}

// generate k8s deployment from model deployment
func ToK8sDeployment(deployment *model.Deployment) *appsv1beta2.Deployment {
	if deployment == nil {
		return nil
	}
	var templ v1.PodTemplateSpec
	if t := ToK8sPodTemplateSpec(&deployment.Spec.Template); t != nil {
		templ = *t
	}
	rep := deployment.Spec.Replicas
	return &appsv1beta2.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1beta2",
		},
		ObjectMeta: ToK8sObjectMeta(deployment.ObjectMeta),
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &rep,
			Selector: &metav1.LabelSelector{
				MatchLabels: deployment.Spec.Selector,
			},
			Template: templ,
			Paused:   deployment.Spec.Paused,
		},
		Status: appsv1beta2.DeploymentStatus{
			Replicas:            deployment.Status.Replicas,
			UpdatedReplicas:     deployment.Status.UpdatedReplicas,
			UnavailableReplicas: deployment.Status.UnavailableReplicas,
			AvailableReplicas:   deployment.Status.AvailableReplicas,
		},
	}
}

// generate k8s pod template spec from model pod template spec
func ToK8sPodTemplateSpec(template *model.PodTemplateSpec) *v1.PodTemplateSpec {
	if template == nil {
		return nil
	}
	var spec v1.PodSpec
	if s := ToK8sPodSpec(&template.Spec); s != nil {
		spec = *s
	}
	return &v1.PodTemplateSpec{
		ObjectMeta: ToK8sObjectMeta(template.ObjectMeta),
		Spec:       spec,
	}
}

// generate k8s replicaset from model replicaset
func ToK8sReplicaSet(rs *model.ReplicaSet) *appsv1beta2.ReplicaSet {
	if rs == nil {
		return nil
	}
	var spec appsv1beta2.ReplicaSetSpec
	if s := ToK8sReplicaSetSpec(&rs.Spec); s != nil {
		spec = *s
	}
	conds := make([]appsv1beta2.ReplicaSetCondition, len(rs.Status.Conditions))
	for i := range rs.Status.Conditions {
		conds[i] = appsv1beta2.ReplicaSetCondition{
			Type:               appsv1beta2.ReplicaSetConditionType(string(rs.Status.Conditions[i].Type)),
			Status:             v1.ConditionStatus(string(rs.Status.Conditions[i].Status)),
			LastTransitionTime: metav1.NewTime(rs.Status.Conditions[i].LastTransitionTime),
			Reason:             rs.Status.Conditions[i].Reason,
			Message:            rs.Status.Conditions[i].Message,
		}
	}
	return &appsv1beta2.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "apps/v1beta2",
		},
		ObjectMeta: ToK8sObjectMeta(rs.ObjectMeta),
		Spec:       spec,
		Status: appsv1beta2.ReplicaSetStatus{
			Replicas:             rs.Status.Replicas,
			FullyLabeledReplicas: rs.Status.FullyLabeledReplicas,
			ReadyReplicas:        rs.Status.ReadyReplicas,
			AvailableReplicas:    rs.Status.AvailableReplicas,
			ObservedGeneration:   rs.Status.ObservedGeneration,
			Conditions:           conds,
		},
	}
}

func ToK8sReplicaSetSpec(spec *model.ReplicaSetSpec) *appsv1beta2.ReplicaSetSpec {
	if spec == nil {
		return nil
	}
	var selector *metav1.LabelSelector
	if spec.Selector != nil {
		selector = &metav1.LabelSelector{
			MatchLabels: spec.Selector.MatchLabels,
		}
	}
	var template v1.PodTemplateSpec
	if t := ToK8sPodTemplateSpec(&spec.Template); t != nil {
		template = *t
	}
	return &appsv1beta2.ReplicaSetSpec{
		Replicas:        spec.Replicas,
		MinReadySeconds: spec.MinReadySeconds,
		Selector:        selector,
		Template:        template,
	}
}

// generate k8s pod from model pod
func ToK8sPod(pod *model.Pod) *v1.Pod {
	if pod == nil {
		return nil
	}
	var spec v1.PodSpec
	if s := ToK8sPodSpec(&pod.Spec); s != nil {
		spec = *s
	}
	var starttime *metav1.Time
	if pod.Status.StartTime != nil {
		t := metav1.NewTime(*pod.Status.StartTime)
		starttime = &t
	}
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(pod.ObjectMeta),
		Spec:       spec,
		Status: v1.PodStatus{
			Phase:     v1.PodPhase(string(pod.Status.Phase)),
			Reason:    pod.Status.Reason,
			HostIP:    pod.Status.HostIP,
			PodIP:     pod.Status.PodIP,
			StartTime: starttime,
		},
	}
}

// generate k8s pod spec from model pod spec
func ToK8sPodSpec(spec *model.PodSpec) *v1.PodSpec {
	if spec == nil {
		return nil
	}
	var volumes []v1.Volume
	for i := range spec.Volumes {
		if v := ToK8sVolume(&spec.Volumes[i]); v != nil {
			volumes = append(volumes, *v)
		}
	}

	var initContainers []v1.Container
	for i := range spec.InitContainers {
		if c := ToK8sContainer(&spec.InitContainers[i]); c != nil {
			initContainers = append(initContainers, *c)
		}
	}

	var containers []v1.Container
	for i := range spec.Containers {
		if c := ToK8sContainer(&spec.Containers[i]); c != nil {
			containers = append(containers, *c)
		}
	}
	return &v1.PodSpec{
		Volumes:        volumes,
		InitContainers: initContainers,
		Containers:     containers,
		NodeSelector:   spec.NodeSelector,
		NodeName:       spec.NodeName,
		HostNetwork:    spec.HostNetwork,
	}
}

func ToK8sVolume(volume *model.Volume) *v1.Volume {
	if volume == nil {
		return nil
	}
	var volumeSource v1.VolumeSource
	if vs := ToK8sVolumeSource(&volume.VolumeSource); vs != nil {
		volumeSource = *vs
	}
	return &v1.Volume{
		Name:         volume.Name,
		VolumeSource: volumeSource,
	}
}

func ToK8sVolumeSource(volumeSource *model.VolumeSource) *v1.VolumeSource {
	if volumeSource == nil {
		return nil
	}
	var hp *v1.HostPathVolumeSource
	if volumeSource.HostPath != nil {
		hp = &v1.HostPathVolumeSource{
			Path: volumeSource.HostPath.Path,
		}
	}
	var nfs *v1.NFSVolumeSource
	if volumeSource.NFS != nil {
		nfs = &v1.NFSVolumeSource{
			Server: volumeSource.NFS.Server,
			Path:   volumeSource.NFS.Path,
		}
	}
	return &v1.VolumeSource{
		HostPath: hp,
		NFS:      nfs,
	}
}

func ToK8sContainer(container *model.K8sContainer) *v1.Container {
	if container == nil {
		return nil
	}
	var ports []v1.ContainerPort
	for i := range container.Ports {
		ports = append(ports, ToK8sContainerPort(container.Ports[i]))
	}
	var envs []v1.EnvVar
	for i := range container.Env {
		envs = append(envs, ToK8sEnvVar(container.Env[i]))
	}

	var mounts []v1.VolumeMount
	for i := range container.VolumeMounts {
		mounts = append(mounts, ToK8sVolumeMount(container.VolumeMounts[i]))
	}

	var resources v1.ResourceRequirements
	resources.Requests = make(v1.ResourceList)
	resources.Limits = make(v1.ResourceList)
	if v, ok := container.Resources.Requests["cpu"]; ok {
		resources.Requests["cpu"] = resource.MustParse(string(v))
	}
	if v, ok := container.Resources.Requests["memory"]; ok {
		resources.Requests["memory"] = resource.MustParse(string(v))
	}
	if v, ok := container.Resources.Limits["cpu"]; ok {
		resources.Limits["cpu"] = resource.MustParse(string(v))
	}
	if v, ok := container.Resources.Limits["memory"]; ok {
		resources.Limits["memory"] = resource.MustParse(string(v))
	}
	return &v1.Container{
		Name:         container.Name,
		Image:        container.Image,
		Command:      container.Command,
		Args:         container.Args,
		WorkingDir:   container.WorkingDir,
		Ports:        ports,
		Env:          envs,
		Resources:    resources,
		VolumeMounts: mounts,
	}
}

func ToK8sContainerPort(containerPort model.ContainerPort) v1.ContainerPort {
	return v1.ContainerPort{
		Name:          containerPort.Name,
		HostPort:      containerPort.HostPort,
		ContainerPort: containerPort.ContainerPort,
		Protocol:      v1.Protocol(string(containerPort.Protocol)),
		HostIP:        containerPort.HostIP,
	}
}

func ToK8sEnvVar(env model.EnvVar) v1.EnvVar {
	return v1.EnvVar{
		Name:  env.Name,
		Value: env.Value,
	}
}

func ToK8sVolumeMount(mount model.VolumeMount) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      mount.Name,
		MountPath: mount.MountPath,
	}
}

func FromK8sObjectMeta(meta metav1.ObjectMeta) model.ObjectMeta {
	var deleteTime *time.Time
	if meta.DeletionTimestamp != nil {
		deleteTime = &meta.DeletionTimestamp.Time
	}
	return model.ObjectMeta{
		Name:              meta.Name,
		Namespace:         meta.Namespace,
		CreationTimestamp: meta.CreationTimestamp.Time,
		DeletionTimestamp: deleteTime,
		Labels:            meta.Labels,
	}
}

// generate model deployment list from k8s deployment list
func FromK8sDeploymentList(deploymentList *appsv1beta2.DeploymentList) *model.DeploymentList {
	if deploymentList == nil {
		return nil
	}
	items := make([]model.Deployment, 0)
	for i := range deploymentList.Items {
		dep := FromK8sDeployment(&deploymentList.Items[i])
		items = append(items, *dep)
	}
	return &model.DeploymentList{
		Items: items,
	}
}

// generate model deployment from k8s deployment
func FromK8sDeployment(deployment *appsv1beta2.Deployment) *model.Deployment {
	if deployment == nil {
		return nil
	}
	var spec model.DeploymentSpec
	if s := FromK8sDeploymentSpec(&deployment.Spec); s != nil {
		spec = *s
	}
	return &model.Deployment{
		ObjectMeta: FromK8sObjectMeta(deployment.ObjectMeta),
		Spec:       spec,
		Status: model.DeploymentStatus{
			Replicas:            deployment.Status.Replicas,
			UpdatedReplicas:     deployment.Status.UpdatedReplicas,
			AvailableReplicas:   deployment.Status.AvailableReplicas,
			UnavailableReplicas: deployment.Status.UnavailableReplicas,
		},
	}
}

func FromK8sDeploymentSpec(spec *appsv1beta2.DeploymentSpec) *model.DeploymentSpec {
	if spec == nil {
		return nil
	}
	var rep int32
	if spec.Replicas != nil {
		rep = *spec.Replicas
	}
	var template model.PodTemplateSpec
	if t := FromK8sPodTemplateSpec(&spec.Template); t != nil {
		template = *t
	}
	return &model.DeploymentSpec{
		Replicas: rep,
		Selector: FromK8sSelector(spec.Selector),
		Template: template,
		Paused:   spec.Paused,
	}
}

func FromK8sSelector(selector *metav1.LabelSelector) map[string]string {
	if selector == nil {
		return nil
	}
	return selector.MatchLabels
}

func FromK8sPodTemplateSpec(template *v1.PodTemplateSpec) *model.PodTemplateSpec {
	if template == nil {
		return nil
	}
	var spec model.PodSpec
	if s := FromK8sPodSpec(&template.Spec); s != nil {
		spec = *s
	}
	return &model.PodTemplateSpec{
		ObjectMeta: FromK8sObjectMeta(template.ObjectMeta),
		Spec:       spec,
	}
}

// generate model replicaset list from k8s replicaset list
func FromK8sReplicaSetList(list *appsv1beta2.ReplicaSetList) *model.ReplicaSetList {
	if list == nil {
		return nil
	}
	items := make([]model.ReplicaSet, len(list.Items))
	for i := range list.Items {
		if r := FromK8sReplicaSet(&list.Items[i]); r != nil {
			items[i] = *r
		}
	}
	return &model.ReplicaSetList{
		Items: items,
	}
}

// generate model replicaset from k8s replicaset
func FromK8sReplicaSet(rs *appsv1beta2.ReplicaSet) *model.ReplicaSet {
	if rs == nil {
		return nil
	}
	conds := make([]model.ReplicaSetCondition, len(rs.Status.Conditions))
	for i := range rs.Status.Conditions {
		conds[i] = model.ReplicaSetCondition{
			Type:               model.ReplicaSetConditionType(string(rs.Status.Conditions[i].Type)),
			Status:             model.ConditionStatus(string(rs.Status.Conditions[i].Status)),
			LastTransitionTime: rs.Status.Conditions[i].LastTransitionTime.Time,
			Reason:             rs.Status.Conditions[i].Reason,
			Message:            rs.Status.Conditions[i].Message,
		}
	}
	var spec model.ReplicaSetSpec
	if s := FromK8sReplicSetSpec(&rs.Spec); s != nil {
		spec = *s
	}
	return &model.ReplicaSet{
		ObjectMeta: FromK8sObjectMeta(rs.ObjectMeta),
		Spec:       spec,
		Status: model.ReplicaSetStatus{
			Replicas:             rs.Status.Replicas,
			FullyLabeledReplicas: rs.Status.FullyLabeledReplicas,
			ReadyReplicas:        rs.Status.ReadyReplicas,
			AvailableReplicas:    rs.Status.AvailableReplicas,
			ObservedGeneration:   rs.Status.ObservedGeneration,
			Conditions:           conds,
		},
	}
}

func FromK8sReplicSetSpec(spec *appsv1beta2.ReplicaSetSpec) *model.ReplicaSetSpec {
	if spec == nil {
		return nil
	}
	var selector *model.LabelSelector
	if spec.Selector != nil {
		selector = &model.LabelSelector{
			MatchLabels: spec.Selector.MatchLabels,
		}
	}
	var template model.PodTemplateSpec
	if t := FromK8sPodTemplateSpec(&spec.Template); t != nil {
		template = *t
	}
	return &model.ReplicaSetSpec{
		Replicas:        spec.Replicas,
		MinReadySeconds: spec.MinReadySeconds,
		Selector:        selector,
		Template:        template,
	}
}

func FromK8sPodList(podList *v1.PodList) *model.PodList {
	if podList == nil {
		return nil
	}
	items := make([]model.Pod, 0)
	for i := range podList.Items {
		if pod := FromK8sPod(&podList.Items[i]); pod != nil {
			items = append(items, *pod)
		}
	}
	return &model.PodList{
		Items: items,
	}
}

// generate model pod from k8s pod
func FromK8sPod(pod *v1.Pod) *model.Pod {
	if pod == nil {
		return nil
	}
	var spec model.PodSpec
	if s := FromK8sPodSpec(&pod.Spec); s != nil {
		spec = *s
	}
	var starttime *time.Time
	if t := pod.Status.StartTime; t != nil {
		starttime = &t.Time
	}
	return &model.Pod{
		ObjectMeta: FromK8sObjectMeta(pod.ObjectMeta),
		Spec:       spec,
		Status: model.PodStatus{
			Phase:     model.PodPhase(string(pod.Status.Phase)),
			Reason:    pod.Status.Reason,
			HostIP:    pod.Status.HostIP,
			PodIP:     pod.Status.PodIP,
			StartTime: starttime,
		},
	}
}

func FromK8sPodSpec(spec *v1.PodSpec) *model.PodSpec {
	if spec == nil {
		return nil
	}
	var volumes []model.Volume = nil
	for i := range spec.Volumes {
		if v := FromK8sVolume(&spec.Volumes[i]); v != nil {
			volumes = append(volumes, *v)
		}
	}

	var initContainers []model.K8sContainer
	for i := range spec.InitContainers {
		if c := FromK8sContainer(&spec.InitContainers[i]); c != nil {
			initContainers = append(initContainers, *c)
		}
	}

	var containers []model.K8sContainer
	for i := range spec.Containers {
		if c := FromK8sContainer(&spec.Containers[i]); c != nil {
			containers = append(containers, *c)
		}
	}
	return &model.PodSpec{
		Volumes:        volumes,
		InitContainers: initContainers,
		Containers:     containers,
		NodeSelector:   spec.NodeSelector,
		NodeName:       spec.NodeName,
		HostNetwork:    spec.HostNetwork,
	}
}

func FromK8sVolume(volume *v1.Volume) *model.Volume {
	if volume == nil {
		return nil
	}
	return &model.Volume{
		Name:         volume.Name,
		VolumeSource: FromK8sVolumeSource(volume.VolumeSource),
	}
}

func FromK8sVolumeSource(volumeSource v1.VolumeSource) model.VolumeSource {
	var hp *model.HostPathVolumeSource
	if volumeSource.HostPath != nil {
		hp = &model.HostPathVolumeSource{
			Path: volumeSource.HostPath.Path,
		}
	}
	var nfs *model.NFSVolumeSource
	if volumeSource.NFS != nil {
		nfs = &model.NFSVolumeSource{
			Server: volumeSource.NFS.Server,
			Path:   volumeSource.NFS.Path,
		}
	}
	return model.VolumeSource{
		HostPath: hp,
		NFS:      nfs,
	}
}

func FromK8sContainer(container *v1.Container) *model.K8sContainer {
	if container == nil {
		return nil
	}
	var ports []model.ContainerPort
	for i := range container.Ports {
		ports = append(ports, FromK8sContainerPort(container.Ports[i]))
	}
	var envs []model.EnvVar
	for i := range container.Env {
		envs = append(envs, FromK8sEnvVar(container.Env[i]))
	}

	var mounts []model.VolumeMount
	for i := range container.VolumeMounts {
		mounts = append(mounts, FromK8sVolumeMount(container.VolumeMounts[i]))
	}

	var resources model.ResourceRequirements
	resources.Requests = make(model.ResourceList)
	resources.Limits = make(model.ResourceList)
	if v, ok := container.Resources.Requests["cpu"]; ok {
		resources.Requests["cpu"] = model.QuantityStr(v.String())
	}
	if v, ok := container.Resources.Requests["memory"]; ok {
		resources.Requests["memory"] = model.QuantityStr(v.String())
	}
	if v, ok := container.Resources.Limits["cpu"]; ok {
		resources.Limits["cpu"] = model.QuantityStr(v.String())
	}
	if v, ok := container.Resources.Limits["memory"]; ok {
		resources.Limits["memory"] = model.QuantityStr(v.String())
	}

	return &model.K8sContainer{
		Name:         container.Name,
		Image:        container.Image,
		Command:      container.Command,
		Args:         container.Args,
		WorkingDir:   container.WorkingDir,
		Ports:        ports,
		Env:          envs,
		Resources:    resources,
		VolumeMounts: mounts,
	}
}

func FromK8sContainerPort(containerPort v1.ContainerPort) model.ContainerPort {
	return model.ContainerPort{
		Name:          containerPort.Name,
		HostPort:      containerPort.HostPort,
		ContainerPort: containerPort.ContainerPort,
		Protocol:      model.Protocol(string(containerPort.Protocol)),
		HostIP:        containerPort.HostIP,
	}
}

func FromK8sEnvVar(env v1.EnvVar) model.EnvVar {
	return model.EnvVar{
		Name:  env.Name,
		Value: env.Value,
	}
}

func FromK8sVolumeMount(mount v1.VolumeMount) model.VolumeMount {
	return model.VolumeMount{
		Name:      mount.Name,
		MountPath: mount.MountPath,
	}
}

//namespace data convert
func ToK8sNamespace(modelNamespace *model.Namespace) *Namespace {
	ns := &Namespace{
		TypeMeta: TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(modelNamespace.ObjectMeta),
	}

	return ns
}

func FromK8sNamespace(typesNamespace *Namespace) *model.Namespace {
	return &model.Namespace{
		ObjectMeta:     FromK8sObjectMeta(typesNamespace.ObjectMeta),
		NamespacePhase: string(typesNamespace.Status.Phase),
	}
}

func FromK8sNamespaceList(typesNamespaceList *NamespaceList) *model.NamespaceList {
	modelNamespaceList := &model.NamespaceList{
		Items: make([]model.Namespace, 0),
	}
	for _, ns := range typesNamespaceList.Items {
		modelNamespaceList.Items = append(modelNamespaceList.Items, *FromK8sNamespace(&ns))
	}

	return modelNamespaceList
}

//service data convert
func ToK8sService(modelService *model.Service) *Service {
	ports := make([]ServicePort, 0)
	for _, port := range modelService.Ports {
		ports = append(ports, ServicePort{
			Name:     port.Name,
			Protocol: Protocol(port.Protocol),
			Port:     port.Port,
			NodePort: port.NodePort,
			TargetPort: IntOrString{
				Type:   Int,
				IntVal: port.TargetPort,
			},
		})
	}
	return &Service{
		TypeMeta: TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(modelService.ObjectMeta),
		Spec: ServiceSpec{
			Ports:        ports,
			Selector:     modelService.Selector,
			ClusterIP:    modelService.ClusterIP,
			Type:         ServiceType(modelService.Type),
			ExternalIPs:  modelService.ExternalIPs,
			ExternalName: modelService.ExternalName,
		},
	}
}

func FromK8sService(typesService *Service) *model.Service {
	ports := make([]model.ServicePort, 0)
	for _, port := range typesService.Spec.Ports {
		ports = append(ports, model.ServicePort{
			Name:       port.Name,
			Protocol:   string(port.Protocol),
			Port:       port.Port,
			NodePort:   port.NodePort,
			TargetPort: port.TargetPort.IntVal,
		})
	}
	return &model.Service{
		ObjectMeta:   FromK8sObjectMeta(typesService.ObjectMeta),
		Ports:        ports,
		Selector:     typesService.Spec.Selector,
		ClusterIP:    typesService.Spec.ClusterIP,
		Type:         string(typesService.Spec.Type),
		ExternalIPs:  typesService.Spec.ExternalIPs,
		ExternalName: typesService.Spec.ExternalName,
	}
}

func FromK8sServiceList(typesServiceList *ServiceList) *model.ServiceList {
	modelServiceList := &model.ServiceList{
		Items: make([]model.Service, 0),
	}
	for _, s := range typesServiceList.Items {
		modelServiceList.Items = append(modelServiceList.Items, *FromK8sService(&s))
	}
	return modelServiceList
}

// generate k8s node status from model node status
func ToK8sNodeStatus(nodestatus model.NodeStatus) v1.NodeStatus {
	capacity := make(map[v1.ResourceName]resource.Quantity)
	for k, v := range nodestatus.Capacity {
		value, _ := strconv.Atoi(string(v))
		q := resource.NewQuantity(int64(value), resource.DecimalExponent)
		capacity[v1.ResourceName(k)] = *q

	}

	allocatable := make(map[v1.ResourceName]resource.Quantity)
	for k, v := range nodestatus.Allocatable {
		value, _ := strconv.Atoi(string(v))
		q := resource.NewQuantity(int64(value), resource.DecimalExponent)
		capacity[v1.ResourceName(k)] = *q

	}

	conditions := make([]v1.NodeCondition, 0)
	for _, v := range nodestatus.Conditions {
		conditions = append(conditions, v1.NodeCondition{
			Type:               v1.NodeConditionType(v.Type),
			Status:             v1.ConditionStatus(v.Status),
			LastHeartbeatTime:  metav1.NewTime(v.LastHeartbeatTime),
			LastTransitionTime: metav1.NewTime(v.LastTransitionTime),
			Reason:             v.Reason,
			Message:            v.Message,
		})

	}

	addresses := make([]v1.NodeAddress, 0)
	for _, v := range nodestatus.Addresses {
		addresses = append(addresses, v1.NodeAddress{
			Type:    v1.NodeAddressType(v.Type),
			Address: v.Address,
		})

	}

	return v1.NodeStatus{
		Capacity:    capacity,
		Allocatable: allocatable,
		Phase:       v1.NodePhase(nodestatus.Phase),
		Conditions:  conditions,
		Addresses:   addresses,
	}
}

func UpdateK8sNodeStatus(k8sNodeStatus *v1.NodeStatus, nodestatus *model.NodeStatus) {
	if nodestatus.Capacity == nil {
		k8sNodeStatus.Capacity = nil
	} else {
		if k8sNodeStatus.Capacity == nil {
			k8sNodeStatus.Capacity = v1.ResourceList(make(map[v1.ResourceName]resource.Quantity))
		}
		for k, v := range nodestatus.Capacity {
			value, _ := strconv.Atoi(string(v))
			q := resource.NewQuantity(int64(value), resource.DecimalExponent)
			k8sNodeStatus.Capacity[v1.ResourceName(k)] = *q
		}
	}

	if nodestatus.Allocatable == nil {
		k8sNodeStatus.Allocatable = nil
	} else {
		if k8sNodeStatus.Allocatable == nil {
			k8sNodeStatus.Allocatable = v1.ResourceList(make(map[v1.ResourceName]resource.Quantity))
		}
		for k, v := range nodestatus.Allocatable {
			value, _ := strconv.Atoi(string(v))
			q := resource.NewQuantity(int64(value), resource.DecimalExponent)
			k8sNodeStatus.Allocatable[v1.ResourceName(k)] = *q

		}
	}

	if nodestatus.Conditions == nil {
		k8sNodeStatus.Conditions = nil
	} else {
		conditions := make([]v1.NodeCondition, 0)
		for _, v := range nodestatus.Conditions {
			conditions = append(conditions, v1.NodeCondition{
				Type:               v1.NodeConditionType(v.Type),
				Status:             v1.ConditionStatus(v.Status),
				LastHeartbeatTime:  metav1.NewTime(v.LastHeartbeatTime),
				LastTransitionTime: metav1.NewTime(v.LastTransitionTime),
				Reason:             v.Reason,
				Message:            v.Message,
			})
		}
		k8sNodeStatus.Conditions = conditions
	}

	if nodestatus.Addresses == nil {
		k8sNodeStatus.Addresses = nil
	} else {
		addresses := make([]v1.NodeAddress, 0)
		for _, v := range nodestatus.Addresses {
			addresses = append(addresses, v1.NodeAddress{
				Type:    v1.NodeAddressType(v.Type),
				Address: v.Address,
			})
		}
		k8sNodeStatus.Addresses = addresses
	}
}

// generate k8s node from model node
func ToK8sNode(node *model.Node) *v1.Node {
	if node == nil {
		return nil
	}

	return &v1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(node.ObjectMeta),
		Spec: v1.NodeSpec{
			Unschedulable: node.Unschedulable,
		},
		Status: ToK8sNodeStatus(node.Status),
	}
}

// update k8s node using model node
func UpdateK8sNode(k8sNode *v1.Node, node *model.Node) {
	if node == nil || k8sNode == nil {
		return
	}
	// just update our attributes.
	k8sNode.Name = node.Name
	k8sNode.Namespace = node.Namespace
	k8sNode.CreationTimestamp = metav1.NewTime(node.CreationTimestamp)
	if node.DeletionTimestamp != nil {
		t := metav1.NewTime(*node.DeletionTimestamp)
		k8sNode.DeletionTimestamp = &t
	} else {
		k8sNode.DeletionTimestamp = nil
	}
	k8sNode.Labels = node.Labels

	k8sNode.Spec.Unschedulable = node.Unschedulable

	UpdateK8sNodeStatus(&k8sNode.Status, &node.Status)
}

// adapt model node.Status from k8s node.Status
func FromK8sNodeStatus(nodestatus v1.NodeStatus) model.NodeStatus {
	capacity := make(map[model.ResourceName]model.QuantityStr)
	for k, v := range nodestatus.Capacity {
		i, _ := v.AsInt64()
		capacity[model.ResourceName(k)] = model.QuantityStr(strconv.Itoa(int(i)))

	}

	allocatable := make(map[model.ResourceName]model.QuantityStr)
	for k, v := range nodestatus.Allocatable {
		i, _ := v.AsInt64()
		allocatable[model.ResourceName(k)] = model.QuantityStr(strconv.Itoa(int(i)))
	}

	conditions := make([]model.NodeCondition, 0)
	for _, v := range nodestatus.Conditions {
		conditions = append(conditions, model.NodeCondition{
			Type:               model.NodeConditionType(v.Type),
			Status:             model.ConditionStatus(v.Status),
			LastHeartbeatTime:  v.LastHeartbeatTime.Time,
			LastTransitionTime: v.LastTransitionTime.Time,
			Reason:             v.Reason,
			Message:            v.Message,
		})

	}

	addresses := make([]model.NodeAddress, 0)
	for _, v := range nodestatus.Addresses {
		addresses = append(addresses, model.NodeAddress{
			Type:    model.NodeAddressType(v.Type),
			Address: v.Address,
		})

	}

	return model.NodeStatus{
		Capacity:    capacity,
		Allocatable: allocatable,
		Phase:       model.NodePhase(nodestatus.Phase),
		Conditions:  conditions,
		Addresses:   addresses,
	}
}

// adapt model node from k8s node
func FromK8sNode(node *v1.Node) *model.Node {

	return &model.Node{
		ObjectMeta:    FromK8sObjectMeta(node.ObjectMeta),
		NodeIP:        node.ObjectMeta.Name,
		Unschedulable: node.Spec.Unschedulable,
		Status:        FromK8sNodeStatus(node.Status),
	}
}

// adapt model nodes from k8s nodes
func FromK8sNodeList(nodeList *v1.NodeList) *model.NodeList {
	if nodeList == nil {
		return nil
	}
	items := make([]model.Node, 0)
	for i := range nodeList.Items {
		if node := FromK8sNode(&nodeList.Items[i]); node != nil {
			items = append(items, *node)
		}
	}
	return &model.NodeList{
		Items: items,
	}
}

func FromK8sScale(scale *v1beta1.Scale) *model.Scale {
	return &model.Scale{
		ObjectMeta: FromK8sObjectMeta(scale.ObjectMeta),
		Spec:       model.ScaleSpec(scale.Spec),
		Status:     model.ScaleStatusK8s(scale.Status),
	}
}

func ToK8sScale(scale *model.Scale) *v1beta1.Scale {
	return &v1beta1.Scale{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Scale",
			APIVersion: "v1beta1",
		},
		ObjectMeta: ToK8sObjectMeta(scale.ObjectMeta),
		Spec:       v1beta1.ScaleSpec(scale.Spec),
		Status:     v1beta1.ScaleStatus(scale.Status),
	}
}

func GenerateDeploymentConfig(deployment *appsv1beta2.Deployment) *appsv1beta2.Deployment {
	containersConfig := []v1.Container{}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		containersConfig = append(containersConfig, v1.Container{
			Name:           container.Name,
			Image:          container.Image,
			Command:        container.Command,
			Args:           container.Args,
			WorkingDir:     container.WorkingDir,
			Ports:          container.Ports,
			EnvFrom:        container.EnvFrom,
			Env:            container.Env,
			Resources:      container.Resources,
			VolumeMounts:   container.VolumeMounts,
			LivenessProbe:  container.LivenessProbe,
			ReadinessProbe: container.ReadinessProbe,
		})
	}
	return &appsv1beta2.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       deploymentKind,
			APIVersion: deploymentAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    deployment.ObjectMeta.Labels,
			Name:      deployment.ObjectMeta.Name,
			Namespace: deployment.ObjectMeta.Namespace,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: deployment.Spec.Replicas,
			Selector: deployment.Spec.Selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deployment.Spec.Template.ObjectMeta.Labels,
					Name:   deployment.Spec.Template.ObjectMeta.Name,
				},
				Spec: v1.PodSpec{
					Affinity:           deployment.Spec.Template.Spec.Affinity,
					Volumes:            deployment.Spec.Template.Spec.Volumes,
					NodeSelector:       deployment.Spec.Template.Spec.NodeSelector,
					ServiceAccountName: deployment.Spec.Template.Spec.ServiceAccountName,
					ImagePullSecrets:   deployment.Spec.Template.Spec.ImagePullSecrets,
					InitContainers:     deployment.Spec.Template.Spec.InitContainers,
					Containers:         containersConfig,
				},
			},
		},
	}
}

func GenerateServiceConfig(service *v1.Service) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       serviceKind,
			APIVersion: serviceAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.ObjectMeta.Name,
			Namespace: service.ObjectMeta.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports:    service.Spec.Ports,
			Selector: service.Spec.Selector,
			Type:     service.Spec.Type,
		},
	}
}

func FromK8sAutoScale(autoscale *autoscalev1.HorizontalPodAutoscaler) *model.AutoScale {
	var lastTime *time.Time
	if autoscale.Status.LastScaleTime != nil {
		lastTime = &autoscale.Status.LastScaleTime.Time
	}
	return &model.AutoScale{
		ObjectMeta: FromK8sObjectMeta(autoscale.ObjectMeta),
		Spec: model.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: model.CrossVersionObjectReference{
				Kind:       autoscale.Spec.ScaleTargetRef.Kind,
				Name:       autoscale.Spec.ScaleTargetRef.Name,
				APIVersion: autoscale.Spec.ScaleTargetRef.APIVersion,
			},
			MinReplicas:                    autoscale.Spec.MinReplicas,
			MaxReplicas:                    autoscale.Spec.MaxReplicas,
			TargetCPUUtilizationPercentage: autoscale.Spec.TargetCPUUtilizationPercentage,
		},
		Status: model.HorizontalPodAutoscalerStatus{
			ObservedGeneration:              autoscale.Status.ObservedGeneration,
			LastScaleTime:                   lastTime,
			CurrentReplicas:                 autoscale.Status.CurrentReplicas,
			DesiredReplicas:                 autoscale.Status.DesiredReplicas,
			CurrentCPUUtilizationPercentage: autoscale.Status.CurrentCPUUtilizationPercentage,
		},
	}
}

func FromK8sAutoScaleList(asList *autoscalev1.HorizontalPodAutoscalerList) *model.AutoScaleList {
	if asList == nil {
		return nil
	}
	items := make([]model.AutoScale, 0)
	for i := range asList.Items {
		if as := FromK8sAutoScale(&asList.Items[i]); as != nil {
			items = append(items, *as)
		}
	}
	return &model.AutoScaleList{
		Items: items,
	}
}

func ToK8sAutoScale(autoscale *model.AutoScale) *autoscalev1.HorizontalPodAutoscaler {
	var lastTime *metav1.Time
	if autoscale.Status.LastScaleTime != nil {
		t := metav1.NewTime(*autoscale.Status.LastScaleTime)
		lastTime = &t
	}
	return &autoscalev1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AutoScaling",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(autoscale.ObjectMeta),
		Spec: autoscalev1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalev1.CrossVersionObjectReference{
				Kind:       autoscale.Spec.ScaleTargetRef.Kind,
				Name:       autoscale.Spec.ScaleTargetRef.Name,
				APIVersion: autoscale.Spec.ScaleTargetRef.APIVersion,
			},
			MinReplicas:                    autoscale.Spec.MinReplicas,
			MaxReplicas:                    autoscale.Spec.MaxReplicas,
			TargetCPUUtilizationPercentage: autoscale.Spec.TargetCPUUtilizationPercentage,
		},
		Status: autoscalev1.HorizontalPodAutoscalerStatus{
			ObservedGeneration:              autoscale.Status.ObservedGeneration,
			LastScaleTime:                   lastTime,
			CurrentReplicas:                 autoscale.Status.CurrentReplicas,
			DesiredReplicas:                 autoscale.Status.DesiredReplicas,
			CurrentCPUUtilizationPercentage: autoscale.Status.CurrentCPUUtilizationPercentage,
		},
	}
}
