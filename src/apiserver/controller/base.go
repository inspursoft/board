package controller

import (
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"git/inspursoft/board/src/apiserver/service"

	"strconv"

	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
)

var tokenServerURL = utils.GetConfig("TOKEN_SERVER_URL")
var tokenExpirtTime = utils.GetConfig("TOKEN_EXPIRE_TIME")
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

type BaseController struct {
	beego.Controller
	currentUser     *model.User
	token           string
	isExternalAuth  bool
	isSysAdmin      bool
	repoName        string
	repoPath        string
	repoServicePath string
	repoImagePath   string
	project         *model.Project
	isRemoved       bool
	operationID     int64
	auditDebug      bool
	auditUser       *model.User
}

func (b *BaseController) Prepare() {
	b.resolveSignedInUser()
	b.recordOperationAudit()
}

func (b *BaseController) Finish() {
	b.updateOperationAudit()
}

func (b *BaseController) recordOperationAudit() {
	b.auditDebug = utils.GetBoolValue("AUDIT_DEBUG")
	audit := b.Ctx.Request.Header.Get("audit")
	if audit == "" && b.auditDebug == false {
		return
	}
	//record data about operation
	operation := service.ParseOperationAudit(b.Ctx)
	err := service.CreateOperationAudit(&operation)
	if err != nil {
		logs.Error("Failed to create operation Audit. Error:%+v", err)
		return
	}
	b.operationID = operation.ID
}

func (b *BaseController) updateOperationAudit() {
	if b.operationID == 0 {
		return
	}

	var err error
	if b.currentUser != nil {
		err = service.UpdateOperationAuditStatus(b.operationID, b.Ctx.ResponseWriter.Status, b.project, b.currentUser)
	} else {
		err = service.UpdateOperationAuditStatus(b.operationID, b.Ctx.ResponseWriter.Status, b.project, b.auditUser)
	}
	if err != nil {
		logs.Error("Failed to update operation Audit. Error:%+v", err)
		return
	}
}

func (b *BaseController) Render() error {
	return nil
}

func (b *BaseController) resolveBody(target interface{}) (err error) {
	err = utils.UnmarshalToJSON(b.Ctx.Request.Body, target)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		b.internalError(err)
		return
	}
	return
}

func (b *BaseController) renderJSON(data interface{}) {
	b.Data["json"] = data
	b.ServeJSON()
}

func (b *BaseController) serveStatus(status int, message string) {
	b.serveJSON(status, struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
	}{
		StatusCode: status,
		Message:    message,
	})
}

func (b *BaseController) serveJSON(status int, data interface{}) {
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.renderJSON(data)
}

func (b *BaseController) internalError(err error) {
	logs.Error("Error occurred: %+v", err)
	b.serveStatus(http.StatusInternalServerError, "Unexpected error occurred.")
}

func (b *BaseController) customAbort(status int, body string) {
	logs.Error("Error of custom aborted: %s", body)
	b.serveStatus(status, body)
}

func parsePostK8sError(message string) int {
	if strings.Contains(message, "No connection could be made") {
		return http.StatusInternalServerError
	}
	return http.StatusBadRequest
}

func parseGetK8sError(message string) int {
	if strings.Contains(message, "not found") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (b *BaseController) parseError(err error, parser func(message string) int) {
	if parser == nil {
		logs.Error("Error in func of parseError,error: parser is nil")
		return
	}
	if err != nil {
		b.customAbort(parser(err.Error()), err.Error())
	}
}

func (b *BaseController) getCurrentUser() *model.User {
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

func (b *BaseController) signOff() error {
	username := b.GetString("username")
	b.auditUser, _ = service.GetUserByName(username)
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

func (b *BaseController) resolveSignedInUser() {
	user := b.getCurrentUser()
	if user == nil {
		b.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	b.currentUser = user
	b.isSysAdmin = (user.SystemAdmin == 1)
}

func (b *BaseController) resolveProject(projectName string) (project *model.Project) {
	var err error
	project, err = service.GetProjectByName(projectName)
	if err != nil {
		b.internalError(err)
		return
	}
	if project == nil {
		b.customAbort(http.StatusNotFound, fmt.Sprintf("Project: %s does not exist.", projectName))
		return
	}
	b.project = project
	return
}

func (b *BaseController) resolveProjectByID(projectID int64) (project *model.Project) {
	var err error
	project, err = service.GetProjectByID(projectID)
	if err != nil {
		b.internalError(err)
		return
	}
	if project == nil {
		b.customAbort(http.StatusNotFound, fmt.Sprintf("Project with ID: %d does not exist.", projectID))
		return
	}
	b.project = project
	return
}

func (b *BaseController) resolveRepoPath(projectName string) {
	username := b.currentUser.Username
	repoName, err := service.ResolveRepoName(projectName, username)
	if err != nil {
		b.customAbort(http.StatusPreconditionFailed, fmt.Sprintf("Failed to generate repo path: %+v", err))
		return
	}
	b.repoPath = service.ResolveRepoPath(repoName, username)
	logs.Debug("Set repo path at file upload: %s", b.repoPath)
}

func (b *BaseController) resolveRepoServicePath(projectName, serviceName string) {
	b.resolveRepoPath(projectName)
	b.repoServicePath = filepath.Join(b.repoPath, serviceName)
}

func (b *BaseController) resolveRepoImagePath(projectName string) {
	b.resolveRepoPath(projectName)
	b.repoImagePath = filepath.Join(b.repoPath, "containers")
}

func (b *BaseController) resolveProjectMember(projectName string) {
	b.resolveUserPrivilege(projectName)
}

func (b *BaseController) resolveProjectMemberByID(projectID int64) (project *model.Project) {
	project = b.resolveProjectByID(projectID)
	b.resolveProjectMember(project.Name)
	return
}

func (b *BaseController) resolveProjectOwnerByID(projectID int64) (project *model.Project) {
	project = b.resolveProjectByID(projectID)
	b.resolveProjectMemberByID(projectID)
	if !(b.isSysAdmin || int64(project.OwnerID) == b.currentUser.ID) {
		b.customAbort(http.StatusForbidden, "User is not the owner of the project.")
		return
	}
	return
}

func (b *BaseController) resolveUserPrivilege(projectName string) {
	b.resolveProject(projectName)
	isMember, err := service.IsProjectMemberByName(projectName, b.currentUser.ID)
	if err != nil {
		b.internalError(err)
		return
	}
	if !(b.isSysAdmin || isMember) {
		b.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
	}
	if b.isSysAdmin && !isMember {
		project := b.resolveProject(projectName)
		isSuccess, err := service.AddOrUpdateProjectMember(project.ID, b.currentUser.ID, 1)
		if err != nil {
			b.internalError(err)
			return
		}
		if !isSuccess {
			logs.Error("Failed to add project: %s with member %s:", projectName, b.currentUser.Username)
			return
		}
		service.ForkRepo(b.currentUser, projectName)
	}
	return
}

func (b *BaseController) resolveUserPrivilegeByID(projectID int64) (project *model.Project) {
	project = b.resolveProjectByID(projectID)
	b.resolveUserPrivilege(project.Name)
	return
}

func (b *BaseController) manipulateRepo(items ...string) error {
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

func (b *BaseController) pushItemsToRepo(items ...string) {
	err := b.manipulateRepo(items...)
	if err != nil {
		logs.Error("Failed to push items to repo: %s, error: %+v", b.repoPath, err)
		b.internalError(err)
	}
}

func (b *BaseController) collaborateWithPullRequest(headBranch, baseBranch string, items ...string) {
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

func (b *BaseController) removeItemsToRepo(items ...string) {
	b.isRemoved = true
	err := b.manipulateRepo(items...)
	if err != nil {
		logs.Error("Failed to remove items to repo: %s, error: %+v", b.repoPath, err)
		b.internalError(err)
	}
}

func signToken(payload map[string]interface{}) (*model.Token, error) {
	var token model.Token
	err := utils.RequestHandle(http.MethodPost, tokenServerURL(), func(req *http.Request) error {
		req.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}
		return nil
	}, payload, func(req *http.Request, resp *http.Response) error {
		return utils.UnmarshalToJSON(resp.Body, &token)
	})
	return &token, err
}

func verifyToken(tokenString string) (map[string]interface{}, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, fmt.Errorf("no token provided")
	}
	var payload map[string]interface{}
	err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s?token=%s", tokenServerURL(), tokenString), nil, nil, func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode == http.StatusUnauthorized {
			logs.Error("Invalid token due to session timeout.")
			return errInvalidToken
		}
		return utils.UnmarshalToJSON(resp.Body, &payload)
	})
	return payload, err
}

func InitController() {
	var err error
	tokenCacheExpireSeconds, err = strconv.Atoi(utils.GetStringValue("TOKEN_CACHE_EXPIRE_SECONDS"))
	if err != nil {
		logs.Error("Failed to get token expire seconds: %+v", err)
	}
	logs.Info("Set token server URL as %s and will expiration time after %d second(s) in cache.", tokenServerURL(), tokenCacheExpireSeconds)

	memoryCache, err = cache.NewCache("memory", `{"interval": 3600}`)
	if err != nil {
		logs.Error("Failed to initialize cache: %+v", err)
	}
	beego.BConfig.MaxMemory = 1 << 22
	logs.Debug("Current auth mode is: %s", authMode())
}
