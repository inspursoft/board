package service

import (
	"errors"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	yaml "yaml-2"

	"github.com/astaxie/beego/logs"
)

var loadPath string

func SetDeploymentPath(Path string) {
	loadPath = Path
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

//build yaml file of service
func BuildServiceYml(reqServiceConfig model.ServiceConfig) error {
	var service model.ServiceStructYml
	var port model.PortsServiceYml
	var selector model.SelectorServiceYml

	serviceLoadPath := GetDeploymentPath()
	err := CheckDeploymentPath(serviceLoadPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return err
	}

	service.ApiVersion = "v1"
	service.Kind = "Service"
	service.Metadata.Name = reqServiceConfig.ServiceYaml.Name
	service.Metadata.Lables.App = reqServiceConfig.ServiceYaml.Name

	if len(reqServiceConfig.ServiceYaml.ExternalPaths) > 0 {
		service.Spec.Tpe = "NodePort"
	}

	for _, nodePort := range reqServiceConfig.ServiceYaml.NodePorts {
		port.Port = "8080"
		port.Tagertport = "8080"
		port.Nodeport = strconv.Itoa(nodePort)
		service.Spec.Ports = append(service.Spec.Ports, port)
	}

	for _, sltor := range reqServiceConfig.ServiceYaml.Selectors {
		selector.App = sltor
		service.Spec.Selector = append(service.Spec.Selector, selector)
	}

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
func BuildDeploymentYml(reqServiceConfig model.ServiceConfig) error {
	var deployment model.DeploymentStructYml
	var nfsvolume model.VolumesDeploymentYml
	var container model.ContainersDeploymentYml
	var port model.PortsDeploymentYml
	var volumemount model.VolumeMountDeploymentYml
	var env model.EnvDeploymentYml

	deploymentLoadPath := GetDeploymentPath()
	err := CheckDeploymentPath(deploymentLoadPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return err
	}

	deployment.ApiVersion = "extensions/v1bata1"
	deployment.Kind = "Deployment"
	deployment.Metadata.Name = reqServiceConfig.DeploymentYaml.Name
	deployment.Spec.Replicas = reqServiceConfig.DeploymentYaml.Replicas
	deployment.Spec.Template.Metadata.Lables.App = reqServiceConfig.DeploymentYaml.Name

	for _, vlme := range reqServiceConfig.DeploymentYaml.VolumeList {
		nfsvolume.Name = vlme.Name
		nfsvolume.Nfs.Path = vlme.Path
		nfsvolume.Nfs.Server = vlme.ServerName
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, nfsvolume)
	}

	for _, ctner := range reqServiceConfig.DeploymentYaml.ContainerList {
		container.Name = ctner.Name
		container.Image = ctner.BaseImage
		container.Workingdir = ctner.WorkDir
		container.Command = ctner.Command
		container.Resource.Request.Cpu = ctner.CPU
		container.Resource.Request.Memory = ctner.Memory

		for _, prt := range ctner.Ports {
			port.ContainerPort = prt
			container.Ports = append(container.Ports, port)
		}

		for _, vlmeMount := range ctner.Volumes {
			volumemount.Name = vlmeMount.TargetStorageName
			volumemount.MountPath = vlmeMount.Dir
			container.VolumeMount = append(container.VolumeMount, volumemount)
		}

		for _, enviroment := range ctner.Envs {
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
