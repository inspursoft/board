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

const jenkinsJobConsoleURL = "http://jenkins:8080/job/{{.JobName}}/{{.BuildSerialID}}/consoleText"

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
		j.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	j.currentUser = user
	j.isProjectAdmin = (j.currentUser.ProjectAdmin == 1)
	if !j.isProjectAdmin {
		j.CustomAbort(http.StatusForbidden, "Insufficient privileges for manipulating Git repos.")
		return
	}
}

func (j *JenkinsJobController) Console() {
	jobName := j.GetString("job_name")
	if jobName == "" {
		j.CustomAbort(http.StatusBadRequest, "No job name found.")
		return
	}
	buildSerialID := j.GetString("build_serial_id", "lastBuild")
	templates := template.Must(template.New("jobConsole").Parse(jenkinsJobConsoleURL))
	var consoleURL bytes.Buffer
	err := templates.Execute(&consoleURL, jobConsole{JobName: jobName, BuildSerialID: buildSerialID})
	if err != nil {
		j.internalError(err)
		return
	}
	logs.Debug("Requested Jenkins console URL: %s", consoleURL.String())

	ws, err := websocket.Upgrade(j.Ctx.ResponseWriter, j.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		j.CustomAbort(http.StatusBadRequest, "Not a websocket handshake.")
		return
	} else if err != nil {
		j.CustomAbort(http.StatusInternalServerError, "Cannot setup websocket connection.")
		return
	}
	defer ws.Close()

	req, err := http.NewRequest("GET", consoleURL.String(), nil)
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
			logs.Debug("WS is writing.")
		case <-done:
			err = ws.Close()
			logs.Debug("WS is being closed.")
		}
		if err != nil {
			logs.Error("Failed to write message: %+v", err)
		}
	}
}
