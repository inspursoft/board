package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/gogs"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

const (
	toBeRemoved   = true
	defaultBranch = "master"
)

type GitRepoController struct {
	baseController
	repoServerURL string
}

type pushObject struct {
	ProjectName string
	UserID      int64
	Items       []string `json:"items"`
	Message     string   `json:"message"`
	JobName     string   `json:"job_name"`
	Value       string   `json:"value"`
	Extras      string   `json:"extras"`
	FileName    string   `json:"file_name"`
}

func (g *GitRepoController) Prepare() {
	user := g.getCurrentUser()
	if user == nil {
		g.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	g.currentUser = user
	g.resolveRepoPath()
}

func (g *GitRepoController) resolveRepoServerURL() {
	projectName := g.GetString("project_name")
	if strings.TrimSpace(projectName) == "" {
		g.CustomAbort(http.StatusBadRequest, "No found for project name.")
		return
	}
	isExists, err := service.ProjectExists(projectName)
	if err != nil {
		g.internalError(err)
		return
	}
	if !isExists {
		g.CustomAbort(http.StatusNotFound, fmt.Sprintf("Project %s does not exist.", projectName))
		return
	}
	username := g.currentUser.Username
	g.repoServerURL = fmt.Sprintf("%s/%s/%s.git", gogitsSSHURL(), username, projectName)
}

func (g *GitRepoController) CreateServeRepo() {
	g.resolveRepoServerURL()
	_, err := service.InitBareRepo(g.repoServerURL)
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to initialize serve repo: %+v\n", err))
		return
	}
}

func (g *GitRepoController) InitUserRepo() {
	g.resolveRepoServerURL()
	_, err := service.InitRepo(g.repoServerURL, g.currentUser.Username, g.repoPath)
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to initialize user's repo: %+v\n", err))
		return
	}
	err = service.CreateMetaConfiguration(make(map[string]string), g.repoPath)
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to initialize user's repo default directories: %+v\n", err))
	}
}

func (g *GitRepoController) PushObjects() {
	var reqPush pushObject
	reqData, err := g.resolveBody()
	if err != nil {
		g.internalError(err)
		return
	}
	err = json.Unmarshal(reqData, &reqPush)
	if err != nil {
		g.internalError(err)
		return
	}

	username := g.currentUser.Username
	email := g.currentUser.Email

	repoPath := filepath.Join(baseRepoPath(), username, reqPush.ProjectName)
	defaultCommitMessage := fmt.Sprintf("added items: %s to repo: %s", strings.Join(reqPush.Items, ","), repoPath)

	if len(reqPush.Message) == 0 {
		reqPush.Message = defaultCommitMessage
	}
	repoHandler, err := service.OpenRepo(repoPath, username)
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to open user's repo: %+v\n", err))
		return
	}
	for _, item := range reqPush.Items {
		repoHandler.Add(item)
	}

	_, err = repoHandler.Commit(reqPush.Message, username, email)
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to commit changes to user's repo: %+v\n", err))
		return
	}
	err = repoHandler.Push()
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to push objects to git repo: %+v\n", err))
	}
}

func (g *GitRepoController) PullObjects() {
	repoHandler, err := service.OpenRepo(g.repoPath, g.currentUser.Username)
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to open user's repo: %+v\n", err))
		return
	}
	err = repoHandler.Pull()
	if err != nil {
		g.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to pull objects from git repo: %+v\n", err))
	}
}

func generateMetaConfiguration(p *pushObject, repoPath string) error {
	conf := make(map[string]string)
	conf["extras"] = p.Extras
	conf["file_name"] = p.FileName
	conf["flag"] = p.JobName
	conf["value"] = p.Value
	conf["apiserver"] = apiServerURL()
	conf["user_id"] = strconv.Itoa(int(p.UserID))
	if p.JobName == imageProcess {
		conf["docker_registry"] = registryBaseURI()
	}
	return service.CreateMetaConfiguration(conf, repoPath)
}

func InternalPushObjects(p *pushObject, g *baseController, actionType ...bool) (int, string, error) {
	username := g.currentUser.Username
	repoPath := filepath.Join(baseRepoPath(), username, p.ProjectName)
	logs.Debug("Repo path for pushing objects: %s", repoPath)

	isRemove := false
	if len(actionType) > 0 && actionType[0] {
		isRemove = true
	}

	actionName := "Added"
	if isRemove {
		actionName = "Removed"
	}
	defaultCommitMessage := fmt.Sprintf("%s items: %s to repo: %s", actionName, strings.Join(p.Items, ","), repoPath)

	if len(p.Message) == 0 {
		p.Message = defaultCommitMessage
	}

	if isRemove {
		p.Message = "[DELETED]" + p.Message
	}

	email := g.currentUser.Email
	repoHandler, err := service.OpenRepo(repoPath, username)
	if err != nil {
		return http.StatusInternalServerError, "Failed to open user's repo", err
	}

	for _, item := range p.Items {
		if isRemove {
			repoHandler.Remove(item)
		} else {
			repoHandler.Add(item)
		}
		logs.Debug(">>>>> pushed item: %s", item)
	}

	_, err = repoHandler.Commit(p.Message, username, email)
	if err != nil {
		return http.StatusInternalServerError, "Failed to commit changes to user's repo", err
	}
	err = repoHandler.Push()
	if err != nil {
		return http.StatusInternalServerError, "Failed to push objects to git repo", err
	}

	project, err := service.GetProject(model.Project{Name: p.ProjectName}, "name")
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Failed to check username: %s to the project: %s", username, p.ProjectName), err
	}

	if project != nil && project.OwnerName != username {
		pullRequestTitle := fmt.Sprintf("Updates from forked repo: %s/%s", username, p.ProjectName)
		pullRequestContent := fmt.Sprintf("Update list: \n\t-\t%s\n", strings.Join(p.Items, "\n\t-\t"))
		pullRequestCompare := fmt.Sprintf("%s...%s:%s", defaultBranch, username, defaultBranch)
		logs.Debug("Pull request info, title: %s, content: %s, compare info: %s", pullRequestTitle, pullRequestContent, pullRequestCompare)
		gogsHandler := gogs.NewGogsHandler(username, g.currentUser.RepoToken)
		prInfo, err := gogsHandler.CreatePullRequest(project.OwnerName, project.Name, pullRequestTitle, pullRequestContent, pullRequestCompare)
		if err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("Failed to create pull request to the repo: %s with username: %s", p.ProjectName, username), err
		}
		if prInfo != nil && prInfo.HasCreated {
			err = gogsHandler.CreateIssueComment(project.OwnerName, project.Name, prInfo.Index, pullRequestContent)
			if err != nil {
				return http.StatusInternalServerError, fmt.Sprintf("Failed to comment issue to the pull request ID: %d, error: %+v", prInfo.IssueID, err), err
			}
		}
	}
	return http.StatusOK, "Internal Push Object successfully", err
}
