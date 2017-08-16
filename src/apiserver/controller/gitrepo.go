package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/logs"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const baseRepoPath = `/repos`

var repoServePath = filepath.Join(baseRepoPath, "board_repo")

type GitRepoController struct {
	baseController
	repoPath string
}

type pushObject struct {
	Items   []string `json:"items"`
	Message string   `json:"message"`
}

func (g *GitRepoController) Prepare() {
	user := g.getCurrentUser()
	if user == nil {
		g.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	g.currentUser = user
	g.repoPath = filepath.Join(baseRepoPath, "board_repo_"+user.Username)
	logs.Debug("Current repo path: %s\n", g.repoPath)
}

func (g *GitRepoController) CreateServeRepo() {
	_, err := service.InitBareRepo(repoServePath)
	if err != nil {
		g.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to initialize serve repo: %+v\n", err))
		return
	}
}

func (g *GitRepoController) InitUserRepo() {
	_, err := service.InitRepo(repoServePath, g.repoPath)
	if err != nil {
		g.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to initialize user's repo: %+v\n", err))
		return
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

	defaultCommitMessage := fmt.Sprintf("Added items: %s to repo: %s", strings.Join(reqPush.Items, ","), g.repoPath)

	if len(reqPush.Message) == 0 {
		reqPush.Message = defaultCommitMessage
	}

	repoHandler, err := service.OpenRepo(g.repoPath)
	if err != nil {
		g.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to open user's repo: %+v\n", err))
		return
	}
	for _, item := range reqPush.Items {
		repoHandler.Add(item)
	}

	username := g.currentUser.Username
	email := g.currentUser.Email

	_, err = repoHandler.Commit(reqPush.Message, &object.Signature{Name: username, Email: email})
	if err != nil {
		g.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to commit changes to user's repo: %+v\n", err))
		return
	}
	err = repoHandler.Push()
	if err != nil {
		g.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to push objects to git repo: %+v\n", err))
	}
}

func (g *GitRepoController) PullObjects() {
	target := g.GetString("target")
	if target == "" {
		g.CustomAbort(http.StatusBadRequest, "No target provided for pulling.")
		return
	}
	targetPath := filepath.Join(baseRepoPath, target)
	repoHandler, err := service.InitRepo(repoServePath, targetPath)
	if err != nil {
		g.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to open user's repo: %+v\n", err))
		return
	}
	err = repoHandler.Pull()
	if err != nil {
		g.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to pull objects from git repo: %+v\n", err))
	}

}
