package service

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/common/model"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var dockerTemplatePath = "templates"
var dockerfileName = "Dockerfile"
var templateNameDefault = "dockerfile-template"

func str2execform(str string) string {
	sli := strings.Split(str, " ")
	for num, node := range sli {
		sli[num] = "\"" + node + "\""
	}
	return strings.Join(sli, ", ")
}

func checkStringHasUpper(str ...string) error {
	for _, node := range str {
		isMatch, err := regexp.MatchString("[A-Z]", node)
		if err != nil {
			return err
		}
		if isMatch {
			errString := fmt.Sprintf("string \"%s\" has upper charactor", node)
			return errors.New(errString)
		}
	}
	return nil
}

func checkStringHasEnter(str ...string) error {
	for _, node := range str {
		isMatch, err := regexp.MatchString("[\n][^$]", node)
		if err != nil {
			return err
		}
		if isMatch {
			errString := fmt.Sprintf("string \"%s\" has enter charactor", node)
			return errors.New(errString)
		}
	}
	return nil
}

func CheckDockerfileConfig(config model.ImageConfig) error {
	err := checkStringHasUpper(config.ImageDockerfile.Base, config.ImageName, config.ImageTag)
	if err != nil {
		return err
	}

	err = checkStringHasEnter(config.ImageDockerfile.EntryPoint, config.ImageDockerfile.Command)
	if err != nil {
		return err
	}
	return nil
}

func BuildDockerfile(reqImageConfig model.ImageConfig, wr ...io.Writer) error {
	var templatename string

	if len(reqImageConfig.ImageTemplate) != 0 {
		templatename = reqImageConfig.ImageTemplate
	} else {
		templatename = templateNameDefault
	}

	tmpl, err := template.New(templatename).Funcs(template.FuncMap{"str2exec": str2execform}).ParseFiles(filepath.Join(dockerTemplatePath, templatename))
	if err != nil {
		return err
	}

	if len(wr) != 0 {
		if err = tmpl.Execute(wr[0], reqImageConfig.ImageDockerfile); err != nil {
			return err
		}
		return nil
	}

	if fi, err := os.Stat(reqImageConfig.ImageDockerfilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(reqImageConfig.ImageDockerfilePath, 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return errors.New("Dockerfile path is not dir")
	}

	dockerfile, err := os.OpenFile(filepath.Join(reqImageConfig.ImageDockerfilePath, dockerfileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
