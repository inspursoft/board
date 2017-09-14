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
			errString := fmt.Sprintf(`string "%s" has upper charactor`, node)
			return errors.New(errString)
		}
	}
	return nil
}

func checkStringHasEnter(str ...string) error {
	for _, node := range str {
		isMatch, err := regexp.MatchString(`^\s*\n?(?:.*[^\n])*\n?\s*$`, node)
		if err != nil {
			return err
		}
		if !isMatch {
			errString := fmt.Sprintf(`string "%s" has enter charactor`, node)
			return errors.New(errString)
		}
	}
	return nil
}

func fixStructEmptyIssue(obj interface{}) {
	if f, ok := obj.(*[]string); ok {
		if len(*f) == 1 && len((*f)[0]) == 0 {
			*f = nil
		}
		return
	}
	if f, ok := obj.(*[]model.CopyStruct); ok {
		if len(*f) == 1 && len((*f)[0].CopyFrom) == 0 && len((*f)[0].CopyTo) == 0 {
			*f = nil
		}
		return
	}
	if f, ok := obj.(*[]model.EnvStruct); ok {
		if len(*f) == 1 && len((*f)[0].EnvName) == 0 && len((*f)[0].EnvValue) == 0 {
			*f = nil
		}
		return
	}
	if f, ok := obj.(*[]int); ok {
		if len(*f) == 1 && (*f)[0] == 0 {
			*f = nil
		}
	}
	return
}

func changeDockerfileStructItem(dockerfile *model.Dockerfile) {
	dockerfile.Base = strings.TrimSpace(dockerfile.Base)
	dockerfile.Author = strings.TrimSpace(dockerfile.Author)
	dockerfile.EntryPoint = strings.TrimSpace(dockerfile.EntryPoint)
	dockerfile.Command = strings.TrimSpace(dockerfile.Command)

	for num, node := range dockerfile.Volume {
		dockerfile.Volume[num] = strings.TrimSpace(node)
	}
	fixStructEmptyIssue(&dockerfile.Volume)

	for num, node := range dockerfile.Copy {
		dockerfile.Copy[num].CopyFrom = strings.TrimSpace(node.CopyFrom)
		dockerfile.Copy[num].CopyTo = strings.TrimSpace(node.CopyTo)
	}
	fixStructEmptyIssue(&dockerfile.Copy)

	for num, node := range dockerfile.RUN {
		dockerfile.RUN[num] = strings.TrimSpace(node)
	}
	fixStructEmptyIssue(&dockerfile.RUN)

	for num, node := range dockerfile.EnvList {
		dockerfile.EnvList[num].EnvName = strings.TrimSpace(node.EnvName)
		dockerfile.EnvList[num].EnvValue = strings.TrimSpace(node.EnvValue)
	}
	fixStructEmptyIssue(&dockerfile.EnvList)

	fixStructEmptyIssue(&dockerfile.ExposePort)
}

func changeImageConfigStructItem(reqImageConfig *model.ImageConfig) {
	reqImageConfig.ImageName = strings.TrimSpace(reqImageConfig.ImageName)
	reqImageConfig.ImageTag = strings.TrimSpace(reqImageConfig.ImageTag)
	reqImageConfig.ProjectName = strings.TrimSpace(reqImageConfig.ProjectName)
	reqImageConfig.ImageTemplate = strings.TrimSpace(reqImageConfig.ImageTemplate)
	reqImageConfig.ImageDockerfilePath = strings.TrimSpace(reqImageConfig.ImageDockerfilePath)

	changeDockerfileStructItem(&reqImageConfig.ImageDockerfile)
}

func CheckDockerfileConfig(config *model.ImageConfig) error {
	changeImageConfigStructItem(config)

	if len(config.ImageDockerfile.Base) == 0 {
		return errors.New("Baseimage in dockerfile should not be empty")
	}

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

func ImageConfigClean(path string) error {
	//remove Tag dir
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	//remove Image dir
	if fi, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		return nil
	} else if !fi.IsDir() {
		errMsg := fmt.Sprintf(`%s is not dir`, filepath.Dir(path))
		return errors.New(errMsg)
	}

	parent, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer parent.Close()

	_, err = parent.Readdirnames(1)
	if err == io.EOF {
		return os.RemoveAll(filepath.Dir(path))
	}

	return nil
}
