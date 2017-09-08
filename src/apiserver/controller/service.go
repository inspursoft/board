package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/common/model"
	//"io/ioutil"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego/logs"
)

type ServiceController struct {
	baseController
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
	fmt.Println(reqData)
	var reqServiceConfig model.ServiceConfig
	err = json.Unmarshal(reqData, &reqServiceConfig)
	if err != nil {
		p.internalError(err)
		return
	}
	fmt.Println(reqServiceConfig)

	//TODO check valid()

	var serviceConfigPath = filepath.Join(repoPath,
		reqServiceConfig.ProjectName, strconv.FormatInt(reqServiceConfig.ServiceID, 10))

	logs.Info("Service config path: %s", serviceConfigPath)
	service.SetDeploymentPath(serviceConfigPath)

	//Build deployment yaml file
	err = service.BuildDeploymentYml(reqServiceConfig)
	if err != nil {
		logs.Info("Build Deployment Yaml failed")
		p.internalError(err)
		return
	}

	//Build deployment yaml file
	err = service.BuildServiceYml(reqServiceConfig)
	if err != nil {
		logs.Info("Build Service Yaml failed")
		p.internalError(err)
		return
	}

	//TODO push service deployment to jenkins
	//push to git
	/*
			var pushobject pushObject

		    pushobject.FileName = "service.yaml and deployment.yaml"
		    pushobject.JobName = "process_service"
		    pushobject.Message = "upload deployment file"
		    pushobject.Value = loadPath

			ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
			if err != nil {
				p.internalError(err)
				return
			}
			logs.Info("Internal push object: %d %s", ret, msg)
			p.CustomAbort(ret, msg)
	*/
}
