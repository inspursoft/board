package controller

import (
	"bytes"
	"context"
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

const maxRetryCount = 300
const buildNumberCacheExpireSecond = time.Duration(maxRetryCount * time.Second)
const toggleBuildingCacheExpireSecond = time.Duration(maxRetryCount * time.Second)

type CIJobCallbackController struct {
	c.BaseController
}

func (j *CIJobCallbackController) Prepare() {
	j.EnableXSRF = false
}

func (j *CIJobCallbackController) BuildNumberCallback() {
	userID := j.Ctx.Input.Param(":userID")
	buildNumber, _ := strconv.Atoi(j.Ctx.Input.Param(":buildNumber"))
	logs.Info("Get build number from CI job callback: %d", buildNumber)
	c.MemoryCache.Put(userID+"_buildNumber", buildNumber, buildNumberCacheExpireSecond)
}

func (j CIJobCallbackController) CustomPushEventPayload() {
	nodeSelection := utils.GetConfig("NODE_SELECTION", "slave")
	data, err := ioutil.ReadAll(j.Ctx.Request.Body)
	if err != nil {
		j.InternalError(err)
	}
	logs.Debug("%s", string(data))
	service.CurrentDevOps().CustomHookPushPayload(data, nodeSelection())
}

func (j CIJobCallbackController) CustomPipelineEventPayload() {
	userID := j.GetString("user_id")
	data, err := ioutil.ReadAll(j.Ctx.Request.Body)
	if err != nil {
		j.InternalError(err)
	}
	logs.Debug("Received pipeline event payload: %s", string(data))
	pipelineID, buildNumber, err := service.CurrentDevOps().CustomHookPipelinePayload(data)
	logs.Info("Got pipeline ID: %d, build number: %d from CI pipeline callback.", pipelineID, buildNumber)
	c.MemoryCache.Put(userID+"_pipelineID", pipelineID, buildNumberCacheExpireSecond)
	c.MemoryCache.Put(userID+"_buildNumber", buildNumber, buildNumberCacheExpireSecond)
}

type CIJobController struct {
	c.BaseController
}

func (j *CIJobController) getStoredID(key string) (int, error) {
	if storedID, ok := c.MemoryCache.Get(strconv.Itoa(int(j.CurrentUser.ID)) + key).(int); ok {
		logs.Info("Get stored ID with key %s from cache is %d", strconv.Itoa(int(j.CurrentUser.ID))+key, storedID)
		return storedID, nil
	}
	return 0, fmt.Errorf("cannot get stored ID from cache currently")
}

type storedContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (j *CIJobController) resolveStoredContext() (ctx context.Context, cancel context.CancelFunc) {
	storedKey := strconv.Itoa(int(j.CurrentUser.ID)) + "_context"
	logs.Info("Resolving context with stored key: %s", storedKey)
	var stored storedContext
	if c.MemoryCache.IsExist(storedKey) {
		logs.Info("Context already exist with stored key: %s", storedKey)
		stored = c.MemoryCache.Get(storedKey).(storedContext)
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		stored = storedContext{ctx, cancel}
		c.MemoryCache.Put(storedKey, stored, toggleBuildingCacheExpireSecond)
		logs.Info("Created and stored context with key: %s", storedKey)
	}
	ctx = stored.ctx
	cancel = stored.cancel
	return
}

func (j *CIJobController) clearBuildNumber() {
	userID := strconv.Itoa(int(j.CurrentUser.ID))
	for _, key := range []string{userID + "_buildNumber", userID + "_pipelineID", userID + "_context"} {
		if c.MemoryCache.IsExist(key) {
			c.MemoryCache.Delete(key)
			logs.Info("Build number stored with key %s has been deleted from cache.", key)
		}
	}
}

func (j *CIJobController) Console() {
	j.clearBuildNumber()
	getBuildNumberRetryCount := 0
	var buildNumber int
	var err error
	for true {
		buildNumber, err = j.getStoredID("_buildNumber")
		if getBuildNumberRetryCount >= maxRetryCount {
			logs.Debug("User canceled current process or exceeded max retry count, will exit.")
			return
		} else if err != nil {
			logs.Debug("Error occurred: %+v", err)
		} else {
			break
		}
		getBuildNumberRetryCount++
		time.Sleep(time.Second * 2)
	}

	jobName := j.GetString("job_name")

	if jobName == "" {
		j.CustomAbortAudit(http.StatusBadRequest, "No job name found.")
		return
	}

	repoName, err := service.ResolveRepoName(jobName, j.CurrentUser.Username)
	if err != nil {
		j.InternalError(err)
		return
	}
	configurations := make(map[string]string)
	configurations["project_name"] = fmt.Sprintf("%s/%s", j.CurrentUser.Username, repoName)
	configurations["job_name"] = repoName
	configurations["build_serial_id"] = strconv.Itoa(buildNumber)
	pipelineID, err := j.getStoredID("_pipelineID")
	if err != nil {
		logs.Warning("Missing to get pipeline ID from store: %+v", err)
	}
	configurations["pipeline_id"] = strconv.Itoa(pipelineID)
	buildConsoleURL, _, _ := service.CurrentDevOps().ResolveHandleURL(configurations)
	logs.Debug("Requested Jenkins build console URL: %s", buildConsoleURL)
	ws, err := websocket.Upgrade(j.Ctx.ResponseWriter, j.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		j.CustomAbortAudit(http.StatusBadRequest, "Not a websocket handshake.")
		return
	} else if err != nil {
		j.CustomAbortAudit(http.StatusInternalServerError, "Cannot setup websocket connection.")
		return
	}
	defer ws.Close()

	req, err := http.NewRequest("GET", buildConsoleURL, nil)
	client := http.Client{}

	buffer := make(chan bytes.Buffer)
	retryCount := 0
	expiryTimer := time.NewTimer(time.Second * 900)
	ticker := time.NewTicker(time.Second * 1)
	var lastPos int

	ctx, cancel := j.resolveStoredContext()
	go func() {
		for range ticker.C {
			resp, err := client.Do(req)
			if err != nil {
				j.InternalError(err)
				logs.Error("Failed to get console response: %+v", err)
				cancel()
			}
			if resp.StatusCode == http.StatusNotFound {
				if retryCount >= maxRetryCount {
					logs.Info("Sent cancel signal to stop WS of console as the retry count has exceeded the maximum.")
					cancel()
				} else {
					retryCount++
					if retryCount%50 == 0 {
						logs.Debug("CI console is not ready at this moment, will retry for next %d request...", retryCount)
					}
					continue
				}
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				j.InternalError(err)
				logs.Error("Sent cancel signal to stop WS of console as failed to read data from response body: %+v", err)
				cancel()
			}
			var partialBuf bytes.Buffer
			partialBuf.Write(bytes.TrimSuffix(data[lastPos:], []byte{'\r', '\n'}))
			buffer <- partialBuf
			partialBuf.Reset()
			lastPos = len(data)
			resp.Body.Close()
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "Finished:") || strings.Contains(line, "Job succeeded") || strings.Contains(line, "Job failed:") {
					logs.Info("Sent cancel signal to stop WS of console as job has reached the end.")
					cancel()
				}
			}
		}
	}()
	for {
		select {
		case content := <-buffer:
			err = ws.WriteMessage(websocket.TextMessage, content.Bytes())
			if err != nil {
				logs.Error("Sent cancel signal to stop WS of console as error occured: %+v", err)
				cancel()
			}
		case <-ctx.Done():
			ticker.Stop()
			logs.Debug("WS is being closed.")
			err = ws.Close()
			if err := ctx.Err(); err != nil {
				logs.Error("Context has canceled with error: %+v", err)
			}
			return
		case <-expiryTimer.C:
			ticker.Stop()
			err = ws.Close()
			logs.Debug("WS is being closed due to timeout.")
			if err != nil {
				logs.Error("Failed to write message: %+v", err)
			}
			return
		}
	}
}

func (j *CIJobController) Stop() {
	_, cancel := j.resolveStoredContext()
	cancel()
	logs.Info("Sent cancel signal to stop WS of console as user prompted.")

	lastBuildNumber, err := j.getStoredID("_buildNumber")
	if err != nil {
		logs.Error("Failed to get job number: %+v", err)
		cancel()
		logs.Info("Sent cancel signal to stop WS of console as miss job build number.")
		return
	}
	jobName := j.GetString("job_name")
	if jobName == "" {
		j.CustomAbortAudit(http.StatusBadRequest, "No job name found.")
		return
	}
	repoName, err := service.ResolveRepoName(jobName, j.CurrentUser.Username)
	if err != nil {
		j.InternalError(err)
		return
	}
	lastPipelineID, err := j.getStoredID("_pipelineID")
	if err != nil {
		logs.Warning("Missing to get pipeline ID from store: %+v", err)
	}
	configurations := make(map[string]string)
	configurations["project_name"] = fmt.Sprintf("%s/%s", j.CurrentUser.Username, repoName)
	configurations["job_name"] = repoName
	configurations["build_serial_id"] = j.GetString("build_serial_id", strconv.Itoa(lastBuildNumber))
	configurations["pipeline_id"] = strconv.Itoa(lastPipelineID)
	configurations["repo_token"] = j.CurrentUser.RepoToken
	_, stopBuildURL, err := service.CurrentDevOps().ResolveHandleURL(configurations)
	if err != nil {
		j.InternalError(err)
		return
	}
	logs.Debug("Requested stop CI build URL: %s", stopBuildURL)
	resp, err := http.PostForm(stopBuildURL, nil)
	if err != nil {
		j.InternalError(err)
		return
	}
	defer func() {
		resp.Body.Close()
		j.clearBuildNumber()
	}()
	logs.Debug("Response status of stopping CI jobs: %d", resp.StatusCode)

	j.ServeStatus(resp.StatusCode, "")
}

func (j *CIJobController) ResetVariable() {
	repoName := j.GetString("repo_name")
	configurations := make(map[string]string)
	configurations["repo_name"] = repoName
	configurations["repo_token"] = j.CurrentUser.RepoToken
	err := service.CurrentDevOps().ResetOpts(configurations)
	if err != nil {
		j.InternalError(err)
		return
	}
}
