package service

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"testing"

	yaml "git/inspursoft/board/src/common/model/yaml"
)

const (
	filePath           = "./test001"
	targetPath         = "/tmp"
	proName            = "oabusi"
	illegalProName     = "OAbusi"
	baseImage          = "oabusi:v1.0"
	proID              = 1
	proPhase           = "yunoa deployment"
	replicas           = 1
	illegalreplicas    = 0
	proPort            = 8080
	proNodePort        = 30080
	illegalProNodePort = 20080
	illegalCommamd     = "ECHO HELLO"
	illegalDir         = "/HOME"
	illegalTarget      = "/Tmp"
	serviceFile        = "yunoaService.yaml"
	deploymentFile     = "yunoaDeployment.yaml"
)

var ServiceConfig = model.ServiceConfig{
	ServiceID:      proID,
	ProjectID:      proID,
	ProjectName:    proName,
	Phase:          proPhase,
	DeploymentYaml: yunoaDeployment,
	ServiceYaml:    yunoaService,
}
var yunoaDeployment = yaml.Deployment{
	Name:          proName,
	Replicas:      replicas,
	ContainerList: yunoaContainers,
}

var yunoaContainers = []yaml.Container{
	yaml.Container{
		Name:      proName,
		BaseImage: baseImage,
		Ports:     []int{proPort},
		Envs:      yunoaEnv,
		Volumes:   yunoaVolume,
	},
}

var yunoaEnv = []yaml.Env{
	yaml.Env{
		Name:  proName,
		Value: proPhase,
	},
}

var yunoaVolume = []yaml.Volume{
	yaml.Volume{
		Dir:               filePath,
		TargetStorageName: targetPath,
	},
}

var yunoaService = yaml.Service{
	Name:      proName,
	External:  ExternalStructList,
	Selectors: []string{proName},
}

var ExternalStructList = []yaml.ExternalStruct{
	yaml.ExternalStruct{
		ContainerPort: proPort,
		NodePort:      proNodePort,
	},
}

var ServiceConfigNil = model.ServiceConfig{
	ServiceID:   proID,
	ProjectID:   proID,
	ProjectName: proName,
	Phase:       proPhase,
}

func TestGetDeploymentPath(t *testing.T) {
	SetDeploymentPath(filePath)
	path := GetDeploymentPath()
	if path != filePath {
		t.Errorf("Error occurred while test GetDeploymentPath,path is wt001, get path is %+v\n", path)
	}
	if path == filePath {
		t.Log("GetDeploymentPath is ok.\n")
	}
	err := os.RemoveAll(filePath)
	if err != nil {
		t.Log("delet dir error:", err)
	}
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
	ServiceConfigNil.DeploymentYaml.Name = proName
	ServiceConfigNil.DeploymentYaml.Replicas = illegalreplicas
	err = CheckDeploymentPara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckDeploymentPara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckDeploymentPara is ok.\n")
	}

	//for upper deployment name
	ServiceConfigNil.DeploymentYaml.Name = illegalProName
	ServiceConfigNil.DeploymentYaml.Replicas = replicas
	err = CheckDeploymentPara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckDeploymentPara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckDeploymentPara is ok.\n")
	}

	//for volumes in complate
	ServiceConfigNil.DeploymentYaml.Name = proName
	ServiceConfigNil.DeploymentYaml.ContainerList = append(ServiceConfigNil.DeploymentYaml.ContainerList,
		yaml.Container{
			Name:      proName,
			BaseImage: baseImage,
			Ports:     []int{proPort},
		})
	ServiceConfigNil.DeploymentYaml.ContainerList[0].Volumes = append(ServiceConfigNil.DeploymentYaml.ContainerList[0].Volumes,
		yaml.Volume{Dir: illegalDir, TargetStorageName: illegalTarget})
	err = CheckDeploymentPara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckDeploymentPara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckDeploymentPara is ok.\n")
	}

	//for command in complate
	ServiceConfigNil.DeploymentYaml.ContainerList[0].Command = []string{illegalCommamd}
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

	ServiceConfigNil.ServiceYaml.Name = proName
	var ExternalStrucTmp = yaml.ExternalStruct{
		ContainerPort: proPort,
		NodePort:      illegalProNodePort,
	}
	ServiceConfigNil.ServiceYaml.External = append(ServiceConfigNil.ServiceYaml.External, ExternalStrucTmp)
	err = CheckServicePara(ServiceConfigNil)
	if err == nil {
		t.Errorf("Error occurred while test CheckServicePara: %+v\n", err)
	}
	if err != nil {
		t.Log("CheckServicePara is ok.\n")
	}

}

func TestBuildServiceYaml(t *testing.T) {
	err := BuildServiceYaml(ServiceConfig, serviceFile)
	if err != nil {
		t.Errorf("Error occurred while test BuildServiceYaml: %+v\n", err)
	}
	if err == nil {
		t.Log("BuildServiceYaml is ok.\n")
	}
	err = os.RemoveAll(filePath)
	if err != nil {
		t.Log("delet dir error:", err)
	}

}

func TestBuildDeploymentYaml(t *testing.T) {
	ServiceConfig.DeploymentYaml.VolumeList = append(ServiceConfig.DeploymentYaml.VolumeList, yaml.NFSVolume{Name: "/root", ServerName: "/home"})
	err := BuildDeploymentYaml(ServiceConfig, deploymentFile)
	if err != nil {
		t.Errorf("Error occurred while test BuildDeploymentYaml: %+v\n", err)
	}
	if err == nil {
		t.Log("BuildDeploymentYaml is ok.\n")
	}
	err = os.RemoveAll(filePath)
	if err != nil {
		t.Log("delet dir error:", err)
	}

}
