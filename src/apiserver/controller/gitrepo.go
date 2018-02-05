package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/logs"
)

type GitRepoController struct {
	baseController
	repoServerURL string
}

type pushObject struct {
	ProjectName string
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
	err = service.CreateBaseDirectory(make(map[string]string), g.repoPath)
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
	conf["docker_registry"] = registryBaseURI()
	return service.CreateBaseDirectory(conf, repoPath)
}

func InternalPushObjects(p *pushObject, g *baseController) (int, string, error) {
	username := g.currentUser.Username
	repoPath := filepath.Join(baseRepoPath(), username, p.ProjectName)
	logs.Debug("Repo path for pushing objects: %s", repoPath)

	defaultCommitMessage := fmt.Sprintf("Added items: %s to repo: %s", strings.Join(p.Items, ","), repoPath)

	if len(p.Message) == 0 {
		p.Message = defaultCommitMessage
	}

	email := g.currentUser.Email
	repoHandler, err := service.OpenRepo(repoPath, username)
	if err != nil {
		return http.StatusInternalServerError, "Failed to open user's repo", err
	}

	generateMetaConfiguration(p, repoPath)
	repoHandler.Add("META.cfg")

	for _, item := range p.Items {
		logs.Debug(">>>>> pushed item: %s", item)
		repoHandler.Add(item)
	}

	_, err = repoHandler.Commit(p.Message, username, email)
	if err != nil {
		return http.StatusInternalServerError, "Failed to commit changes to user's repo", err
	}
	err = repoHandler.Push()
	if err != nil {
		return http.StatusInternalServerError, "Failed to push objects to git repo", err
	}
	return 0, "Internal Push Object successfully", err
}

// Clean git repo after remove config files
func InternalCleanObjects(p *pushObject, g *baseController) (int, string, error) {

	repoPath := filepath.Join(baseRepoPath(), p.ProjectName, p.Value)
	logs.Debug("Repo path for pushing objects: %s", repoPath)

	defaultCommitMessage := fmt.Sprintf("Removed items: %s from repo: %s", strings.Join(p.Items, ","), repoPath)
	if len(p.Message) == 0 {
		p.Message = defaultCommitMessage
	}

	username := g.currentUser.Username
	email := g.currentUser.Email
	repoHandler, err := service.OpenRepo(repoPath, username)
	if err != nil {
		return http.StatusInternalServerError, "Failed to open user's repo", err
	}
	for _, item := range p.Items {
		repoHandler.Remove(item)
	}

	_, err = repoHandler.Commit(p.Message, username, email)
	if err != nil {
		return http.StatusInternalServerError, "Failed to commit changes to user's repo", err
	}
	err = repoHandler.Push()
	if err != nil {
		return http.StatusInternalServerError, "Failed to push objects to git repo", err
	}
	return 0, "Internal Push Object for cleaning successfully", nil
}
