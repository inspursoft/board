package service

import (
	"errors"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/astaxie/beego/logs"
)

var loadPath string

const (
	serviceAPIVersion    = "v1"
	serviceKind          = "Service"
	nodePort             = "NodePort"
	deploymentAPIVersion = "extensions/v1beta1"
	deploymentKind       = "Deployment"
	maxPort              = 65535
	minPort              = 30000
)

func SetDeploymentPath(path string) {
	loadPath = path
}

func GetDeploymentPath() string {
	return loadPath
}

func CheckDeploymentPath(loadPath string) error {
	if fi, err := os.Stat(loadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(loadPath, 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return errors.New("Deployment path is not directory.")
	}

	return nil
}

//check parameter of service yaml file
func CheckServicePara(reqServiceConfig model.ServiceConfig) error {
	//check empty
	if reqServiceConfig.ServiceYaml.Name == "" {
		return errors.New("ServiceYaml.Name is empty.")
	}

	for _, external := range reqServiceConfig.ServiceYaml.External {
		if external.NodePort > maxPort {
			return errors.New("Service_nodeport exceed maximum limit.")
		} else if external.NodePort < minPort {
			return errors.New("Service_nodeport exceed minimum limit.")
		}
	}

	//check upper
	err := checkStringHasUpper(reqServiceConfig.ServiceYaml.Name,
		reqServiceConfig.ServiceYaml.Selectors[0])
	if err != nil {
		return err
	}

	for _, ext := range reqServiceConfig.ServiceYaml.External {
		err := checkStringHasUpper(ext.ContainerName, ext.ExternalPath)
		if err != nil {
			return err
		}
	}

	return nil
}

//check parameter of deployment yaml file
func CheckDeploymentPara(reqServiceConfig model.ServiceConfig) error {
	//check empty
	if reqServiceConfig.DeploymentYaml.Name == "" {
		return errors.New("Deployment_Name is empty.")
	}

	if reqServiceConfig.DeploymentYaml.Replicas < 1 {
		return errors.New("Deployment_Replicas < 1 is invaild.")
	}

	if len(reqServiceConfig.DeploymentYaml.ContainerList) < 1 {
		return errors.New("Container_List is empty.")
	}

	//check upper
	err := checkStringHasUpper(reqServiceConfig.DeploymentYaml.Name)
	if err != nil {
		return err
	}

	for _, cont := range reqServiceConfig.DeploymentYaml.ContainerList {

		err := checkStringHasUpper(cont.Name, cont.BaseImage, cont.WorkDir, cont.CPU, cont.Memory)
		if err != nil {
			return err
		}

		for _, com := range cont.Command {
			err := checkStringHasUpper(com)
			if err != nil {
				return err
			}
		}

		for _, volMount := range cont.Volumes {
			err := checkStringHasUpper(volMount.Dir, volMount.TargetStorageName)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

//check request massage parameters
func CheckReqPara(reqServiceConfig model.ServiceConfig) error {
	var err error
	//check upper of project name and phase
	err = checkStringHasUpper(reqServiceConfig.ProjectName, reqServiceConfig.Phase)
	if err != nil {
		return err
	}
	// Check deployment parameters
	err = CheckDeploymentPara(reqServiceConfig)
	if err != nil {
		return err
	}

	// Check service parameters
	err = CheckServicePara(reqServiceConfig)
	if err != nil {
		return err
	}

	return err
}

//build yaml file of service
func BuildServiceYaml(reqServiceConfig model.ServiceConfig, yamlFileName string) error {
	var service model.ServiceStructYaml
	var port model.PortsServiceYaml

	serviceLoadPath := GetDeploymentPath()
	err := CheckDeploymentPath(serviceLoadPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return err
	}

	service.ApiVersion = serviceAPIVersion
	service.Kind = serviceKind
	service.Metadata.Name = reqServiceConfig.ServiceYaml.Name
	service.Metadata.Labels.App = reqServiceConfig.ServiceYaml.Selectors[0]

	if len(reqServiceConfig.ServiceYaml.External) > 0 {
		service.Spec.Tpe = nodePort
	}

	for _, external := range reqServiceConfig.ServiceYaml.External {
		port.Port = external.ContainerPort
		port.TargetPort = external.ContainerPort
		port.NodePort = external.NodePort
		service.Spec.Ports = append(service.Spec.Ports, port)
	}

	service.Spec.Selector.App = reqServiceConfig.ServiceYaml.Selectors[0]

	context, err := yaml.Marshal(&service)
	if err != nil {
		logs.Error("Failed to Marshal service yaml file: %+v\n", err)
		return err
	}

	fileName := filepath.Join(serviceLoadPath, yamlFileName)
	err = ioutil.WriteFile(fileName, context, 0644)
	if err != nil {
		logs.Error("Failed to build service yaml file: %+v\n", err)
		return err
	}
	return nil
}

//build yaml file of deployment
func BuildDeploymentYaml(reqServiceConfig model.ServiceConfig, yamlFileName string) error {
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

	deployment.ApiVersion = deploymentAPIVersion
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

	fileName := filepath.Join(deploymentLoadPath, yamlFileName)

	err = ioutil.WriteFile(fileName, context, 0644)
	if err != nil {
		logs.Error("Failed to build deployment yaml file: %+v\n", err)
		return err
	}
	return nil
}
