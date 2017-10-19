package controller

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

const jenkinsBuildConsoleURL = "http://jenkins:8080/job/{{.JobName}}/{{.BuildSerialID}}/consoleText"
const jenkinsLastBuildNumberURL = "http://jenkins:8080/job/{{.JobName}}/lastBuild/buildNumber"

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
	j.isProjectAdmin = (j.currentUser.ProjectAdmin == 1)
	if !j.isProjectAdmin {
		j.customAbort(http.StatusForbidden, "Insufficient privileges for manipulating Git repos.")
		return
	}
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

func (j *JenkinsJobController) Console() {
	jobName := j.GetString("job_name")
	if jobName == "" {
		j.customAbort(http.StatusBadRequest, "No job name found.")
		return
	}
	buildSerialID := j.GetString("build_serial_id", "lastBuild")

	lastBuildNumberURL, err := generateURL(jenkinsLastBuildNumberURL, jobConsole{JobName: jobName})
	if err != nil {
		j.internalError(err)
		return
	}

	respBuildNumber, err := http.Get(lastBuildNumberURL)
	if err != nil {
		j.internalError(err)
		return
	}
	defer respBuildNumber.Body.Close()
	respBuildNumberData, err := ioutil.ReadAll(respBuildNumber.Body)
	if err != nil {
		j.internalError(err)
		return
	}
	if respBuildNumberData != nil {
		buildSerialID = string(respBuildNumberData)
	}

	buildConsoleURL, err := generateURL(jenkinsBuildConsoleURL, jobConsole{JobName: jobName, BuildSerialID: buildSerialID})
	if err != nil {
		j.internalError(err)
		return
	}

	logs.Debug("Requested Jenkins build console URL: %s", buildConsoleURL)
	logs.Debug("Requested Jenkins build number URL: %s", lastBuildNumberURL)

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
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for range ticker.C {
			resp, err := client.Do(req)
			if err != nil {
				j.internalError(err)
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				j.internalError(err)
				return
			}
			buffer <- data
			resp.Body.Close()

			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "Finished:") {
					ticker.Stop()
					done <- true
				}
			}
		}
	}()

	for {
		select {
		case content := <-buffer:
			err = ws.WriteMessage(websocket.TextMessage, content)
		case <-done:
			err = ws.Close()
			logs.Debug("WS is being closed.")
		case <-time.After(time.Second * 80):
			err = ws.Close()
			logs.Debug("WS is being closed due to timeout.")
		}
		if err != nil {
			logs.Error("Failed to write message: %+v", err)
		}
	}
}
