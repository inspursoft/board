package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/astaxie/beego/logs"
)

type ServiceRollingUpdateController struct {
	baseController
}

func (p *ServiceRollingUpdateController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
}

func (p *ServiceRollingUpdateController) GetRollingUpdateServiceConfigAction() {
	serviceConfig, _ := p.getServiceConfig()
	if serviceConfig.Spec.Template == nil || len(serviceConfig.Spec.Template.Spec.Containers) < 1 {
		p.customAbort(http.StatusBadRequest, "Requested service's config is invalid.")
	}

	var imageList []model.ImageIndex
	for _, container := range serviceConfig.Spec.Template.Spec.Containers {
		indexProject := strings.IndexByte(container.Image, '/')
		indexImage := strings.LastIndexByte(container.Image, '/')
		indexTag := strings.LastIndexByte(container.Image, ':')
		imageList = append(imageList, model.ImageIndex{ImageName: container.Image[indexProject+1 : indexTag],
			ImageTag:    container.Image[indexTag+1:],
			ProjectName: container.Image[indexProject+1 : indexImage]})
	}

	p.Data["json"] = imageList
	p.ServeJSON()
}
func (p *ServiceRollingUpdateController) PostRollingUpdateServiceConfigAction() {
	var imageList []model.ImageIndex
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
	}
	err = json.Unmarshal(reqData, &imageList)
	if err != nil {
		p.internalError(err)
	}
	logs.Debug("Image list info: %+v\n", imageList)

	serviceConfig, projectName := p.getServiceConfig()
	if len(serviceConfig.Spec.Template.Spec.Containers) != len(imageList) {
		p.customAbort(http.StatusConflict, "Image's config is invalid.")
	}

	var rollingUpdateConfig v1.ReplicationController
	rollingUpdateConfig.Spec.Template = &v1.PodTemplateSpec{}
	for index, container := range serviceConfig.Spec.Template.Spec.Containers {
		image := registryBaseURI() + "/" + imageList[index].ImageName + ":" + imageList[index].ImageTag
		if serviceConfig.Spec.Template.Spec.Containers[index].Image != image {
			rollingUpdateConfig.Spec.Template.Spec.Containers = append(rollingUpdateConfig.Spec.Template.Spec.Containers, v1.Container{
				Name:  container.Name,
				Image: image,
			})
		}
	}

	serviceRollConfig, err := json.Marshal(rollingUpdateConfig)
	if err != nil {
		p.internalError(err)
	}
	extras := fmt.Sprintf("%s%s", kubeMasterURL(), filepath.Join(deploymentAPI, projectName, "deployments", serviceConfig.Name))
	logs.Debug("Requested rolling update jenkins with extras: %s", extras)

	req, err := http.NewRequest("POST", `http://jenkins:8080/job/rolling_update/buildWithParameters`, nil)
	if err != nil {
		p.internalError(err)
	}
	form := url.Values{}
	form.Add("value", string(serviceRollConfig))
	form.Add("extras", extras)
	req.URL.RawQuery = form.Encode()
	logs.Debug("Request object to Jenkins is %+v", req)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.internalError(err)
	}
	if resp != nil {
		defer resp.Body.Close()
		respData, _ := ioutil.ReadAll(resp.Body)
		logs.Debug("Response with requested Jenkins rolling update job: %s", string(respData))
	}

}

func (p *ServiceRollingUpdateController) getServiceConfig() (*v1.ReplicationController, string) {
	projectName := p.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		p.internalError(err)
	}
	if isExistence != true {
		p.customAbort(http.StatusBadRequest, "Project don't exist.")
	}

	serviceName := p.GetString("service_name")
	serviceStatus, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		p.internalError(err)
	}
	if serviceStatus == nil {
		p.customAbort(http.StatusBadRequest, "Service name don't exist.")
	}

	absFileName := filepath.Join(repoPath(), projectName, strconv.Itoa(int(serviceStatus.ID)), deploymentFilename)
	logs.Info("User: %s get deployment.yaml images info from %s.", p.currentUser.Username, absFileName)

	yamlFile, err := ioutil.ReadFile(absFileName)
	if err != nil {
		p.internalError(err)
	}

	var rcConfig v1.ReplicationController
	err = yaml.Unmarshal(yamlFile, &rcConfig)
	if err != nil {
		p.internalError(err)
	}

	return &rcConfig, projectName
}

func (p *ServiceRollingUpdateController) PatchRollingUpdateServiceAction() {

	var imageList []model.ImageIndex
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
	}
	//logs.Debug("reqData %+v\n", string(reqData))
	err = json.Unmarshal(reqData, &imageList)
	if err != nil {
		p.internalError(err)
	}
	logs.Debug("Image list info: %+v\n", imageList)

	serviceConfig, projectName := p.getServiceConfig()
	if len(serviceConfig.Spec.Template.Spec.Containers) != len(imageList) {
		p.customAbort(http.StatusConflict, "Image's config is invalid.")
	}

	var rollingUpdateConfig v1.ReplicationController
	rollingUpdateConfig.Spec.Template = &v1.PodTemplateSpec{}
	for index, container := range serviceConfig.Spec.Template.Spec.Containers {
		image := registryBaseURI() + "/" + imageList[index].ImageName + ":" + imageList[index].ImageTag
		if serviceConfig.Spec.Template.Spec.Containers[index].Image != image {
			rollingUpdateConfig.Spec.Template.Spec.Containers = append(rollingUpdateConfig.Spec.Template.Spec.Containers, v1.Container{
				Name:  container.Name,
				Image: image,
			})
		}
	}

	if len(rollingUpdateConfig.Spec.Template.Spec.Containers) == 0 {
		logs.Info("Nothing to be updated")
		return
	}

	serviceRollConfig, err := json.Marshal(rollingUpdateConfig)
	if err != nil {
		logs.Debug("rollingUpdateConfig %+v\n", rollingUpdateConfig)
		p.internalError(err)
	}
	//logs.Debug("Marshal serviceRollConfig %+v\n", string(serviceRollConfig))

	cli, err := service.K8sCliFactory("", kubeMasterURL(), "v1beta1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		p.internalError(err)
	}

	d := apiSet.Deployments(projectName)
	patchType := api.StrategicMergePatchType
	deployData, err := d.Patch(serviceConfig.Name, patchType, serviceRollConfig)
	if err != nil {
		logs.Error("Failed to update service %+v\n", err)
		p.internalError(err)
	}
	logs.Debug("New updated deployment: %+v\n", deployData)

	//TODO update deployment yaml file
	//err = service.rollingUpdateYaml(serviceID, deployData)

}
