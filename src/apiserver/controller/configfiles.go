package controller

import (
	"errors"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego/logs"
)

const (
	deploymentType = "deployment"
	serviceType    = "service"
)

type ConfigFilesController struct {
	baseController
}

func (f *ConfigFilesController) Prepare() {
	user := f.getCurrentUser()
	if user == nil {
		f.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	f.currentUser = user
}

func (f *ConfigFilesController) UploadDeploymentYamlFileAction() {
	projectName := f.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExistence != true {
		f.customAbort(http.StatusBadRequest, "Project don't exist.")
		return
	}

	serviceName := f.GetString("service_name")
	serviceID, err := getServiceID(serviceName, projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if serviceID == "" {
		serviceID, err = createUploadedServiceConfig(projectName, serviceName, f.currentUser.ID, f.currentUser.Username)
		if err != nil {
			f.internalError(err)
			return
		}
		os.MkdirAll(filepath.Join(repoPath(), projectName, serviceID), 0755)
	}

	yamlType := f.GetString("yaml_type")
	fileName := getYamlFileName(yamlType)
	if fileName == "" {
		f.customAbort(http.StatusBadRequest, "yaml type invalid.")
		return
	}
	targetFilePath := filepath.Join(repoPath(), projectName, serviceID)
	logs.Info("User: %s uploaded %s yaml file to %s.", f.currentUser.Username, yamlType, targetFilePath)

	err = f.SaveToFile("upload_file", filepath.Join(targetFilePath, fileName))
	if err != nil {
		f.internalError(err)
	}

	f.Data["json"] = serviceID
	f.ServeJSON()
}

func (f *ConfigFilesController) DownloadDeploymentYamlFileAction() {
	projectName := f.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExistence != true {
		f.customAbort(http.StatusBadRequest, "Project name invalid.")
		return
	}

	serviceName := f.GetString("service_name")
	serviceID, err := getServiceID(serviceName, projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if serviceID == "" {
		f.customAbort(http.StatusBadRequest, "Service name invalid.")
		return
	}

	//get paras
	yamlType := f.GetString("yaml_type")
	fileName := getYamlFileName(yamlType)
	if fileName == "" {
		f.customAbort(http.StatusBadRequest, "yaml type invalid.")
		return
	}
	absFileName := filepath.Join(repoPath(), projectName, serviceID, fileName)
	logs.Info("User: %s download %s yaml file from %s.", f.currentUser.Username, yamlType, absFileName)

	f.Ctx.Output.Download(absFileName, fileName)
}

func (f *ConfigFilesController) UploadDockerfileFileAction() {
	projectName := f.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExistence != true {
		f.customAbort(http.StatusBadRequest, "Project don't exist.")
		return
	}

	imageName := f.GetString("image_name")
	tagName := f.GetString("tag_name")
	targetFilePath := filepath.Join(repoPath(), projectName, imageName, tagName)
	err = os.MkdirAll(targetFilePath, 0755)
	if err != nil {
		f.internalError(err)
		return
	}
	logs.Info("User: %s uploaded Dockerfile file to %s.", f.currentUser.Username, targetFilePath)

	err = f.SaveToFile("upload_file", filepath.Join(targetFilePath, "Dockerfile"))
	if err != nil {
		f.internalError(err)
	}

}

func (f *ConfigFilesController) DownloadDockerfileFileAction() {
	projectName := f.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExistence != true {
		f.customAbort(http.StatusBadRequest, "Project name invalid.")
		return
	}

	imageName := f.GetString("image_name")
	tagName := f.GetString("tag_name")
	targetFilePath := filepath.Join(repoPath(), projectName, imageName, tagName)
	if _, err := os.Stat(targetFilePath); os.IsNotExist(err) {
		f.customAbort(http.StatusBadRequest, "image Name and  tag name are invalid.")
		return
	}

	absFileName := filepath.Join(repoPath(), projectName, imageName, tagName, "Dockerfile")
	logs.Info("User: %s download Dockerfile file from %s.", f.currentUser.Username, absFileName)

	f.Ctx.Output.Download(absFileName, "Dockerfile")
}

func getServiceID(serviceName string, projectName string) (string, error) {
	var servicequery model.ServiceStatus
	servicequery.Name = serviceName
	servicequery.ProjectName = projectName
	service, err := service.GetService(servicequery, "name", "project_name")
	if err != nil {
		return "", err
	}
	if service == nil {
		logs.Info("service doesn't exist")
		return "", nil
	}
	return strconv.Itoa(int(service.ID)), nil
}

func createUploadedServiceConfig(projectName string, serviceName string, ID int64, userName string) (string, error) {
	//get ProjectID
	query := model.Project{Name: projectName}
	project, err := service.GetProject(query, "name")
	if err != nil {
		return "", err
	}
	if project == nil {
		logs.Info("project doesn't exist")
		return "", errors.New("project is invalid")
	}

	var newservice model.ServiceStatus
	newservice.Name = serviceName
	newservice.ProjectName = projectName
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = ID
	newservice.OwnerName = userName
	newservice.ProjectID = project.ID
	serviceID, err := service.CreateServiceConfig(newservice)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(serviceID)), err
}

func getYamlFileName(yamlType string) string {
	var fileName string
	if yamlType == deploymentType {
		fileName = deploymentFilename
	} else if yamlType == serviceType {
		fileName = serviceFilename
	} else {
		return ""
	}
	return fileName
}
