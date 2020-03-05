package nodeService

import (
	"bufio"
	"bytes"
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/dao/nodeDao"
	"git/inspursoft/board/src/adminserver/models/nodeModel"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetNodeStatusList(nodeStatusList *[]nodeModel.NodeStatus) error {
	return nodeDao.GetNodeStatusList(nodeStatusList)
}

func GetPaginatedNodeLogList(v *nodeModel.PaginatedNodeLogList) error {
	offset := (v.Pagination.PageIndex - 1) * v.Pagination.PageSize
	if err := nodeDao.GetNodeLogList(v.LogList, v.Pagination.PageSize, offset); err != nil {
		return err
	}

	totalCount, err := nodeDao.GetLogTotalRecordCount()
	if err != nil {
		return err
	}
	v.Pagination.TotalCount = totalCount
	v.Pagination.PageCount = int(v.Pagination.TotalCount)/v.Pagination.PageSize + 1
	return nil
}

func GetNodeLogDetail(logTimestamp int64, nodeIp string, nodeLogDetail *[]nodeModel.NodeLogDetail) error {
	var reader *bufio.Reader
	if CheckExistsInCache(nodeIp) {
		logCache := dao.GlobalCache.Get(nodeIp).(*nodeModel.NodeLogCache)
		reader = bufio.NewReader(strings.NewReader(logCache.DetailBuffer.String()))
	} else {
		detailInfo, err := nodeDao.GetNodeLogDetail(logTimestamp)
		if err != nil {
			return err
		}
		reader = bufio.NewReader(bytes.NewBufferString(detailInfo.Detail))
	}
	for {
		line, err := reader.ReadString('\n')
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

func ExecuteCommand(nodeLog *nodeModel.NodeLog, yamlPathFile string, shellPathFile string) error {
	if _, err := os.Stat(yamlPathFile); err != nil {
		return err
	}
	if _, err := os.Stat(shellPathFile); err != nil {
		return err
	}

	cmd := exec.Command("nohup", "sh", shellPathFile, yamlPathFile)

	var nodeLogCache = nodeModel.NodeLogCache{}
	nodeLogCache.NodeLogPtr = nodeLog

	if nodeLog.LogType == nodeModel.ActionTypeAddNode {
		nodeLogCache.DetailBuffer.WriteString(fmt.Sprintf("---Begin add node:%s----", nodeLog.Ip))
	} else {
		nodeLogCache.DetailBuffer.WriteString(fmt.Sprintf("---Begin remove node:%s----", nodeLog.Ip))
	}

	cmd.Stdout = &nodeLogCache.DetailBuffer
	cmd.Stderr = &nodeLogCache.DetailBuffer

	if err := dao.GlobalCache.Put(nodeLog.Ip, &nodeLogCache, 3600*time.Second); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	insertData(cmd, nodeLog)
	return nil
}

func insertData(cmd *exec.Cmd, nodeLog *nodeModel.NodeLog) {
	go func() {
		cmd.Wait();
		nodeLog.Success = cmd.ProcessState.Success()
		nodeLog.Pid = cmd.ProcessState.Pid()

		logCache := dao.GlobalCache.Get(nodeLog.Ip).(*nodeModel.NodeLogCache)

		var endStr string
		if nodeLog.LogType == nodeModel.ActionTypeAddNode {
			endStr = fmt.Sprintf("---Add node completed:%s----\n", nodeLog.Ip)
		} else {
			endStr = fmt.Sprintf("---Remove node completed:%s----\n", nodeLog.Ip)
		}

		if _, err := logCache.DetailBuffer.WriteString(endStr); err != nil {
			logs.Info(err)
		}
		detail := nodeModel.NodeLogDetailInfo{CreationTime: nodeLog.CreationTime,
			Detail: logCache.DetailBuffer.String()}

		if _, err := nodeDao.InsertNodeLogDetail(&detail); err != nil {
			logs.Info(err.Error())
		}

		if err := nodeDao.InsertNodeLog(nodeLog); err != nil {
			logs.Info(err.Error())
		}

		if nodeLog.Success {
			if nodeLog.LogType == nodeModel.ActionTypeAddNode {
				if err := nodeDao.InsertNodeStatus(&nodeModel.NodeStatus{
					Ip:           nodeLog.Ip,
					CreationTime: nodeLog.CreationTime}); err != nil {
					logs.Info(err.Error())
				}
			}
			if nodeLog.LogType == nodeModel.ActionTypeDeleteNode {
				if err := nodeDao.DeleteNodeStatus(&nodeModel.NodeStatus{Ip: nodeLog.Ip}); err != nil {
					logs.Info(err.Error())
				}
			}
		}
		dao.GlobalCache.Delete(nodeLog.Ip)
	}()
}

func GenerateHostFile(masterIp, nodeIp, registryIp, nodePathFile string) error {
	addHosts, err := os.Create(nodePathFile)
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

func CheckExistsInCache(nodeIp string) bool {
	return dao.GlobalCache.IsExist(nodeIp)
}

func GetLogInfoInCache(nodeIp string) *nodeModel.NodeLog {
	logCache := dao.GlobalCache.Get(nodeIp).(*nodeModel.NodeLogCache)
	return logCache.NodeLogPtr
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
	} else if strings.Index(log, "---") == 0 {
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseStart}
	} else if checkIsEndingLog(log) {
		if checkIsSuccessExecuted(log) {
			return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseSuccess}
		}
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseFailed}
	} else {
		return nodeModel.NodeLogDetail{Message: log, Status: nodeModel.NodeLogResponseNormal}
	}
}
