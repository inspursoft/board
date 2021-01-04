package service_test

import (
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"testing"

	"github.com/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var path = "./"
var configServiceStep = model.ConfigServiceStep{
	ProjectID:   1,
	Instance:    1,
	ServiceName: "testService",
	ContainerList: []model.Container{
		{
			Name: "testService",
			Image: model.ImageIndex{
				ImageName:   "library/demooanginx",
				ImageTag:    "1.0",
				ProjectName: "library",
			},
		},
	},
	ExternalServiceList: []model.ExternalService{
		{
			ContainerName: "testService",
			NodeConfig: model.NodeType{
				TargetPort: 80,
				NodePort:   32080,
			},
		},
	},
}

var configStatefulSet = model.ConfigServiceStep{
	ProjectID:   1,
	ProjectName: "library",
	Instance:    1,
	ServiceName: "unitteststatefulset001",
	ServiceType: model.ServiceTypeStatefulSet,
	// ClusterIP:   "None",
	ContainerList: []model.Container{
		{
			Name: "nginx",
			Image: model.ImageIndex{
				ImageName:   "library/nginx",
				ImageTag:    "1.11.5",
				ProjectName: "library",
			},
		},
	},
	ExternalServiceList: []model.ExternalService{
		{
			ContainerName: "nginx",
			NodeConfig: model.NodeType{
				TargetPort: 80,
				Port:       80,
			},
		},
	},
}

var statufulsetName = "unitteststatefulset001"

// TODO: unit test case later
// TestDeployStatefulSet
func TestDeployStatefulSet(t *testing.T) {
	assert := assert.New(t)
	t.Log("Check KubeMaster")
	masterIP := utils.GetStringValue("KUBE_MASTER_IP")
	logs.Info("KUBE_MASTER_IP %s", masterIP)

	registryURI := utils.GetStringValue("REGISTRY_BASE_URI")
	logs.Info("REGISTRY_URI %s", registryURI)

	deployStatefulSetInfo, err := service.DeployStatefulSet(&configStatefulSet, registryURI)
	assert.Nil(err, "Failed, err when create test StatefulSet")
	assert.Equal(statufulsetName, deployStatefulSetInfo.Service.Name, "Failed to create StatefulSet")
	logs.Info("Created statefulset %v %v", deployStatefulSetInfo.Service, deployStatefulSetInfo.StatefulSet)

	//clean test
	t.Log("Clean TestDeployStatefulSet")
	cleanStatefulSet(configStatefulSet.ProjectName, configStatefulSet.ServiceName)
	t.Log("Tested TestDeployStatefulSet")
}

// clean test
func cleanStatefulSet(projectName string, serviceName string) {
	logs.Info("cleanStatefulSet %s %s", projectName, serviceName)
	err := service.StopStatefulSetK8s(&model.ServiceStatus{
		Name:        serviceName,
		ProjectName: projectName,
	})
	if err != nil {
		logs.Info("cleanStatefulSet failed %v", err)
		return
	}
	logs.Info("cleaned StatefulSet")
}
