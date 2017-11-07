package service

import (
	"git/inspursoft/board/src/common/model"
	"testing"

	"k8s.io/client-go/pkg/api/v1"
)

const (
	filePath                 = "./test001"
	targetPath               = "/tmp"
	proName                  = "oabusi"
	illegalProName           = "OAbusi"
	baseImage                = "oabusi:v1.0"
	proID                    = 1
	proPhase                 = "yunoa deployment"
	proPort            int32 = 8080
	proNodePort              = 30080
	illegalProNodePort       = 20080
	illegalCommamd           = "ECHO HELLO"
	illegalDir               = "/HOME"
	illegalTarget            = "/Tmp"
	serviceFile              = "yunoaService.yaml"
	deploymentFile           = "yunoaDeployment.yaml"
)

var ServiceConfig = model.ServiceConfig2{
	Project:    yunoaProject,
	Deployment: yunoaDeployment,
	Service:    yunoaService,
}

var yunoaProject = model.ProjectInfo{
	ServiceID:   proID,
	ProjectID:   proID,
	ProjectName: proName,
}

var yunoaDeployment = v1.ReplicationController{
	ObjectMeta: yunoaObjectMeta,
	Spec:       yunoaSpec,
}

var yunoaObjectMeta = v1.ObjectMeta{
	Name: proName,
}

var replicase int32 = 1
var yunoaSpec = v1.ReplicationControllerSpec{
	Replicas: &replicase,
	Template: &yunoaTemplate,
}

var yunoaTemplate = v1.PodTemplateSpec{
	ObjectMeta: yunoaObjectMeta,
	Spec:       yunoaPodSpec,
}

var yunoaPodSpec = v1.PodSpec{
	Containers: yunoaContainers,
}

var yunoaContainers = []v1.Container{
	v1.Container{
		Name:  proName,
		Image: baseImage,
		Ports: []v1.ContainerPort{
			v1.ContainerPort{ContainerPort: proPort}},
		Env:          yunoaEnv,
		VolumeMounts: yunoaVolume,
	},
}

var yunoaEnv = []v1.EnvVar{
	v1.EnvVar{
		Name:  proName,
		Value: proPhase,
	},
}

var yunoaVolume = []v1.VolumeMount{
	v1.VolumeMount{
		MountPath: targetPath,
	},
}

var yunoaService = v1.Service{
	ObjectMeta: yunoaObjectMeta,
	Spec:       yunoaServiceSpec,
}
var yunoaServiceSpec = v1.ServiceSpec{}

var ServiceConfigNil = model.ServiceConfig2{
	Project: yunoaProject2,
}

var yunoaProject2 = model.ProjectInfo{
	ServiceID:   proID,
	ProjectID:   proID,
	ProjectName: proName,
}

func TestCheckDeploymentPath(t *testing.T) {
	err := CheckDeploymentPath(filePath)
	if err != nil {
		t.Errorf("Error occurred while test CheckDeploymentPath: %+v\n", err)
	}
	if err == nil {
		t.Log("CheckDeploymentPath is ok.\n")
	}

}

func TestCheckReqPara(t *testing.T) {
	err := CheckReqPara(ServiceConfig)
	if err != nil {
		t.Errorf("Error occurred while test CheckReqPara: %+v\n", err)
	}
	if err == nil {
		t.Log("CheckReqPara is ok.\n")
	}

}

func TestCheckDeploymentPara(t *testing.T) {
	var err error

	//for null name
	err = CheckDeploymentPara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckDeploymentPara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckDeploymentPara is ok.\n")
	}

	//for replicas <1
	var illegalreplicas int32 = 0
	var replicas int32 = 1
	ServiceConfigNil.Deployment.ObjectMeta.Name = proName
	ServiceConfigNil.Deployment.Spec.Replicas = &illegalreplicas
	err = CheckDeploymentPara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckDeploymentPara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckDeploymentPara is ok.\n")
	}

	//for upper deployment name
	ServiceConfigNil.Deployment.ObjectMeta.Name = proName
	ServiceConfigNil.Deployment.Spec.Replicas = &replicas
	err = CheckDeploymentPara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckDeploymentPara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckDeploymentPara is ok.\n")
	}
}

func TestCheckServicePara(t *testing.T) {
	err := CheckServicePara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckServicePara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckServicePara is ok.\n")
	}

	ServiceConfigNil.Service.ObjectMeta.Name = proName
	var servicePort = v1.ServicePort{
		Port:     proPort,
		NodePort: illegalProNodePort,
	}
	ServiceConfigNil.Service.Spec.Ports = append(ServiceConfigNil.Service.Spec.Ports, servicePort)
	err = CheckServicePara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckServicePara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckServicePara is ok.\n")
	}

}
