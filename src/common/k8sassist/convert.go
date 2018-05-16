package k8sassist

import (
	"time"

	"git/inspursoft/board/src/common/model"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// generate k8s objectmeta from model objectmeta
func toK8sObjectMeta(meta model.ObjectMeta) metav1.ObjectMeta {
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
func toK8sDeployment(deployment *model.Deployment) *appsv1beta2.Deployment {
	if deployment == nil {
		return nil
	}
	var templ v1.PodTemplateSpec
	if t := toK8sPodTemplateSpec(&deployment.Spec.Template); t != nil {
		templ = *t
	}
	rep := deployment.Spec.Replicas
	return &appsv1beta2.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1beta2",
		},
		ObjectMeta: toK8sObjectMeta(deployment.ObjectMeta),
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
func toK8sPodTemplateSpec(template *model.PodTemplateSpec) *v1.PodTemplateSpec {
	if template == nil {
		return nil
	}
	var spec v1.PodSpec
	if s := toK8sPodSpec(&template.Spec); s != nil {
		spec = *s
	}
	return &v1.PodTemplateSpec{
		ObjectMeta: toK8sObjectMeta(template.ObjectMeta),
		Spec:       spec,
	}
}

// generate k8s pod from model pod
func toK8sPod(pod *model.Pod) *v1.Pod {
	if pod == nil {
		return nil
	}
	var spec v1.PodSpec
	if s := toK8sPodSpec(&pod.Spec); s != nil {
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
		ObjectMeta: toK8sObjectMeta(pod.ObjectMeta),
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
func toK8sPodSpec(spec *model.PodSpec) *v1.PodSpec {
	if spec == nil {
		return nil
	}
	var volumes []v1.Volume
	for i := range spec.Volumes {
		if v := toK8sVolume(&spec.Volumes[i]); v != nil {
			volumes = append(volumes, *v)
		}
	}

	var initContainers []v1.Container
	for i := range spec.InitContainers {
		if c := toK8sContainer(&spec.InitContainers[i]); c != nil {
			initContainers = append(initContainers, *c)
		}
	}

	var containers []v1.Container
	for i := range spec.Containers {
		if c := toK8sContainer(&spec.Containers[i]); c != nil {
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

func toK8sVolume(volume *model.Volume) *v1.Volume {
	if volume == nil {
		return nil
	}
	var volumeSource v1.VolumeSource
	if vs := toK8sVolumeSource(&volume.VolumeSource); vs != nil {
		volumeSource = *vs
	}
	return &v1.Volume{
		Name:         volume.Name,
		VolumeSource: volumeSource,
	}
}

func toK8sVolumeSource(volumeSource *model.VolumeSource) *v1.VolumeSource {
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

func toK8sContainer(container *model.K8sContainer) *v1.Container {
	if container == nil {
		return nil
	}
	var ports []v1.ContainerPort
	for i := range container.Ports {
		ports = append(ports, toK8sContainerPort(container.Ports[i]))
	}
	var envs []v1.EnvVar
	for i := range container.Env {
		envs = append(envs, toK8sEnvVar(container.Env[i]))
	}

	var mounts []v1.VolumeMount
	for i := range container.VolumeMounts {
		mounts = append(mounts, toK8sVolumeMount(container.VolumeMounts[i]))
	}
	return &v1.Container{
		Name:         container.Name,
		Image:        container.Image,
		Command:      container.Command,
		Args:         container.Args,
		WorkingDir:   container.WorkingDir,
		Ports:        ports,
		Env:          envs,
		VolumeMounts: mounts,
	}
}

func toK8sContainerPort(containerPort model.ContainerPort) v1.ContainerPort {
	return v1.ContainerPort{
		Name:          containerPort.Name,
		HostPort:      containerPort.HostPort,
		ContainerPort: containerPort.ContainerPort,
		Protocol:      v1.Protocol(string(containerPort.Protocol)),
		HostIP:        containerPort.HostIP,
	}
}

func toK8sEnvVar(env model.EnvVar) v1.EnvVar {
	return v1.EnvVar{
		Name:  env.Name,
		Value: env.Value,
	}
}

func toK8sVolumeMount(mount model.VolumeMount) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      mount.Name,
		MountPath: mount.MountPath,
	}
}

func fromK8sObjectMeta(meta metav1.ObjectMeta) model.ObjectMeta {
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
func fromK8sDeploymentList(deploymentList *appsv1beta2.DeploymentList) *model.DeploymentList {
	if deploymentList == nil {
		return nil
	}
	items := make([]model.Deployment, 0)
	for i := range deploymentList.Items {
		dep := fromK8sDeployment(&deploymentList.Items[i])
		items = append(items, *dep)
	}
	return &model.DeploymentList{
		Items: items,
	}
}

// generate model deployment from k8s deployment
func fromK8sDeployment(deployment *appsv1beta2.Deployment) *model.Deployment {
	if deployment == nil {
		return nil
	}
	var spec model.DeploymentSpec
	if s := fromK8sDeploymentSpec(&deployment.Spec); s != nil {
		spec = *s
	}
	return &model.Deployment{
		ObjectMeta: fromK8sObjectMeta(deployment.ObjectMeta),
		Spec:       spec,
		Status: model.DeploymentStatus{
			Replicas:            deployment.Status.Replicas,
			UpdatedReplicas:     deployment.Status.UpdatedReplicas,
			AvailableReplicas:   deployment.Status.AvailableReplicas,
			UnavailableReplicas: deployment.Status.UnavailableReplicas,
		},
	}
}

func fromK8sDeploymentSpec(spec *appsv1beta2.DeploymentSpec) *model.DeploymentSpec {
	if spec == nil {
		return nil
	}
	var rep int32
	if spec.Replicas != nil {
		rep = *spec.Replicas
	}
	var template model.PodTemplateSpec
	if t := fromK8sPodTemplateSpec(&spec.Template); t != nil {
		template = *t
	}
	return &model.DeploymentSpec{
		Replicas: rep,
		Selector: fromK8sSelector(spec.Selector),
		Template: template,
		Paused:   spec.Paused,
	}
}

func fromK8sSelector(selector *metav1.LabelSelector) map[string]string {
	if selector == nil {
		return nil
	}
	return selector.MatchLabels
}

func fromK8sPodTemplateSpec(template *v1.PodTemplateSpec) *model.PodTemplateSpec {
	if template == nil {
		return nil
	}
	var spec model.PodSpec
	if s := fromK8sPodSpec(&template.Spec); s != nil {
		spec = *s
	}
	return &model.PodTemplateSpec{
		ObjectMeta: fromK8sObjectMeta(template.ObjectMeta),
		Spec:       spec,
	}
}

func fromK8sPodList(podList *v1.PodList) *model.PodList {
	if podList == nil {
		return nil
	}
	items := make([]model.Pod, 0)
	for i := range podList.Items {
		if pod := fromK8sPod(&podList.Items[i]); pod != nil {
			items = append(items, *pod)
		}
	}
	return &model.PodList{
		Items: items,
	}
}

// generate model pod from k8s pod
func fromK8sPod(pod *v1.Pod) *model.Pod {
	if pod == nil {
		return nil
	}
	var spec model.PodSpec
	if s := fromK8sPodSpec(&pod.Spec); s != nil {
		spec = *s
	}
	var starttime *time.Time
	if t := pod.Status.StartTime; t != nil {
		starttime = &t.Time
	}
	return &model.Pod{
		ObjectMeta: fromK8sObjectMeta(pod.ObjectMeta),
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

func fromK8sPodSpec(spec *v1.PodSpec) *model.PodSpec {
	if spec == nil {
		return nil
	}
	var volumes []model.Volume = nil
	for i := range spec.Volumes {
		if v := fromK8sVolume(&spec.Volumes[i]); v != nil {
			volumes = append(volumes, *v)
		}
	}

	var initContainers []model.K8sContainer
	for i := range spec.InitContainers {
		if c := fromK8sContainer(&spec.InitContainers[i]); c != nil {
			initContainers = append(initContainers, *c)
		}
	}

	var containers []model.K8sContainer
	for i := range spec.Containers {
		if c := fromK8sContainer(&spec.Containers[i]); c != nil {
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

func fromK8sVolume(volume *v1.Volume) *model.Volume {
	if volume == nil {
		return nil
	}
	return &model.Volume{
		Name:         volume.Name,
		VolumeSource: fromK8sVolumeSource(volume.VolumeSource),
	}
}

func fromK8sVolumeSource(volumeSource v1.VolumeSource) model.VolumeSource {
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

func fromK8sContainer(container *v1.Container) *model.K8sContainer {
	if container == nil {
		return nil
	}
	var ports []model.ContainerPort
	for i := range container.Ports {
		ports = append(ports, fromK8sContainerPort(container.Ports[i]))
	}
	var envs []model.EnvVar
	for i := range container.Env {
		envs = append(envs, fromK8sEnvVar(container.Env[i]))
	}

	var mounts []model.VolumeMount
	for i := range container.VolumeMounts {
		mounts = append(mounts, fromK8sVolumeMount(container.VolumeMounts[i]))
	}
	return &model.K8sContainer{
		Name:         container.Name,
		Image:        container.Image,
		Command:      container.Command,
		Args:         container.Args,
		WorkingDir:   container.WorkingDir,
		Ports:        ports,
		Env:          envs,
		VolumeMounts: mounts,
	}
}

func fromK8sContainerPort(containerPort v1.ContainerPort) model.ContainerPort {
	return model.ContainerPort{
		Name:          containerPort.Name,
		HostPort:      containerPort.HostPort,
		ContainerPort: containerPort.ContainerPort,
		Protocol:      model.Protocol(string(containerPort.Protocol)),
		HostIP:        containerPort.HostIP,
	}
}

func fromK8sEnvVar(env v1.EnvVar) model.EnvVar {
	return model.EnvVar{
		Name:  env.Name,
		Value: env.Value,
	}
}

func fromK8sVolumeMount(mount v1.VolumeMount) model.VolumeMount {
	return model.VolumeMount{
		Name:      mount.Name,
		MountPath: mount.MountPath,
	}
}
