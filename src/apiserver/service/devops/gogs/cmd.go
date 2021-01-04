package gogs

import (
	"fmt"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"
)

var gogitsBaseURL = utils.GetConfig("GOGITS_BASE_URL")
var JenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")
var maxRetryCount = 30

type createAccessTokenOption struct {
	Name string `json:"name" binding:"Required"`
}

type createKeyOption struct {
	Title string `json:"title" binding:"Required"`
	Key   string `json:"key" binding:"Required"`
}

type createRepoOption struct {
	Name        string `json:"name" binding:"Required;AlphaDashDot;MaxSize(100)"`
	Description string `json:"description" binding:"MaxSize(255)"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
	Gitignores  string `json:"gitignores"`
	License     string `json:"license"`
	Readme      string `json:"readme"`
}

type createIssueOption struct {
	Title      string  `json:"title" binding:"Required"`
	Body       string  `json:"body"`
	Assignee   string  `json:"assignee"`
	Milestone  int64   `json:"milestone"`
	Labels     []int64 `json:"labels"`
	Closed     bool    `json:"closed"`
	LabelIDs   string  `json:"label_ids"`
	AssigneeID int64   `json:"assignee_id"`
}

type createIssueCommentOption struct {
	Body string `json:"body" binding:"Required"`
}

type createHookOption struct {
	Type   string            `json:"type" binding:"Required"`
	Config map[string]string `json:"config" binding:"Required"`
	Events []string          `json:"events"`
	Active bool              `json:"active"`
}

type AccessToken struct {
	Name string `json:"name"`
	Sha1 string `json:"sha1"`
}

type gogsHandler struct {
	username string
	password string
	token    string
}

type PullRequestInfo struct {
	IssueID    int64 `json:"issue_id"`
	Index      int64 `json:"index"`
	HasCreated bool  `json:"has_created"`
}

func pingGogitsService() {
	pingURL := fmt.Sprintf("%s", gogitsBaseURL())
	for i := 0; i < maxRetryCount; i++ {
		logs.Debug("Ping Gogits server %d time(s)...", i+1)
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
			logs.Info("Successful connected to the Gogits service.")
			break
		}
		time.Sleep(time.Second)
	}
}

func NewGogsHandler(username, token string) *gogsHandler {
	pingGogitsService()
	return &gogsHandler{
		username: username,
		token:    token,
	}
}

func userExists(username string) (bool, error) {
	pingGogitsService()
	logs.Info("Requesting Gogits API of user exists ...")
	err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/api/v1/users/%s", gogitsBaseURL(), username), nil, nil, func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("user: %s already exists", username)
		}
		return nil
	})
	if err != nil {
		return true, nil
	}
	return false, nil
}

func SignUp(user model.User) error {
	pingGogitsService()
	userExists, err := userExists(user.Username)
	if err != nil {
		logs.Error("Error occurred while checking user existing: %+v", err)
		return nil
	}
	if userExists {
		logs.Info("User: %s already exists in Gogits.", user.Username)
		return nil
	}
	logs.Info("Requesting Gogits API of sign up ...")
	return utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/user/sign_up", gogitsBaseURL()), func(req *http.Request) error {
		req.Header = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
		formData := url.Values{}
		formData.Set("user_name", user.Username)
		formData.Set("password", user.Password)
		formData.Set("retype", user.Password)
		formData.Set("email", user.Email)
		req.URL.RawQuery = formData.Encode()
		return nil
	}, nil, nil)
}

func CreateAccessToken(username, password string) (*AccessToken, error) {
	opt := createAccessTokenOption{Name: "ACCESS-TOKEN"}
	var token AccessToken
	logs.Info("Requesting Gogits API of create access token ...")
	err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/users/%s/tokens", gogitsBaseURL(), username), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json", "application/form-data"},
			"Authorization": []string{"Basic " + utils.BasicAuthEncode(username, password)},
		}
		return nil
	}, &opt, func(req *http.Request, resp *http.Response) error {
		return utils.UnmarshalToJSON(resp.Body, &token)
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (g *gogsHandler) getAccessHeader() http.Header {
	return http.Header{
		"content-type":  []string{"application/json"},
		"Authorization": []string{"token " + g.token},
	}
}

func (g *gogsHandler) CreatePublicKey(title, publicKey string) error {
	opt := createKeyOption{
		Title: title,
		Key:   publicKey,
	}
	return utils.SimplePostRequestHandle(fmt.Sprintf("%s/api/v1/user/keys", gogitsBaseURL()), g.getAccessHeader(), &opt)
}

func (g *gogsHandler) DeletePublicKey(userID int64) error {
	return utils.SimpleDeleteRequestHandle(fmt.Sprintf("%s/api/v1/user/keys/%d", gogitsBaseURL(), userID), nil)
}

func (g *gogsHandler) CreateRepo(repoName string) error {
	var opt = createRepoOption{
		Name:        repoName,
		Description: "Created by Board API for DevOps.",
	}
	return utils.SimplePostRequestHandle(fmt.Sprintf("%s/api/v1/user/repos", gogitsBaseURL()), g.getAccessHeader(), &opt)
}

func (g *gogsHandler) DeleteRepo(username, repoName string) error {
	return utils.SimpleDeleteRequestHandle(fmt.Sprintf("%s/api/v1/repos/%s/%s", gogitsBaseURL(), username, repoName), g.getAccessHeader())
}

func (g *gogsHandler) ForkRepo(ownerName, baseRepoName, forkRepoName, description string) error {
	logs.Info("Requesting Gogits API of fork repo ...")
	return utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/repos/%s/%s/forks", gogitsBaseURL(), ownerName, baseRepoName), func(req *http.Request) error {
		req.Header = http.Header{
			"Authorization": []string{"token " + g.token},
		}
		formData := url.Values{}
		formData.Set("repo_name", forkRepoName)
		formData.Set("description", description)
		req.URL.RawQuery = formData.Encode()
		return nil
	}, nil, nil)
}

func (g *gogsHandler) CreatePullRequest(ownerName, baseRepoName, title, content, compareInfo string) (*PullRequestInfo, error) {
	var opt = createIssueOption{
		Title: title,
		Body:  content,
	}
	var prInfo PullRequestInfo
	logs.Info("Requesting Gogits API of create pull request ...")
	err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/repos/%s/%s/pull-request/%s", gogitsBaseURL(), ownerName, baseRepoName, compareInfo), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"token " + g.token},
		}
		return nil
	}, &opt, func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("unexpected error occurred with response status code: %d", resp.StatusCode)
		}
		err := utils.UnmarshalToJSON(resp.Body, &prInfo)
		if err != nil {
			return err
		}
		if &prInfo != nil {
			prInfo.HasCreated = (resp.StatusCode == http.StatusConflict)
		}
		logs.Info("Requested Gogits create pull request with response status code: %d", resp.StatusCode)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &prInfo, nil
}

func (g *gogsHandler) CreateIssueComment(ownerName string, baseRepoName string, issueIndex int64, comment string) error {
	var opt = createIssueCommentOption{
		Body: comment,
	}
	logs.Info("Requesting Gogits comment issue ...")
	return utils.SimplePostRequestHandle(fmt.Sprintf("%s/api/v1/repos/%s/%s/issues/%d/comments", gogitsBaseURL(), ownerName, baseRepoName, issueIndex), g.getAccessHeader(), &opt)
}

func (g *gogsHandler) CreateHook(ownerName string, repoName string, hookURL string) error {
	config := make(map[string]string)
	config["url"] = hookURL
	config["content_type"] = "json"

	opt := createHookOption{
		Type:   "gogs",
		Config: config,
		Events: []string{"push"},
		Active: true,
	}
	logs.Info("Requesting Gogits API of create hook ...")
	return utils.SimplePostRequestHandle(fmt.Sprintf("%s/api/v1/repos/%s/%s/hooks", gogitsBaseURL(), ownerName, repoName), g.getAccessHeader(), &opt)
}

func (g *gogsHandler) DeleteUser(username string) error {
	return utils.SimpleDeleteRequestHandle(fmt.Sprintf("%s/api/v1/admin/users/%s", gogitsBaseURL(), username), g.getAccessHeader())
}
