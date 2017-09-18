package service

import (
	"errors"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/astaxie/beego/logs"
)

var loadPath string

const serviceApiVersion = "v1"
const serviceKind = "Service"
const nodePort = "NodePort"
const deploymentApiVersion = "extensions/v1beta1"
const deploymentKind = "Deployment"
const MaxPort = 65535
const MinPort = 30000

func SetDeploymentPath(Path string) {
	loadPath = strings.Replace(Path, " ", "", -1)
}

func GetDeploymentPath() string {
	return loadPath
}

func CheckDeploymentPath(loadPath string) error {
	if len(loadPath) == 0 {
		return errors.New("loadPath is Null.")
	}

	if fi, err := os.Stat(loadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(loadPath, 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return errors.New("Doployment path is not directory.")
	}

	return nil
}

//check parameter of service yaml file
func CheckServiceYamlPara(reqServiceConfig model.ServiceConfig) error {
	if reqServiceConfig.DeploymentYaml.Name == "" {
		return errors.New("Deployment_Name is NULL.")
	}

	if reqServiceConfig.DeploymentYaml.Replicas < 1 {
		return errors.New("Deployment_Replicas < 1 is invaild.")
	}

	if len(reqServiceConfig.DeploymentYaml.ContainerList) < 1 {
		return errors.New("Container_List is NULL.")
	}

	return nil
}

//check parameter of deployment yaml file
func CheckDeploymentYamlPara(reqServiceConfig model.ServiceConfig) error {
	if reqServiceConfig.ServiceYaml.Name == "" {
		return errors.New("ServiceYaml.Name is NULL.")
	}

	for _, nodePort := range reqServiceConfig.ServiceYaml.NodePorts {
		if nodePort.ExternalPort > MaxPort {
			return errors.New("Service_nodeports exceed maximum limit.")
		} else if nodePort.ExternalPort < MinPort {
			return errors.New("Service_nodeports exceed minimum limit.")
		}
	}

	return nil
}

//build yaml file of service
func BuildServiceYaml(reqServiceConfig model.ServiceConfig) error {
	var service model.ServiceStructYaml
	var port model.PortsServiceYaml

	serviceLoadPath := GetDeploymentPath()
	err := CheckDeploymentPath(serviceLoadPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return err
	}

	service.ApiVersion = serviceApiVersion
	service.Kind = serviceKind
	service.Metadata.Name = reqServiceConfig.ServiceYaml.Name
	service.Metadata.Labels.App = reqServiceConfig.ServiceYaml.Name

	if len(reqServiceConfig.ServiceYaml.NodePorts) > 0 {
		service.Spec.Tpe = nodePort
	}

	for _, nodePort := range reqServiceConfig.ServiceYaml.NodePorts {
		port.Port = nodePort.ContainerPort
		port.TargetPort = nodePort.ContainerPort
		port.NodePort = nodePort.ExternalPort
		service.Spec.Ports = append(service.Spec.Ports, port)
	}

	// for _, sltor := range reqServiceConfig.ServiceYaml.Selectors {
	// 	selector.App = sltor
	// 	service.Spec.Selector = append(service.Spec.Selector, selector)
	// }
	service.Spec.Selector.App = reqServiceConfig.ServiceYaml.Selectors[0]

	context, err := yaml.Marshal(&service)
	if err != nil {
		logs.Error("Failed to Marshal service yaml file: %+v\n", err)
		return err
	}

	fileName := filepath.Join(serviceLoadPath, "service.yaml")
	err = ioutil.WriteFile(fileName, context, 0644)
	if err != nil {
		logs.Error("Failed to build service yaml file: %+v\n", err)
		return err
	}
	return nil
}

//build yaml file of deployment
func BuildDeploymentYaml(reqServiceConfig model.ServiceConfig) error {
	var deployment model.DeploymentStructYaml
	var nfsvolume model.VolumesDeploymentYaml
	var container model.ContainersDeploymentYaml
	var port model.PortsDeploymentYaml
	var volumeMount model.VolumeMountDeploymentYaml
	var env model.EnvDeploymentYaml

	deploymentLoadPath := GetDeploymentPath()
	err := CheckDeploymentPath(deploymentLoadPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return err
	}

	deployment.ApiVersion = deploymentApiVersion
	deployment.Kind = deploymentKind
	deployment.Metadata.Name = reqServiceConfig.DeploymentYaml.Name
	deployment.Spec.Replicas = reqServiceConfig.DeploymentYaml.Replicas
	deployment.Spec.Template.Metadata.Labels.App = reqServiceConfig.DeploymentYaml.Name

	for _, vol := range reqServiceConfig.DeploymentYaml.VolumeList {
		nfsvolume.Name = vol.Name

		if vol.ServerName == "" {
			nfsvolume.HostPath.Path = vol.Path
		} else {
			nfsvolume.Nfs.Path = vol.Path
			nfsvolume.Nfs.Server = vol.ServerName
		}

		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, nfsvolume)
	}

	for _, cont := range reqServiceConfig.DeploymentYaml.ContainerList {

		container.Name = cont.Name
		container.Image = cont.BaseImage
		container.Workingdir = cont.WorkDir
		container.Command = cont.Command
		container.Resource.Request.Cpu = cont.CPU
		container.Resource.Request.Memory = cont.Memory

		container.Ports = make([]model.PortsDeploymentYaml, 0)
		for _, por := range cont.Ports {
			port.ContainerPort = por
			container.Ports = append(container.Ports, port)
		}

		container.VolumeMount = make([]model.VolumeMountDeploymentYaml, 0)
		for _, volMount := range cont.Volumes {
			volumeMount.Name = volMount.TargetStorageName
			volumeMount.MountPath = volMount.Dir
			container.VolumeMount = append(container.VolumeMount, volumeMount)
		}

		container.Env = make([]model.EnvDeploymentYaml, 0)
		for _, enviroment := range cont.Envs {
			env.Name = enviroment.Name
			env.Value = enviroment.Value
			container.Env = append(container.Env, env)
		}

		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)
	}

	context, err := yaml.Marshal(&deployment)
	if err != nil {
		logs.Error("Failed to Marshal deployment yaml file: %+v\n", err)
		return err
	}

	fileName := filepath.Join(deploymentLoadPath, "deployment.yaml")

	err = ioutil.WriteFile(fileName, context, 0644)
	if err != nil {
		logs.Error("Failed to build deployment yaml file: %+v\n", err)
		return err
	}
	return nil
}
