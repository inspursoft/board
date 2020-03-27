package nodeService

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/dao/nodeDao"
	"git/inspursoft/board/src/adminserver/models/nodeModel"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/adminserver/utils"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func AddRemoveNodeByContainer(nodePostData *nodeModel.AddNodePostData,
	actionType nodeModel.ActionType, yamlFile string) (*nodeModel.NodeLog, error) {
	configuration, statusMessage := service.GetAllCfg("")
	if statusMessage == "BadRequest" {
		return nil, fmt.Errorf("failed to get the configuration")
	}
	hostName := configuration.Apiserver.Hostname
	masterIp := configuration.Apiserver.KubeMasterIP
	registryIp := configuration.Apiserver.RegistryIP

	hostFilePath := path.Join(nodeModel.BasePath, nodeModel.HostFileDir)
	if _, err := os.Stat(hostFilePath); os.IsNotExist(err) {
		os.MkdirAll(hostFilePath, os.ModePerm)
	}

	hostFileName := fmt.Sprintf("%s/%s@%s", hostFilePath, nodeModel.NodeHostsFile, nodePostData.NodeIp)

	if err := GenerateHostFile(masterIp, nodePostData.NodeIp, registryIp, hostFileName); err != nil {
		return nil, err
	}

	var newLogId int64
	nodeLog := nodeModel.NodeLog{
		Ip: nodePostData.NodeIp, Completed: false, Success: false, CreationTime: time.Now().Unix(), LogType: actionType}
	if id, err := InsertLog(&nodeLog); err != nil {
		return nil, err
	} else {
		newLogId = id
	}

	var containerEnv = nodeModel.ContainerEnv{NodeIp: nodePostData.NodeIp,
		NodePassword:   nodePostData.NodePassword,
		HostIp:         hostName,
		HostPassword:   nodePostData.HostPassword,
		HostUserName:   nodePostData.HostUsername,
		HostFile:       fmt.Sprintf("%s@%s", nodeModel.NodeHostsFile, nodePostData.NodeIp),
		MasterIp:       masterIp,
		MasterPassword: nodePostData.MasterPassword,
		InstallFile:    yamlFile,
		LogId:          newLogId,
		LogTimestamp:   nodeLog.CreationTime}
	if err := LaunchAnsibleContainer(&containerEnv); err != nil {
		return nil, err
	}
	return &nodeLog, nil
}

func LaunchAnsibleContainer(env *nodeModel.ContainerEnv) error {
	logCache := dao.GlobalCache.Get(env.NodeIp).(*nodeModel.NodeLogCache)
	var secureShell *utils.SecureShell
	var err error
	secureShell, err = utils.NewSecureShell(&logCache.DetailBuffer, env.HostIp, env.HostUserName, env.HostPassword)
	if err != nil {
		return err
	}

	envStr := fmt.Sprintf(""+
		"--env MASTER_PASS=\"%s\" \\\n"+
		"--env MASTER_IP=\"%s\" \\\n"+
		"--env NODE_IP=\"%s\" \\\n"+
		"--env NODE_PASS=\"%s\" \\\n"+
		"--env LOG_ID=\"%d\" \\\n"+
		"--env ADMIN_SERVER_IP=\"%s\" \\\n"+
		"--env ADMIN_SERVER_PORT=\"%d\" \\\n"+
		"--env INSTALL_FILE=\"%s\" \\\n"+
		"--env LOG_TIMESTAMP=\"%d\" \\\n"+
		"--env HOSTS_FILE=\"%s\" ",
		env.MasterPassword,
		env.MasterIp,
		env.NodeIp,
		env.NodePassword,
		env.LogId,
		env.HostIp,
		8081,
		env.InstallFile,
		env.LogTimestamp,
		env.HostFile)

	LogFilePath := path.Join(nodeModel.BasePath, nodeModel.LogFileDir)
	HostDirPath := path.Join(nodeModel.BasePath, nodeModel.HostFileDir)
	PreEnvPath := path.Join(nodeModel.BasePath, nodeModel.PreEnvDir)
	cmdStr := fmt.Sprintf("docker run -td \\\n "+
		"-v %s:/tmp/log \\\n "+
		"-v %s:/tmp/hosts_dir \\\n"+
		"-v %s:/ansible_k8s/pre-env \\\n "+
		"%s \\\n k8s_install:1",
		LogFilePath, HostDirPath, PreEnvPath, envStr)
	err = secureShell.ExecuteCommand(cmdStr)

	if err != nil {
		return err
	}

	return nil
}

func UpdateLog(putLogData *nodeModel.UpdateNodeLog) error {
	var logData *nodeModel.NodeLog
	var err error
	logData, err = nodeDao.GetNodeLog(putLogData.LogId);
	if err != nil {
		return err
	}

	logData.Success = putLogData.Success == 0
	logData.Completed = true
	if errUpdate := nodeDao.UpdateNodeLog(logData); errUpdate != nil {
		return errUpdate
	}

	logFilePath := path.Join(nodeModel.BasePath, nodeModel.LogFileDir)
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.MkdirAll(logFilePath, os.ModePerm)
	}
	logFileName := fmt.Sprintf("%s/%s", logFilePath, putLogData.LogFile)
	if errInsert := InsertLogDetail(logData.Ip, logFileName, logData.CreationTime); errInsert != nil {
		return errInsert
	}

	if putLogData.Success == 0 {
		if putLogData.InstallFile == nodeModel.AddNodeYamlFile {
			if err := nodeDao.InsertNodeStatus(&nodeModel.NodeStatus{
				Ip:           putLogData.Ip,
				CreationTime: logData.CreationTime}); err != nil {
				logs.Info(err.Error())
			}
		}
		if putLogData.InstallFile == nodeModel.RemoveNodeYamlFile {
			if err := nodeDao.DeleteNodeStatus(&nodeModel.NodeStatus{Ip: putLogData.Ip}); err != nil {
				logs.Info(err.Error())
			}
		}
	}
	return RemoveCacheData(putLogData.Ip)
}

func InsertLog(nodeLog *nodeModel.NodeLog) (int64, error) {
	var nodeLogCache = nodeModel.NodeLogCache{}
	nodeLogCache.NodeLogPtr = nodeLog
	if err := dao.GlobalCache.Put(nodeLog.Ip, &nodeLogCache, 3600*time.Second); err != nil {
		return 0, err
	}

	if nodeLog.LogType == nodeModel.ActionTypeAddNode {
		nodeLogCache.DetailBuffer.WriteString(fmt.Sprintf("---Begin add node:%s----\n", nodeLog.Ip))
	} else {
		nodeLogCache.DetailBuffer.WriteString(fmt.Sprintf("---Begin remove node:%s----\n", nodeLog.Ip))
	}

	if newId, err := nodeDao.InsertNodeLog(nodeLog); err != nil {
		return 0, err
	} else {
		return newId, nil
	}
}

func InsertLogDetail(ip, logFileName string, creationTime int64) error {
	var logCache = nodeModel.NodeLogCache{}
	if dao.GlobalCache.IsExist(ip) {
		logCache = *dao.GlobalCache.Get(ip).(*nodeModel.NodeLogCache)
	} else {
		logs.Info(fmt.Sprintf("No cache data for node:%s", ip))
	}

	if _, err := os.Stat(logFileName); os.IsExist(err) {
		filePtr, _ := os.Open(logFileName)
		defer filePtr.Close()
		if fileContent, err := ioutil.ReadFile(logFileName); err != nil {
			return err
		} else {
			if _, writeErr := logCache.DetailBuffer.Write(fileContent); writeErr != nil {
				return writeErr
			}
		}
	}
	logCache.DetailBuffer.WriteString("---End log---\n")
	detail := nodeModel.NodeLogDetailInfo{
		CreationTime: creationTime,
		Detail:       logCache.DetailBuffer.String()}
	if _, err := nodeDao.InsertNodeLogDetail(&detail); err != nil {
		logs.Info(err.Error())
	}
	return nil
}

func GetNodeResponseList(nodeListResponse *[]nodeModel.NodeListResponse) error {
	var apiServerNodeList []nodeModel.ApiServerNodeListResult;
	if err := getNodeListFromApiServer(&apiServerNodeList); err != nil {
		logs.Info(err)
	}

	var nodeStatusList []nodeModel.NodeStatus
	if err := nodeDao.GetNodeStatusList(&nodeStatusList); err != nil {
		return err
	}

	for _, item := range nodeStatusList {
		* nodeListResponse = append(*nodeListResponse, nodeModel.NodeListResponse{
			Ip:           item.Ip,
			CreationTime: item.CreationTime,
			Origin:       0})
	}

	for _, apiItem := range apiServerNodeList {
		var existsInDb = false
		for _, item := range nodeStatusList {
			if item.Ip == apiItem.NodeIP {
				existsInDb = true
			}
		}
		if existsInDb == false {
			* nodeListResponse = append(*nodeListResponse, nodeModel.NodeListResponse{
				Ip:           apiItem.NodeIP,
				CreationTime: time.Now().Unix(),
				Origin:       1})
		}
	}

	return nil
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

func CheckNodeLogInfoInUse(logTimestamp int64) bool {
	return nodeDao.CheckNodeStatusExists(logTimestamp)
}

func DeleteNodeLogInfo(logTimestamp int64) error {
	if err := nodeDao.DeleteNodeLog(logTimestamp); err != nil {
		return err
	}
	if err := nodeDao.DeleteNodeLogDetail(logTimestamp); err != nil {
		return err
	}
	return nil
}

func GetNodeLogDetail(logTimestamp int64, nodeIp string, nodeLogDetail *[]nodeModel.NodeLogDetail) error {
	var reader *bufio.Reader
	if CheckExistsInCache(nodeIp) {
		logCache := dao.GlobalCache.Get(nodeIp).(*nodeModel.NodeLogCache)
		logFilePath := path.Join(nodeModel.BasePath, nodeModel.LogFileDir)
		logFileName := fmt.Sprintf("%s/%d.log", logFilePath, logTimestamp)
		if _, err := os.Stat(logFileName); err == nil {
			filePtr, _ := os.Open(logFileName)
			defer filePtr.Close()

			if fileContent, err := ioutil.ReadFile(logFileName); err != nil {
				return err
			} else {
				if _, writeErr := logCache.DetailBuffer.Write(fileContent); writeErr != nil {
					return writeErr
				}
			}
		}
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

func RemoveCacheData(nodeIp string) error {
	if dao.GlobalCache.IsExist(nodeIp) {
		return dao.GlobalCache.Delete(nodeIp)
	}
	return nil
}

func GetLogInfoInCache(nodeIp string) *nodeModel.NodeLog {
	logCache := dao.GlobalCache.Get(nodeIp).(*nodeModel.NodeLogCache)
	return logCache.NodeLogPtr
}

func getNodeListFromApiServer(nodeList *[]nodeModel.ApiServerNodeListResult) error {
	allConfig, statusMessage := service.GetAllCfg("")
	if statusMessage == "BadRequest" {
		return fmt.Errorf("failed to get the configuration")
	}
	host := allConfig.Apiserver.Hostname
	url := fmt.Sprintf("http://%s:%d/api/v1/nodes?skip=AMS", host, 8088)
	err := utils.RequestHandle(http.MethodGet, url, func(req *http.Request) error {
		req.Header = http.Header{"Content-Type": []string{"application/json"}}
		return nil
	}, nil, func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode == 200 {
			return UnmarshalToJSON(resp.Body, nodeList)
		}
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to get nodes from apiserver.status:%d;message:%s",
			resp.StatusCode, string(data))
	})
	return err
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

func UnmarshalToJSON(in io.ReadCloser, target interface{}) error {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}
