package service

import (
	"encoding/json"
	"errors"

	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
)

var loadPath string

const (
	serviceAPIVersion    = "v1"
	serviceKind          = "Service"
	nodePort             = "NodePort"
	deploymentAPIVersion = "extensions/v1beta1"
	deploymentKind       = "Deployment"
	maxPort              = 32765
	minPort              = 30000
)

var (
	pathErr             = errors.New("ERR_DEPLOYMENT_PATH_NOT_DIRECTORY")
	emptyServiceNameErr = errors.New("ERR_NO_SERVICE_NAME")
	portMaxErr          = errors.New("ERR_SERVICE_NODEPORT_EXCEED_MAX_LIMIT")
	portMinErr          = errors.New("ERR_SERVICE_NODEPORT_EXCEED_MIN_LIMIT")
	emptyDeployErr      = errors.New("ERR_NO_DEPLOYMENT_NAME")
	invalidErr          = errors.New("ERR_DEPLOYMENT_REPLICAS_INVAILD")
	emptyContainerErr   = errors.New("ERR_NO_CONTAINER")
)

func CheckDeploymentPath(loadPath string) error {
	if fi, err := os.Stat(loadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(loadPath, 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return pathErr
	}

	return nil
}

//check parameter of service yaml file
func CheckServicePara(reqServiceConfig model.ServiceConfig2) error {
	//check empty
	if reqServiceConfig.Service.ObjectMeta.Name == "" {
		return emptyServiceNameErr
	}

	for _, external := range reqServiceConfig.Service.Spec.Ports {
		if external.NodePort > maxPort {
			return portMaxErr
		} else if external.NodePort < minPort {
			return portMinErr
		}
	}

	//check upper
	err := checkStringHasUpper(reqServiceConfig.Service.ObjectMeta.Name)
	if err != nil {
		return err
	}

	for _, extPath := range reqServiceConfig.Project.ServiceExternalPath {
		err := checkStringHasUpper(extPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func IsServiceNameDuplicated(serviceName string, projectName string) (bool, error) {
	var servicequery model.ServiceStatus
	servicequery.Name = serviceName
	servicequery.ProjectName = projectName
	s, err := GetService(servicequery, "name", "project_name")

	return s != nil, err
}

//check parameter of deployment yaml file
func CheckDeploymentPara(reqServiceConfig model.ServiceConfig2) error {
	//check empty
	if reqServiceConfig.Deployment.ObjectMeta.Name == "" {
		return emptyDeployErr
	}

	if *reqServiceConfig.Deployment.Spec.Replicas < 1 {
		return invalidErr
	}

	if reqServiceConfig.Deployment.Spec.Template == nil {
		return emptyContainerErr
	}

	if len(reqServiceConfig.Deployment.Spec.Template.Spec.Containers) < 1 {
		return emptyContainerErr
	}

	//check upper
	err := checkStringHasUpper(reqServiceConfig.Deployment.ObjectMeta.Name)
	if err != nil {
		return err
	}

	for _, cont := range reqServiceConfig.Deployment.Spec.Template.Spec.Containers {

		err := checkStringHasUpper(cont.Name, cont.Image)
		if err != nil {
			return err
		}

		for _, com := range cont.Command {
			err := checkStringHasUpper(com)
			if err != nil {
				return err
			}
		}

		for _, volMount := range cont.VolumeMounts {
			err := checkStringHasUpper(volMount.Name, volMount.MountPath)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

//check request massage parameters
func CheckReqPara(reqServiceConfig model.ServiceConfig2) error {
	var err error
	//check upper of project name and phase
	err = checkStringHasUpper(reqServiceConfig.Project.ProjectName, reqServiceConfig.Project.Phase)
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

func GenerateYamlFile(name string, structdata interface{}) error {
	info, err := json.Marshal(structdata)
	if err != nil {
		logs.Error("Marhal json failed, err:%+v\n", err)
		return err
	}

	context, err := yaml.JSONToYAML(info)
	if err != nil {
		logs.Error("Generate yaml data failed, err:%+v\n", err)
		return err
	}

	err = ioutil.WriteFile(name, context, 0644)
	if err != nil {
		logs.Error("Generate yaml file failed, err:%+v\n", err)
		return err
	}
	return nil
}

func UnmarshalServiceConfigYaml(serviceConfig *model.ServiceConfig2, serviceConfigPath string) error {
	err := CheckDeploymentPath(serviceConfigPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return err
	}

	err = getYamlFileData(&serviceConfig.Service, serviceConfigPath, "service.yaml")
	if err != nil {
		return err
	}

	err = getYamlFileData(&serviceConfig.Deployment, serviceConfigPath, "deployment.yaml")
	if err != nil {
		return err
	}

	return nil
}

func UpdateServiceConfigYaml(reqServiceConfig model.ServiceConfig2, serviceConfigPath string) error {
	err := CheckDeploymentPath(serviceConfigPath)
	if err != nil {
		logs.Error("Failed to check deployment path: %+v\n", err)
		return err
	}

	deploymentFileName := filepath.Join(serviceConfigPath, "deployment.yaml")
	err = GenerateYamlFile(deploymentFileName, reqServiceConfig.Deployment)
	if err != nil {
		return err
	}

	serviceFileName := filepath.Join(serviceConfigPath, "service.yaml")
	err = GenerateYamlFile(serviceFileName, reqServiceConfig.Service)
	if err != nil {
		return err
	}
	return nil
}

func DeleteServiceConfigYaml(serviceConfigPath string) error {
	err := os.RemoveAll(serviceConfigPath)
	if err != nil {
		logs.Error("Failed to delete deployment files: %+v\n", err)
		return err
	}

	return nil
}

func getYamlFileData(serviceConfig interface{}, serviceConfigPath string, fileName string) error {
	serviceFileName := filepath.Join(serviceConfigPath, fileName)
	yamlData, err := ioutil.ReadFile(serviceFileName)
	if err != nil {
		return err
	}

	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, serviceConfig)
	if err != nil {
		return err
	}

	return nil
}
