package service

import (
	"errors"
	"git/inspursoft/board/src/common/model"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var dockerfilePath = filepath.Join("/", "repos", "board_repo", "library")
var dockerTemplatePath = "templates"
var dockerfileName = "Dockerfile"
var templateNameDefault = "dockerfile-template"
var copyFromPath = "upload"

func SetDockerfilePath(path string) {
	dockerfilePath = path
}

func GetDockerfilePath() string {
	return dockerfilePath
}

func SetCopyFromPath(path string) {
	copyFromPath = path
}

func GetCopyFromPath() string {
	return copyFromPath
}

func str2execform(str string) string {
	sli := strings.Split(str, " ")
	for num, node := range sli {
		sli[num] = "\"" + node + "\""
	}
	return strings.Join(sli, ", ")
}

func CheckDockerfileConfig(config model.ImageConfig) error {
	isMatch, err := regexp.MatchString("[A-Z]", config.ImageDockerfile.Base)
	if err != nil {
		return err
	}
	if isMatch {
		return errors.New("dockerfile's baseimage name shouldn't contain upper character")
	}

	isMatch, err = regexp.MatchString("[A-Z]", config.ImageName)
	if err != nil {
		return err
	}
	if isMatch {
		return errors.New("docker image's name shouldn't contain upper character")
	}

	isMatch, err = regexp.MatchString("[A-Z]", config.ImageTag)
	if err != nil {
		return err
	}
	if isMatch {
		return errors.New("docker image's tag shouldn't contain upper character")
	}

	return nil
}

func changeDockerfileStructPath(dockerfile model.Dockerfile) error {
	if len(GetCopyFromPath()) == 0 {
		return nil
	}

	for num, node := range dockerfile.Copy {
		dockerfile.Copy[num].CopyFrom = filepath.Join(GetCopyFromPath(), node.CopyFrom)
	}

	return nil
}

func BuildDockerfile(reqImageConfig model.ImageConfig, wr io.Writer) error {
	var templatename string

	if err := changeDockerfileStructPath(reqImageConfig.ImageDockerfile); err != nil {
		return err
	}

	if len(reqImageConfig.ImageTemplate) != 0 {
		templatename = reqImageConfig.ImageTemplate
	} else {
		templatename = templateNameDefault
	}

	tmpl, err := template.New(templatename).Funcs(template.FuncMap{"str2exec": str2execform}).ParseFiles(filepath.Join(dockerTemplatePath, templatename))
	if err != nil {
		return err
	}

	if wr != nil {
		if err = tmpl.Execute(wr, reqImageConfig.ImageDockerfile); err != nil {
			return err
		}
		return nil
	}

	if fi, err := os.Stat(GetDockerfilePath()); os.IsNotExist(err) {
		if err := os.MkdirAll(GetDockerfilePath(), 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return errors.New("Dockerfile path is not dir")
	}

	dockerfile, err := os.OpenFile(filepath.Join(GetDockerfilePath(), dockerfileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dockerfile.Close()

	err = tmpl.Execute(dockerfile, reqImageConfig.ImageDockerfile)
	if err != nil {
		return err
	}

	return nil
}
