package service

import (
	"errors"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

const (
	hostPath           = "hostpath"
	nfs                = "nfs"
	emptyDir           = ""
	deploymentFilename = "deployment.yaml"
	serviceFilename    = "service.yaml"
)

type DeployInfo struct {
	Service            *model.Service
	ServiceFileInfo    []byte
	Deployment         *model.Deployment
	DeploymentFileInfo []byte
}

func DeployService(serviceConfig *model.ConfigServiceStep, K8sMasterURL string, registryURI string) (*DeployInfo, error) {
	clusterConfig := &k8sassist.K8sAssistConfig{K8sMasterURL: K8sMasterURL}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)
	deploymentConfig := MarshalDeployment(serviceConfig, registryURI)
	deploymentInfo, deploymentFileInfo, err := cli.AppV1().Deployment(serviceConfig.ProjectName).Create(deploymentConfig)
	if err != nil {
		logs.Error("Deploy deployment object of %s failed. error: %+v\n", serviceConfig.ServiceName, err)
		return nil, err
	}

	svcConfig := MarshalService(serviceConfig)
	serviceInfo, serviceFileInfo, err := cli.AppV1().Service(serviceConfig.ProjectName).Create(svcConfig)
	if err != nil {
		logs.Error("Deploy service object of %s failed. error: %+v\n", serviceConfig.ServiceName, err)
		return nil, err
	}

	return &DeployInfo{
		Service:            serviceInfo,
		ServiceFileInfo:    serviceFileInfo,
		Deployment:         deploymentInfo,
		DeploymentFileInfo: deploymentFileInfo,
	}, nil
}

func GenerateDeployYamlFiles(deployInfo *DeployInfo, loadPath string) error {
	if deployInfo == nil {
		logs.Error("Deploy info is empty.")
		return errors.New("Deploy info is empty.")
	}
	err := GenerateServiceYamlFile(deployInfo.ServiceFileInfo, loadPath)
	if err != nil {
		return err
	}
	err = GenerateDeploymentYamlFile(deployInfo.DeploymentFileInfo, loadPath)
	if err != nil {
		return err
	}

	return nil
}

func GenerateDeploymentYamlFile(deploymentInfo []byte, loadPath string) error {
	deploymentAbsName := filepath.Join(loadPath, deploymentFilename)
	err := ioutil.WriteFile(deploymentAbsName, deploymentInfo, 0644)
	if err != nil {
		logs.Error("Generate deployment object yaml file failed, err:%+v\n", err)
		return err
	}

	return nil
}

func GenerateServiceYamlFile(serviceInfo []byte, loadPath string) error {
	ServiceAbsName := filepath.Join(loadPath, serviceFilename)
	err := ioutil.WriteFile(ServiceAbsName, serviceInfo, 0644)
	if err != nil {
		logs.Error("Generate service object yaml file failed, err:%+v\n", err)
		return err
	}

	return nil
}

func DeployServiceByYaml(projectName, K8sMasterURL, loadPath string) error {
	clusterConfig := &k8sassist.K8sAssistConfig{K8sMasterURL: K8sMasterURL}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)

	deploymentAbsName := filepath.Join(loadPath, deploymentFilename)
	deploymentFile, err := os.Open(deploymentAbsName)
	if err != nil {
		return err
	}

	defer deploymentFile.Close()
	_, err = cli.AppV1().Deployment(projectName).CreateByYaml(deploymentFile)
	if err != nil {
		logs.Error("Deploy deployment object by deployment.yaml failed, err:%+v\n", err)
		return err
	}

	ServiceAbsName := filepath.Join(loadPath, serviceFilename)
	serviceFile, err := os.Open(ServiceAbsName)
	if err != nil {
		return err
	}
	defer serviceFile.Close()
	_, err = cli.AppV1().Service(projectName).CreateByYaml(serviceFile)
	if err != nil {
		logs.Error("Deploy service object by service.yaml failed, err:%+v\n", err)
		return err
	}
	return nil
}

//check yaml file config
func CheckDeployYamlConfig(serviceFile, deploymentFile io.Reader, projectName, K8sMasterURL string) (*DeployInfo, error) {
	clusterConfig := &k8sassist.K8sAssistConfig{K8sMasterURL: K8sMasterURL}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)

	deploymentInfo, err := cli.AppV1().Deployment(projectName).CheckYaml(deploymentFile)
	if err != nil {
		logs.Error("Check deployment object by deployment.yaml failed, err:%+v\n", err)
		return nil, err
	}

	serviceInfo, err := cli.AppV1().Service(projectName).CheckYaml(serviceFile)
	if err != nil {
		logs.Error("Check service object by service.yaml failed, err:%+v\n", err)
		return nil, err
	}
	return &DeployInfo{
		Service:    serviceInfo,
		Deployment: deploymentInfo,
	}, nil
}
