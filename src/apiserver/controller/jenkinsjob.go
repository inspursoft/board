package controller

import (
	"fmt"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

const jenkinsBuildConsoleTemplateURL = "%s/job/{{.JobName}}/{{.BuildSerialID}}/consoleText"
const jenkinsStopBuildTemplateURL = "%s/job/{{.JobName}}/{{.BuildSerialID}}/stop"
const maxRetryCount = 600
const buildNumberCacheExpireSecond = time.Duration(maxRetryCount * time.Second)
const toggleBuildingCacheExpireSecond = time.Duration(maxRetryCount * time.Second)

type jobConsole struct {
	JobName       string `json:"job_name"`
	BuildSerialID string `json:"build_serial_id"`
}

type JenkinsJobCallbackController struct {
	baseController
}

func (j *JenkinsJobCallbackController) BuildNumberCallback() {
	userID := j.Ctx.Input.Param(":userID")
	buildNumber, _ := strconv.Atoi(j.Ctx.Input.Param(":buildNumber"))
	logs.Info("Get build number from Jenkins job callback: %d", buildNumber)
	memoryCache.Put(userID+"_buildNumber", buildNumber, buildNumberCacheExpireSecond)
}

type JenkinsJobController struct {
	baseController
}

func (j *JenkinsJobController) getBuildNumber() (int, error) {
	if buildNumber, ok := memoryCache.Get(strconv.Itoa(int(j.currentUser.ID)) + "_buildNumber").(int); ok {
		logs.Info("Get build number with key %s from cache is %d", strconv.Itoa(int(j.currentUser.ID))+"_lastBuildNumber", buildNumber)
		return buildNumber, nil
	}
	return 0, fmt.Errorf("cannot get build number from cache currently")
}

func (j *JenkinsJobController) clearBuildNumber() {
	key := strconv.Itoa(int(j.currentUser.ID)) + "_buildNumber"
	if memoryCache.IsExist(key) {
		memoryCache.Delete(key)
		logs.Info("Build number stored with key %s has been deleted from cache.", key)
	}
}

func (j *JenkinsJobController) toggleBuild(status bool) {
	logs.Info("Set building signal as %+v currently.", status)
	memoryCache.Put(strconv.Itoa(int(j.currentUser.ID))+"_buildSignal", status, toggleBuildingCacheExpireSecond)
}

func (j *JenkinsJobController) getBuildSignal() bool {
	key := strconv.Itoa(int(j.currentUser.ID)) + "_buildSignal"
	if buildingSignal, ok := memoryCache.Get(key).(bool); ok {
		return buildingSignal
	}
	return false
}

func (j *JenkinsJobController) Console() {
	j.clearBuildNumber()
	j.toggleBuild(true)
	getBuildNumberRetryCount := 0
	var buildNumber int
	var err error
	for true {
		buildNumber, err = j.getBuildNumber()
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
		j.customAbort(http.StatusBadRequest, "No job name found.")
		return
	}
	query := jobConsole{JobName: jobName}
	query.BuildSerialID = strconv.Itoa(buildNumber)
	buildConsoleURL, err := utils.GenerateURL(fmt.Sprintf(jenkinsBuildConsoleTemplateURL, jenkinsBaseURL()), query)
	if err != nil {
		j.internalError(err)
		return
	}
	logs.Debug("Requested Jenkins build console URL: %s", buildConsoleURL)
	ws, err := websocket.Upgrade(j.Ctx.ResponseWriter, j.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		j.customAbort(http.StatusBadRequest, "Not a websocket handshake.")
		return
	} else if err != nil {
		j.customAbort(http.StatusInternalServerError, "Cannot setup websocket connection.")
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
	go func() {
		for range ticker.C {
			resp, err := client.Do(req)
			if err != nil {
				j.internalError(err)
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
						logs.Debug("Jenkins console is not ready at this moment, will retry for next %d request...", retryCount)
					}
					continue
				}
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				j.internalError(err)
				logs.Error("Failed to read data from response body: %+v", err)
				done <- true
				return
			}
			buffer <- data
			resp.Body.Close()
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "Finished:") {
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

func (j *JenkinsJobController) Stop() {
	lastBuildNumber, err := j.getBuildNumber()
	if err != nil {
		logs.Error("Failed to get job number: %+v", err)
		j.toggleBuild(false)
		return
	}

	jobName := j.GetString("job_name")
	if jobName == "" {
		j.customAbort(http.StatusBadRequest, "No job name found.")
		return
	}
	query := jobConsole{JobName: jobName}
	query.BuildSerialID = j.GetString("build_serial_id", strconv.Itoa(lastBuildNumber))
	stopBuildURL, err := utils.GenerateURL(fmt.Sprintf(jenkinsStopBuildTemplateURL, jenkinsBaseURL()), query)
	if err != nil {
		j.internalError(err)
		return
	}
	logs.Debug("Requested stop Jenkins build URL: %s", stopBuildURL)
	resp, err := http.PostForm(stopBuildURL, nil)
	if err != nil {
		j.internalError(err)
		return
	}
	defer func() {
		resp.Body.Close()
		j.clearBuildNumber()
	}()
	logs.Debug("Response status of stopping Jenkins jobs: %d", resp.StatusCode)
	j.serveStatus(resp.StatusCode, "")
}
