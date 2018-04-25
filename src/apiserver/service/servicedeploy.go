package service

import (
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"path/filepath"
	"strings"

	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

type DeploymentConfig v1.ReplicationController
type ServiceConfig v1.Service

var registryBaseURI = utils.GetConfig("REGISTRY_BASE_URI")

const (
	hostPath           = "hostpath"
	nfs                = "nfs"
	emptyDir           = ""
	deploymentFilename = "deployment.yaml"
	serviceFilename    = "service.yaml"
)

func NewDeployment() *DeploymentConfig {
	deployConfig := DeploymentConfig{
		TypeMeta: unversioned.TypeMeta{
			Kind:       deploymentKind,
			APIVersion: deploymentAPIVersion,
		},
		Spec: v1.ReplicationControllerSpec{
			Template: &v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: make(map[string]string),
				},
			},
		},
	}
	return &deployConfig
}

func (d *DeploymentConfig) setDeploymentName(name string) {
	d.ObjectMeta.Name = name
	d.Spec.Template.ObjectMeta.Labels["app"] = name
}

func (d *DeploymentConfig) setDeploymentInstance(Instance *int32) {
	d.Spec.Replicas = Instance
}

func (d *DeploymentConfig) setDeploymentNamespace(name string) {
	d.ObjectMeta.Namespace = name
}

func (d *DeploymentConfig) setDeploymentNodeSelector(nodeOrNodeGroupName string) {
	if nodeOrNodeGroupName == "" {
		return
	}
	d.Spec.Template.Spec.NodeSelector = make(map[string]string)
	nodeGroupExists, _ := NodeGroupExists(nodeOrNodeGroupName)
	if nodeGroupExists {
		d.Spec.Template.Spec.NodeSelector[nodeOrNodeGroupName] = "true"
	} else {
		d.Spec.Template.Spec.NodeSelector["kubernetes.io/hostname"] = nodeOrNodeGroupName
	}
}

func (d *DeploymentConfig) setDeploymentContainers(ContainerList []model.Container) {
	for _, cont := range ContainerList {
		container := v1.Container{}
		container.Name = cont.Name

		if cont.WorkingDir != "" {
			container.WorkingDir = cont.WorkingDir
		}

		if len(cont.Command) > 0 {
			container.Command = append(container.Command, "/bin/sh")
			container.Args = append(container.Args, "-c", cont.Command)
		}

		if cont.VolumeMounts.VolumeName != "" {
			container.VolumeMounts = append(container.VolumeMounts, v1.VolumeMount{
				Name:      cont.VolumeMounts.VolumeName,
				MountPath: cont.VolumeMounts.ContainerPath,
			})
		}

		if len(cont.Env) > 0 {
			for _, enviroment := range cont.Env {
				if enviroment.EnvName != "" {
					container.Env = append(container.Env, v1.EnvVar{
						Name:  enviroment.EnvName,
						Value: enviroment.EnvValue,
					})
				}
			}
		}

		for _, port := range cont.ContainerPort {
			container.Ports = append(container.Ports, v1.ContainerPort{
				ContainerPort: int32(port),
			})
		}

		container.Image = registryBaseURI() + "/" + cont.Image.ImageName + ":" + cont.Image.ImageTag

		d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)
	}
}

func (d *DeploymentConfig) setDeploymentVolumes(ContainerList []model.Container) {
	for _, cont := range ContainerList {
		if strings.ToLower(cont.VolumeMounts.TargetStorageService) == hostPath {
			d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, v1.Volume{
				Name: cont.VolumeMounts.VolumeName,
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: cont.VolumeMounts.TargetPath,
					},
				},
			})
		} else if strings.ToLower(cont.VolumeMounts.TargetStorageService) == nfs {
			index := strings.IndexByte(cont.VolumeMounts.TargetPath, '/')
			d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, v1.Volume{
				Name: cont.VolumeMounts.VolumeName,
				VolumeSource: v1.VolumeSource{
					NFS: &v1.NFSVolumeSource{
						Server: cont.VolumeMounts.TargetPath[:index],
						Path:   cont.VolumeMounts.TargetPath[index:],
					},
				},
			})
		}
	}
}

func NewService() *ServiceConfig {
	serviceConfig := ServiceConfig{
		TypeMeta: unversioned.TypeMeta{
			Kind:       serviceKind,
			APIVersion: serviceAPIVersion,
		},
		Spec: v1.ServiceSpec{
			Selector: make(map[string]string),
		},
	}
	return &serviceConfig
}

func (s *ServiceConfig) setServiceName(name string) {
	s.ObjectMeta.Name = name
}

func (s *ServiceConfig) setServiceSelector(name string) {
	s.Spec.Selector["app"] = name
}

func (s *ServiceConfig) setServiceNamespace(name string) {
	s.ObjectMeta.Namespace = name
}

func (s *ServiceConfig) setServicePort(ExternalServiceList []model.ExternalService) {
	for _, extService := range ExternalServiceList {
		s.Spec.Ports = append(s.Spec.Ports, v1.ServicePort{
			Port:     int32(extService.NodeConfig.TargetPort),
			NodePort: int32(extService.NodeConfig.NodePort),
		})
	}
	s.Spec.Type = nodePort
}

func AssembleDeploymentYaml(serviceConfig *model.ConfigServiceStep, loadPath string) error {
	//build struct
	instance := (int32)(serviceConfig.Instance)
	deployConfig := NewDeployment()
	deployConfig.setDeploymentName(serviceConfig.ServiceName)
	deployConfig.setDeploymentNamespace(serviceConfig.ProjectName)
	deployConfig.setDeploymentInstance(&instance)
	deployConfig.setDeploymentNodeSelector(serviceConfig.NodeSelector)
	deployConfig.setDeploymentContainers(serviceConfig.ContainerList)
	deployConfig.setDeploymentVolumes(serviceConfig.ContainerList)
	deploymentAbsName := filepath.Join(loadPath, deploymentFilename)
	err := GenerateYamlFile(deploymentAbsName, deployConfig)
	if err != nil {
		return err
	}

	return nil
}

func AssembleServiceYaml(serviceConfig *model.ConfigServiceStep, loadPath string) error {
	//build struct
	svcConfig := NewService()
	svcConfig.setServiceName(serviceConfig.ServiceName)
	svcConfig.setServiceNamespace(serviceConfig.ProjectName)
	svcConfig.setServiceSelector(serviceConfig.ServiceName)
	svcConfig.setServicePort(serviceConfig.ExternalServiceList)
	ServiceAbsName := filepath.Join(loadPath, serviceFilename)
	err := GenerateYamlFile(ServiceAbsName, svcConfig)
	if err != nil {
		return err
	}

	return nil
}
