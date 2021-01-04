package service

import (
	"errors"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"io"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

const (
	hostPath             = "hostpath"
	nfs                  = "nfs"
	emptyDir             = ""
	deploymentFilename   = "deployment.yaml"
	statefulsetFilename  = "statefulset.yaml"
	serviceFilename      = "service.yaml"
	serviceStoppedStatus = 2
)

// DeployStatefulSetInfo is the data for yaml files of statefulset and its service
type DeployStatefulSetInfo struct {
	Service             *model.Service
	ServiceFileInfo     []byte
	StatefulSet         *model.StatefulSet
	StatefulSetFileInfo []byte
}

type DeployInfo struct {
	Service            *model.Service
	ServiceFileInfo    []byte
	Deployment         *model.Deployment
	DeploymentFileInfo []byte
}

func DeployService(serviceConfig *model.ConfigServiceStep, registryURI string) (*DeployInfo, error) {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)
	deploymentConfig := MarshalDeployment(serviceConfig, registryURI)
	//logs.Debug("Marshaled deployment: ", deploymentConfig)
	if serviceConfig.ServiceType == model.ServiceTypeEdgeComputing {
		deploymentConfig.Spec.Template.Spec.HostNetwork = true
		deploymentConfig.Spec.Template.Spec.Affinity.NodeAffinity = model.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &model.NodeSelector{
				NodeSelectorTerms: []model.NodeSelectorTerm{
					model.NodeSelectorTerm{
						MatchExpressions: []model.NodeSelectorRequirement{
							model.NodeSelectorRequirement{
								Key:      "node-role.kubernetes.io/edge",
								Operator: model.NodeSelectorOpExists,
							}}}}}}
	}
	deploymentInfo, deploymentFileInfo, err := cli.AppV1().Deployment(serviceConfig.ProjectName).Create(deploymentConfig)
	if err != nil {
		logs.Error("Deploy deployment object of %s failed. error: %+v\n", serviceConfig.ServiceName, err)
		return nil, err
	}
	logs.Debug("Created deployment: ", deploymentInfo)

	var serviceInfo *model.Service
	var serviceFileInfo []byte
	if serviceConfig.ServiceType != model.ServiceTypeEdgeComputing {
		svcConfig := MarshalService(serviceConfig)
		serviceInfo, serviceFileInfo, err = cli.AppV1().Service(serviceConfig.ProjectName).Create(svcConfig)
		if err != nil {
			cli.AppV1().Deployment(serviceConfig.ProjectName).Delete(serviceConfig.ServiceName)
			logs.Error("Deploy service object of %s failed. error: %+v\n", serviceConfig.ServiceName, err)
			return nil, err
		}
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
	if deployInfo.ServiceFileInfo != nil {
		err := utils.GenerateFile(deployInfo.ServiceFileInfo, loadPath, serviceFilename)
		if err != nil {
			return err
		}
	} else {
		logs.Warning("The file of deployInfo.ServiceFileInfo is nil.")
	}

	err := utils.GenerateFile(deployInfo.DeploymentFileInfo, loadPath, deploymentFilename)
	if err != nil {
		return err
	}

	return nil
}

// TODO: this func should be refactored with GenerateDeployYamlFiles
// GenerateStatefulSetYamlFiles
func GenerateStatefulSetYamlFiles(deployInfo *DeployStatefulSetInfo, loadPath string) error {
	if deployInfo == nil {
		logs.Error("Deploy info is empty.")
		return errors.New("Deploy info is empty.")
	}
	err := utils.GenerateFile(deployInfo.ServiceFileInfo, loadPath, serviceFilename)
	if err != nil {
		return err
	}
	err = utils.GenerateFile(deployInfo.StatefulSetFileInfo, loadPath, statefulsetFilename)
	if err != nil {
		return err
	}

	return nil
}

func DeployServiceByYaml(projectName, loadPath string) error {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)

	deploymentAbsName := filepath.Join(loadPath, deploymentFilename)
	deploymentFile, err := FetchFileContentByDevOpsOpt("master", deploymentAbsName)
	if err != nil {
		return err
	}
	defer func() {
		if h, ok := deploymentFile.(*os.File); ok {
			h.Close()
		}
	}()
	deploymentInfo, err := cli.AppV1().Deployment(projectName).CreateByYaml(deploymentFile)
	if err != nil {
		logs.Error("Deploy deployment object by deployment.yaml failed, err:%+v\n", err)
		return err
	}
	serviceAbsName := filepath.Join(loadPath, serviceFilename)
	serviceFile, err := FetchFileContentByDevOpsOpt("master", serviceAbsName)
	if err != nil {
		return err
	}
	defer func() {
		if h, ok := serviceFile.(*os.File); ok {
			h.Close()
		}
	}()
	_, err = cli.AppV1().Service(projectName).CreateByYaml(serviceFile)
	if err != nil {
		cli.AppV1().Deployment(projectName).Delete(deploymentInfo.Name)
		logs.Error("Deploy service object by service.yaml failed, err:%+v\n", err)
		return err
	}
	return nil
}

//check yaml file config
func CheckDeployYamlConfig(serviceFile, deploymentFile io.Reader, projectName string) (*DeployInfo, error) {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
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

func GetStoppedSeviceNodePorts() ([]int32, error) {
	stoppedServices, err := dao.GetServices("status", serviceStoppedStatus)
	if err != nil {
		logs.Error("Failed to get the services when get NodePorts.")
		return nil, err
	}
	ports := []int32{}
	type config struct {
		Spec struct {
			Ports []struct {
				NodePort int `yaml:"nodePort,flow"`
			}
		}
	}
	for _, serviceConfig := range stoppedServices {
		err := utils.UnmarshalYamlData([]byte(serviceConfig.ServiceYaml), &config{}, func(in interface{}) error {
			if c, ok := in.(*config); ok {
				for _, port := range c.Spec.Ports {
					ports = append(ports, int32(port.NodePort))
				}
			}
			return nil
		})
		if err != nil {
			logs.Error("Failed to Unmarshal data of the service.")
			return nil, err
		}
	}
	return ports, nil
}

// DeployStatefulSet is to deploy the statefulset service in k8s
func DeployStatefulSet(serviceConfig *model.ConfigServiceStep, registryURI string) (*DeployStatefulSetInfo, error) {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)
	statefulsetConfig := MarshalStatefulSet(serviceConfig, registryURI)
	//logs.Debug("Marshaled deployment: ", deploymentConfig)
	statefulsetInfo, statefulsetFileInfo, err := cli.AppV1().StatefulSet(serviceConfig.ProjectName).Create(statefulsetConfig)
	if err != nil {
		logs.Error("Deploy statefulset object of %s failed. error: %+v\n", serviceConfig.ServiceName, err)
		return nil, err
	}
	logs.Debug("Created statefulset: ", statefulsetInfo)
	svcConfig := MarshalService(serviceConfig)
	serviceInfo, serviceFileInfo, err := cli.AppV1().Service(serviceConfig.ProjectName).Create(svcConfig)
	if err != nil {
		cli.AppV1().StatefulSet(serviceConfig.ProjectName).Delete(serviceConfig.ServiceName)
		logs.Error("Deploy service object of %s failed. error: %+v\n", serviceConfig.ServiceName, err)
		return nil, err
	}

	return &DeployStatefulSetInfo{
		Service:             serviceInfo,
		ServiceFileInfo:     serviceFileInfo,
		StatefulSet:         statefulsetInfo,
		StatefulSetFileInfo: statefulsetFileInfo,
	}, nil
}
