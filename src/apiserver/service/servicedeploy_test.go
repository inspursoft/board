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
