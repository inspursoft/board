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
	f.isProjectAdmin = (user.ProjectAdmin == 1)
	if !f.isProjectAdmin {
		f.customAbort(http.StatusForbidden, "Insufficient privileges.")
		return
	}
}

func (f *ConfigFilesController) UploadDeploymentYamlFileAction() {
	projectName := f.GetString("project_name")
	isExist, err := isProjectName(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExist != true {
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
		serviceID, err = CreateServiceConfig(projectName, f.currentUser.ID)
		if err != nil {
			f.internalError(err)
			return
		}
		os.MkdirAll(filepath.Join(repoPath, projectName, serviceID), 0755)
	}

	yamlType := f.GetString("yaml_type")
	fileName, flag := GetYamlFileName(yamlType)
	if flag == false {
		f.customAbort(http.StatusBadRequest, "yaml type invalid.")
		return
	}
	targetFilePath := filepath.Join(repoPath, projectName, serviceID)
	logs.Info("User: %s uploaded deployment yaml file to %s.", f.currentUser.Username, targetFilePath)

	err = f.SaveToFile("upload_file", filepath.Join(targetFilePath, fileName))
	if err != nil {
		f.internalError(err)
	}

}

func (f *ConfigFilesController) DownloadDeploymentYamlFileAction() {
	projectName := f.GetString("project_name")
	isExist, err := isProjectName(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExist != true {
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
	fileName, flag := GetYamlFileName(yamlType)
	if flag == false {
		f.customAbort(http.StatusBadRequest, "yaml type invalid.")
		return
	}
	absFileName := filepath.Join(repoPath, projectName, serviceID, fileName)
	logs.Info("User: %s download deployment yaml file from %s.", f.currentUser.Username, absFileName)

	f.Ctx.Output.Download(absFileName, fileName)
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
		logs.Info("service is not exist")
		return "", nil
	}
	return strconv.Itoa(int(service.ID)), nil
}

func CreateServiceConfig(projectName string, ID int64) (string, error) {
	//get ProjectID
	ProjectID, err := GetProjectID(projectName)
	if err != nil {
		return "", err
	}
	if ProjectID == 0 {
		return "", errors.New("project is invalid")
	}

	var newservice model.ServiceStatus
	newservice.ProjectName = projectName
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = ID
	newservice.ProjectID = ProjectID
	serviceID, err := service.CreateServiceConfig(newservice)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(serviceID)), err
}

func GetProjectID(projectName string) (int64, error) {
	query := model.Project{Name: projectName}
	project, err := service.GetProject(query, "name")
	if err != nil {
		return 0, err
	}
	if project == nil {
		logs.Info("project is not exist")
		return 0, nil
	}
	return project.ID, nil
}

func isProjectName(projectName string) (bool, error) {
	projectExists, err := service.ProjectExists(projectName)
	if err != nil {
		return false, err
	}
	return projectExists, nil
}

func CreateProject(projectName string, ID int, userName string) error {
	var reqProject model.Project
	reqProject.Name = projectName
	reqProject.OwnerID = ID
	reqProject.OwnerName = userName
	isSuccess, err := service.CreateProject(reqProject)
	if err != nil {
		return err
	}
	if !isSuccess {
		return errors.New("Project name is illegal.")
	}
	return nil
}

func GetYamlFileName(yamlType string) (string, bool) {
	var fileName string
	if yamlType == deploymentType {
		fileName = deploymentFilename
	} else if yamlType == serviceType {
		fileName = serviceFilename
	} else {
		return "", false
	}
	return fileName, true
}
