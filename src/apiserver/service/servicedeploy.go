package service

import (
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
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
	deploymentConfig := utils.MarshalDeployment(serviceConfig, registryURI)
	deploymentInfo, deploymentFileInfo, err := cli.AppV1().Deployment(serviceConfig.ProjectName).Create(deploymentConfig)
	if err != nil {
		logs.Error("Deploy deployment object of %s failed.", serviceConfig.ServiceName)
		return nil, err
	}

	svcConfig := utils.MarshalService(serviceConfig)
	serviceInfo, serviceFileInfo, err := cli.AppV1().Service(serviceConfig.ProjectName).Create(svcConfig)
	if err != nil {
		logs.Error("Deploy service object of %s failed.", serviceConfig.ServiceName)
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
	ServiceAbsName := filepath.Join(loadPath, serviceFilename)
	err := ioutil.WriteFile(ServiceAbsName, deployInfo.ServiceFileInfo, 0644)
	if err != nil {
		logs.Error("Generate service object yaml file failed, err:%+v\n", err)
		return err
	}

	deploymentAbsName := filepath.Join(loadPath, deploymentFilename)
	err = ioutil.WriteFile(deploymentAbsName, deployInfo.DeploymentFileInfo, 0644)
	if err != nil {
		logs.Error("Generate deployment object yaml file failed, err:%+v\n", err)
		return err
	}

	return nil
}


func DeployServiceByYaml(projectName,K8sMasterURL,loadPath string)( error){
	clusterConfig := &k8sassist.K8sAssistConfig{K8sMasterURL: K8sMasterURL}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)

	deploymentAbsName := filepath.Join(loadPath, deploymentFilename)
	deploymentFile, err := os.Open(deploymentAbsName)
	if err != nil {
		return  err
	}
	defer deploymentFile.Close()
	deploymentInfo, err := cli.AppV1().Deployment(projectName).CreateByYaml(deploymentFile)
	if err != nil {
		logs.Error("Deploy deployment object by deployment.yaml failed.")
		return  err
	}

	ServiceAbsName := filepath.Join(loadPath, serviceFilename)
	serviceFile, err := os.Open(ServiceAbsName)
	if err != nil {
		return  err
	}
	defer serviceFile.Close()
	serviceInfo, err := cli.AppV1().Service(projectName).CreateByYaml()
	if err != nil {
		logs.Error("Deploy service object by service.yaml failed.")
		return  err
	}
	return nil
}