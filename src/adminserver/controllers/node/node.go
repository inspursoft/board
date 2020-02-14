package node

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/adminserver/models/node"
	"github.com/astaxie/beego"
	"net/http"
	"os"
)

type Controller struct {
	beego.Controller
}

// @Title Get node list
// @Description Get node list
// @Success 200 {object} []node.NodeListType  success
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /list [get]
func (controller *Controller) GetNodeListAction() {
	var nodeListJson []node.NodeListType
	var filePtr *os.File
	if _, err := os.Stat(node.AddDeleteNodeJsonFileName); os.IsNotExist(err) {
		controller.CustomAbort(http.StatusNotFound, "The file of addNode.json is not exists")
		return
	} else {
		filePtr, _ = os.Open(node.AddDeleteNodeJsonFileName)
		decoder := json.NewDecoder(filePtr)
		readErr := decoder.Decode(&nodeListJson)
		if readErr != nil {
			errorMsg := fmt.Sprintf("Unexpected error occurred.%s", readErr.Error())
			controller.CustomAbort(http.StatusInternalServerError, errorMsg)
			return
		} else {
			controller.Data["json"] = nodeListJson
			controller.ServeJSON()
			return
		}
	}
}

// @Title add node
// @Description Get add node
// @Success 200
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /add [get]
func (controller *Controller) AddNodeAction() {
	controller.AddDeleteNode(node.ActionTypeAddNode, node.AddNodeYamlFileName)
}

// @Title delete node
// @Description Get delete node
// @Success 200
// @Failure 400 bad request
// @Failure 500 Internal Server Error
// @router /add [get]
func (controller *Controller) DeleteNodeAction() {
	controller.AddDeleteNode(node.ActionTypeDeleteNode, node.DeleteNodeYamlFileName)
}
