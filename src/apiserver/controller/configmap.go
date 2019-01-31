package controller

import (
	//"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	//"io/ioutil"
	"net/http"
	//"strconv"
	//"strings"

	"github.com/astaxie/beego/logs"
)

type ConfigMapController struct {
	BaseController
}

func (n *ConfigMapController) AddConfigMapAction() {
	var reqCM model.ConfigMapStruct
	var err error
	err = n.resolveBody(&reqCM)
	if err != nil {
		return
	}

	if reqCM.Name == "" || reqCM.Namespace == "" {
		n.customAbort(http.StatusBadRequest, "ConfigMap Name and project should not null")
		return
	}

	configmap, err := service.CreateConfigMapK8s(&reqCM)
	if err != nil {
		logs.Debug("Failed to add configmap %v", reqCM)
		n.internalError(err)
		return
	}
	logs.Info("Added configmap %v", configmap)
}

func (n *ConfigMapController) RemoveConfigMapAction() {
	cmName := n.Ctx.Input.Param(":configmapname")

	if cmName == "" {
		n.customAbort(http.StatusBadRequest, "ConfigMap Name should not null")
		return
	}

	projectName := n.GetString("project_name")
	if projectName == "" {
		n.customAbort(http.StatusBadRequest, "project should not null")
		return
	}
	//TODO check ConfigMap existing

	err := service.DeleteConfigMapK8s(cmName, projectName)
	if err != nil {
		logs.Info("Delete ConfigMap %s from K8s Failed %v", cmName, err)
		n.internalError(err)
	}
	logs.Info("Delete ConfigMap %s from K8s Successful %v", cmName, err)

}

func (n *ConfigMapController) GetConfigMapListAction() {

	projectName := n.GetString("project_name")
	if projectName == "" {
		// TODO support all later
		n.customAbort(http.StatusBadRequest, "project should not null")
		return
	}

	res, err := service.GetConfigMapListByProject(projectName)
	if err != nil {
		logs.Debug("Failed to get ConfigMap List")
		n.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.renderJSON(res)
}
