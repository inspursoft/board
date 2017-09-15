package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/common/model"
	//"io/ioutil"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
)

type ServiceController struct {
	baseController
}

var KubeMasterIp string
var KubeMasterStatus bool
var deploymentFilename = "deployment.yaml"
var serviceFilename = "service.yaml"
var serviceProcess = "process_service"

var apiheader = "Content-Type: application/yaml"
var deploymentAPI = "/apis/extensions/v1beta1/namespaces/"
var serviceAPI = "/api/v1/namespaces/"

var serviceNamespace = "default" //TODO create in project post

func init() {
	var masterip = os.Getenv("KUBEMASTER_IP")
	var masterport = os.Getenv("KUBEMASTER_PORT")
	KubeMasterIp = masterip + ":" + masterport

	logs.Info("Service api started KubeMaster %s %s", KubeMasterIp, time.Now())
}

//  Checking the user priviledge by token
func (p *ServiceController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}

// API to deploy service
func (p *ServiceController) DeployServiceAction() {
	var err error

	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}
	var reqServiceConfig model.ServiceConfig
	err = json.Unmarshal(reqData, &reqServiceConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	// Check deployment parameters
	err = service.CheckDeploymentYmlPara(reqServiceConfig)
	if err != nil {
		logs.Info("Deployment config parameters error")
		p.internalError(err)
		return
	}

	// Check service parameters
	err = service.CheckServiceYmlPara(reqServiceConfig)
	if err != nil {
		logs.Info("Service config parameters error")
		p.internalError(err)
		return
	}

	var serviceConfigPath = filepath.Join(repoPath,
		reqServiceConfig.ProjectName, strconv.FormatInt(reqServiceConfig.ServiceID, 10))

	logs.Info("Service config path: %s", serviceConfigPath)
	service.SetDeploymentPath(serviceConfigPath)

	//Add registry to container images for deployment
	registryprefix := os.Getenv("REGISTRY_HOST") + ":" + os.Getenv("REGISTRY_PORT")
	for index, container := range reqServiceConfig.DeploymentYaml.ContainerList {
		reqServiceConfig.DeploymentYaml.ContainerList[index].BaseImage =
			filepath.Join(registryprefix, container.BaseImage)
	}
	logs.Info(reqServiceConfig)

	//Build deployment yaml file
	err = service.BuildDeploymentYml(reqServiceConfig)
	if err != nil {
		logs.Info("Build Deployment Yaml failed")
		p.internalError(err)
		return
	}

	//Build service yaml file
	err = service.BuildServiceYml(reqServiceConfig)
	if err != nil {
		logs.Info("Build Service Yaml failed")
		p.internalError(err)
		return
	}

	//serviceNamespace = reqServiceConfig.ProjectName TODO in project

	// Push deployment to jenkins
	var pushobject pushObject
	pushobject.FileName = deploymentFilename
	pushobject.JobName = "process_service"
	pushobject.Value = filepath.Join(reqServiceConfig.ProjectName,
		strconv.FormatInt(reqServiceConfig.ServiceID, 10))
	pushobject.Message = fmt.Sprintf("Create deployment for project %s service %d",
		reqServiceConfig.ProjectName, reqServiceConfig.ServiceID)
	pushobject.Extras = KubeMasterIp + deploymentAPI + serviceNamespace + "/deployments"

	// Add deployment file
	pushobject.Items = []string{filepath.Join(pushobject.Value, deploymentFilename)}

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		logs.Info("Create deployment failed %s", pushobject.Extras)
		p.internalError(err)
		return
	}
	logs.Info("Internal push deployment object: %d %s", ret, msg)

	//TODO: If fail to create deployment, should not continue to create service

	//Push service to jenkins
	pushobject.FileName = serviceFilename
	pushobject.JobName = "process_service"
	pushobject.Value = filepath.Join(reqServiceConfig.ProjectName,
		strconv.FormatInt(reqServiceConfig.ServiceID, 10))
	pushobject.Message = fmt.Sprintf("Create service for project %s service %d",
		reqServiceConfig.ProjectName, reqServiceConfig.ServiceID)
	pushobject.Extras = KubeMasterIp + serviceAPI + serviceNamespace + "/services"

	// Add deployment file
	pushobject.Items = []string{filepath.Join(pushobject.Value, serviceFilename)}

	ret, msg, err = InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		logs.Info("Create service failed %s", pushobject.Extras)
		p.internalError(err)
		return
	}
	logs.Info("Internal push service object: %d %s", ret, msg)
	p.CustomAbort(ret, msg)

}

// TODO API to create service config
func (p *ServiceController) CreateServiceConfigAction() {
	//TODO: Assign and return Service ID with mysql
	var serviceID = "1"
	p.Data["json"] = serviceID
	p.ServeJSON()
}
