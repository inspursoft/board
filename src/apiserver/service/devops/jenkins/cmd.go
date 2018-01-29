package jenkins

import (
	"fmt"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

type jenkinsHandler struct{}

func NewJenkinsHandler() *jenkinsHandler {
	return &jenkinsHandler{}
}

func (j *jenkinsHandler) CreateJob(projectName string) error {
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/createItem?name=%s&mode=copy&from=base", jenkinsBaseURL, projectName), func(req *http.Request) error {
		req.Header = http.Header{
			"Authorization": []string{"token " + utils.BasicAuthEncode("admin", "admin")},
		}
		return nil
	}, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("Internal error: %+v", err)
		}
		logs.Info("Requested Jenkins clone job with response status code: %d", resp.StatusCode)
	}
	return nil
}
