package service

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"

	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/astaxie/beego/logs"
)

const (
	dockerTemplatePath  = "templates"
	templateNameDefault = "dockerfile-template"
)

func str2execform(str string) string {
	sli := strings.Split(str, " ")
	for num, node := range sli {
		sli[num] = "\"" + node + "\""
	}
	return strings.Join(sli, ", ")
}

func exec2str(str string) string {
	line := strings.TrimSpace(str)
	line = strings.TrimLeft(strings.TrimRight(line, "]"), "[")
	split := strings.Split(line, ",")
	for num, node := range split {
		node = strings.TrimSpace(node)
		split[num] = strings.TrimLeft(strings.TrimRight(node, "\""), "\"")
	}
	return strings.Join(split, " ")
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
			*f = make([]string, 0, 0)
		}
		return
	}
	if f, ok := obj.(*[]model.CopyStruct); ok {
		if len(*f) == 1 && len((*f)[0].CopyFrom) == 0 && len((*f)[0].CopyTo) == 0 {
			*f = make([]model.CopyStruct, 0, 0)
		}
		return
	}
	if f, ok := obj.(*[]model.EnvStruct); ok {
		if len(*f) == 1 && len((*f)[0].EnvName) == 0 && len((*f)[0].EnvValue) == 0 {
			*f = make([]model.EnvStruct, 0, 0)
		}
		return
	}
	return
}

func changeDockerfileStructItem(dockerfile *model.Dockerfile, relPath string) {
	dockerfile.Base = strings.TrimSpace(dockerfile.Base)
	dockerfile.Author = strings.TrimSpace(dockerfile.Author)
	dockerfile.EntryPoint = strings.TrimSpace(dockerfile.EntryPoint)
	dockerfile.Command = strings.TrimSpace(dockerfile.Command)

	for num, node := range dockerfile.Volume {
		dockerfile.Volume[num] = strings.TrimSpace(node)
	}
	fixStructEmptyIssue(&dockerfile.Volume)

	for num, node := range dockerfile.Copy {
		fromPath := filepath.Join(relPath, strings.TrimSpace(node.CopyFrom))
		dockerfile.Copy[num].CopyFrom = fromPath
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

	for num, node := range dockerfile.ExposePort {
		dockerfile.ExposePort[num] = strings.TrimSpace(node)
	}
	fixStructEmptyIssue(&dockerfile.ExposePort)
}

func changeImageConfigStructItem(reqImageConfig *model.ImageConfig) {
	reqImageConfig.ImageName = strings.TrimSpace(reqImageConfig.ImageName)
	reqImageConfig.ImageTag = strings.TrimSpace(reqImageConfig.ImageTag)
	reqImageConfig.ProjectName = strings.TrimSpace(reqImageConfig.ProjectName)
	reqImageConfig.ImageTemplate = strings.TrimSpace(reqImageConfig.ImageTemplate)
	reqImageConfig.ImageDockerfilePath = strings.TrimSpace(reqImageConfig.ImageDockerfilePath)
}

func CheckDockerfileItem(dockerfile *model.Dockerfile, relPath string) error {
	changeDockerfileStructItem(dockerfile, relPath)

	if len(dockerfile.Base) == 0 {
		return errors.New("Baseimage in dockerfile should not be empty")
	}

	if err := checkStringHasUpper(dockerfile.Base); err != nil {
		return err
	}

	if err := checkStringHasEnter(dockerfile.EntryPoint, dockerfile.Command); err != nil {
		return err
	}

	for _, node := range dockerfile.ExposePort {
		if _, err := strconv.Atoi(node); err != nil {
			return err
		}
	}

	return nil
}

func CheckDockerfileConfig(config *model.ImageConfig) error {
	changeImageConfigStructItem(config)

	if err := checkStringHasUpper(config.ImageName, config.ImageTag); err != nil {
		return err
	}

	return CheckDockerfileItem(&config.ImageDockerfile, "upload")
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

	dockerfileName := ResolveDockerfileName(reqImageConfig.ImageName, reqImageConfig.ImageTag)
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

func GetDockerfileInfo(dockerfilePath, imageName, tag string) (*model.Dockerfile, error) {
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		logs.Error("Failed to find Dockerfile path: %+v", err)
		return nil, err
	}

	var Dockerfile model.Dockerfile
	var fulline string
	dockerfileName := ResolveDockerfileName(imageName, tag)
	dockerfile, err := os.Open(filepath.Join(dockerfilePath, dockerfileName))
	if err != nil {
		logs.Error("Failed to find Dockerfile : %+v with name: %s", err, dockerfileName)
		return nil, err
	}
	defer dockerfile.Close()

	scanner := bufio.NewScanner(dockerfile)
	for scanner.Scan() {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), "#") {
			continue
		}
		fulline += string(scanner.Text())
		if strings.HasSuffix(scanner.Text(), "\\") {
			fulline = fulline[:len(fulline)-1]
			continue
		}
		split := strings.SplitN(strings.TrimSpace(fulline), " ", 2)
		fulline = ""

		// ignore empty line and lines with only one field
		if len(split) < 2 {
			continue
		}
		split[1] = strings.TrimSpace(split[1])
		switch split[0] {
		case "FROM":
			Dockerfile.Base = split[1]
		case "MAINTAINER":
			Dockerfile.Author = split[1]
		case "VOLUME":
			Dockerfile.Volume = append(Dockerfile.Volume, exec2str(split[1]))
		case "COPY":
			{
				var node model.CopyStruct
				var copyfrom, copyto string
				copystring := exec2str(split[1])
				split_copy := strings.Split(strings.TrimSpace(copystring), " ")
				copyfrom = strings.Join(split_copy[:len(split_copy)-1], " ")
				copyto = split_copy[len(split_copy)-1]
				node.CopyFrom = copyfrom
				node.CopyTo = copyto
				Dockerfile.Copy = append(Dockerfile.Copy, node)
			}
		case "RUN":
			Dockerfile.RUN = append(Dockerfile.RUN, split[1])
		case "ENTRYPOINT":
			Dockerfile.EntryPoint = exec2str(split[1])
		case "CMD":
			Dockerfile.Command = exec2str(split[1])
		case "ENV":
			{
				var node model.EnvStruct
				envstring := split[1]
				split_env := strings.SplitN(envstring, " ", 2)
				node.EnvName = split_env[0]
				node.EnvValue = strings.TrimSpace(split_env[1])
				Dockerfile.EnvList = append(Dockerfile.EnvList, node)
			}
		case "EXPOSE":
			Dockerfile.ExposePort = append(Dockerfile.ExposePort, split[1])
		}
	}

	return &Dockerfile, nil
}

// Image in database
func CreateImage(image model.Image) (int64, error) {
	imageID, err := dao.AddImage(image)
	if err != nil {
		return 0, err
	}
	return imageID, nil
}

func GetImage(image model.Image, selectedFields ...string) (*model.Image, error) {
	m, err := dao.GetImage(image, selectedFields...)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetImageByName(imageName string) (*model.Image, error) {
	return GetImage(model.Image{ImageName: imageName, ImageDeleted: 0}, "name")
}

func UpdateImage(image model.Image, fieldNames ...string) (bool, error) {
	if image.ImageID == 0 {
		return false, errors.New("no Image ID provided")
	}
	_, err := dao.UpdateImage(image, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetImageTag(imageTag model.ImageTag, selectedFields ...string) (*model.ImageTag, error) {
	mt, err := dao.GetImageTag(imageTag, selectedFields...)
	if err != nil {
		return nil, err
	}
	return mt, nil
}

func UpdateImageTag(imageTag model.ImageTag, fieldNames ...string) (bool, error) {
	if imageTag.ImageTagID == 0 {
		return false, errors.New("no Image ID provided")
	}
	_, err := dao.UpdateImageTag(imageTag, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteImage(image model.Image) error {
	//TODO delete registry image

	//Mark image deleted in db
	image.ImageDeleted = 1
	_, err := dao.UpdateImage(image, "deleted")
	if err != nil {
		return err
	}
	return nil
}

func DeleteImageTag(imageTag model.ImageTag) error {
	//TODO delete registry image tag

	//Mark image tag deleted in db
	imageTag.ImageTagDeleted = 1
	_, err := dao.UpdateImageTag(imageTag, "deleted")
	if err != nil {
		return err
	}
	return nil
}

func CreateImageTag(imageTag model.ImageTag) (int64, error) {
	imageTagID, err := dao.AddImageTag(imageTag)
	if err != nil {
		return 0, err
	}
	return imageTagID, nil
}

func UpdateDockerfileCopyCommand(repoImagePath, dockerfileName string) ([]byte, error) {
	dockerfile, err := os.OpenFile(filepath.Join(repoImagePath, dockerfileName), os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer dockerfile.Close()

	pattern := "^(COPY|ADD)\\s+([\\[\"\\s]*)([\\w./-]+)(.*)"
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(dockerfile)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	for scanner.Scan() {
		content := scanner.Text()
		matches := re.FindStringSubmatch(content)
		if len(matches) > 0 {
			// replace the source path of ADD or COPY line in Dockerfile; eg: ADD [ "./abc/r.sql", "/tmp/my.cnf" ]
			// will be replaced by ADD [ "update/r.sql", "/tmp/my.cnf" ]
			_, filename := filepath.Split(matches[3])
			content = fmt.Sprintf("%s %s %s %s", matches[1], matches[2], filepath.Join("upload", strings.TrimSpace(filename)), matches[4])
		}
		writer.WriteString(fmt.Sprintf("%s\n", content))
	}
	writer.Flush()
	dockerfile.Truncate(0)
	dockerfile.Seek(0, 0)
	bufferInfo := buffer.Bytes()
	buffer.WriteTo(dockerfile)
	return bufferInfo, nil
}

func ExistRegistry(projectName string, imageName string, imageTag string) (bool, error) {
	currentName := filepath.Join(projectName, imageName)
	//check image
	repoList, err := GetRegistryCatalog()
	if err != nil {
		logs.Error("Failed to unmarshal repoList body %+v", err)
		return false, err
	}
	for _, imageRegistry := range repoList.Names {
		if imageRegistry == currentName {
			//check tag
			tagList, err := GetRegistryImageTags(currentName)
			if err != nil {
				logs.Error("Failed to unmarshal body %+v", err)
				return false, err
			}
			for _, tagID := range tagList.Tags {
				if imageTag == tagID {
					logs.Info("Image tag existing %s:%s", currentName, tagID)
					return true, nil
				}
			}
		}
	}
	return false, err
}
