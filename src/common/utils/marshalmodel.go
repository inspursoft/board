package utils

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"strings"
)

func MarshalService(serviceConfig *model.ConfigServiceStep) *model.Service {
	if serviceConfig == nil {
		return nil
	}
	ports := make([]model.ServicePort, 0)
	for _, port := range serviceConfig.ExternalServiceList {
		ports = append(ports, model.ServicePort{
			Port:     int32(port.NodeConfig.TargetPort),
			NodePort: int32(port.NodeConfig.NodePort),
		})
	}

	return &model.Service{
		ObjectMeta: model.ObjectMeta{Name: serviceConfig.ServiceName},
		Ports:      ports,
		Selector:   map[string]string{"app": serviceConfig.ServiceName},
		Type:       "NodePort",
	}
}

func setDeploymentNodeSelector(nodeOrNodeGroupName string) map[string]string {
	if nodeOrNodeGroupName == "" {
		return nil
	}
	nodegroup, _ := dao.GetNodeGroup(model.NodeGroup{GroupName: nodeOrNodeGroupName}, "name")
	if nodegroup != nil && nodegroup.ID != 0 {
		return map[string]string{nodeOrNodeGroupName: "true"}
	} else {
		return map[string]string{"kubernetes.io/hostname": nodeOrNodeGroupName}
	}
}

func setDeploymentContainers(containerList []model.Container, registryURI string) []model.K8sContainer {
	if containerList == nil {
		return nil
	}
	k8sContainerList := make([]model.K8sContainer, 0)
	for _, cont := range containerList {
		container := model.K8sContainer{}
		container.Name = cont.Name

		if cont.WorkingDir != "" {
			container.WorkingDir = cont.WorkingDir
		}

		if len(cont.Command) > 0 {
			container.Command = append(container.Command, "/bin/sh")
			container.Args = append(container.Args, "-c", cont.Command)
		}

		if cont.VolumeMounts.VolumeName != "" {
			container.VolumeMounts = append(container.VolumeMounts, model.VolumeMount{
				Name:      cont.VolumeMounts.VolumeName,
				MountPath: cont.VolumeMounts.ContainerPath,
			})
		}

		if len(cont.Env) > 0 {
			for _, enviroment := range cont.Env {
				if enviroment.EnvName != "" {
					container.Env = append(container.Env, model.EnvVar{
						Name:  enviroment.EnvName,
						Value: enviroment.EnvValue,
					})
				}
			}
		}

		for _, port := range cont.ContainerPort {
			container.Ports = append(container.Ports, model.ContainerPort{
				ContainerPort: int32(port),
			})
		}

		container.Image = registryURI + "/" + cont.Image.ImageName + ":" + cont.Image.ImageTag

		k8sContainerList = append(k8sContainerList, container)
	}
	return k8sContainerList
}

func setDeploymentVolumes(containerList []model.Container) []model.Volume {
	if containerList == nil {
		return nil
	}
	volumes := make([]model.Volume, 0)
	for _, cont := range containerList {
		if strings.ToLower(cont.VolumeMounts.TargetStorageService) == "hostpath" {
			volumes = append(volumes, model.Volume{
				Name: cont.VolumeMounts.VolumeName,
				VolumeSource: model.VolumeSource{
					HostPath: &model.HostPathVolumeSource{
						Path: cont.VolumeMounts.TargetPath,
					},
				},
			})
		} else if strings.ToLower(cont.VolumeMounts.TargetStorageService) == "nfs" {
			index := strings.IndexByte(cont.VolumeMounts.TargetPath, '/')
			volumes = append(volumes, model.Volume{
				Name: cont.VolumeMounts.VolumeName,
				VolumeSource: model.VolumeSource{
					NFS: &model.NFSVolumeSource{
						Server: cont.VolumeMounts.TargetPath[:index],
						Path:   cont.VolumeMounts.TargetPath[index:],
					},
				},
			})
		}
	}
	return volumes
}

func MarshalDeployment(serviceConfig *model.ConfigServiceStep, registryURI string) *model.Deployment {
	if serviceConfig == nil {
		return nil
	}
	podTemplate := model.PodTemplateSpec{
		ObjectMeta: model.ObjectMeta{
			Name:   serviceConfig.ServiceName,
			Labels: map[string]string{"app": serviceConfig.ServiceName},
		},
		Spec: model.PodSpec{
			Volumes:      setDeploymentVolumes(serviceConfig.ContainerList),
			Containers:   setDeploymentContainers(serviceConfig.ContainerList, registryURI),
			NodeSelector: setDeploymentNodeSelector(serviceConfig.NodeSelector),
		},
	}

	return &model.Deployment{
		ObjectMeta: model.ObjectMeta{
			Name:      serviceConfig.ServiceName,
			Namespace: serviceConfig.ProjectName,
		},
		Spec: model.DeploymentSpec{
			Replicas: int32(serviceConfig.Instance),
			Selector: map[string]string{"app": serviceConfig.ServiceName},
			Template: podTemplate,
		},
	}
}

func MarshalNode(nodeName, labelKey string, schedFlag bool) *model.Node {
	label := make(map[string]string)
	if labelKey != "" {
		label[labelKey] = "true"
	}
	return &model.Node{
		ObjectMeta: model.ObjectMeta{
			Name:   nodeName,
			Labels: label,
		},
		Unschedulable: schedFlag,
	}
}

func MarshalNamespace(namespace string) *model.Namespace {
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{Name: namespace},
	}
}
