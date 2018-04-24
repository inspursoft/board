package gogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"
)

var gogitsBaseURL = utils.GetConfig("GOGITS_BASE_URL")
var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")
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

func NewGogsHandler(username, token string) *gogsHandler {
	for i := 0; i < maxRetryCount; i++ {
		logs.Debug("Ping Gogits server %d time(s)...", i+1)
		resp, err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s", gogitsBaseURL()), nil, nil)
		if err != nil {
			logs.Error("Failed to request Gogits server: %+v", err)
		}
		if resp != nil {
			if resp.StatusCode <= 400 {
				break
			}
		} else if i == maxRetryCount-1 {
			logs.Warn("Failed to ping Gogits due to exceed max retry count.")
		}
		time.Sleep(time.Second)
	}
	return &gogsHandler{
		username: username,
		token:    token,
	}
}

func userExists(username string) (bool, error) {
	resp, err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/api/v1/users/%s", gogitsBaseURL(), username), nil, nil)
	if err != nil {
		return false, err
	}
	if resp != nil {
		defer resp.Body.Close()
		return (resp.StatusCode != http.StatusNotFound), nil
	}
	return false, nil
}

func SignUp(user model.User) error {
	userExists, err := userExists(user.Username)
	if err != nil {
		logs.Error("Error occurred while checking user existing: %+v", err)
		return nil
	}
	if userExists {
		logs.Info("User: %s already exists in Gogits.", user.Username)
		return nil
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/user/sign_up", gogitsBaseURL()), func(req *http.Request) error {
		req.Header = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
		formData := url.Values{}
		formData.Set("user_name", user.Username)
		formData.Set("password", user.Password)
		formData.Set("retype", user.Password)
		formData.Set("email", user.Email)
		req.URL.RawQuery = formData.Encode()
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
		logs.Info("Requested Gogits sign up with response status code: %d", resp.StatusCode)
	}
	return nil
}

func CreateAccessToken(username, password string) (*AccessToken, error) {
	opt := createAccessTokenOption{Name: "ACCESS-TOKEN"}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/users/%s/tokens", gogitsBaseURL(), username), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"Basic " + utils.BasicAuthEncode(username, password)},
		}
		return nil
	}, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if resp != nil {
		output, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		logs.Debug("token output: %s", string(output))
		var token AccessToken
		err = json.Unmarshal(output, &token)
		if err != nil {
			return nil, err
		}
		return &token, nil
	}
	return nil, nil
}

func (g *gogsHandler) CreatePublicKey(title, publicKey string) error {
	opt := createKeyOption{
		Title: title,
		Key:   publicKey,
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/user/keys", gogitsBaseURL()), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"token " + g.token},
		}
		return nil
	}, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
		logs.Info("Requested Gogits create public key with response status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *gogsHandler) DeletePublicKey(userID int64) error {
	resp, err := utils.RequestHandle(http.MethodDelete, fmt.Sprintf("%s/api/v1/user/keys/%d", gogitsBaseURL(), userID), nil, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
		logs.Info("Requested Gogits delete public key with response status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *gogsHandler) CreateRepo(repoName string) error {
	var opt = createRepoOption{
		Name:        repoName,
		Description: "Created by Board API for DevOps.",
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/user/repos", gogitsBaseURL()), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"token " + g.token},
		}
		return nil
	}, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("Internal error: %+v", err)
		}
		logs.Info("Requested Gogits create repo with response status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *gogsHandler) DeleteRepo(repoName string) error {
	resp, err := utils.RequestHandle(http.MethodDelete, fmt.Sprintf("%s/api/v1/repos/%s/%s", gogitsBaseURL(), g.username, repoName), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"token " + g.token},
		}
		return nil
	}, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("Internal error: %+v", err)
		}
		logs.Info("Requested Gogits delete repo with response status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *gogsHandler) ForkRepo(ownerName, baseRepoName, forkRepoName, description string) error {
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/repos/%s/%s/forks", gogitsBaseURL(), ownerName, baseRepoName), func(req *http.Request) error {
		req.Header = http.Header{
			"Authorization": []string{"token " + g.token},
		}
		formData := url.Values{}
		formData.Set("repo_name", forkRepoName)
		formData.Set("description", description)
		req.URL.RawQuery = formData.Encode()
		return nil
	}, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("Internal error: %+v", err)
		}
		data, _ := ioutil.ReadAll(resp.Body)
		logs.Debug(string(data))
		logs.Info("Requested Gogits fork with response status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *gogsHandler) CreatePullRequest(ownerName, baseRepoName, title, content, compareInfo string) (*PullRequestInfo, error) {
	var opt = createIssueOption{
		Title: title,
		Body:  content,
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/repos/%s/%s/pull-request/%s", gogitsBaseURL(), ownerName, baseRepoName, compareInfo), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"token " + g.token},
		}
		return nil
	}, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if resp != nil {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, fmt.Errorf("Internal error: %+v", err)
		}
		data, _ := ioutil.ReadAll(resp.Body)
		var prInfo PullRequestInfo
		json.Unmarshal(data, &prInfo)
		if &prInfo != nil {
			prInfo.HasCreated = (resp.StatusCode == http.StatusConflict)
		}
		logs.Info("Requested Gogits create pull request with response status code: %d", resp.StatusCode)
		return &prInfo, nil
	}
	return nil, nil
}

func (g *gogsHandler) CreateIssueComment(ownerName string, baseRepoName string, issueIndex int64, comment string) error {
	var opt = createIssueCommentOption{
		Body: comment,
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/repos/%s/%s/issues/%d/comments", gogitsBaseURL(), ownerName, baseRepoName, issueIndex), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"token " + g.token},
		}
		return nil
	}, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("Internal error: %+v", err)
		}
		data, _ := ioutil.ReadAll(resp.Body)
		logs.Debug(string(data))
		logs.Info("Requested Gogits comment issue with response status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *gogsHandler) CreateHook(ownerName string, baseRepoName string) error {

	config := make(map[string]string)
	config["url"] = fmt.Sprintf("%s/generic-webhook-trigger/invoke?token=123456", jenkinsBaseURL())
	config["content_type"] = "json"

	opt := createHookOption{
		Type:   "gogs",
		Config: config,
		Events: []string{"push"},
		Active: true,
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/repos/%s/%s/hooks", gogitsBaseURL(), ownerName, baseRepoName), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"token " + g.token},
		}
		return nil
	}, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("Internal error: %+v", err)
		}
		data, _ := ioutil.ReadAll(resp.Body)
		logs.Debug(string(data))
		logs.Info("Requested Gogits create hook with response status code: %d", resp.StatusCode)
	}
	return nil
}
