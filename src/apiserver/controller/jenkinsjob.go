package controller

import (
	"bytes"
	"io"
	"net/http"
	"text/template"

	"github.com/astaxie/beego/logs"
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
	req, err := http.NewRequest("GET", consoleURL.String(), nil)
	if err != nil {
		j.internalError(err)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		j.internalError(err)
		return
	}
	_, err = io.Copy(j.Ctx.ResponseWriter, resp.Body)
	if err != nil {
		j.internalError(err)
		return
	}
}
