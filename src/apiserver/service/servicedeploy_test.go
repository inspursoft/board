package service_test

import (
	"git/inspursoft/board/src/common/model"
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

// TODO: unit test case later
// TestDeployStatefulSet
func TestDeployStatefulSet(t *testing.T) {
	//assert := assert.New(t)
	//id, err := service.DeployStatefulSet(&configServiceStep, )
	//assert.Nil(err, "Failed, err when create test image.")
	//assert.NotEqual(0, id, "Failed to assign a image id")
	//testImageid = id
	t.Log(configServiceStep)
}

// TODO: clean test
func cleanStatefulSetByID(id int64) {

}
