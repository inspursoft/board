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
var jenkinsfileRepoURL = utils.GetConfig("JENKINSFILE_REPO_URL")
var maxRetryCount = 240
var seedJobName = "base"
var jenkinsHostIP = utils.GetConfig("JENKINS_HOST_IP")
var jenkinsHostPort = utils.GetConfig("JENKINS_HOST_PORT")
var jenkinsNodeIP = utils.GetConfig("JENKINS_NODE_IP")
var kvmRegistryPort = utils.GetConfig("KVM_REGISTRY_PORT")

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
	return utils.SimpleGetRequestHandle(fmt.Sprintf("%s/job/%s/buildWithParameters?F00=%s&&F01=%s&F02=%s&F03=%s&F04=%s&F05=%s&F06=%s&F07=%s", jenkinsBaseURL(), seedJobName, projectName, jenkinsfileRepoURL(), username, email, jenkinsHostIP(), jenkinsHostPort(), jenkinsNodeIP(), kvmRegistryPort()))
}
