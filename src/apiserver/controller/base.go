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

var registryBaseURI = utils.GetConfig("REGISTRY_BASE_URI")
var authMode = utils.GetConfig("AUTH_MODE")

var baseRepoPath = utils.GetConfig("BASE_REPO_PATH")
var boardAPIBaseURL = utils.GetConfig("BOARD_API_BASE_URL")
var gogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

type commonController struct {
	beego.Controller
	token          string
	isExternalAuth bool
}

func (c *commonController) Render() error {
	return nil
}

func (c *commonController) resolveBody(target interface{}) {
	err := utils.UnmarshalToJSON(c.Ctx.Request.Body, target)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		c.internalError(err)
	}
}

func (c *commonController) serveStatus(status int, message string) {
	c.Ctx.ResponseWriter.WriteHeader(status)
	c.Data["json"] = struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
	}{
		StatusCode: status,
		Message:    message,
	}
	c.ServeJSON()
}

func (c *commonController) internalError(err error) {
	logs.Error("Error occurred: %+v", err)
	c.CustomAbort(http.StatusInternalServerError, "Unexpected error occurred.")
}

func (c *commonController) customAbort(status int, body string) {
	logs.Error("Error occurred: %s", body)
	c.CustomAbort(status, body)
}

func (c *commonController) getCurrentUser() *model.User {
	token := c.Ctx.Request.Header.Get("token")
	if token == "" {
		token = c.GetString("token")
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
	c.token = token

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
			c.Ctx.ResponseWriter.Header().Set("token", token)
		}
		return user
	}
	return nil
}

func (c *commonController) signOff() error {
	username := c.GetString("username")
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

type baseController struct {
	commonController
	currentUser *model.User

	isSysAdmin bool
	repoPath   string
	project    *model.Project
	isRemoved  bool
}

func (b *baseController) Prepare() {
	b.resolveSignedInUser()
}

func (b *baseController) resolveSignedInUser() {
	user := b.getCurrentUser()
	if user == nil {
		b.customAbort(http.StatusUnauthorized, "Need to login first.")
	}
	b.currentUser = user
	b.isSysAdmin = (user.SystemAdmin == 1)
}

func (b *baseController) resolveRepoPath(projectName string) {
	if projectName == "" {
		b.customAbort(http.StatusBadRequest, "No found project name.")
	}
	var err error
	b.project, err = service.GetProject(model.Project{Name: projectName}, "name")
	if err != nil {
		b.internalError(err)
		return
	}
	if b.project == nil {
		b.customAbort(http.StatusNotFound, fmt.Sprintf("Project: %s does not exist.", projectName))
		return
	}
	b.repoPath = filepath.Join(baseRepoPath(), b.currentUser.Username, projectName)
	logs.Debug("Set repo path at file upload: %s", b.repoPath)
}

func (b *baseController) resolveProjectMember(projectName string) {
	isMember, err := service.IsProjectMemberByName(projectName, b.currentUser.ID)
	if err != nil {
		b.internalError(err)
		return
	}
	if !isMember {
		b.customAbort(http.StatusForbidden, fmt.Sprintf("Project %s is not the member to the current user.", b.currentUser.Username))
		return
	}
}

func (b *baseController) resolveUserPrivilege(projectName string) (isMember bool) {
	var err error
	isMember, err = service.IsProjectMemberByName(projectName, b.currentUser.ID)
	if err != nil {
		b.internalError(err)
	}
	if !(b.isSysAdmin || isMember) {
		b.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
	}
	return
}

func (b *baseController) resolveUserPrivilegeByID(projectID int64) (project *model.Project) {
	var err error
	project, err = service.GetProjectByID(projectID)
	if err != nil {
		b.internalError(err)
	}
	if project == nil {
		b.customAbort(http.StatusNotFound, fmt.Sprintf("Project ID: %+d does not exist.", projectID))
	}
	b.resolveUserPrivilege(project.Name)
	return
}

func (b *baseController) manipulateRepo(items ...string) error {
	if b.repoPath == "" {
		return fmt.Errorf("repo path cannot be empty")
	}
	username := b.currentUser.Username
	email := b.currentUser.Email
	repoHandler, err := service.OpenRepo(b.repoPath, username, email)
	if err != nil {
		logs.Error("Failed to open repo: %+v", err)
		return err
	}
	if b.isRemoved {
		repoHandler.ToRemove()
	}
	return repoHandler.SimplePush(items...)
}

func (b *baseController) pushItemsToRepo(items ...string) {
	err := b.manipulateRepo(items...)
	if err != nil {
		logs.Error("Failed to push items to repo: %s, error: %+v", b.repoPath, err)
		b.internalError(err)
	}
}

func (b *baseController) collaborateWithPullRequest(headBranch, baseBranch string, items ...string) {
	if b.repoPath == "" {
		b.customAbort(http.StatusPreconditionFailed, "Repo path cannot be empty.")
		return
	}
	if b.project == nil {
		b.customAbort(http.StatusPreconditionFailed, "Project info cannot be nil.")
		return
	}
	username := b.currentUser.Username
	repoName := b.project.Name
	ownerName := b.project.OwnerName
	if ownerName == username {
		logs.Info("User %s is the owner to the current repo: %s", username, repoName)
		return
	}

	title := fmt.Sprintf("Updates from forked repo: %s/%s", username, repoName)
	content := fmt.Sprintf("Update list: \n\t-\t%s\n", strings.Join(items, "\n\t-\t"))
	compareInfo := fmt.Sprintf("%s...%s:%s", headBranch, username, baseBranch)
	logs.Debug("Pull request info, title: %s, content: %s, compare info: %s", title, content, compareInfo)

	repoToken := b.currentUser.RepoToken
	err := service.CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, content)
	if err != nil {
		logs.Error("Failed to create pull request and comment: %+v", err)
		b.internalError(err)
	}
}

func (b *baseController) removeItemsToRepo(items ...string) {
	b.isRemoved = true
	err := b.manipulateRepo(items...)
	if err != nil {
		logs.Error("Failed to remove items to repo: %s, error: %+v", b.repoPath, err)
		b.internalError(err)
	}
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
