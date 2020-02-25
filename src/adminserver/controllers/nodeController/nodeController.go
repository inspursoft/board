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
	"os"
	"path/filepath"
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
// @Success 200 {object} []node.NodeListType  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /list [get]
func (controller *Controller) GetNodeListAction() {
	var nodeListJson []nodeModel.NodeListType
	err := nodeService.GetArrayJsonByFile(nodeModel.AddNodeListJson, &nodeListJson)
	if err != nil {
		errorMsg := fmt.Sprintf("Bad request.%s", err.Error())
		controller.CustomAbort(http.StatusBadRequest, errorMsg)
		return
	}
	controller.Data["json"] = nodeListJson
	controller.ServeJSON()
}

// @Title Get node log history
// @Description Get node log history
// @Success 200 {object} []nodeModel.NodeLogDetail  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /logs [get]
func (controller *Controller) GetNodeLogHistory() {
	var logHistoryList []nodeModel.LogHistory
	err := nodeService.GetArrayJsonByFile(nodeModel.AddNodeHistoryJson, &logHistoryList)
	if err != nil {
		errorMsg := fmt.Sprintf("Bad request.%s", err.Error())
		controller.CustomAbort(http.StatusBadRequest, errorMsg)
		return
	}
	controller.Data["json"] = logHistoryList
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
	logFileName := controller.Ctx.Input.Query("file_name")
	if _, err := os.Stat(filepath.Join(nodeModel.AddNodeLogPath, logFileName)); os.IsNotExist(err) {
		controller.CustomAbort(http.StatusBadRequest, "The file of "+logFileName+" is not exists")
		return
	}
	var nodeLogDetail []nodeModel.NodeLogDetail
	err := nodeService.GetNodeLogDetail(filepath.Join(nodeModel.AddNodeLogPath, logFileName), &nodeLogDetail)
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
	if logHistory := nodeService.CheckExecuting(nodeIp); logHistory != nil {
		controller.Data["json"] = *logHistory
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

	logFileJson := nodeModel.LogHistory{
		Ip: nodeIp, Success: false, Pid: 0, CreationTime: time.Now().Unix(), Type: actionType}
	if err := nodeService.ExecuteCommand(&logFileJson, yamlFile); err != nil {
		controller.CustomAbort(http.StatusBadRequest, err.Error())
		return
	}
	controller.Data["json"] = logFileJson
	controller.ServeJSON()
}

func (controller *Controller) resolveBody(target interface{}) (err error) {
	err = utils.UnmarshalToJSON(controller.Ctx.Request.Body, target)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		return
	}
	return
}
