package nodeController

import (
	"fmt"
	"git/inspursoft/board/src/adminserver/models/nodeModel"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/adminserver/service/nodeService"
	"git/inspursoft/board/src/common/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	beego.Controller
}

func (controller *Controller) Render() error {
	return nil
}

// @Title Get node list
// @Description Get node list
// @Success 200 {object} []nodeModel.NodeStatus  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /list [get]
func (controller *Controller) GetNodeListAction() {
	var nodeStatusList []nodeModel.NodeStatus
	err := nodeService.GetNodeStatusList(&nodeStatusList)
	if err != nil {
		errorMsg := fmt.Sprintf("Bad request.%s", err.Error())
		controller.CustomAbort(http.StatusBadRequest, errorMsg)
		return
	}
	controller.Data["json"] = nodeStatusList
	controller.ServeJSON()
}

// @Title Get node log list
// @Description Get node log list
// @Success 200 {object} nodeModel.PaginatedNodeLogList  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /logs [get]
func (controller *Controller) GetNodeLogList() {
	var paginatedNodeLogList = nodeModel.PaginatedNodeLogList{}
	var nodeList []nodeModel.NodeLog
	pageIndex, _ := strconv.Atoi(controller.Ctx.Input.Query("page_index"))
	pageSize, _ := strconv.Atoi(controller.Ctx.Input.Query("page_size"))
	paginatedNodeLogList.Pagination = &nodeModel.Pagination{PageIndex: pageIndex, PageSize: pageSize}
	paginatedNodeLogList.LogList = &nodeList

	err := nodeService.GetPaginatedNodeLogList(&paginatedNodeLogList)
	if err != nil {
		errorMsg := fmt.Sprintf("Bad request.%s", err.Error())
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
// @Param	file_name	query 	string	true	""
// @router /log [get]
func (controller *Controller) GetNodeLogDetail() {
	nodeIp := controller.Ctx.Input.Query("node_ip")
	creationTime, _ := strconv.ParseInt(controller.Ctx.Input.Query("creation_time"), 10, 64)
	var nodeLogDetail []nodeModel.NodeLogDetail
	err := nodeService.GetNodeLogDetail(creationTime, nodeIp, &nodeLogDetail)
	if err != nil {
		controller.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	controller.Data["json"] = nodeLogDetail
	controller.ServeJSON()
	return
}

// @Title add nodeModel
// @Description Get add nodeModel
// @Success 200
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /add [post]
func (controller *Controller) AddNodeAction() {
	var postData nodeModel.AddNodePostData
	controller.resolveBody(&postData)
	controller.AddRemoveNode(postData.NodeIp, nodeModel.ActionTypeAddNode, nodeModel.AddNodeYamlFile)
}

// @Title remove node
// @Description remove node
// @Success 200
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /remove [delete]
func (controller *Controller) RemoveNodeAction() {
	nodeIp := controller.Ctx.Input.Query("node_ip")
	controller.AddRemoveNode(nodeIp, nodeModel.ActionTypeDeleteNode, nodeModel.RemoveNodeYamlFile)
}

func (controller *Controller) AddRemoveNode(nodeIp string, actionType nodeModel.ActionType, yamlFile string) {
	if nodeService.CheckExistsInCache(nodeIp) {
		controller.Data["json"] = *nodeService.GetLogInfoInCache(nodeIp)
		controller.ServeJSON()
		return
	}

	configuration, statusMessage := service.GetAllCfg("")
	if statusMessage == "BadRequest" {
		controller.CustomAbort(http.StatusBadRequest, "Failed to get the configuration.")
		return
	}
	masterIp := configuration.Apiserver.KubeMasterIP
	registryIp := configuration.Apiserver.RegistryIP

	if err := nodeService.GenerateHostFile(masterIp, nodeIp, registryIp); err != nil {
		controller.CustomAbort(http.StatusBadRequest, err.Error())
		return
	}

	nodeLog := nodeModel.NodeLog{
		Ip: nodeIp, Success: false, Pid: 0, CreationTime: time.Now().Unix(), LogType: actionType}
	if err := nodeService.ExecuteCommand(&nodeLog, yamlFile); err != nil {
		controller.CustomAbort(http.StatusBadRequest, err.Error())
		return
	}
	controller.Data["json"] = nodeLog
	controller.ServeJSON()
	return
}

func (controller *Controller) resolveBody(target interface{}) (err error) {
	err = utils.UnmarshalToJSON(controller.Ctx.Request.Body, target)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		return
	}
	return
}
