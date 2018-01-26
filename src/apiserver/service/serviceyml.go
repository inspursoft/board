package service

import (
	"encoding/json"
	"errors"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
	modelK8s "k8s.io/client-go/pkg/api/v1"
	modelK8sExt "k8s.io/client-go/pkg/apis/extensions/v1beta1"
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
	pathErr                        = errors.New("ERR_DEPLOYMENT_PATH_NOT_DIRECTORY")
	emptyServiceNameErr            = errors.New("ERR_NO_SERVICE_NAME")
	portMaxErr                     = errors.New("ERR_SERVICE_NODEPORT_EXCEED_MAX_LIMIT")
	portMinErr                     = errors.New("ERR_SERVICE_NODEPORT_EXCEED_MIN_LIMIT")
	emptyDeployErr                 = errors.New("ERR_NO_DEPLOYMENT_NAME")
	invalidErr                     = errors.New("ERR_DEPLOYMENT_REPLICAS_INVAILD")
	emptyContainerErr              = errors.New("ERR_NO_CONTAINER")
	NameInconsistentErr            = errors.New("ERR_SERVICE_NAME_AND_DEPLOYMENT_NAME_INCONSISTENT")
	ServiceNameInconsistentErr     = errors.New("ERR_SERVICE_NAME_INCONSISTENT_WITH_YAML_FILE")
	ProjectNameInconsistentErr     = errors.New("ERR_PROJECT_NAME_INCONSISTENT_WITH_YAML_FILE")
	DeploymentNotFoundErr          = errors.New("ERR_DEPLOYMENT_NOT_FOUND")
	ServiceNotFoundErr             = errors.New("ERR_SERVICE_NOT_FOUND")
	ServiceYamlFileUnmarshalErr    = errors.New("ERR_SERVICE_YAML_FILE_UNMARSHAL")
	DeploymentYamlFileUnmarshalErr = errors.New("ERR_DEPLOYMENT_YAML_FILE_UNMARSHAL")
	deploymentKindErr              = errors.New("ERR_DEPLOYMENT_YAML_KIND")
	serviceKindErr                 = errors.New("ERR_SERVICE_YAML_KIND")
	deploymentAPIVersionErr        = errors.New("ERR_DEPLOYMENT_YAML_API_VERSION")
	serviceAPIVersionErr           = errors.New("ERR_SERVICE_YAML_API_VERSION")
)

type Service struct {
	modelK8s.Service
}

type Deployment struct {
	modelK8sExt.Deployment
}

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

func GetYamlFileServiceName(file multipart.File, fileName string) (string, error) {
	config, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	var service modelK8s.Service
	var deployment modelK8s.ReplicationController
	var serviceName string
	if fileName == deploymentFilename {
		err = yaml.Unmarshal(config, &deployment)
		if err != nil {
			return "", err
		}
		serviceName = deployment.ObjectMeta.Name
	} else if fileName == serviceFilename {
		err = yaml.Unmarshal(config, &service)
		if err != nil {
			return "", err
		}
		serviceName = service.ObjectMeta.Name
	}

	return serviceName, nil
}

func ServiceExists(serviceName string, projectName string) (bool, error) {
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

func GenerateDeploymentYamlFileFromK8S(deployConfigURL string, absFileName string) error {
	deployConfig, err := GetDeployConfig(deployConfigURL)
	if err != nil {
		return err
	}

	return GenerateYamlFile(absFileName, &deployConfig)
}

func GenerateServiceYamlFileFromK8S(serviceConfigURL string, absFileName string) error {
	serviceConfig, err := GetServiceStatus(serviceConfigURL)
	if err != nil {
		return err
	}

	return GenerateYamlFile(absFileName, &serviceConfig)
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

//check parameter of service yaml file
func CheckServiceConfig(projectName string, serviceConfig Service) error {
	//check empty
	if serviceConfig.Kind != serviceKind {
		return deploymentKindErr
	}
	if serviceConfig.APIVersion != serviceAPIVersion {
		return deploymentAPIVersionErr
	}
	if serviceConfig.ObjectMeta.Name == "" {
		return emptyServiceNameErr
	}
	if serviceConfig.ObjectMeta.Namespace != projectName {
		return ProjectNameInconsistentErr
	}

	for _, external := range serviceConfig.Spec.Ports {
		if external.NodePort > maxPort {
			return portMaxErr
		} else if external.NodePort < minPort {
			return portMinErr
		}
	}

	return nil
}

func CheckDeploymentConfig(projectName string, deploymentConfig Deployment) error {
	//check empty
	if deploymentConfig.Kind != deploymentKind {
		return deploymentKindErr
	}
	if deploymentConfig.APIVersion != deploymentAPIVersion {
		return deploymentAPIVersionErr
	}
	if deploymentConfig.ObjectMeta.Name == "" {
		return emptyDeployErr
	}
	if deploymentConfig.ObjectMeta.Namespace != projectName {
		return ProjectNameInconsistentErr
	}
	if *deploymentConfig.Spec.Replicas < 1 {
		return invalidErr
	}
	if len(deploymentConfig.Spec.Template.Spec.Containers) < 1 {
		return emptyContainerErr
	}

	for _, cont := range deploymentConfig.Spec.Template.Spec.Containers {

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

func RollingUpdateDeploymentYaml(repoPath string, deploydata *modelK8sExt.Deployment) error {
	projectName := deploydata.ObjectMeta.Namespace
	serviceName := deploydata.ObjectMeta.Name
	serviceInfo, err := GetServiceByProject(serviceName, projectName)
	if err != nil {
		return err
	}
	if serviceInfo == nil {
		return ServiceNotFoundErr
	}

	fileName := filepath.Join(repoPath, projectName, strconv.Itoa(int(serviceInfo.ID)), deploymentFilename)
	err = GenerateYamlFile(fileName, deploydata)
	if err != nil {
		return err
	}

	return nil
}
