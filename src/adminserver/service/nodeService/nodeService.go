package nodeService

import (
	"bufio"
	"bytes"
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/dao/nodeDao"
	"git/inspursoft/board/src/adminserver/models/nodeModel"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/adminserver/tools/secureShell"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/token"
	"git/inspursoft/board/src/common/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

func AddRemoveNodeByContainer(nodePostData *nodeModel.AddNodePostData,
	actionType nodeModel.ActionType, yamlFile string) (*nodeModel.NodeLog, error) {
	configuration, err := service.GetAllCfg("", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get the configuration")
	}
	hostName := configuration.Board.Hostname
	registryIp := configuration.K8s.RegistryIP

	hostFilePath := path.Join(nodeModel.BasePath, nodeModel.HostFileDir)
	if _, err := os.Stat(hostFilePath); os.IsNotExist(err) {
		dirErr := os.MkdirAll(hostFilePath, os.ModePerm)
		if dirErr != nil {
			return nil, dirErr
		}
	}

	hostFileName := fmt.Sprintf("%s/%s@%s", hostFilePath, nodeModel.NodeHostsFile, nodePostData.NodeIp)

	if err := GenerateHostFile(nodePostData.MasterIp, nodePostData.NodeIp, registryIp, hostFileName); err != nil {
		return nil, err
	}

	var nodeLogCache = nodeModel.NodeLogCache{}
	if err := dao.GlobalCache.Put(nodePostData.NodeIp, &nodeLogCache, 3600*time.Second); err != nil {
		return nil, err
	}

	var secure *secureShell.SecureShell
	var secureErr error
	secure, secureErr = secureShell.NewSecureShell(&nodeLogCache.DetailBuffer,
		hostName, nodePostData.HostUsername, nodePostData.HostPassword)
	if secureErr != nil {
		RemoveCacheData(nodePostData.NodeIp)
		return nil, secureErr
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
		MasterIp:       nodePostData.MasterIp,
		MasterPassword: nodePostData.MasterPassword,
		InstallFile:    yamlFile,
		LogId:          newLogId,
		LogTimestamp:   nodeLog.CreationTime}
	if err := LaunchAnsibleContainer(&containerEnv, secure); err != nil {
		logFileName := fmt.Sprintf("%d.log", nodeLog.CreationTime)
		updateLogInfo := &nodeModel.UpdateNodeLog{LogId: int(newLogId),
			Ip:          nodePostData.NodeIp,
			InstallFile: yamlFile,
			LogFile:     logFileName,
			ExitCode:    1}
		UpdateLog(updateLogInfo)
		return nil, err
	}
	return &nodeLog, nil
}

func LaunchAnsibleContainer(env *nodeModel.ContainerEnv, secure *secureShell.SecureShell) error {
	if currentToken, ok := dao.GlobalCache.Get("boardadmin").(string); ok {
		envStr := fmt.Sprintf(`--env MASTER_PASS=%s \
--env MASTER_IP=%s \
--env NODE_IP=%s \
--env NODE_PASS=%s \
--env LOG_ID=%d \
--env ADMIN_SERVER_IP=%s \
--env ADMIN_SERVER_PORT=%d \
--env INSTALL_FILE=%s \
--env LOG_TIMESTAMP=%d \
--env HOSTS_FILE=%s \
--env TOKEN=%s`,
			env.MasterPassword,
			env.MasterIp,
			env.NodeIp,
			env.NodePassword,
			env.LogId,
			env.HostIp,
			8081,
			env.InstallFile,
			env.LogTimestamp,
			env.HostFile,
			currentToken)

		LogFilePath := path.Join(nodeModel.BasePath, nodeModel.LogFileDir)
		HostDirPath := path.Join(nodeModel.BasePath, nodeModel.HostFileDir)
		cmdStr := fmt.Sprintf(`docker run --rm -d \
-v %s:/tmp/log \
-v %s:/tmp/hosts_dir \
-v %s:/ansible_k8s/pre-env \
%s k8s_install:1.19 `, LogFilePath, HostDirPath, nodeModel.PreEnvDir, envStr)

		if err := secure.ExecuteCommand(cmdStr); err != nil {
			return err
		}

		return nil
	} else {
		return common.ErrInvalidToken
	}
}

func UpdateLog(putLogData *nodeModel.UpdateNodeLog) error {
	var logData *nodeModel.NodeLog
	var err error
	logData, err = nodeDao.GetNodeLog(putLogData.LogId)
	if err != nil {
		return err
	}

	logData.Success = putLogData.ExitCode == 0
	logData.Completed = true
	if errUpdate := nodeDao.UpdateNodeLog(logData); errUpdate != nil {
		return errUpdate
	}

	logFilePath := path.Join(nodeModel.BasePath, nodeModel.LogFileDir)
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		dirErr := os.MkdirAll(logFilePath, os.ModePerm)
		if dirErr != nil {
			return dirErr
		}
	}
	logFileName := fmt.Sprintf("%s/%s", logFilePath, putLogData.LogFile)
	if errInsert := InsertLogDetail(logData.Ip, logFileName, logData.CreationTime); errInsert != nil {
		return errInsert
	}

	if putLogData.ExitCode == 0 {
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
	nodeLogCache := dao.GlobalCache.Get(nodeLog.Ip).(*nodeModel.NodeLogCache)

	nodeLogCache.NodeLogPtr = nodeLog
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
		errStr := fmt.Sprintf("No cache data for node:%s \n", ip)
		logs.Info(errStr)
		logCache.DetailBuffer.WriteString(errStr)
	}

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
	logCache.DetailBuffer.WriteString("---End log---\n")
	detail := nodeModel.NodeLogDetailInfo{
		CreationTime: creationTime,
		Detail:       logCache.DetailBuffer.String()}
	if nodeDao.CheckNodeLogDetailExists(creationTime) {
		if err := nodeDao.UpdateNodeLogDetail(&detail); err != nil {
			return err
		}
	} else {
		if _, err := nodeDao.InsertNodeLogDetail(&detail); err != nil {
			return err
		}
	}
	return nil
}

func GetNodeResponseList(nodeListResponse *[]nodeModel.NodeListResponse) error {
	var apiServerNodeList []nodeModel.ApiServerNodeListResult
	if err := getNodeListFromApiServer(&apiServerNodeList); err != nil {
		return err
	}

	var nodeStatusList []nodeModel.NodeStatus
	if err := nodeDao.GetNodeStatusList(&nodeStatusList); err != nil {
		return err
	}

	for _, item := range apiServerNodeList {
		_, isMaster := item.Labels["node-role.kubernetes.io/master"]
		var origin = 0
		var logTime = item.CreateTime
		for _, adminItem := range nodeStatusList {
			if item.NodeIP == adminItem.Ip {
				origin = 1
				logTime = adminItem.CreationTime
			}
		}
		*nodeListResponse = append(*nodeListResponse, nodeModel.NodeListResponse{
			Ip:           item.NodeIP,
			CreationTime: item.CreateTime,
			Status:       item.Status,
			NodeName:     item.NodeName,
			IsMaster:     isMaster,
			IsEdge:       item.NodeType == "edge",
			LogTime:      logTime,
			Origin:       origin})
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
	var cacheBuffer = bytes.NewBuffer([]byte{})
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
				cacheBuffer.Write(logCache.DetailBuffer.Bytes())
				if _, writeErr := cacheBuffer.Write(fileContent); writeErr != nil {
					return writeErr
				}
			}
		}
		reader = bufio.NewReader(strings.NewReader(cacheBuffer.String()))
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
	masters := strings.Split(masterIp, "_")
	for _, master := range masters {
		addHosts.WriteString(fmt.Sprintf("%s\n", master))
	}
	addHosts.WriteString("[etcd]\n")
	for _, master := range masters {
		addHosts.WriteString(fmt.Sprintf("%s\n", master))
	}
	addHosts.WriteString("[nodes]\n")
	nodes := strings.Split(nodeIp, "_")
	for _, node := range nodes {
		addHosts.WriteString(fmt.Sprintf("%s\n", node))
	}
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

func GetNodeControlStatusFromApiServer(nodeControlStatus *model.NodeControlStatus) error {
	url := fmt.Sprintf("api/v1/nodes/%s", nodeControlStatus.NodeName)
	return getResponseJsonFromApiServer(url, nodeControlStatus)
}

func DeleteNode(nodeIp string) error {
	urlPath := fmt.Sprintf("api/v1/nodes/%s?node_ip=%s", nodeIp, nodeIp)
	return deleteActionFromApiServer(urlPath)
}

func getNodeListFromApiServer(nodeList *[]nodeModel.ApiServerNodeListResult) error {
	return getResponseJsonFromApiServer("api/v1/nodes", nodeList)
}

func deleteActionFromApiServer(urlPath string) error {
	allConfig, errCfg := service.GetAllCfg("", false)
	if errCfg != nil {
		return fmt.Errorf("failed to get the configuration")
	}
	host := allConfig.Board.Hostname
	port := allConfig.Board.APIServerPort
	url := fmt.Sprintf("http://%s:%s/%s", host, port, urlPath)

	if currentToken, ok := dao.GlobalCache.Get("boardadmin").(string); ok {
		err := utils.RequestHandle(http.MethodDelete, url, func(req *http.Request) error {
			req.Header = http.Header{
				"Content-Type": []string{"application/json"},
				"token":        []string{currentToken},
			}
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode == 200 {
				return nil
			}
			if resp.StatusCode == 401 {
				return common.ErrInvalidToken
			}
			data, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("failed to request apiserver.status:%d;message:%s",
				resp.StatusCode, string(data))
		})
		return err
	} else {
		return common.ErrInvalidToken
	}
}

func getResponseJsonFromApiServer(urlPath string, res interface{}) error {
	allConfig, errCfg := service.GetAllCfg("", false)
	if errCfg != nil {
		return fmt.Errorf("failed to get the configuration")
	}
	host := allConfig.Board.Hostname
	port := allConfig.Board.APIServerPort
	url := fmt.Sprintf("http://%s:%s/%s", host, port, urlPath)

	if currentToken, ok := dao.GlobalCache.Get("boardadmin").(string); ok {
		err := utils.RequestHandle(http.MethodGet, url, func(req *http.Request) error {
			req.Header = http.Header{
				"Content-Type": []string{"application/json"},
				"token":        []string{currentToken},
			}
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode == 200 {
				return utils.UnmarshalToJSON(resp.Body, res)
			}
			if resp.StatusCode == 401 {
				return common.ErrInvalidToken
			}
			data, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("failed to request apiserver.status:%d;message:%s",
				resp.StatusCode, string(data))
		})
		return err
	} else {
		return common.ErrInvalidToken
	}
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
