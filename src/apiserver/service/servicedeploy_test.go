package service

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
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

func TestAssembleDeploymentYaml(t *testing.T) {
	assert := assert.New(t)
	err := AssembleDeploymentYaml(&configServiceStep, path)
	assert.Nil(err, "Error occurred while testing AssembleDeploymentYaml.")
	deleteFile(filepath.Join(loadPath, deploymentFilename))
}

func TestAssembleServiceYaml(t *testing.T) {
	assert := assert.New(t)
	err := AssembleServiceYaml(&configServiceStep, path)
	assert.Nil(err, "Error occurred while testing AssembleServiceYaml.")
	deleteFile(filepath.Join(loadPath, serviceFilename))
}

func deleteFile(file string) {
	err := os.Remove(file)
	if err != nil {
		logs.Error("Error occurred while removing file %s\n", file)
	}
}
