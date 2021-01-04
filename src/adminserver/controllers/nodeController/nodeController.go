package nodeController

import (
	"fmt"
	"github.com/inspursoft/board/src/adminserver/controllers"
	"github.com/inspursoft/board/src/adminserver/models/nodeModel"
	"github.com/inspursoft/board/src/adminserver/service"
	"github.com/inspursoft/board/src/adminserver/service/nodeService"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/token"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

type NodeController struct {
	controllers.BaseController
}

func (controller *NodeController) Render() error {
	return nil
}

// @Title Get node list
// @Description Get node list
// @Success 200 {object} []nodeModel.NodeStatus  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router / [get]
func (controller *NodeController) GetNodeListAction() {
	var nodeResponseList []nodeModel.NodeListResponse
	err := nodeService.GetNodeResponseList(&nodeResponseList)
	if err != nil {
		if err == common.ErrInvalidToken {
			errorMsg := fmt.Sprintf("Token was expired.%s", err.Error())
			controller.CustomAbort(http.StatusUnauthorized, errorMsg)
			return
		}
		errorMsg := fmt.Sprintf("Bad request.%s", err.Error())
		logs.Error(errorMsg)
		controller.CustomAbort(http.StatusBadRequest, errorMsg)
		return
	}
	controller.Data["json"] = nodeResponseList
	controller.ServeJSON()
}

// @Title Get node log list
// @Description Get node log list
// @Success 200 {object} nodeModel.PaginatedNodeLogList  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /logs [get]
func (controller *NodeController) GetNodeLogList() {
	var paginatedNodeLogList = nodeModel.PaginatedNodeLogList{}
	var nodeList []nodeModel.NodeLog
	pageIndex, _ := strconv.Atoi(controller.Ctx.Input.Query("page_index"))
	pageSize, _ := strconv.Atoi(controller.Ctx.Input.Query("page_size"))
	paginatedNodeLogList.Pagination = &nodeModel.Pagination{PageIndex: pageIndex, PageSize: pageSize}
	paginatedNodeLogList.LogList = &nodeList

	err := nodeService.GetPaginatedNodeLogList(&paginatedNodeLogList)
	if err != nil {
		errorMsg := fmt.Sprintf("Bad request.%s", err.Error())
		logs.Error(errorMsg)
		controller.CustomAbort(http.StatusBadRequest, errorMsg)
		return
	}
	controller.Data["json"] = paginatedNodeLogList
	controller.ServeJSON()
}

// @Title Get detail of history log info
// @Description Get detail of history log info
// @Success 200 {object} []nodeModel.NodeLogDetail  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @Param	node_ip	query 	string	true	""
// @Param	creation_time	query 	string	true	""
// @router /log [get]
func (controller *NodeController) GetNodeLogDetail() {
	nodeIp := controller.Ctx.Input.Query("node_ip")
	creationTime, _ := strconv.ParseInt(controller.Ctx.Input.Query("creation_time"), 10, 64)
	var nodeLogDetail []nodeModel.NodeLogDetail
	err := nodeService.GetNodeLogDetail(creationTime, nodeIp, &nodeLogDetail)
	if err != nil {
		logs.Error(err)
		controller.CustomAbort(http.StatusNotFound, err.Error())
		return
	}

	controller.Data["json"] = nodeLogDetail
	controller.ServeJSON()
	return
}

// @Title Delete node log
// @Description Delete node log info from node_log table and node_log_detail_info table
// @Success 200 success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @Param	creation_time	query 	string	true	""
// @router /log [delete]
func (controller *NodeController) DeleteNodeLog() {
	creationTime, _ := strconv.ParseInt(controller.Ctx.Input.Query("creation_time"), 10, 64)

	if nodeService.CheckNodeLogInfoInUse(creationTime) {
		controller.CustomAbort(http.StatusConflict, "Log info in used.")
		return
	}

	err := nodeService.DeleteNodeLogInfo(creationTime)
	if err != nil {
		controller.CustomAbort(http.StatusBadRequest, err.Error())
		logs.Error(err)
		return
	}

	return
}

// @Title get preparation data
// @Description get preparation data
// @Success 200 {object} nodeModel.PreparationData  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /preparation [get]
func (controller *NodeController) PreparationAction() {
	configuration, err := service.GetAllCfg("", false)
	if err != nil {
		logs.Error(err)
		controller.CustomAbort(http.StatusBadRequest, "Failed to get the configuration.")
		return
	}
	hostName := configuration.Board.Hostname
	masterIp := configuration.K8s.KubeMasterIP

	var preparationData = nodeModel.PreparationData{HostIp: hostName, MasterIp: masterIp}
	controller.Data["json"] = preparationData
	controller.ServeJSON()
	return
}

// @Title Update node log
// @Description Update node log
// @Param	body	body	nodeModel.UpdateNodeLog	true	""
// @Success 200
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /callback [put]
func (controller *NodeController) CallBackAction() {
	var putData nodeModel.UpdateNodeLog
	controller.resolveBody(&putData)
	if err := nodeService.UpdateLog(&putData); err != nil {
		errMsg := fmt.Sprintf("Failed to update node log: %v", err)
		logs.Error(errMsg)
		controller.CustomAbort(http.StatusBadRequest, errMsg)
		return
	}
	if putData.ExitCode == 0 && putData.InstallFile == nodeModel.RemoveNodeYamlFile {
		if err := nodeService.DeleteNode(putData.Ip); err != nil {
			if err == common.ErrInvalidToken {
				errMsg := fmt.Sprintf("Token was expired.%s", err.Error())
				controller.CustomAbort(http.StatusUnauthorized, errMsg)
				return
			}
			errMsg := fmt.Sprintf("Failed to delete node: %v", err)
			logs.Error(errMsg)
			controller.CustomAbort(http.StatusBadRequest, errMsg)
			return
		}
	}
	return
}

// @Title add nodeModel
// @Description Get add nodeModel
// @Param	body	body	nodeModel.AddNodePostData	true	""
// @Success 200
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router / [post]
func (controller *NodeController) AddNodeAction() {
	var postData nodeModel.AddNodePostData
	controller.resolveBody(&postData)
	controller.AddRemoveNode(&postData, nodeModel.ActionTypeAddNode, nodeModel.AddNodeYamlFile)
}

// @Title remove node
// @Description remove node
// @Success 200
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @Param	node_ip	        query	string	true	""
// @Param	node_password	query	string	true	""
// @Param	host_password	query	string	true	""
// @Param	host_username	query	string	true	"root"
// @Param	master_password	query	string	true	""
// @router / [delete]
func (controller *NodeController) RemoveNodeAction() {
	nodeIp := controller.Ctx.Input.Query("node_ip")
	nodePassword := controller.Ctx.Input.Query("node_password")
	hostPassword := controller.Ctx.Input.Query("host_password")
	hostUsername := controller.Ctx.Input.Query("host_username")
	masterPassword := controller.Ctx.Input.Query("master_password")
	controller.AddRemoveNode(&nodeModel.AddNodePostData{
		NodePassword:   nodePassword,
		HostUsername:   hostUsername,
		HostPassword:   hostPassword,
		NodeIp:         nodeIp,
		MasterPassword: masterPassword},
		nodeModel.ActionTypeDeleteNode, nodeModel.RemoveNodeYamlFile)
}

// @Title Get node control status
// @Description Get node control status
// @Param	node_name	        path	string	true	""
// @Success 200 {object} model.NodeControlStatus  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /:node_name [get]
func (controller *NodeController) ControlStatusAction() {
	nodeName := strings.TrimSpace(controller.Ctx.Input.Param(":node_name"))
	var nodeControlStatus = model.NodeControlStatus{NodeName: nodeName}
	if err := nodeService.GetNodeControlStatusFromApiServer(&nodeControlStatus); err != nil {
		if err == common.ErrInvalidToken {
			errMsg := fmt.Sprintf("Token was expired.%s", err.Error())
			controller.CustomAbort(http.StatusUnauthorized, errMsg)
			return
		}
		errMsg := fmt.Sprintf("Failed to get node control status: %v", err)
		logs.Error(errMsg)
		controller.CustomAbort(http.StatusBadRequest, errMsg)
		return
	}
	controller.Data["json"] = nodeControlStatus
	controller.ServeJSON()
	return
}

func (controller *NodeController) AddRemoveNode(nodePostData *nodeModel.AddNodePostData,
	actionType nodeModel.ActionType, yamlFile string) {
	if nodeService.CheckExistsInCache(nodePostData.NodeIp) {
		controller.CustomAbort(http.StatusNotAcceptable, "node was locked.")
		return
	}

	if nodeLog, err := nodeService.AddRemoveNodeByContainer(nodePostData, actionType, yamlFile); err != nil {
		nodeService.RemoveCacheData(nodePostData.NodeIp)
		logs.Error(err)
		controller.CustomAbort(http.StatusBadRequest, err.Error())
		return
	} else {
		controller.Data["json"] = *nodeLog
		controller.ServeJSON()
		return
	}
}

func (controller *NodeController) resolveBody(target interface{}) (err error) {
	err = utils.UnmarshalToJSON(controller.Ctx.Request.Body, target)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		return
	}
	return
}
