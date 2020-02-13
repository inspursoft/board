package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/adminserver/models/node"
	"git/inspursoft/board/src/common/utils"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func (controller *Controller) AddDeleteNode(actionType node.ActionType, yamlFile string) {
	nodeIp := controller.Ctx.Input.Query("node_ip")
	//masterIp := utils.GetStringValue("KUBE_MASTER_IP")
	//registryIp := utils.GetStringValue("REGISTRY_IP")
	masterIp := "192.168.122.44"
	registryIp := "192.168.122.44"
	actionName := "Add node ";
	if actionType == node.ActionTypeDeleteNode {
		actionName = "Delete node"
	}
	ws, err := websocket.Upgrade(controller.Ctx.ResponseWriter, controller.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		controller.CustomAbort(http.StatusBadRequest, "Not a webSocket handshake.")
		return
	} else if err != nil {
		controller.CustomAbort(http.StatusInternalServerError, "Cannot setup webSocket connection.")
		return
	}
	defer ws.Close()
	if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
		fileNotExists := fmt.Sprintf("File [%s] not exists", yamlFile)
		controller.sendMessage(ws, fileNotExists, node.WsNodeResponseError)
		return
	}
	controller.generateHostFile(masterIp, nodeIp, registryIp)
	command := fmt.Sprintf("ansible-playbook -i " + node.AddDeleteNodeFileName + " " + yamlFile)
	errTip := fmt.Sprintf("Start execute command.%s", command)
	controller.sendMessage(ws, errTip, node.WsNodeResponseStart)
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		controller.sendMessage(ws, err.Error(), node.WsNodeResponseError)
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		controller.sendMessage(ws, line, node.WsNodeResponseNormal)
		if controller.checkIsEndingLog(line, nodeIp) {
			if controller.checkIsSuccessExecuted(line) {
				msgTip := fmt.Sprintf(actionName+"successed!%s", nodeIp)
				controller.sendMessage(ws, msgTip, node.WsNodeResponseSuccess)
				if actionType == node.ActionTypeAddNode {
					updateErr := controller.appendNodeInfo(nodeIp)
					if updateErr != nil {
						updateTip := fmt.Sprintf("Update the node list failed.%s", updateErr)
						controller.sendMessage(ws, updateTip, node.WsNodeResponseWarning)
					}
				} else {
					updateErr := controller.removeNodeInfo(nodeIp)
					if updateErr != nil {
						updateTip := fmt.Sprintf("Update the node list failed.%s", updateErr)
						controller.sendMessage(ws, updateTip, node.WsNodeResponseWarning)
					}
				}
			} else {
				errTip := fmt.Sprintf(actionName+"failed.%s", nodeIp)
				controller.sendMessage(ws, errTip, node.WsNodeResponseFailed)
			}
		}
		time.Sleep(100 * time.Microsecond)
	}
	cmd.Wait()
}

func (controller *Controller) Render() error {
	return nil
}

func (controller *Controller) generateHostFile(masterIp, nodeIp, registryIp string) error {
	addHosts, err := os.Create(node.AddDeleteNodeFileName)
	defer addHosts.Close()
	if err != nil {
		return err
	}
	addHosts.WriteString("[masters]\n")
	addHosts.WriteString(fmt.Sprintf("%s\n", masterIp))
	addHosts.WriteString("[etcd]\n")
	addHosts.WriteString(fmt.Sprintf("%s\n", masterIp))
	addHosts.WriteString("[nodes]\n")
	addHosts.WriteString(fmt.Sprintf("%s\n", nodeIp))
	addHosts.WriteString("[registry]\n")
	addHosts.WriteString(fmt.Sprintf("%s\n", registryIp))
	return nil
}

func (controller *Controller) sendMessage(ws *websocket.Conn, msg string, status node.WsNodeResponseStatus) {
	send := node.WsNodeResponse{Status: status, Message: msg}
	data, _ := json.Marshal(send)
	ws.WriteMessage(websocket.TextMessage, data)
}

func (controller *Controller) checkIsEndingLog(log, nodeIp string) bool {
	return len(log) > 0 &&
		strings.Contains(log, nodeIp) &&
		strings.Contains(log, "ok") &&
		strings.Contains(log, "changed") &&
		strings.Contains(log, "unreachable") &&
		strings.Contains(log, "failed")
}

func (controller *Controller) checkIsSuccessExecuted(log string) bool {
	return len(log) > 0 &&
		log[strings.Index(log, "unreachable")+12:strings.Index(log, "unreachable")+13] == "0" &&
		log[strings.Index(log, "failed")+7:strings.Index(log, "failed")+8] == "0"
}

func (controller *Controller) removeNodeInfo(nodeIp string) error {
	var nodeListJson []node.NodeListType
	var filePtr *os.File
	_, err := os.Stat(node.AddDeleteNodeJsonFileName);
	if err != nil {
		return err
	} else {
		filePtr, _ = os.Open(node.AddDeleteNodeJsonFileName)
		decoder := json.NewDecoder(filePtr)
		readErr := decoder.Decode(&nodeListJson)
		if readErr != nil {
			return readErr
		}
	}
	defer filePtr.Close();
	if len(nodeListJson) > 0 {
		for index, nodeList := range nodeListJson {
			if nodeList.Ip == nodeIp {
				nodeListJson = append(nodeListJson[:index], nodeListJson[index+1:]...)
			}
		}
	}
	nodeListJsonBytes, _ := json.Marshal(nodeListJson)
	writeErr := ioutil.WriteFile(node.AddDeleteNodeJsonFileName, nodeListJsonBytes, os.ModeType)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func (controller *Controller) appendNodeInfo(nodeIp string) error {
	var nodeListJson []node.NodeListType
	var filePtr *os.File
	if _, err := os.Stat(node.AddDeleteNodeJsonFileName); os.IsNotExist(err) {
		filePtr, _ = os.Create(node.AddDeleteNodeJsonFileName)
	} else {
		filePtr, _ = os.Open(node.AddDeleteNodeJsonFileName)
		decoder := json.NewDecoder(filePtr)
		readErr := decoder.Decode(&nodeListJson)
		if readErr != nil {
			return readErr
		}
	}
	defer filePtr.Close();
	nodeListJson = append(nodeListJson, node.NodeListType{Ip: nodeIp, CreationTime: time.Now()})
	nodeListJsonBytes, _ := json.Marshal(nodeListJson)
	writeErr := ioutil.WriteFile(node.AddDeleteNodeJsonFileName, nodeListJsonBytes, os.ModeType)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func (controller *Controller) resolveBody(target interface{}) (err error) {
	err = utils.UnmarshalToJSON(controller.Ctx.Request.Body, target)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		return
	}
	return
}
