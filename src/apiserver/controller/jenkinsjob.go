package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

const jenkinsBuildConsoleTemplateURL = "http://jenkins:8080/job/{{.JobName}}/{{.BuildSerialID}}/consoleText"
const jenkinsStopBuildTemplateURL = "http://jenkins:8080/job/{{.JobName}}/{{.BuildSerialID}}/stop"
const maxRetryCount = 600

type jobConsole struct {
	JobName       string `json:"job_name"`
	BuildSerialID string `json:"build_serial_id"`
}

type JenkinsJobController struct {
	baseController
}

func (j *JenkinsJobController) Prepare() {
	user := j.getCurrentUser()
	if user == nil {
		j.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	j.currentUser = user
}

func generateURL(rawTemplate string, data interface{}) (string, error) {
	t := template.Must(template.New("").Parse(rawTemplate))
	var targetURL bytes.Buffer
	err := t.Execute(&targetURL, data)
	if err != nil {
		return "", err
	}
	return targetURL.String(), nil
}

func (j *JenkinsJobController) getBuildNumber() (int, error) {
	if buildNumber, ok := memoryCache.Get(strconv.Itoa(int(j.currentUser.ID)) + "_lastBuildNumber").(int); ok {
		logs.Info("Get build number with key %s from cache is %d", strconv.Itoa(int(j.currentUser.ID))+"_lastBuildNumber", buildNumber)
		return buildNumber, nil
	}
	return 0, fmt.Errorf("failed to get build number from cache")
}

func (j *JenkinsJobController) GetLastBuildNumber() {
	lastBuildNumber, err := j.getBuildNumber()
	if err != nil {
		j.internalError(err)
		return
	}
	j.Data["json"] = lastBuildNumber
	j.ServeJSON()
}

func (j *JenkinsJobController) Console() {
	jobName := j.GetString("job_name")
	if jobName == "" {
		j.customAbort(http.StatusBadRequest, "No job name found.")
		return
	}
	query := jobConsole{JobName: jobName}
	buildNumber, err := j.getBuildNumber()
	if err != nil {
		j.internalError(err)
		return
	}
	query.BuildSerialID = strconv.Itoa(buildNumber)
	buildConsoleURL, err := generateURL(jenkinsBuildConsoleTemplateURL, query)
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

	go func(done chan bool) {
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
	}(done)

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
	jobName := j.GetString("job_name")
	if jobName == "" {
		j.customAbort(http.StatusBadRequest, "No job name found.")
		return
	}
	query := jobConsole{JobName: jobName}
	lastBuildNumber, err := j.getBuildNumber()
	if err != nil {
		j.customAbort(http.StatusInternalServerError, "Failed to get job number.")
		return
	}
	query.BuildSerialID = j.GetString("build_serial_id", strconv.Itoa(lastBuildNumber))
	stopBuildURL, err := generateURL(jenkinsStopBuildTemplateURL, query)
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
	defer resp.Body.Close()
	logs.Debug("Response status of stopping Jenkins jobs: %d", resp.StatusCode)
	j.serveStatus(resp.StatusCode, "")
}
