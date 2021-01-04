package controller

import (
	//"encoding/json"
	"fmt"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"

	//"io/ioutil"
	"net/http"
	//"strconv"
	//"strings"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
)

type ConfigMapController struct {
	c.BaseController
}

func (n *ConfigMapController) AddConfigMapAction() {
	var reqCM model.ConfigMapStruct
	var err error
	err = n.ResolveBody(&reqCM)
	if err != nil {
		return
	}

	if reqCM.Name == "" || reqCM.Namespace == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap Name and project should not null")
		return
	}

	if !utils.ValidateWithPattern("configmapname", reqCM.Name) {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap name content is illegal.")
		return
	}

	configmap, err := service.CreateConfigMapK8s(&reqCM)
	if err != nil {
		logs.Debug("Failed to add configmap %v", reqCM)
		n.InternalError(err)
		return
	}
	logs.Info("Added configmap %v", configmap)
}

func (n *ConfigMapController) RemoveConfigMapAction() {
	cmName := n.Ctx.Input.Param(":configmapname")

	if cmName == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap Name should not null")
		return
	}

	projectName := n.GetString("project_name")
	if projectName == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "project should not null")
		return
	}
	//TODO check ConfigMap existing

	err := service.DeleteConfigMapK8s(cmName, projectName)
	if err != nil {
		logs.Info("Delete ConfigMap %s from K8s Failed %v", cmName, err)
		n.InternalError(err)
	}
	logs.Info("Delete ConfigMap %s from K8s Successful %v", cmName, err)

}

func (n *ConfigMapController) GetConfigMapListAction() {
	projectName := n.GetString("project_name")
	configmapName := n.GetString("configmap_name")
	if projectName == "" {
		res, err := service.GetConfigMapListByUser(n.CurrentUser.ID)
		if err != nil {
			logs.Debug("Failed to get ConfigMap List from User")
			n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		n.RenderJSON(res)
	} else if configmapName == "" {
		res, err := service.GetConfigMapListByProject(projectName)
		if err != nil {
			logs.Debug("Failed to get ConfigMap List")
			n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		n.RenderJSON(res)
	} else {
		cm, err := service.GetConfigMapK8s(configmapName, projectName)
		if err != nil {
			logs.Debug("Failed to get ConfigMap")
			n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		n.RenderJSON(cm)
	}
}

func (n *ConfigMapController) GetConfigMapAction() {

	projectName := n.GetString("project_name")
	if projectName == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "project should not null")
		return
	}
	cmName := n.Ctx.Input.Param(":configmapname")
	if cmName == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap Name should not null")
		return
	}

	cm, err := service.GetConfigMapK8s(cmName, projectName)
	if err != nil {
		logs.Debug("Failed to get ConfigMap")
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.RenderJSON(cm)

}

func (n *ConfigMapController) UpdateConfigMapAction() {
	var reqCM model.ConfigMapStruct
	var err error
	err = n.ResolveBody(&reqCM)
	if err != nil {
		return
	}

	if reqCM.Name == "" || reqCM.Namespace == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap Name and project should not null")
		return
	}

	configmap, err := service.UpdateConfigMapK8s(&reqCM)
	if err != nil {
		logs.Debug("Failed to update configmap %v", reqCM)
		n.InternalError(err)
		return
	}
	logs.Info("Updated configmap %v", configmap)
	n.RenderJSON(configmap)
}

//Remove a configmap by name and project
func (n *ConfigMapController) RemoveConfigMapByName() {
	projectName := n.GetString("project_name")
	configmapName := n.GetString("configmap_name")
	if projectName == "" || configmapName == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap Name and project should not null")
		return
	}
	//TODO check ConfigMap existing
	err := service.DeleteConfigMapK8s(configmapName, projectName)
	if err != nil {
		logs.Info("Delete ConfigMap %s from K8s Failed %v", configmapName, err)
		n.InternalError(err)
	}
	logs.Info("Delete ConfigMap %s from K8s Successful %v", configmapName, err)
}

//Update a configmap by name and project
func (n *ConfigMapController) UpdateConfigMapByName() {
	projectName := n.GetString("project_name")
	configmapName := n.GetString("configmap_name")
	if projectName == "" || configmapName == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap Name and project should not null")
		return
	}

	cm, err := service.GetConfigMapK8s(configmapName, projectName)
	if err != nil {
		logs.Debug("Failed to get ConfigMap")
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	} else if cm == nil {
		n.CustomAbortAudit(http.StatusNotFound, "ConfigMap Name not found")
		return
	}
	var reqCM model.ConfigMapStruct
	err = n.ResolveBody(&reqCM)
	if err != nil {
		return
	}

	if reqCM.Name == "" || reqCM.Namespace == "" {
		n.CustomAbortAudit(http.StatusBadRequest, "ConfigMap Name and project should not null")
		return
	}

	configmap, err := service.UpdateConfigMapK8s(&reqCM)
	if err != nil {
		logs.Debug("Failed to update configmap %v", reqCM)
		n.InternalError(err)
		return
	}
	logs.Info("Updated configmap %v", configmap)
	n.RenderJSON(configmap)
}
