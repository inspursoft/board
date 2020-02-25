package nodeService

import (
	"bufio"
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/adminserver/models/nodeModel"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GetArrayJsonByFile(fileName string, v interface{}) error {
	var filePtr *os.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil
	}
	filePtr, _ = os.Open(fileName)
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	if err := decoder.Decode(v); err != nil {
		errorMsg := fmt.Sprintf("Unexpected error occurred.%s", err.Error())
		return fmt.Errorf(errorMsg)
	}
	return nil
}

func GetNodeLogDetail(fileName string, nodeLogDetail *[]nodeModel.NodeLogDetail) error {
	filePtr, _ := os.Open(fileName)
	defer filePtr.Close()
	fileBuf := bufio.NewReader(filePtr)
	for {
		line, err := fileBuf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			errorMsg := fmt.Sprintf("Unexpected error occurred.%s", err.Error())
			return fmt.Errorf(errorMsg)
		}
		line = strings.TrimSpace(line)
		detail := setLogStatus(line)
		*nodeLogDetail = append(*nodeLogDetail, detail)
	}
	return nil
}

func ExecuteCommand(nodeLog *nodeModel.LogHistory, yamlFile string) error {
	if _, err := os.Stat(yamlFile); err != nil {
		return err
	}

	if _, err := os.Stat(nodeModel.AddRemoveShellFile); err != nil {
		return err
	}

	if _, err := os.Stat(nodeModel.AddNodeLogPath); os.IsNotExist(err) {
		os.MkdirAll(nodeModel.AddNodeLogPath, os.ModePerm)
	}

	logFileName := fmt.Sprintf("%s@%d.txt", nodeLog.Ip, time.Now().Unix())
	logFile, createErr := os.Create(filepath.Join(nodeModel.AddNodeLogPath, logFileName))
	if createErr != nil {
		return createErr
	}
	cmd := exec.Command("nohup", "sh", nodeModel.AddRemoveShellFile, yamlFile)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := updateHistoryLog(nodeLog); err != nil {
		return err
	}
	updateList(cmd, nodeLog)
	return nil
}

func GenerateHostFile(masterIp, nodeIp, registryIp string) error {
	addHosts, err := os.Create(nodeModel.AddRemoveNodeFile)
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

func CheckExecuting(nodeIp string) *nodeModel.LogHistory {
	var logHistoryList []nodeModel.LogHistory
	err := GetArrayJsonByFile(nodeModel.AddNodeHistoryJson, &logHistoryList)
	if err != nil {
		return nil
	}
	for _, logHistory := range logHistoryList {
		if logHistory.Completed == false && logHistory.Ip == nodeIp {
			return &logHistory
		}
	}
	return nil
}

func updateList(cmd *exec.Cmd, logHistory *nodeModel.LogHistory) {
	go func() {
		cmd.Wait();
		logHistory.Success = cmd.ProcessState.Success()
		logHistory.Pid = cmd.ProcessState.Pid()
		logHistory.Completed = true

		if err := updateHistoryLog(logHistory); err != nil {
			logs.Info(err.Error())
		}

		if logHistory.Success {
			if logHistory.Type == nodeModel.ActionTypeAddNode {
				if err := appendNodeInfo(logHistory.Ip, logHistory.CreationTime); err != nil {
					logs.Info(err.Error())
				}
			}
			if logHistory.Type == nodeModel.ActionTypeDeleteNode {
				if err := removeNodeInfo(logHistory.Ip); err != nil {
					logs.Info(err.Error())
				}
			}
		}

	}()
}

func removeNodeInfo(nodeIp string) error {
	var nodeListJson []nodeModel.NodeListType
	var filePtr *os.File
	_, err := os.Stat(nodeModel.AddNodeListJson);
	if err != nil {
		return err
	} else {
		filePtr, _ = os.Open(nodeModel.AddNodeListJson)
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
				break
			}
		}
	}
	nodeListJsonBytes, _ := json.Marshal(nodeListJson)
	writeErr := ioutil.WriteFile(nodeModel.AddNodeListJson, nodeListJsonBytes, os.ModeType)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func appendNodeInfo(nodeIp string, creationTime int64) error {
	var nodeListJson []nodeModel.NodeListType
	var filePtr *os.File
	if _, err := os.Stat(nodeModel.AddNodeListJson); os.IsNotExist(err) {
		filePtr, _ = os.Create(nodeModel.AddNodeListJson)
	} else {
		filePtr, _ = os.Open(nodeModel.AddNodeListJson)
		decoder := json.NewDecoder(filePtr)
		readErr := decoder.Decode(&nodeListJson)
		if readErr != nil {
			return readErr
		}
	}
	defer filePtr.Close();
	nodeListJson = append(nodeListJson, nodeModel.NodeListType{Ip: nodeIp, CreationTime: creationTime})
	nodeListJsonBytes, _ := json.Marshal(nodeListJson)
	writeErr := ioutil.WriteFile(nodeModel.AddNodeListJson, nodeListJsonBytes, os.ModeType)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func updateHistoryLog(nodeLog *nodeModel.LogHistory) error {
	var nodeLogHistoryList []nodeModel.LogHistory
	var filePtr *os.File
	if _, err := os.Stat(nodeModel.AddNodeHistoryJson); os.IsNotExist(err) {
		filePtr, _ = os.Create(nodeModel.AddNodeHistoryJson)
	} else {
		filePtr, _ = os.Open(nodeModel.AddNodeHistoryJson)
		decoder := json.NewDecoder(filePtr)
		readErr := decoder.Decode(&nodeLogHistoryList)
		if readErr != nil {
			return readErr
		}
	}
	defer filePtr.Close();
	if len(nodeLogHistoryList) > 0 {
		for index, nodeLogHistory := range nodeLogHistoryList {
			if nodeLogHistory.Ip == nodeLog.Ip && nodeLogHistory.CreationTime == nodeLog.CreationTime {
				nodeLogHistoryList = append(nodeLogHistoryList[:index], nodeLogHistoryList[index+1:]...)
				break
			}
		}
	}
	nodeLogHistoryList = append(nodeLogHistoryList, *nodeLog)
	nodeListJsonBytes, _ := json.Marshal(nodeLogHistoryList)
	writeErr := ioutil.WriteFile(nodeModel.AddNodeHistoryJson, nodeListJsonBytes, os.ModeType)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func checkIsEndingLog(log string) bool {
	return len(log) > 0 &&
		strings.Contains(log, "ok") &&
		strings.Contains(log, "changed") &&
		strings.Contains(log, "unreachable") &&
		strings.Contains(log, "failed")
}

func checkIsSuccessExecuted(log string) bool {
	return len(log) > 0 &&
		log[strings.Index(log, "unreachable")+12:strings.Index(log, "unreachable")+13] == "0" &&
		log[strings.Index(log, "failed")+7:strings.Index(log, "failed")+8] == "0"
}

func setLogStatus(log string) nodeModel.NodeLogDetail {
	if strings.Index(log, "fatal") == 0 || strings.Index(log, "failed") == 0 {
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseError}
	} else if strings.Index(log, "TASK") == 0 {
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseStart}
	} else if strings.Index(log, "ok") == 0 || strings.Index(log, "changed") == 0 {
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseWarning}
	} else if checkIsEndingLog(log) {
		if checkIsSuccessExecuted(log) {
			return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseSuccess}
		}
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseFailed}
	} else {
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseNormal}
	}
}
