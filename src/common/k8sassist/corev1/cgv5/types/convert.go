package types

import (
	"time"

	"git/inspursoft/board/src/common/model"

	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	autoscalev1 "k8s.io/api/autoscaling/v1"
	autoscalingapi "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/tools/remotecommand"
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
func ToK8sDeployment(deployment *model.Deployment) *appsv1.Deployment {
	if deployment == nil {
		return nil
	}
	var templ v1.PodTemplateSpec
	if t := ToK8sPodTemplateSpec(&deployment.Spec.Template); t != nil {
		templ = *t
	}
	rep := deployment.Spec.Replicas
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: ToK8sObjectMeta(deployment.ObjectMeta),
		Spec: appsv1.DeploymentSpec{
			Replicas: &rep,
			Selector: &metav1.LabelSelector{
				MatchLabels: deployment.Spec.Selector,
			},
			Template: templ,
			Paused:   deployment.Spec.Paused,
		},
		Status: appsv1.DeploymentStatus{
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
func ToK8sReplicaSet(rs *model.ReplicaSet) *appsv1.ReplicaSet {
	if rs == nil {
		return nil
	}
	var spec appsv1.ReplicaSetSpec
	if s := ToK8sReplicaSetSpec(&rs.Spec); s != nil {
		spec = *s
	}
	conds := make([]appsv1.ReplicaSetCondition, len(rs.Status.Conditions))
	for i := range rs.Status.Conditions {
		conds[i] = appsv1.ReplicaSetCondition{
			Type:               appsv1.ReplicaSetConditionType(string(rs.Status.Conditions[i].Type)),
			Status:             v1.ConditionStatus(string(rs.Status.Conditions[i].Status)),
			LastTransitionTime: metav1.NewTime(rs.Status.Conditions[i].LastTransitionTime),
			Reason:             rs.Status.Conditions[i].Reason,
			Message:            rs.Status.Conditions[i].Message,
		}
	}
	return &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: ToK8sObjectMeta(rs.ObjectMeta),
		Spec:       spec,
		Status: appsv1.ReplicaSetStatus{
			Replicas:             rs.Status.Replicas,
			FullyLabeledReplicas: rs.Status.FullyLabeledReplicas,
			ReadyReplicas:        rs.Status.ReadyReplicas,
			AvailableReplicas:    rs.Status.AvailableReplicas,
			ObservedGeneration:   rs.Status.ObservedGeneration,
			Conditions:           conds,
		},
	}
}

func ToK8sReplicaSetSpec(spec *model.ReplicaSetSpec) *appsv1.ReplicaSetSpec {
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
	return &appsv1.ReplicaSetSpec{
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

	affinity := &v1.Affinity{}
	if spec.Affinity.PodAffinity != nil {
		affinity.PodAffinity = &v1.PodAffinity{}
		for _, term := range spec.Affinity.PodAffinity {
			affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
				affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution, ToK8sAffinityTerm(term),
			)
		}
	}
	if spec.Affinity.PodAntiAffinity != nil {
		affinity.PodAntiAffinity = &v1.PodAntiAffinity{}
		for _, term := range spec.Affinity.PodAntiAffinity {
			affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
				affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, ToK8sAffinityTerm(term),
			)
		}
	}
	if spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		affinity.NodeAffinity = &v1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
				NodeSelectorTerms: ToK8sNodeSelectorTerms(spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms),
			},
		}
	}

	return &v1.PodSpec{
		Volumes:        volumes,
		InitContainers: initContainers,
		Containers:     containers,
		NodeSelector:   spec.NodeSelector,
		NodeName:       spec.NodeName,
		HostNetwork:    spec.HostNetwork,
		Affinity:       affinity,
		RestartPolicy:  v1.RestartPolicy(string(spec.RestartPolicy)),
		Tolerations:    ToK8sTolerations(spec.Tolerations),
	}
}

func ToK8sTolerations(modeltols []model.Toleration) []v1.Toleration {
	var tols []v1.Toleration
	if modeltols != nil {
		for _, tol := range modeltols {
			tols = append(tols, ToK8sToleration(tol))
		}
	}
	return tols
}

func ToK8sToleration(tol model.Toleration) v1.Toleration {
	return v1.Toleration{
		Key:               tol.Key,
		Operator:          v1.TolerationOperator(string(tol.Operator)),
		Value:             tol.Value,
		Effect:            v1.TaintEffect(string(tol.Effect)),
		TolerationSeconds: tol.TolerationSeconds,
	}
}

func ToK8sNodeSelectorTerms(terms []model.NodeSelectorTerm) []v1.NodeSelectorTerm {
	var nodeSelectorTerms []v1.NodeSelectorTerm
	for _, term := range terms {
		nodeSelectorTerms = append(nodeSelectorTerms, v1.NodeSelectorTerm{
			MatchExpressions: ToK8sNodeSelectorRequirements(term.MatchExpressions),
		})
	}
	return nodeSelectorTerms
}

func ToK8sNodeSelectorRequirements(NodeSelectorRequirements []model.NodeSelectorRequirement) []v1.NodeSelectorRequirement {
	var K8sNodeSelectorRequirements []v1.NodeSelectorRequirement
	for _, NodeSelectorRequirement := range NodeSelectorRequirements {
		K8sNodeSelectorRequirements = append(K8sNodeSelectorRequirements, v1.NodeSelectorRequirement{
			Key:      NodeSelectorRequirement.Key,
			Operator: v1.NodeSelectorOperator(NodeSelectorRequirement.Operator),
			Values:   NodeSelectorRequirement.Values,
		})
	}
	return K8sNodeSelectorRequirements
}

func ToK8sAffinityTerm(term model.PodAffinityTerm) v1.PodAffinityTerm {
	return v1.PodAffinityTerm{
		LabelSelector: &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      term.LabelSelector.MatchExpressions[0].Key,
					Operator: metav1.LabelSelectorOperator(term.LabelSelector.MatchExpressions[0].Operator),
					Values:   term.LabelSelector.MatchExpressions[0].Values,
				},
			},
		},
		Namespaces:  term.Namespaces,
		TopologyKey: term.TopologyKey,
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

	var pvc *v1.PersistentVolumeClaimVolumeSource
	if volumeSource.PersistentVolumeClaim != nil {
		pvc = &v1.PersistentVolumeClaimVolumeSource{
			ClaimName: volumeSource.PersistentVolumeClaim.ClaimName,
			ReadOnly:  volumeSource.PersistentVolumeClaim.ReadOnly,
		}
	}

	var configmap *v1.ConfigMapVolumeSource
	if volumeSource.ConfigMap != nil {
		configmap = &v1.ConfigMapVolumeSource{
			LocalObjectReference: (v1.LocalObjectReference)(volumeSource.ConfigMap.LocalObjectReference),
			//			Items:                volumeSource.ConfigMap.Items,
			//TODO: suppurt Item later
		}
	}

	return &v1.VolumeSource{
		HostPath:              hp,
		NFS:                   nfs,
		PersistentVolumeClaim: pvc,
		ConfigMap:             configmap,
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
	if v, ok := container.Resources.Limits["nvidia.com/gpu"]; ok {
		resources.Limits["nvidia.com/gpu"] = resource.MustParse(string(v))
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

func ToK8sConfigMapKeySelector(cmks *model.ConfigMapKeySelector) *v1.ConfigMapKeySelector {
	if cmks == nil {
		return nil
	}
	return &v1.ConfigMapKeySelector{
		Key:      cmks.Key,
		Optional: cmks.Optional,
		LocalObjectReference: v1.LocalObjectReference{
			Name: cmks.LocalObjectReference.Name,
		},
	}
}

func ToK8sEnvVar(env model.EnvVar) v1.EnvVar {
	var valuefrom *v1.EnvVarSource
	if env.ValueFrom == nil {
		valuefrom = nil
	} else {
		valuefrom = &v1.EnvVarSource{
			ConfigMapKeyRef: ToK8sConfigMapKeySelector(env.ValueFrom.ConfigMapKeyRef),
		}
	}
	return v1.EnvVar{
		Name:      env.Name,
		Value:     env.Value,
		ValueFrom: valuefrom,
	}
}

func ToK8sVolumeMount(mount model.VolumeMount) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      mount.Name,
		MountPath: mount.MountPath,
		SubPath:   mount.SubPath,
	}
}

func ToK8sGroupResource(gr model.GroupResource) schema.GroupResource {
	return schema.GroupResource{
		Group:    gr.Group,
		Resource: gr.Resource,
	}
}

func ToK8sTerminalSizeQueue(tsq model.TerminalSizeQueue) remotecommand.TerminalSizeQueue {
	return TerminalSizeQueueFunc(func() *remotecommand.TerminalSize {
		size := tsq.Next()
		return &remotecommand.TerminalSize{
			Width:  size.Width,
			Height: size.Height,
		}
	})
}

func ToK8sDaemonSet(daemonset *model.DaemonSet) *appsv1.DaemonSet {
	if daemonset == nil {
		return nil
	}
	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       daemonsetKind,
			APIVersion: daemonsetAPIVersion,
		},
		ObjectMeta: ToK8sObjectMeta(daemonset.ObjectMeta),
		Spec:       ToK8sDaemonSetSpec(daemonset.Spec),
		Status:     ToK8sDaemonSetStatus(daemonset.Status),
	}
}

func ToK8sDaemonSetSpec(spec model.DaemonSetSpec) appsv1.DaemonSetSpec {
	var templ v1.PodTemplateSpec
	if t := ToK8sPodTemplateSpec(&spec.Template); t != nil {
		templ = *t
	}
	return appsv1.DaemonSetSpec{
		Selector:             ToK8sLabelSelector(spec.Selector),
		Template:             templ,
		UpdateStrategy:       ToK8sDaemonSetUpdateStrategy(spec.UpdateStrategy),
		MinReadySeconds:      spec.MinReadySeconds,
		RevisionHistoryLimit: spec.RevisionHistoryLimit,
	}
}

func ToK8sDaemonSetUpdateStrategy(strategy model.DaemonSetUpdateStrategy) appsv1.DaemonSetUpdateStrategy {
	return appsv1.DaemonSetUpdateStrategy{
		Type:          appsv1.DaemonSetUpdateStrategyType(string(strategy.Type)),
		RollingUpdate: ToK8sRollingUpdateDaemonSet(strategy.RollingUpdate),
	}
}

func ToK8sRollingUpdateDaemonSet(rollingupdate *model.RollingUpdateDaemonSet) *appsv1.RollingUpdateDaemonSet {
	if rollingupdate == nil {
		return nil
	}
	var max *intstr.IntOrString
	if modelmax := rollingupdate.MaxUnavailable; modelmax != nil {
		m := intstr.Parse(modelmax.String())
		max = &m
	}
	return &appsv1.RollingUpdateDaemonSet{
		MaxUnavailable: max,
	}
}

func ToK8sDaemonSetStatus(status model.DaemonSetStatus) appsv1.DaemonSetStatus {
	return appsv1.DaemonSetStatus{
		CurrentNumberScheduled: status.CurrentNumberScheduled,
		NumberMisscheduled:     status.NumberMisscheduled,
		DesiredNumberScheduled: status.DesiredNumberScheduled,
		NumberReady:            status.NumberReady,
		ObservedGeneration:     status.ObservedGeneration,
		UpdatedNumberScheduled: status.UpdatedNumberScheduled,
		NumberAvailable:        status.NumberAvailable,
		NumberUnavailable:      status.NumberUnavailable,
		CollisionCount:         status.CollisionCount,
		Conditions:             ToK8sDaemonSetConditions(status.Conditions),
	}
}

func ToK8sDaemonSetConditions(list []model.DaemonSetCondition) []appsv1.DaemonSetCondition {
	if list == nil {
		return nil
	}
	conds := make([]appsv1.DaemonSetCondition, 0, len(list))
	for _, cond := range list {
		conds = append(conds, ToK8sDaemonSetCondition(cond))
	}
	return conds
}

func ToK8sDaemonSetCondition(cond model.DaemonSetCondition) appsv1.DaemonSetCondition {
	return appsv1.DaemonSetCondition{
		Type:               appsv1.DaemonSetConditionType(string(cond.Type)),
		Status:             v1.ConditionStatus(cond.Status),
		LastTransitionTime: metav1.NewTime(cond.LastTransitionTime),
		Reason:             cond.Reason,
		Message:            cond.Message,
	}
}

func FromK8sInfo(info *version.Info) *model.KubernetesInfo {
	return &model.KubernetesInfo{
		Major:        info.Major,
		Minor:        info.Minor,
		GitVersion:   info.GitVersion,
		GitCommit:    info.GitCommit,
		GitTreeState: info.GitTreeState,
		BuildDate:    info.BuildDate,
		GoVersion:    info.GoVersion,
		Compiler:     info.Compiler,
		Platform:     info.Platform,
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
func FromK8sDeploymentList(deploymentList *appsv1.DeploymentList) *model.DeploymentList {
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
func FromK8sDeployment(deployment *appsv1.Deployment) *model.Deployment {
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

func FromK8sDeploymentSpec(spec *appsv1.DeploymentSpec) *model.DeploymentSpec {
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
func FromK8sReplicaSetList(list *appsv1.ReplicaSetList) *model.ReplicaSetList {
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
func FromK8sReplicaSet(rs *appsv1.ReplicaSet) *model.ReplicaSet {
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

func FromK8sReplicSetSpec(spec *appsv1.ReplicaSetSpec) *model.ReplicaSetSpec {
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
		RestartPolicy:  model.RestartPolicy(string(spec.RestartPolicy)),
		Tolerations:    FromK8sTolerations(spec.Tolerations),
	}
}

func FromK8sTolerations(tols []v1.Toleration) []model.Toleration {
	var k8stols []model.Toleration
	if tols != nil {
		for _, tol := range tols {
			k8stols = append(k8stols, FromK8sToleration(tol))
		}
	}
	return k8stols
}

func FromK8sToleration(tol v1.Toleration) model.Toleration {
	return model.Toleration{
		Key:               tol.Key,
		Operator:          model.TolerationOperator(string(tol.Operator)),
		Value:             tol.Value,
		Effect:            model.TaintEffect(string(tol.Effect)),
		TolerationSeconds: tol.TolerationSeconds,
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
	var pvc *model.PersistentVolumeClaimVolumeSource
	if volumeSource.PersistentVolumeClaim != nil {
		pvc = &model.PersistentVolumeClaimVolumeSource{
			ClaimName: volumeSource.PersistentVolumeClaim.ClaimName,
			ReadOnly:  volumeSource.PersistentVolumeClaim.ReadOnly,
		}
	}

	var configmap *model.ConfigMapVolumeSource
	if volumeSource.ConfigMap != nil {
		configmap = &model.ConfigMapVolumeSource{
			LocalObjectReference: (model.LocalObjectReference)(volumeSource.ConfigMap.LocalObjectReference),
		}
	}

	return model.VolumeSource{
		HostPath:              hp,
		NFS:                   nfs,
		PersistentVolumeClaim: pvc,
		ConfigMap:             configmap,
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
	if v, ok := container.Resources.Limits["nvidia.com/gpu"]; ok {
		resources.Limits["nvidia.com/gpu"] = model.QuantityStr(v.String())
	}
	return &model.K8sContainer{
		Name:            container.Name,
		Image:           container.Image,
		Command:         container.Command,
		Args:            container.Args,
		WorkingDir:      container.WorkingDir,
		Ports:           ports,
		Env:             envs,
		Resources:       resources,
		VolumeMounts:    mounts,
		SecurityContext: FromK8sSecurityContext(container.SecurityContext),
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
	var valuefrom *model.EnvVarSource
	if env.ValueFrom == nil {
		valuefrom = nil
	} else if env.ValueFrom.ConfigMapKeyRef != nil {
		// configmap env type
		valuefrom = &model.EnvVarSource{
			ConfigMapKeyRef: &model.ConfigMapKeySelector{
				LocalObjectReference: model.LocalObjectReference(env.ValueFrom.ConfigMapKeyRef.LocalObjectReference),
				Key:                  env.ValueFrom.ConfigMapKeyRef.Key,
				Optional:             env.ValueFrom.ConfigMapKeyRef.Optional,
			},
		}
		// TODO support other types
	}
	return model.EnvVar{
		Name:      env.Name,
		Value:     env.Value,
		ValueFrom: valuefrom,
	}
}

func FromK8sVolumeMount(mount v1.VolumeMount) model.VolumeMount {
	return model.VolumeMount{
		Name:      mount.Name,
		MountPath: mount.MountPath,
		SubPath:   mount.SubPath,
	}
}

func FromK8sSecurityContext(context *v1.SecurityContext) *model.SecurityContext {
	var privileged *bool
	if context != nil && context.Privileged != nil {
		privileged = context.Privileged
	}
	return &model.SecurityContext{
		Privileged: privileged,
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
	var sessionAffinity SessionAffinity
	var sessionAffinityConfig *SessionAffinityConfig
	var timeoutSecond int32 = 0
	if modelService.SessionAffinityFlag != 0 {
		if modelService.SessionAffinityTime == 0 {
			timeoutSecond = 1800
		} else {
			timeoutSecond = int32(modelService.SessionAffinityTime)
		}
		sessionAffinity = SessionAffinityClientIP
		sessionAffinityConfig = &SessionAffinityConfig{
			ClientIP: &ClientIPConfig{TimeoutSeconds: &timeoutSecond},
		}
	} else {
		sessionAffinity = ServiceAffinityNone
		sessionAffinityConfig = &SessionAffinityConfig{
			ClientIP: &ClientIPConfig{TimeoutSeconds: &timeoutSecond},
		}
	}
	return &Service{
		TypeMeta: TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(modelService.ObjectMeta),
		Spec: ServiceSpec{
			Ports:                 ports,
			Selector:              modelService.Selector,
			ClusterIP:             modelService.ClusterIP,
			Type:                  ServiceType(modelService.Type),
			ExternalIPs:           modelService.ExternalIPs,
			ExternalName:          modelService.ExternalName,
			SessionAffinity:       sessionAffinity,
			SessionAffinityConfig: sessionAffinityConfig,
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
	var sessionAffinityFlag int
	var sessionAffinityTime int
	if typesService.Spec.SessionAffinity == SessionAffinityClientIP {
		sessionAffinityFlag = 1
		sessionAffinityTime = int(*typesService.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds)
	}
	return &model.Service{
		ObjectMeta:          FromK8sObjectMeta(typesService.ObjectMeta),
		Ports:               ports,
		Selector:            typesService.Spec.Selector,
		ClusterIP:           typesService.Spec.ClusterIP,
		Type:                string(typesService.Spec.Type),
		ExternalIPs:         typesService.Spec.ExternalIPs,
		ExternalName:        typesService.Spec.ExternalName,
		SessionAffinityFlag: sessionAffinityFlag,
		SessionAffinityTime: sessionAffinityTime,
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

func ToK8sNodeTaints(taints []model.Taint) []v1.Taint {
	var k8staints []v1.Taint
	if taints != nil {
		for _, t := range taints {
			k8staints = append(k8staints, ToK8sNodeTaint(t))
		}
	}
	return k8staints
}

func ToK8sNodeTaint(taint model.Taint) v1.Taint {
	var added *metav1.Time
	if taint.TimeAdded != nil {
		added = &metav1.Time{
			*taint.TimeAdded,
		}
	}
	return v1.Taint{
		Key:       taint.Key,
		Value:     taint.Value,
		Effect:    v1.TaintEffect(string(taint.Effect)),
		TimeAdded: added,
	}
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
			Taints:        ToK8sNodeTaints(node.Taints),
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
	k8sNode.Spec.Taints = ToK8sNodeTaints(node.Taints)

	UpdateK8sNodeStatus(&k8sNode.Status, &node.Status)
}

func FromK8sNodeTaints(taints []v1.Taint) []model.Taint {
	var ts []model.Taint
	if taints != nil {
		for _, t := range taints {
			ts = append(ts, FromK8sNodeTaint(t))
		}
	}
	return ts
}

func FromK8sNodeTaint(taint v1.Taint) model.Taint {
	var t *time.Time
	if taint.TimeAdded != nil {
		t = &taint.TimeAdded.Time
	}
	return model.Taint{
		Key:       taint.Key,
		Value:     taint.Value,
		Effect:    model.TaintEffect(string(taint.Effect)),
		TimeAdded: t,
	}
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
	nodeip := node.ObjectMeta.Name
	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			nodeip = addr.Address
			break
		}
	}
	return &model.Node{
		ObjectMeta:    FromK8sObjectMeta(node.ObjectMeta),
		NodeIP:        nodeip,
		Unschedulable: node.Spec.Unschedulable,
		Taints:        FromK8sNodeTaints(node.Spec.Taints),
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

func FromK8sScale(scale *autoscalingapi.Scale) *model.Scale {
	return &model.Scale{
		ObjectMeta: FromK8sObjectMeta(scale.ObjectMeta),
		Spec:       model.ScaleSpec(scale.Spec),
		Status:     model.ScaleStatusK8s(scale.Status),
	}
}

func ToK8sScale(scale *model.Scale) *autoscalingapi.Scale {
	return &autoscalingapi.Scale{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Scale",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(scale.ObjectMeta),
		Spec:       autoscalingapi.ScaleSpec(scale.Spec),
		Status:     autoscalingapi.ScaleStatus(scale.Status),
	}
}

func GenerateDeploymentConfig(deployment *appsv1.Deployment) *appsv1.Deployment {
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
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       deploymentKind,
			APIVersion: deploymentAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    deployment.ObjectMeta.Labels,
			Name:      deployment.ObjectMeta.Name,
			Namespace: deployment.ObjectMeta.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
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
					HostNetwork:        deployment.Spec.Template.Spec.HostNetwork,
					Tolerations:        deployment.Spec.Template.Spec.Tolerations,
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
			SessionAffinity:       service.Spec.SessionAffinity,
			SessionAffinityConfig: service.Spec.SessionAffinityConfig,
			ClusterIP:             service.Spec.ClusterIP,
			Ports:                 service.Spec.Ports,
			Selector:              service.Spec.Selector,
			Type:                  service.Spec.Type,
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

// update k8s autoscale using model autosacle
func UpdateK8sAutoScale(k8sHPA *autoscalev1.HorizontalPodAutoscaler, autoscale *model.AutoScale) {
	if k8sHPA == nil || autoscale == nil {
		return
	}
	// just update our attributes.
	k8sHPA.Spec = autoscalev1.HorizontalPodAutoscalerSpec{
		ScaleTargetRef: autoscalev1.CrossVersionObjectReference{
			Kind:       autoscale.Spec.ScaleTargetRef.Kind,
			Name:       autoscale.Spec.ScaleTargetRef.Name,
			APIVersion: autoscale.Spec.ScaleTargetRef.APIVersion,
		},
		MinReplicas:                    autoscale.Spec.MinReplicas,
		MaxReplicas:                    autoscale.Spec.MaxReplicas,
		TargetCPUUtilizationPercentage: autoscale.Spec.TargetCPUUtilizationPercentage,
	}
}

func IsNotFoundError(err error) bool {
	return errors.IsNotFound(err)
}

func IsAlreadyExistError(err error) bool {
	return errors.IsAlreadyExists(err)
}

func FromK8sRBD(rbd *v1.RBDPersistentVolumeSource) *model.RBDPersistentVolumeSource {
	if rbd == nil {
		return nil
	}
	return &model.RBDPersistentVolumeSource{
		CephMonitors: rbd.CephMonitors,
		RBDImage:     rbd.RBDImage,
		FSType:       rbd.FSType,
		RBDPool:      rbd.RBDPool,
		RadosUser:    rbd.RadosUser,
		Keyring:      rbd.Keyring,
		SecretRef:    (*model.SecretReference)(rbd.SecretRef),
		ReadOnly:     rbd.ReadOnly,
	}
}

func FromK8sPVAccessMode(am []v1.PersistentVolumeAccessMode) []model.PersistentVolumeAccessMode {
	items := make([]model.PersistentVolumeAccessMode, 0)
	for _, valueK8s := range am {
		items = append(items, (model.PersistentVolumeAccessMode)(valueK8s))
	}
	return items
}

func FromK8sPVObjectReference(or *v1.ObjectReference) *model.ObjectReference {
	if or == nil {
		return nil
	}
	return &model.ObjectReference{
		Kind:            or.Kind,
		Namespace:       or.Namespace,
		Name:            or.Name,
		UID:             model.UID(or.UID),
		APIVersion:      or.APIVersion,
		ResourceVersion: or.ResourceVersion,
		FieldPath:       or.FieldPath,
	}
}

func FromK8sPV(pv *v1.PersistentVolume) *model.PersistentVolumeK8scli {
	//var lastTime *time.Time
	//if autoscale.Status.LastScaleTime != nil {
	//	lastTime = &autoscale.Status.LastScaleTime.Time
	//}
	capacity := make(map[model.ResourceName]model.QuantityStr)
	for k, v := range pv.Spec.Capacity {
		i, _ := v.AsInt64()
		capacity[model.ResourceName(k)] = model.QuantityStr(strconv.Itoa(int(i)))

	}

	return &model.PersistentVolumeK8scli{
		ObjectMeta: FromK8sObjectMeta(pv.ObjectMeta),
		Spec: model.PersistentVolumeSpec{
			Capacity: capacity,
			PersistentVolumeSource: model.PersistentVolumeSource{
				NFS: (*model.NFSVolumeSource)(pv.Spec.PersistentVolumeSource.NFS),
				RBD: FromK8sRBD(pv.Spec.PersistentVolumeSource.RBD),
			},
			AccessModes:                   FromK8sPVAccessMode(pv.Spec.AccessModes),
			ClaimRef:                      FromK8sPVObjectReference(pv.Spec.ClaimRef),
			PersistentVolumeReclaimPolicy: (model.PersistentVolumeReclaimPolicy)(pv.Spec.PersistentVolumeReclaimPolicy),
			StorageClassName:              pv.Spec.StorageClassName,
			MountOptions:                  pv.Spec.MountOptions,
		},
		Status: model.PersistentVolumeStatus{
			Phase:   (model.PersistentVolumePhase)(pv.Status.Phase),
			Message: pv.Status.Message,
			Reason:  pv.Status.Reason,
		},
	}
}

func FromK8sPVList(pvList *v1.PersistentVolumeList) *model.PersistentVolumeList {
	if pvList == nil {
		return nil
	}
	items := make([]model.PersistentVolumeK8scli, 0)
	for i := range pvList.Items {
		if pv := FromK8sPV(&pvList.Items[i]); pv != nil {
			items = append(items, *pv)
		}
	}
	return &model.PersistentVolumeList{
		Items: items,
	}
}

func ToK8sRBD(rbd *model.RBDPersistentVolumeSource) *v1.RBDPersistentVolumeSource {
	if rbd == nil {
		return nil
	}
	return &v1.RBDPersistentVolumeSource{
		CephMonitors: rbd.CephMonitors,
		RBDImage:     rbd.RBDImage,
		FSType:       rbd.FSType,
		RBDPool:      rbd.RBDPool,
		RadosUser:    rbd.RadosUser,
		Keyring:      rbd.Keyring,
		SecretRef:    (*v1.SecretReference)(rbd.SecretRef),
		ReadOnly:     rbd.ReadOnly,
	}
}

func ToK8sPVAccessMode(am []model.PersistentVolumeAccessMode) []v1.PersistentVolumeAccessMode {
	items := make([]v1.PersistentVolumeAccessMode, 0)
	for _, valueK8s := range am {
		items = append(items, (v1.PersistentVolumeAccessMode)(valueK8s))
	}
	return items
}

func ToK8sPVObjectReference(or *model.ObjectReference) *v1.ObjectReference {
	if or == nil {
		return nil
	}
	return &v1.ObjectReference{
		Kind:            or.Kind,
		Namespace:       or.Namespace,
		Name:            or.Name,
		UID:             types.UID(or.UID),
		APIVersion:      or.APIVersion,
		ResourceVersion: or.ResourceVersion,
		FieldPath:       or.FieldPath,
	}
}

func ToK8sPV(pv *model.PersistentVolumeK8scli) *v1.PersistentVolume {
	//var lastTime *metav1.Time
	//if autoscale.Status.LastScaleTime != nil {
	//	t := metav1.NewTime(*autoscale.Status.LastScaleTime)
	//	lastTime = &t
	//}
	capacity := make(map[v1.ResourceName]resource.Quantity)
	if v, ok := pv.Spec.Capacity["storage"]; ok {
		capacity["storage"] = resource.MustParse(string(v))
	}

	return &v1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(pv.ObjectMeta),
		Spec: v1.PersistentVolumeSpec{
			Capacity: capacity,
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: (*v1.NFSVolumeSource)(pv.Spec.PersistentVolumeSource.NFS),
				RBD: ToK8sRBD(pv.Spec.PersistentVolumeSource.RBD),
			},
			AccessModes:                   ToK8sPVAccessMode(pv.Spec.AccessModes),
			ClaimRef:                      ToK8sPVObjectReference(pv.Spec.ClaimRef),
			PersistentVolumeReclaimPolicy: (v1.PersistentVolumeReclaimPolicy)(pv.Spec.PersistentVolumeReclaimPolicy),
			StorageClassName:              pv.Spec.StorageClassName,
			MountOptions:                  pv.Spec.MountOptions,
		},
		Status: v1.PersistentVolumeStatus{
			Phase:   (v1.PersistentVolumePhase)(pv.Status.Phase),
			Message: pv.Status.Message,
			Reason:  pv.Status.Reason,
		},
	}
}

//TODO implement update later， only support capacity now
func UpdateK8sPV(k8sPV *v1.PersistentVolume, pv *model.PersistentVolumeK8scli) {
	if k8sPV == nil || pv == nil {
		return
	}
	capacity := make(map[v1.ResourceName]resource.Quantity)
	for k, v := range pv.Spec.Capacity {
		value, _ := strconv.Atoi(string(v))
		q := resource.NewQuantity(int64(value), resource.DecimalExponent)
		capacity[v1.ResourceName(k)] = *q

	}
	// just update our attributes.
	k8sPV.Spec.Capacity = capacity
}

// PVC convert
func FromK8sPVC(pvc *v1.PersistentVolumeClaim) *model.PersistentVolumeClaimK8scli {

	var resources model.ResourceRequirements
	resources.Requests = make(model.ResourceList)
	resources.Limits = make(model.ResourceList)
	if v, ok := pvc.Spec.Resources.Requests["storage"]; ok {
		resources.Requests["storage"] = model.QuantityStr(v.String())
	}

	if v, ok := pvc.Spec.Resources.Limits["storage"]; ok {
		resources.Limits["storage"] = model.QuantityStr(v.String())
	}

	capacity := make(map[model.ResourceName]model.QuantityStr)
	for k, v := range pvc.Status.Capacity {
		i, _ := v.AsInt64()
		capacity[model.ResourceName(k)] = model.QuantityStr(strconv.Itoa(int(i)))

	}

	return &model.PersistentVolumeClaimK8scli{
		ObjectMeta: FromK8sObjectMeta(pvc.ObjectMeta),
		Spec: model.PersistentVolumeClaimSpec{
			AccessModes:      FromK8sPVAccessMode(pvc.Spec.AccessModes),
			VolumeName:       pvc.Spec.VolumeName,
			Resources:        resources,
			StorageClassName: pvc.Spec.StorageClassName,
		},
		Status: model.PersistentVolumeClaimStatus{
			Phase:       (model.PersistentVolumeClaimPhase)(pvc.Status.Phase),
			AccessModes: FromK8sPVAccessMode(pvc.Status.AccessModes),
			Capacity:    capacity,
		},
	}
}

func FromK8sPVCList(pvcList *v1.PersistentVolumeClaimList) *model.PersistentVolumeClaimList {
	if pvcList == nil {
		return nil
	}
	items := make([]model.PersistentVolumeClaimK8scli, 0)
	for i := range pvcList.Items {
		if pvc := FromK8sPVC(&pvcList.Items[i]); pvc != nil {
			items = append(items, *pvc)
		}
	}
	return &model.PersistentVolumeClaimList{
		Items: items,
	}
}

func ToK8sPVC(pvc *model.PersistentVolumeClaimK8scli) *v1.PersistentVolumeClaim {

	capacity := make(map[v1.ResourceName]resource.Quantity)
	for k, v := range pvc.Status.Capacity {
		value, _ := strconv.Atoi(string(v))
		q := resource.NewQuantity(int64(value), resource.DecimalExponent)
		capacity[v1.ResourceName(k)] = *q

	}

	var resources v1.ResourceRequirements
	resources.Requests = make(v1.ResourceList)
	resources.Limits = make(v1.ResourceList)
	if v, ok := pvc.Spec.Resources.Requests["storage"]; ok {
		resources.Requests["storage"] = resource.MustParse(string(v))
	}
	if v, ok := pvc.Spec.Resources.Limits["storage"]; ok {
		resources.Limits["storage"] = resource.MustParse(string(v))
	}

	return &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(pvc.ObjectMeta),
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes:      ToK8sPVAccessMode(pvc.Spec.AccessModes),
			VolumeName:       pvc.Spec.VolumeName,
			Resources:        resources,
			StorageClassName: pvc.Spec.StorageClassName,
		},
		Status: v1.PersistentVolumeClaimStatus{
			Phase:       (v1.PersistentVolumeClaimPhase)(pvc.Status.Phase),
			AccessModes: ToK8sPVAccessMode(pvc.Status.AccessModes),
			Capacity:    capacity,
		},
	}
}

//TODO implement update later， only support capacity now
func UpdateK8sPVC(k8sPVC *v1.PersistentVolumeClaim, pvc *model.PersistentVolumeClaimK8scli) {
	//	if k8sPV == nil || pv == nil {
	//		return
	//	}
	//	capacity := make(map[v1.ResourceName]resource.Quantity)
	//	for k, v := range pv.Spec.Capacity {
	//		value, _ := strconv.Atoi(string(v))
	//		q := resource.NewQuantity(int64(value), resource.DecimalExponent)
	//		capacity[v1.ResourceName(k)] = *q

	//	}
	//	// just update our attributes.
	//	k8sPV.Spec.Capacity = capacity
}

// ConfigMap convert
func FromK8sConfigMap(configmap *v1.ConfigMap) *model.ConfigMap {

	return &model.ConfigMap{
		ObjectMeta: FromK8sObjectMeta(configmap.ObjectMeta),
		Data:       configmap.Data,
	}
}

func FromK8sConfigMapList(configmapList *v1.ConfigMapList) *model.ConfigMapList {
	if configmapList == nil {
		return nil
	}
	items := make([]model.ConfigMap, 0)
	for i := range configmapList.Items {
		if configmap := FromK8sConfigMap(&configmapList.Items[i]); configmap != nil {
			items = append(items, *configmap)
		}
	}
	return &model.ConfigMapList{
		Items: items,
	}
}

func ToK8sConfigMap(configmap *model.ConfigMap) *v1.ConfigMap {

	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: ToK8sObjectMeta(configmap.ObjectMeta),
		Data:       configmap.Data,
	}
}

//TODO implement update later
func UpdateK8sConfigMap(k8sCM *v1.ConfigMap, cm *model.ConfigMap) {
	if k8sCM == nil || cm == nil {
		return
	}
	k8sCM.Data = cm.Data

}

// ToK8sStatefulSet is to generate k8s statefulset from model statefulset
func ToK8sStatefulSet(statefulset *model.StatefulSet) *appsv1.StatefulSet {
	if statefulset == nil {
		return nil
	}
	var templ v1.PodTemplateSpec
	if t := ToK8sPodTemplateSpec(&statefulset.Spec.Template); t != nil {
		templ = *t
	}
	//rep := deployment.Spec.Replicas

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: ToK8sObjectMeta(statefulset.ObjectMeta),
		Spec: appsv1.StatefulSetSpec{
			Replicas: statefulset.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: statefulset.Spec.Selector.MatchLabels,
			},
			Template: templ,
			//TODO: support storage class later
			//VolumeClaimTemplates:   statefulset.Spec.VolumeClaimTemplates,
			ServiceName: statefulset.Spec.ServiceName,
			//UpdateStrategy:  statefulset.Spec.UpdateStrategy,
			//RevisionHistoryLimit:  statefulset.Spec.RevisionHistoryLimit,
		},
		Status: appsv1.StatefulSetStatus{
			Replicas:        statefulset.Status.Replicas,
			ReadyReplicas:   statefulset.Status.ReadyReplicas,
			CurrentReplicas: statefulset.Status.CurrentReplicas,
			UpdatedReplicas: statefulset.Status.UpdatedReplicas,
			CurrentRevision: statefulset.Status.CurrentRevision,
			UpdateRevision:  statefulset.Status.UpdateRevision,
		},
	}
}
func ToK8sJob(job *model.Job) *Job {
	if job == nil {
		return nil
	}
	var templ v1.PodTemplateSpec
	if t := ToK8sPodTemplateSpec(&job.Spec.Template); t != nil {
		templ = *t
	}
	conditions := make([]JobCondition, 0)
	for _, condition := range job.Status.Conditions {
		conditions = append(conditions, JobCondition{
			Type:   JobConditionType(condition.Type),
			Status: ConditionStatus(condition.Status),
			LastProbeTime: metav1.Time{
				condition.LastProbeTime,
			},
			LastTransitionTime: metav1.Time{
				condition.LastTransitionTime,
			},
			Reason:  condition.Reason,
			Message: condition.Message,
		})
	}
	var starttime *metav1.Time
	if job.Status.StartTime != nil {
		starttime = &metav1.Time{
			*job.Status.StartTime,
		}
	}
	var completiontime *metav1.Time
	if job.Status.CompletionTime != nil {
		completiontime = &metav1.Time{
			*job.Status.CompletionTime,
		}
	}
	return &Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: ToK8sObjectMeta(job.ObjectMeta),
		Spec: JobSpec{
			Parallelism:           job.Spec.Parallelism,
			Completions:           job.Spec.Completions,
			ActiveDeadlineSeconds: job.Spec.ActiveDeadlineSeconds,
			BackoffLimit:          job.Spec.BackoffLimit,
			Selector:              ToK8sLabelSelector(job.Spec.Selector),
			ManualSelector:        job.Spec.ManualSelector,
			Template:              templ,
		},
		Status: JobStatus{
			Conditions:     conditions,
			StartTime:      starttime,
			CompletionTime: completiontime,
			Active:         job.Status.Active,
			Succeeded:      job.Status.Succeeded,
			Failed:         job.Status.Failed,
		},
	}
}

func FromK8sJob(job *Job) *model.Job {
	if job == nil {
		return nil
	}
	var templ model.PodTemplateSpec
	if t := FromK8sPodTemplateSpec(&job.Spec.Template); t != nil {
		templ = *t
	}
	conditions := make([]model.JobCondition, 0)
	for _, condition := range job.Status.Conditions {
		conditions = append(conditions, model.JobCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastProbeTime:      condition.LastProbeTime.Time,
			LastTransitionTime: condition.LastTransitionTime.Time,
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	var starttime *time.Time
	if job.Status.StartTime != nil {
		starttime = &job.Status.StartTime.Time
	}
	var completiontime *time.Time
	if job.Status.CompletionTime != nil {
		completiontime = &job.Status.CompletionTime.Time
	}
	return &model.Job{
		ObjectMeta: FromK8sObjectMeta(job.ObjectMeta),
		Spec: model.JobSpec{
			Parallelism:           job.Spec.Parallelism,
			Completions:           job.Spec.Completions,
			ActiveDeadlineSeconds: job.Spec.ActiveDeadlineSeconds,
			BackoffLimit:          job.Spec.BackoffLimit,
			Selector:              FromK8sLabelSelector(job.Spec.Selector),
			ManualSelector:        job.Spec.ManualSelector,
			Template:              templ,
		},
		Status: model.JobStatus{
			Conditions:     conditions,
			StartTime:      starttime,
			CompletionTime: completiontime,
			Active:         job.Status.Active,
			Succeeded:      job.Status.Succeeded,
			Failed:         job.Status.Failed,
		},
	}
}

func GenerateJobConfig(job *Job) *Job {
	containersConfig := []v1.Container{}
	for _, container := range job.Spec.Template.Spec.Containers {
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
	return &Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       job.TypeMeta.Kind,
			APIVersion: job.TypeMeta.APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    job.ObjectMeta.Labels,
			Name:      job.ObjectMeta.Name,
			Namespace: job.ObjectMeta.Namespace,
		},
		Spec: JobSpec{
			Parallelism:           job.Spec.Parallelism,
			Completions:           job.Spec.Completions,
			ActiveDeadlineSeconds: job.Spec.ActiveDeadlineSeconds,
			BackoffLimit:          job.Spec.BackoffLimit,
			Selector:              job.Spec.Selector,
			ManualSelector:        job.Spec.ManualSelector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: job.Spec.Template.ObjectMeta.Labels,
					Name:   job.Spec.Template.ObjectMeta.Name,
				},
				Spec: v1.PodSpec{
					Affinity:           job.Spec.Template.Spec.Affinity,
					Volumes:            job.Spec.Template.Spec.Volumes,
					NodeSelector:       job.Spec.Template.Spec.NodeSelector,
					ServiceAccountName: job.Spec.Template.Spec.ServiceAccountName,
					ImagePullSecrets:   job.Spec.Template.Spec.ImagePullSecrets,
					InitContainers:     job.Spec.Template.Spec.InitContainers,
					Containers:         containersConfig,
					RestartPolicy:      job.Spec.Template.Spec.RestartPolicy,
				},
			},
		},
	}
}

func FromK8sJobList(jobList *JobList) *model.JobList {
	if jobList == nil {
		return nil
	}
	items := make([]model.Job, 0)
	for i := range jobList.Items {
		job := FromK8sJob(&jobList.Items[i])
		items = append(items, *job)
	}
	return &model.JobList{
		Items: items,
	}
}

func FromK8sLabelSelector(selector *metav1.LabelSelector) *model.LabelSelector {
	if selector == nil {
		return nil
	}
	var expretions []model.LabelSelectorRequirement
	for i := range selector.MatchExpressions {
		expretions = append(expretions, model.LabelSelectorRequirement{
			Key:      selector.MatchExpressions[i].Key,
			Operator: string(selector.MatchExpressions[i].Operator),
			Values:   selector.MatchExpressions[i].Values,
		})
	}
	return &model.LabelSelector{
		MatchLabels:      selector.MatchLabels,
		MatchExpressions: expretions,
	}
}

func ToK8sLabelSelector(selector *model.LabelSelector) *metav1.LabelSelector {
	if selector == nil {
		return nil
	}
	var expretions []metav1.LabelSelectorRequirement
	for i := range selector.MatchExpressions {
		expretions = append(expretions, metav1.LabelSelectorRequirement{
			Key:      selector.MatchExpressions[i].Key,
			Operator: metav1.LabelSelectorOperator(selector.MatchExpressions[i].Operator),
			Values:   selector.MatchExpressions[i].Values,
		})
	}
	return &metav1.LabelSelector{
		MatchLabels:      selector.MatchLabels,
		MatchExpressions: expretions,
	}
}

func ToK8sPodLogOptions(opt *model.PodLogOptions) *v1.PodLogOptions {
	if opt == nil {
		return nil
	}
	var since *metav1.Time
	if opt.SinceTime != nil {
		since = &metav1.Time{
			*opt.SinceTime,
		}
	}
	return &v1.PodLogOptions{
		Container:    opt.Container,
		Follow:       opt.Follow,
		Previous:     opt.Previous,
		SinceSeconds: opt.SinceSeconds,
		SinceTime:    since,
		Timestamps:   opt.Timestamps,
		TailLines:    opt.TailLines,
		LimitBytes:   opt.LimitBytes,
	}
}

func ToK8sListOptions(opts model.ListOptions) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector:   opts.LabelSelector,
		FieldSelector:   opts.FieldSelector,
		Watch:           opts.Watch,
		ResourceVersion: opts.ResourceVersion,
		TimeoutSeconds:  opts.TimeoutSeconds,
	}
}

func LabelSelectorToString(selector *model.LabelSelector) string {
	return metav1.FormatLabelSelector(ToK8sLabelSelector(selector))
}

// FromK8sStatefulSet is to generate model StatefulSet from k8s StatefulSet
func FromK8sStatefulSet(statefulset *appsv1.StatefulSet) *model.StatefulSet {
	if statefulset == nil {
		return nil
	}

	var template model.PodTemplateSpec
	if t := FromK8sPodTemplateSpec(&statefulset.Spec.Template); t != nil {
		template = *t
	}
	return &model.StatefulSet{
		ObjectMeta: FromK8sObjectMeta(statefulset.ObjectMeta),
		Spec: model.StatefulSetSpec{
			Replicas: statefulset.Spec.Replicas,
			Selector: &model.LabelSelector{
				MatchLabels: statefulset.Spec.Selector.MatchLabels,
			},
			Template: template,
			//TODO: support storage class later
			//VolumeClaimTemplates:
			ServiceName: statefulset.Spec.ServiceName,
			//UpdateStrategy:
			//RevisionHistoryLimit: statefulset.Spec.RevisionHistoryLimit,
		},
		Status: model.StatefulSetStatus{
			ObservedGeneration: statefulset.Status.ObservedGeneration,
			Replicas:           statefulset.Status.Replicas,
			ReadyReplicas:      statefulset.Status.ReadyReplicas,
			UpdatedReplicas:    statefulset.Status.UpdatedReplicas,
			CurrentReplicas:    statefulset.Status.CurrentReplicas,
			CurrentRevision:    statefulset.Status.CurrentRevision,
			UpdateRevision:     statefulset.Status.UpdateRevision,
			CollisionCount:     statefulset.Status.CollisionCount,
		},
	}
}

// FromK8sStatefulSetList is to generate model StatefulSetList from k8s StatefulSetList
func FromK8sStatefulSetList(statefulsetList *appsv1.StatefulSetList) *model.StatefulSetList {
	if statefulsetList == nil {
		return nil
	}
	items := make([]model.StatefulSet, 0)
	for i := range statefulsetList.Items {
		if statefulset := FromK8sStatefulSet(&statefulsetList.Items[i]); statefulset != nil {
			items = append(items, *statefulset)
		}
	}
	return &model.StatefulSetList{
		Items: items,
	}
}

// GenerateStatefulSetConfig is to generate stateful config
func GenerateStatefulSetConfig(statefulset *appsv1.StatefulSet) *appsv1.StatefulSet {
	containersConfig := []v1.Container{}
	for _, container := range statefulset.Spec.Template.Spec.Containers {
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
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    statefulset.ObjectMeta.Labels,
			Name:      statefulset.ObjectMeta.Name,
			Namespace: statefulset.ObjectMeta.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: statefulset.Spec.Replicas,
			Selector: statefulset.Spec.Selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: statefulset.Spec.Template.ObjectMeta.Labels,
					Name:   statefulset.Spec.Template.ObjectMeta.Name,
				},
				Spec: v1.PodSpec{
					Affinity:           statefulset.Spec.Template.Spec.Affinity,
					Volumes:            statefulset.Spec.Template.Spec.Volumes,
					NodeSelector:       statefulset.Spec.Template.Spec.NodeSelector,
					ServiceAccountName: statefulset.Spec.Template.Spec.ServiceAccountName,
					ImagePullSecrets:   statefulset.Spec.Template.Spec.ImagePullSecrets,
					InitContainers:     statefulset.Spec.Template.Spec.InitContainers,
					Containers:         containersConfig,
					Tolerations:        statefulset.Spec.Template.Spec.Tolerations,
				},
			},
			ServiceName: statefulset.Spec.ServiceName,
		},
	}
}

// generate model daemonset list from k8s daemonset list
func FromK8sDaemonSetList(daemonsetList *appsv1.DaemonSetList) *model.DaemonSetList {
	if daemonsetList == nil {
		return nil
	}
	items := make([]model.DaemonSet, 0)
	for i := range daemonsetList.Items {
		ds := FromK8sDaemonSet(&daemonsetList.Items[i])
		items = append(items, *ds)
	}
	return &model.DaemonSetList{
		Items: items,
	}
}

func FromK8sDaemonSet(daemonset *appsv1.DaemonSet) *model.DaemonSet {
	if daemonset == nil {
		return nil
	}
	return &model.DaemonSet{
		ObjectMeta: FromK8sObjectMeta(daemonset.ObjectMeta),
		Spec:       FromK8sDaemonSetSpec(daemonset.Spec),
		Status:     FromK8sDaemonSetStatus(daemonset.Status),
	}
}

func FromK8sDaemonSetSpec(spec appsv1.DaemonSetSpec) model.DaemonSetSpec {
	var template model.PodTemplateSpec
	if t := FromK8sPodTemplateSpec(&spec.Template); t != nil {
		template = *t
	}
	return model.DaemonSetSpec{
		Selector:             FromK8sLabelSelector(spec.Selector),
		Template:             template,
		UpdateStrategy:       FromK8sDaemonSetUpdateStrategy(spec.UpdateStrategy),
		MinReadySeconds:      spec.MinReadySeconds,
		RevisionHistoryLimit: spec.RevisionHistoryLimit,
	}
}

func FromK8sDaemonSetUpdateStrategy(strategy appsv1.DaemonSetUpdateStrategy) model.DaemonSetUpdateStrategy {
	return model.DaemonSetUpdateStrategy{
		Type:          model.DaemonSetUpdateStrategyType(string(strategy.Type)),
		RollingUpdate: FromK8sRollingUpdateDaemonSet(strategy.RollingUpdate),
	}
}

func FromK8sRollingUpdateDaemonSet(rollingupdate *appsv1.RollingUpdateDaemonSet) *model.RollingUpdateDaemonSet {
	if rollingupdate == nil {
		return nil
	}
	var max *model.IntOrString
	if k8smax := rollingupdate.MaxUnavailable; k8smax != nil {
		m := model.Parse(k8smax.String())
		max = &m
	}
	return &model.RollingUpdateDaemonSet{
		MaxUnavailable: max,
	}
}

func FromK8sDaemonSetStatus(status appsv1.DaemonSetStatus) model.DaemonSetStatus {
	return model.DaemonSetStatus{
		CurrentNumberScheduled: status.CurrentNumberScheduled,
		NumberMisscheduled:     status.NumberMisscheduled,
		DesiredNumberScheduled: status.DesiredNumberScheduled,
		NumberReady:            status.NumberReady,
		ObservedGeneration:     status.ObservedGeneration,
		UpdatedNumberScheduled: status.UpdatedNumberScheduled,
		NumberAvailable:        status.NumberAvailable,
		NumberUnavailable:      status.NumberUnavailable,
		CollisionCount:         status.CollisionCount,
		Conditions:             FromK8sDaemonSetConditions(status.Conditions),
	}
}

func FromK8sDaemonSetConditions(list []appsv1.DaemonSetCondition) []model.DaemonSetCondition {
	if list == nil {
		return nil
	}
	conds := make([]model.DaemonSetCondition, 0, len(list))
	for _, cond := range list {
		conds = append(conds, FromK8sDaemonSetCondition(cond))
	}
	return conds
}

func FromK8sDaemonSetCondition(cond appsv1.DaemonSetCondition) model.DaemonSetCondition {
	return model.DaemonSetCondition{
		Type:               model.DaemonSetConditionType(string(cond.Type)),
		Status:             model.ConditionStatus(cond.Status),
		LastTransitionTime: cond.LastTransitionTime.Time,
		Reason:             cond.Reason,
		Message:            cond.Message,
	}
}

func GenerateDaemonSetConfig(daemonset *appsv1.DaemonSet) *appsv1.DaemonSet {
	containersConfig := []v1.Container{}
	for _, container := range daemonset.Spec.Template.Spec.Containers {
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
	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       daemonsetKind,
			APIVersion: daemonsetAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    daemonset.ObjectMeta.Labels,
			Name:      daemonset.ObjectMeta.Name,
			Namespace: daemonset.ObjectMeta.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: daemonset.Spec.Selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: daemonset.Spec.Template.ObjectMeta.Labels,
					Name:   daemonset.Spec.Template.ObjectMeta.Name,
				},
				Spec: v1.PodSpec{
					Affinity:           daemonset.Spec.Template.Spec.Affinity,
					Volumes:            daemonset.Spec.Template.Spec.Volumes,
					NodeSelector:       daemonset.Spec.Template.Spec.NodeSelector,
					ServiceAccountName: daemonset.Spec.Template.Spec.ServiceAccountName,
					ImagePullSecrets:   daemonset.Spec.Template.Spec.ImagePullSecrets,
					InitContainers:     daemonset.Spec.Template.Spec.InitContainers,
					Containers:         containersConfig,
					Tolerations:        daemonset.Spec.Template.Spec.Tolerations,
				},
			},
			UpdateStrategy:       daemonset.Spec.UpdateStrategy,
			MinReadySeconds:      daemonset.Spec.MinReadySeconds,
			RevisionHistoryLimit: daemonset.Spec.RevisionHistoryLimit,
		},
	}
}
