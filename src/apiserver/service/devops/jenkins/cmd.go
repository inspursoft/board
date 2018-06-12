package jenkins

import (
	"fmt"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"
)

var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")
var gogitsBaseURL = utils.GetConfig("GOGITS_BASE_URL")
var maxRetryCount = 120
var seedJobName = "base"

type jenkinsHandler struct{}

func NewJenkinsHandler() *jenkinsHandler {
	pingURL := fmt.Sprintf("%s/job/%s", jenkinsBaseURL(), seedJobName)
	for i := 0; i < maxRetryCount; i++ {
		logs.Debug("Ping Jenkins server %d time(s)...", i+1)
		if i == maxRetryCount-1 {
			logs.Warn("Failed to ping Gogits due to exceed max retry count.")
			break
		}
		err := utils.RequestHandle(http.MethodGet, pingURL, nil, nil,
			func(req *http.Request, resp *http.Response) error {
				if resp.StatusCode <= 400 {
					return nil
				}
				return fmt.Errorf("Requested URL %s with unexpected response: %d", pingURL, resp.StatusCode)
			})
		if err == nil {
			logs.Info("Successful connected to the Jenkins service.")
			break
		}
		time.Sleep(time.Second)
	}
	return &jenkinsHandler{}
}

func (j *jenkinsHandler) CreateJobWithParameter(projectName, username, email string) error {
	repoURL := fmt.Sprintf("%s/%s/%s.git", gogitsBaseURL(), username, projectName)
	return utils.SimpleGetRequestHandle(fmt.Sprintf("%s/job/%s/buildWithParameters?F00=%s&&F01=%s&F02=%s&F03=%s", jenkinsBaseURL(), seedJobName, projectName, repoURL, username, email))
}
