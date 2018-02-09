package controller

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"encoding/json"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"bytes"

	"git/inspursoft/board/src/apiserver/service"

	"strconv"

	"net/url"

	"strings"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var conf config.Configer
var tokenServerURL *url.URL
var tokenCacheExpireSeconds int
var memoryCache cache.Cache

var errInvalidToken = errors.New("error for invalid token")

var apiServerURL = utils.GetConfig("API_SERVER_URL")

var kubeMasterURL = utils.GetConfig("KUBE_MASTER_URL")
var registryURL = utils.GetConfig("REGISTRY_URL")
var registryBaseURI = utils.GetConfig("REGISTRY_BASE_URI")
var authMode = utils.GetConfig("AUTH_MODE")

var baseRepoPath = utils.GetConfig("BASE_REPO_PATH")
var gogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

type baseController struct {
	beego.Controller
	currentUser    *model.User
	token          string
	isSysAdmin     bool
	isExternalAuth bool
	repoPath       string
}

func (b *baseController) Render() error {
	return nil
}

func (b *baseController) resolveBody() ([]byte, error) {
	data, err := ioutil.ReadAll(b.Ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *baseController) resolveRepoPath() {
	projectName := b.GetString("project_name")
	if strings.TrimSpace(projectName) == "" {
		b.customAbort(http.StatusBadRequest, "No found project name.")
		return
	}
	isExists, err := service.ProjectExists(projectName)
	if err != nil {
		b.internalError(err)
		return
	}
	if !isExists {
		b.customAbort(http.StatusNotFound, "Project name does not exist.")
		return
	}
	b.repoPath = filepath.Join(baseRepoPath(), b.currentUser.Username, projectName)
	logs.Debug("Set repo path at file upload: %s", b.repoPath)
}

type messageStatus struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func (b *baseController) serveStatus(status int, message string) {
	ms := messageStatus{
		StatusCode: status,
		Message:    message,
	}
	b.Data["json"] = ms
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.ServeJSON()
}

func (b *baseController) internalError(err error) {
	logs.Error("Error occurred: %+v", err)
	b.CustomAbort(http.StatusInternalServerError, "Unexpected error occurred.")
}

func (b *baseController) customAbort(status int, body string) {
	logs.Error("Error occurred: %s", body)
	b.CustomAbort(status, body)
}

func (b *baseController) getCurrentUser() *model.User {
	token := b.Ctx.Request.Header.Get("token")
	if token == "" {
		token = b.GetString("token")
	}
	if isTokenExists := memoryCache.IsExist(token); !isTokenExists {
		logs.Info("Token stored in cache has expired.")
		return nil
	}
	var hasResignedToken bool
	payload, err := verifyToken(token)
	if err != nil {
		if err == errInvalidToken {
			if lastPayload, ok := memoryCache.Get(token).(map[string]interface{}); ok {
				newToken, err := signToken(lastPayload)
				if err != nil {
					logs.Error("failed to sign token: %+v\n", err)
					return nil
				}
				hasResignedToken = true
				token = newToken.TokenString
				payload = lastPayload
				logs.Info("Token has been re-signed due to timeout.")
			}
		} else {
			logs.Error("failed to verify token: %+v\n", err)
		}
	}

	memoryCache.Put(token, payload, time.Second*time.Duration(tokenCacheExpireSeconds))
	b.token = token

	if strID, ok := payload["id"].(string); ok {
		userID, err := strconv.Atoi(strID)
		if err != nil {
			logs.Error("Error occurred on converting userID: %+v\n", err)
			return nil
		}
		user, err := service.GetUserByID(int64(userID))
		if err != nil {
			logs.Error("Error occurred while getting user by ID: %d\n", err)
			return nil
		}
		if currentToken, ok := memoryCache.Get(user.Username).(string); ok {
			if !hasResignedToken && currentToken != "" && currentToken != token {
				logs.Info("Another same name user has signed in other places.")
				return nil
			}
			memoryCache.Put(user.Username, token, time.Second*time.Duration(tokenCacheExpireSeconds))
			b.Ctx.ResponseWriter.Header().Set("token", token)
		}
		return user
	}
	return nil
}

func (b *baseController) signOff() error {
	username := b.GetString("username")
	var err error
	if token, ok := memoryCache.Get(username).(string); ok {
		if payload, ok := memoryCache.Get(token).(map[string]interface{}); ok {
			if userID, ok := payload["id"].(int); ok {
				err = memoryCache.Delete(strconv.Itoa(userID))
				if err != nil {
					logs.Error("Failed to delete by userID from memory cache: %+v", err)
				}
			}
		}
		err = memoryCache.Delete(token)
		if err != nil {
			logs.Error("Failed to delete by token from memory cache: %+v", err)
		}
	}
	err = memoryCache.Delete(username)
	if err != nil {
		logs.Error("Failed to delete by username from memory cache: %+v", err)
	}
	logs.Info("Successful signed off from API server.")
	return nil
}

func signToken(payload map[string]interface{}) (*model.Token, error) {
	var err error
	reqData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(tokenServerURL.String(), "application/json", bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var token model.Token
	err = json.Unmarshal(respData, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func verifyToken(tokenString string) (map[string]interface{}, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, fmt.Errorf("no token was provided")
	}
	resp, err := http.Get(tokenServerURL.String() + "?token=" + tokenString)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		logs.Error("Invalid token due to session timeout.")
		return nil, errInvalidToken
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload map[string]interface{}
	err = json.Unmarshal(respData, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func InitController() {

	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		logs.Error("Failed to load config file: %+v\n", err)
	}
	rawURL := conf.String("tokenServerURL")
	tokenServerURL, err = url.Parse(rawURL)
	if err != nil {
		logs.Error("Failed to parse token server URL: %+v\n", err)
	}
	tokenCacheExpireSeconds, err = conf.Int("tokenCacheExpireSeconds")
	if err != nil {
		logs.Error("Failed to parse token expire seconds: %+v\n", err)
	}

	logs.Info("Set token server URL as %s and will expiration time after %d second(s) in cache", tokenServerURL.String(), tokenCacheExpireSeconds)
	memoryCache, err = cache.NewCache("memory", `{"interval": 3600}`)
	if err != nil {
		logs.Error("Failed to initialize cache: %+v\n", err)
	}

	beego.BConfig.MaxMemory = 1 << 22
	logs.Debug("Current auth mode is: %s", authMode())
}
