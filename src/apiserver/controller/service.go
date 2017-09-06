package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/common/model"
	//"io/ioutil"
)

type ServiceController struct {
	baseController
}

// API to deploy service
func (p *ServiceController) DeployServiceAction() {
	var err error
	//var basepath = repoPath
	//var projectname = ""
	//var content = []byte("deployment yaml\n")
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
	//filename := basepath + projectname + "deploymentyaml"
	//ioutil.WriteFile(filename, content, 0666)
	//TODO implement service deployment api
	p.serveStatus(200, "deploy service successfully")
}
