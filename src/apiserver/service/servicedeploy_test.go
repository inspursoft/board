package service_test

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"testing"

	"github.com/astaxie/beego/logs"
)

var path = "./"
var configServiceStep = model.ConfigServiceStep{
	ProjectID:   1,
	Instance:    1,
	ServiceName: "testService",
	ContainerList: []model.Container{
		model.Container{
			Name: "testService",
			Image: model.ImageIndex{
				ImageName:   "library/demooanginx",
				ImageTag:    "1.0",
				ProjectName: "library",
			},
		},
	},
	ExternalServiceList: []model.ExternalService{
		model.ExternalService{
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
	Instance:    1,
	ServiceName: "teststatefulset001",
	ServiceType: model.ServiceTypeStatefulSet,
	ClusterIP:   "None",
	ContainerList: []model.Container{
		model.Container{
			Name: "nginx",
			Image: model.ImageIndex{
				ImageName:   "library/nginx",
				ImageTag:    "1.11.5",
				ProjectName: "library",
			},
		},
	},
	ExternalServiceList: []model.ExternalService{
		model.ExternalService{
			ContainerName: "nginx",
			NodeConfig: model.NodeType{
				TargetPort: 80,
				Port:       80,
			},
		},
	},
}

// TODO: unit test case later
// TestDeployStatefulSet
func TestDeployStatefulSet(t *testing.T) {
	assert := assert.New(t)
	//deployStatefulSetInfo, err := service.DeployStatefulSet(&configServiceStep, registryBaseURI())
	//assert.Nil(err, "Failed, err when create test image.")
	//assert.NotEqual(0, id, "Failed to assign a image id")
	//testImageid = id
	//t.Log(deployStatefulSetInfo)
	t.Log("Test KubeMaster")
	masterIP = utils.GetConfig("KUBE_MASTER_IP")
	registryIP = utils.GetConfig("REGISTRY_IP")
	t.Log("KUBE_MASTER_IP %s  REGISTRY_IP %s", masterIP, registryIP)

	t.Log("Tested TestDeployStatefulSet")
}

// TODO: clean test
func cleanStatefulSetByID(id int64) {
	logs.Debug("cleanStatefulSetByID")
}
