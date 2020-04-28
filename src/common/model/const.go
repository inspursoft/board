package model

import (
	"time"
)

const (
	ProjectAdmin = int64(iota + 1)
	Developer
	Visitor
	ServiceStart = int64(iota + 1)
	ServiceStop
)

const (
	Preparing = iota
	Running
	Stopped
	Uncompleted
	Warning
	Deploying
	Completed
	Failed
)

const (
	DockerfileName         = "Dockerfile"
	DeploymentFilename     = "deployment.yaml"
	StatefulsetFilename    = "statefulset.yaml"
	ServiceFilename        = "service.yaml"
	RollingUpdateFilename  = "rollingUpdateDeployment.yaml"
	DeploymentTestFilename = "testdeployment.yaml"
	ServiceTestFilename    = "testservice.yaml"

	Apiheader        = "Content-Type: application/yaml"
	DeploymentAPI    = "/apis/extensions/v1beta1/namespaces/"
	ServiceAPI       = "/api/v1/namespaces/"
	Test             = "test"
	ServiceNamespace = "default" //TODO create in project post
	K8sServices      = "kubernetes"
	DeploymentType   = "deployment"
	ServiceType      = "service"
	StatefulsetType  = "statefulset"
	StartingDuration = 300 * time.Second //300 seconds
)
