package service

import (
	"git/inspursoft/board/src/common/model"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var Copy = model.CopyStruct{
	CopyFrom: "from",
	CopyTo:   "to",
}

var Env = model.EnvStruct{
	EnvName:  "key",
	EnvValue: "value",
}

var Dockerfile model.Dockerfile
var imageConfig model.ImageConfig

func TestCheckDockerfileConfig(t *testing.T) {
	err := CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check nil dockerfile config should fail")
	} else {
		t.Log("Check nil config fail as expect.")
	}

	imageConfig.ImageDockerfile.Base = "Test:1.0"
	err = CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check dockerfile with upper charactor in base showld fail")
	} else {
		t.Log("Check dockerfile with upper charactor in base successfully.")
	}

	imageConfig.ImageDockerfile.Base = "test:1.0"
	imageConfig.ImageDockerfile.EntryPoint = "a\nb"
	err = CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check dockerfile with enter in Entrypoint showld fail")
	} else {
		t.Log("Check dockerfile with enter in Entrypoint successfully.")
	}

	imageConfig.ImageDockerfile.EntryPoint = "ab"
	imageConfig.ImageDockerfile.ExposePort = append(imageConfig.ImageDockerfile.ExposePort, "0s")
	err = CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check dockerfile port showld fail")
	} else {
		t.Log("Check dockerfile port successfully.")
	}

	imageConfig.ImageDockerfile.ExposePort = nil
	imageConfig.ImageDockerfile.ExposePort = append(imageConfig.ImageDockerfile.ExposePort, "8888")
	imageConfig.ImageDockerfile.Volume = append(imageConfig.ImageDockerfile.Volume, "volume")
	//imageConfig.ImageDockerfile.Copy = append(imageConfig.ImageDockerfile.Copy, Copy)
	imageConfig.ImageDockerfile.RUN = append(imageConfig.ImageDockerfile.RUN, "run")
	imageConfig.ImageDockerfile.EnvList = append(imageConfig.ImageDockerfile.EnvList, Env)
	err = CheckDockerfileConfig(&imageConfig)
	if err != nil {
		t.Errorf("Check dockerfile error: %+v", err)
	} else {
		t.Log("Check dockerfile successfully.")
	}
}

func TestBuildDockerfile(t *testing.T) {
	imageConfig.ImageDockerfile.Base = "test:1.0"
	imageConfig.ImageDockerfile.Copy = append(imageConfig.ImageDockerfile.Copy, Copy)
	imageConfig.ImageDockerfilePath = "path"
	err := BuildDockerfile(imageConfig)
	if err != nil {
		t.Errorf("Build dockerfile fail: %+v", err)
	} else {
		t.Log("Build dockerfile successfully.")
	}
}

func TestGetDockerfileInfo(t *testing.T) {
	dockerfile, err := GetDockerfileInfo("path")
	if err != nil {
		t.Errorf("Get dockerfile info error: %+v", err)
	}
	if dockerfile.Base == imageConfig.ImageDockerfile.Base &&
		dockerfile.EntryPoint == imageConfig.ImageDockerfile.EntryPoint {
		t.Log("Get dockerfile info successfully.")
	}
}

func TestImageConfigClean(t *testing.T) {
	err := ImageConfigClean("path")
	if err != nil {
		t.Errorf("Clean config error: %+v", err)
	} else {
		t.Log("Clean config successfully.")
	}
}
