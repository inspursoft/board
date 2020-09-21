package controller

import (
	"bytes"
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

func (j *CIJobController) clearBuildNumber() {
	userID := strconv.Itoa(int(j.CurrentUser.ID))
	for _, key := range []string{userID + "_buildNumber", userID + "_pipelineID"} {
		if c.MemoryCache.IsExist(key) {
			c.MemoryCache.Delete(key)
			logs.Info("Build number stored with key %s has been deleted from cache.", key)
		}
	}
}

func (j *CIJobController) toggleBuild(status bool) {
	logs.Info("Set building signal as %+v currently.", status)
	c.MemoryCache.Put(strconv.Itoa(int(j.CurrentUser.ID))+"_buildSignal", status, toggleBuildingCacheExpireSecond)
}

func (j *CIJobController) getBuildSignal() bool {
	key := strconv.Itoa(int(j.CurrentUser.ID)) + "_buildSignal"
	if buildingSignal, ok := c.MemoryCache.Get(key).(bool); ok {
		return buildingSignal
	}
	return false
}

func (j *CIJobController) Console() {
	j.clearBuildNumber()
	j.toggleBuild(true)
	getBuildNumberRetryCount := 0
	var buildNumber int
	var err error
	for true {
		buildNumber, err = j.getStoredID("_buildNumber")
		if j.getBuildSignal() == false || getBuildNumberRetryCount >= maxRetryCount {
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

	buffer := make(chan []byte, 1024)
	done := make(chan bool)

	retryCount := 0
	expiryTimer := time.NewTimer(time.Second * 900)
	ticker := time.NewTicker(time.Second * 1)
	var lastPos int
	go func() {
		for range ticker.C {
			resp, err := client.Do(req)
			if err != nil {
				j.InternalError(err)
				logs.Error("Failed to get console response: %+v", err)
				done <- true
				return
			}
			if resp.StatusCode == http.StatusNotFound {
				if retryCount >= maxRetryCount {
					done <- true
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
				logs.Error("Failed to read data from response body: %+v", err)
				done <- true
				return
			}
			buffer <- bytes.TrimSuffix(data[lastPos:], []byte{'\r', '\n'})
			lastPos = len(data)
			resp.Body.Close()
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "Finished:") || strings.Contains(line, "Job succeeded") || strings.Contains(line, "Job failed:") {
					done <- true
				}
			}
		}
	}()

	for {
		select {
		case content := <-buffer:
			err = ws.WriteMessage(websocket.TextMessage, content)
			if err != nil {
				done <- true
			}
		case <-done:
			ticker.Stop()
			err = ws.Close()
			logs.Debug("WS is being closed.")
		case <-expiryTimer.C:
			ticker.Stop()
			err = ws.Close()
			logs.Debug("WS is being closed due to timeout.")
		}
		if err != nil {
			logs.Error("Failed to write message: %+v", err)
		}
	}
}

func (j *CIJobController) Stop() {
	lastBuildNumber, err := j.getStoredID("_buildNumber")
	if err != nil {
		logs.Error("Failed to get job number: %+v", err)
		j.toggleBuild(false)
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
