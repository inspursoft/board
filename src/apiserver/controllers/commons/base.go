package commons

import (
	"git/inspursoft/board/src/common/model"
	t "git/inspursoft/board/src/common/token"
	"git/inspursoft/board/src/common/utils"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/gitlab"

	"strconv"

	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils/captcha"
)

var MemoryCache cache.Cache
var DefaultCacheDuration = time.Second * time.Duration(TokenCacheExpireSeconds)
var Cpt *captcha.Captcha
var TokenServerURL = utils.GetConfig("TOKEN_SERVER_URL")
var TokenExpireTime = utils.GetConfig("TOKEN_EXPIRE_TIME")
var TokenCacheExpireSeconds int

var APIServerURL = utils.GetConfig("API_SERVER_URL")

var KubeMasterURL = utils.GetConfig("KUBE_MASTER_URL")

var RegistryBaseURI = utils.GetConfig("REGISTRY_BASE_URI")
var AuthMode = utils.GetConfig("AUTH_MODE")

var BaseRepoPath = utils.GetConfig("BASE_REPO_PATH")
var BoardAPIBaseURL = utils.GetConfig("BOARD_API_BASE_URL")
var GogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var JenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

type BaseController struct {
	beego.Controller
	CurrentUser     *model.User
	Token           string
	IsExternalAuth  bool
	IsSysAdmin      bool
	RepoName        string
	RepoPath        string
	RepoServicePath string
	RepoImagePath   string
	Project         *model.Project
	IsRemoved       bool
	OperationID     int64
	AuditDebug      bool
	AuditUser       *model.User
}

func (b *BaseController) Prepare() {
	b.EnableXSRF = false
	b.ResolveSignedInUser()
	b.RecordOperationAudit()
}

func (b *BaseController) Finish() {
	b.UpdateOperationAudit(b.Ctx.ResponseWriter.Status)
}

func (b *BaseController) RecordOperationAudit() {
	b.AuditDebug = utils.GetBoolValue("AUDIT_DEBUG")
	audit := b.Ctx.Request.Header.Get("audit")
	if audit == "" && b.AuditDebug == false {
		return
	}
	//record data about operation
	operation := service.ParseOperationAudit(b.Ctx)
	err := service.CreateOperationAudit(&operation)
	if err != nil {
		logs.Error("Failed to create operation Audit. Error:%+v", err)
		return
	}
	b.OperationID = operation.ID
}

func (b *BaseController) UpdateOperationAudit(statusCode int) {
	if b.OperationID == 0 {
		return
	}
	user := b.CurrentUser
	if b.CurrentUser == nil {
		user = b.AuditUser
	}
	err := service.UpdateOperationAuditStatus(b.OperationID, statusCode, b.Project, user)
	if err != nil {
		logs.Error("Failed to update operation Audit. Error:%+v", err)
		return
	}
}

func (b *BaseController) Render() error {
	return nil
}

func (b *BaseController) ResolveBody(target interface{}) (err error) {
	err = utils.UnmarshalToJSON(b.Ctx.Request.Body, target)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		b.InternalError(err)
		return
	}
	return
}

func (b *BaseController) RenderJSON(data interface{}) {
	b.Data["json"] = data
	b.ServeJSON()
}

func (b *BaseController) ServeStatus(statusCode int, message string) {
	b.ServeJSONOutput(statusCode, struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
	}{
		StatusCode: statusCode,
		Message:    message,
	})
}

func (b *BaseController) ServeJSONOutput(statusCode int, data interface{}) {
	b.Ctx.Output.SetStatus(statusCode)
	b.RenderJSON(data)
}

func (b *BaseController) InternalError(err error) {
	logs.Error("Error occurred: %+v", err)
	b.CustomAbortAudit(http.StatusInternalServerError, "Unexpected error occurred.")
}

func (b *BaseController) CustomAbortAudit(statusCode int, body string) {
	logs.Error("Error of custom aborted: %s", body)
	b.UpdateOperationAudit(statusCode)
	b.CustomAbort(statusCode, template.HTMLEscapeString(body))
}

func ParsePostK8sError(message string) int {
	if strings.Contains(message, "No connection could be made") {
		return http.StatusInternalServerError
	}
	return http.StatusBadRequest
}

func ParseGetK8sError(message string) int {
	if strings.Contains(message, "not found") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (b *BaseController) ParseError(err error, parser func(message string) int) {
	if parser == nil {
		logs.Error("Error in func of parseError,error: parser is nil")
		return
	}
	if err != nil {
		b.CustomAbortAudit(parser(err.Error()), err.Error())
	}
}

func (b *BaseController) GetCurrentUser() *model.User {
	token := b.Ctx.Request.Header.Get("token")
	if token == "" {
		token = b.GetString("token")
	}
	var hasResignedToken bool
	payload, err := t.VerifyToken(TokenServerURL(), token)
	if err != nil {
		if err == t.ErrInvalidToken {
			if lastPayload, ok := MemoryCache.Get(token).(map[string]interface{}); ok {
				newToken, err := t.SignToken(TokenServerURL(), lastPayload)
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
	MemoryCache.Put(token, payload, DefaultCacheDuration)
	b.Token = token

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
		if currentToken, ok := MemoryCache.Get(user.Username).(string); ok {
			if !hasResignedToken && currentToken != "" && currentToken != token {
				logs.Info("Another same name user has signed in other places, removing previous token ...")
				MemoryCache.Delete(token)
				return nil
			}
			MemoryCache.Put(user.Username, token, DefaultCacheDuration)
			b.Ctx.ResponseWriter.Header().Set("token", token)
			if hasResignedToken && MemoryCache.IsExist(currentToken) {
				logs.Info("Deleting stale token stored in cache...")
				MemoryCache.Delete(currentToken)
			}
		}
		user.Password = ""
		return user
	}
	return nil
}

func (b *BaseController) SignOff() error {
	username := b.GetString("username")
	b.AuditUser, _ = service.GetUserByName(username)
	var err error
	if token, ok := MemoryCache.Get(username).(string); ok {
		if payload, ok := MemoryCache.Get(token).(map[string]interface{}); ok {
			if userID, ok := payload["id"].(int); ok {
				err = MemoryCache.Delete(strconv.Itoa(userID))
				if err != nil {
					logs.Error("Failed to delete by userID from memory cache: %+v", err)
				}
			}
		}
		err = MemoryCache.Delete(token)
		if err != nil {
			logs.Error("Failed to delete by token from memory cache: %+v", err)
		}
	}
	err = MemoryCache.Delete(username)
	if err != nil {
		logs.Error("Failed to delete by username from memory cache: %+v", err)
	}
	logs.Info("Successful signed off from API server.")
	return nil
}

func (b *BaseController) ResolveSignedInUser() {
	user := b.GetCurrentUser()
	if user == nil {
		b.CustomAbortAudit(http.StatusUnauthorized, "Need to login first.")
		return
	}
	b.CurrentUser = user
	b.IsSysAdmin = (user.SystemAdmin == 1)
}

func (b *BaseController) ResolveProject(projectName string) (project *model.Project) {
	var err error
	project, err = service.GetProjectByName(projectName)
	if err != nil {
		b.InternalError(err)
		return
	}
	if project == nil {
		b.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("Project: %s does not exist.", projectName))
		return
	}
	b.Project = project
	return
}

func (b *BaseController) ResolveProjectByID(projectID int64) (project *model.Project) {
	var err error
	project, err = service.GetProjectByID(projectID)
	if err != nil {
		b.InternalError(err)
		return
	}
	if project == nil {
		b.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("Project with ID: %d does not exist.", projectID))
		return
	}
	b.Project = project
	return
}

func (b *BaseController) ResolveRepoPath(projectName string) {
	username := b.CurrentUser.Username
	repoName, err := service.ResolveRepoName(projectName, username)
	if err != nil {
		b.CustomAbortAudit(http.StatusPreconditionFailed, fmt.Sprintf("Failed to generate repo path: %+v", err))
		return
	}
	b.RepoPath = service.ResolveRepoPath(repoName, username)
	b.RepoName = repoName
	logs.Debug("Set repo path at file upload: %s and repo name: %s", b.RepoPath, b.RepoName)
}

func (b *BaseController) ResolveRepoServicePath(projectName, serviceName string) {
	b.ResolveRepoPath(projectName)
	b.RepoServicePath = filepath.Join(b.RepoPath, serviceName)
}

func (b *BaseController) ResolveRepoImagePath(projectName string) {
	b.ResolveRepoPath(projectName)
	b.RepoImagePath = filepath.Join(b.RepoPath, "containers")
}

func (b *BaseController) ResolveProjectMember(projectName string) {
	b.ResolveUserPrivilege(projectName)
}

func (b *BaseController) ResolveProjectMemberByID(projectID int64) (project *model.Project) {
	project = b.ResolveProjectByID(projectID)
	b.ResolveProjectMember(project.Name)
	return
}

func (b *BaseController) ResolveProjectOwnerByID(projectID int64) (project *model.Project) {
	project = b.ResolveProjectByID(projectID)
	b.ResolveProjectMemberByID(projectID)
	if !(b.IsSysAdmin || int64(project.OwnerID) == b.CurrentUser.ID) {
		b.CustomAbortAudit(http.StatusForbidden, "User is not the owner of the project.")
		return
	}
	return
}

func (b *BaseController) ResolveUserPrivilege(projectName string) {
	b.ResolveProject(projectName)
	isMember, err := service.IsProjectMemberByName(projectName, b.CurrentUser.ID)
	if err != nil {
		b.InternalError(err)
		return
	}
	if !(b.IsSysAdmin || isMember) {
		b.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to operation.")
		return
	}
	if b.IsSysAdmin && !isMember {
		project := b.ResolveProject(projectName)
		isSuccess, err := service.AddOrUpdateProjectMember(project.ID, b.CurrentUser.ID, 1)
		if err != nil {
			b.InternalError(err)
			return
		}
		if !isSuccess {
			logs.Error("Failed to add project: %s with member %s:", projectName, b.CurrentUser.Username)
			return
		}
		service.CurrentDevOps().ForkRepo(*b.CurrentUser, projectName)
	}
	return
}

func (b *BaseController) ResolveUserPrivilegeByID(projectID int64) (project *model.Project) {
	project = b.ResolveProjectByID(projectID)
	b.ResolveUserPrivilege(project.Name)
	return
}

func (b *BaseController) ManipulateRepo(items ...string) error {
	username := b.CurrentUser.Username
	email := b.CurrentUser.Email
	b.ResolveRepoPath(b.Project.Name)
	commitItems := []service.CommitItem{}
	for _, item := range items {
		commitItems = append(commitItems, service.CommitItem{
			PathWithName: item,
			Content:      utils.GetContentFromFile(filepath.Join(b.RepoPath, item)),
		})
	}
	return service.CurrentDevOps().CommitAndPush(b.RepoName, b.IsRemoved, username, email, commitItems...)
}

func (b *BaseController) PushItemsToRepo(items ...string) {
	err := b.ManipulateRepo(items...)
	if err != nil {
		logs.Error("Failed to push items to repo: %s, error: %+v", b.RepoPath, err)
		b.InternalError(err)
	}
}

func (b *BaseController) CollaborateWithPullRequest(headBranch, baseBranch string, items ...string) {
	if b.RepoPath == "" {
		b.CustomAbortAudit(http.StatusPreconditionFailed, "Repo path cannot be empty.")
		return
	}
	if b.Project == nil {
		b.CustomAbortAudit(http.StatusPreconditionFailed, "Project info cannot be nil.")
		return
	}
	username := b.CurrentUser.Username
	repoName := b.Project.Name
	ownerName := b.Project.OwnerName
	if ownerName == username {
		logs.Info("Skip to create Merge Request as the user %s is the source repo: %s", username, repoName)
		return
	}

	title := fmt.Sprintf("Updates from forked repo: %s/%s", username, repoName)
	content := fmt.Sprintf("Update list: \n\t-\t%s\n", strings.Join(items, "\n\t-\t"))
	compareInfo := fmt.Sprintf("%s...%s:%s", headBranch, username, baseBranch)
	logs.Debug("Pull request info, title: %s, content: %s, compare info: %s", title, content, compareInfo)

	repoToken := b.CurrentUser.RepoToken
	err := service.CurrentDevOps().CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, content)
	if err != nil {
		logs.Error("Failed to create pull request and comment: %+v", err)
		b.InternalError(err)
	}
}

func (b *BaseController) RemoveItemsToRepo(items ...string) {
	b.IsRemoved = true
	err := b.ManipulateRepo(items...)
	if err != nil {
		logs.Error("Failed to remove items to repo: %s, error: %+v", b.RepoPath, err)
		b.InternalError(err)
	}
}

func (b *BaseController) MergeCollaborativePullRequest() {
	if b.CurrentUser.Username != b.Project.OwnerName {
		logs.Info("Skip to merge as the user %s is not the source repo: %s", b.CurrentUser.Username, b.Project.Name)
		return
	}
	projectOwner, err := service.GetUserByName(b.Project.OwnerName)
	if err != nil {
		logs.Error("Failed to get project owner by user ID: %d, error: %+v", b.Project.OwnerID, err)
		b.InternalError(err)
	}
	err = service.CurrentDevOps().MergePullRequest(b.Project.Name, projectOwner.RepoToken)
	if err != nil {
		logs.Warning("Failed to merge request with error: %+v", err)
		if err == gitlab.ErrBranchCannotBeMerged {
			b.CustomAbort(http.StatusNotAcceptable, "Branch has conflicts than cannot be merged.")
		}
	}
}

func (b *BaseController) GeneratePodLogOptions() *model.PodLogOptions {
	var err error
	opt := &model.PodLogOptions{}
	opt.Container = b.GetString("container")
	opt.Follow, err = b.GetBool("follow", false)
	if err != nil {
		logs.Warn("Follow parameter %s is invalid: %+v", b.GetString("follow"), err)
	}
	opt.Previous, err = b.GetBool("previous", false)
	if err != nil {
		logs.Warn("Privious parameter %s is invalid: %+v", b.GetString("privious"), err)
	}
	opt.Timestamps, err = b.GetBool("timestamps", false)
	if err != nil {
		logs.Warn("Timestamps parameter %s is invalid: %+v", b.GetString("timestamps"), err)
	}

	if b.GetString("since_seconds") != "" {
		since, err := b.GetInt64("since_seconds")
		if err != nil {
			logs.Warn("SinceSeconds parameter %s is invalid: %+v", b.GetString("since_seconds"), err)
		} else {
			opt.SinceSeconds = &since
		}
	}

	since := b.GetString("since_time")
	if since != "" {
		sinceTime, err := time.Parse(time.RFC3339, since)
		if err != nil {
			logs.Warn("since_time parameter %s is invalid: %+v", since, err)
		} else {
			opt.SinceTime = &sinceTime
		}
	}

	tail, err := b.GetInt64("tail_lines", -1)
	if err != nil {
		logs.Warn("tail_lines parameter %s is invalid: %+v", b.GetString("tail_lines"), err)
	} else if tail != -1 {
		opt.TailLines = &tail
	}

	limit, err := b.GetInt64("limit_bytes", -1)
	if err != nil {
		logs.Warn("limit_bytes parameter %s is invalid: %+v", b.GetString("limit_bytes"), err)
	} else if limit != -1 {
		opt.LimitBytes = &limit
	}

	return opt
}

func InitController() {
	var err error
	TokenCacheExpireSeconds, err = strconv.Atoi(utils.GetStringValue("TOKEN_CACHE_EXPIRE_SECONDS"))
	if err != nil {
		logs.Error("Failed to get token expire seconds: %+v", err)
	}
	logs.Info("Set token server URL as %s and will expiration time after %d second(s) in cache.", TokenServerURL(), TokenCacheExpireSeconds)
	beego.BConfig.MaxMemory = 1 << 22
	MemoryCache, err = cache.NewCache("memory", `{"interval":3600}`)
	if err != nil {
		logs.Error("Failed to initialize cache: %+v", err)
	}
	Cpt = captcha.NewWithFilter("/captcha/", MemoryCache)
	logs.Debug("Current auth mode is: %s", AuthMode())
}
