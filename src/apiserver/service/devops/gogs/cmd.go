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

	"github.com/astaxie/beego/logs"
)

var gogsBaseURL = utils.GetConfig("GOGS_BASE_URL")

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

type accessToken struct {
	Name string `json:"name"`
	Sha1 string `json:"sha1"`
}

type gogsHandler struct {
	username string
	token    string
}

func NewGogsHandler(username, token string) *gogsHandler {
	return &gogsHandler{
		username: username,
		token:    token,
	}
}

func userExists(username string) (bool, error) {
	resp, err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/api/v1/users/%s", gogsBaseURL(), username), nil, nil)
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
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/user/sign_up", gogsBaseURL()), func(req *http.Request) error {
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

func CreateAccessToken(username, password string) (*accessToken, error) {
	opt := createAccessTokenOption{Name: "ACCESS-TOKEN"}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/users/%s/tokens", gogsBaseURL(), username), func(req *http.Request) error {
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
		var token accessToken
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
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/user/keys", gogsBaseURL()), func(req *http.Request) error {
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
	resp, err := utils.RequestHandle(http.MethodDelete, fmt.Sprintf("%s/api/v1/user/keys/%d", gogsBaseURL(), userID), nil, nil)
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
	resp, err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/user/repos", gogsBaseURL()), func(req *http.Request) error {
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
	resp, err := utils.RequestHandle(http.MethodDelete, fmt.Sprintf("%s/api/v1/repos/%s/%s", gogsBaseURL(), g.username, repoName), func(req *http.Request) error {
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
