package service

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/apiserver/service/devops/gogs"
	"git/inspursoft/board/src/apiserver/service/devops/jenkins"
	"git/inspursoft/board/src/common/utils"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

var baseRepoPath = utils.GetConfig("BASE_REPO_PATH")
var gogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

func CreateRepoAndJob(userID int64, projectName string) error {

	user, err := GetUserByID(userID)
	if err != nil {
		logs.Error("Failed to get user: %+v", err)
		return err
	}
	if user == nil {
		return fmt.Errorf("user with ID: %d is nil", userID)
	}

	username := user.Username
	email := user.Email
	accessToken := user.RepoToken

	logs.Info("Create repo and job with username: %s, project name: %s.", username, projectName)

	repoName, err := ResolveRepoName(projectName, username)
	if err != nil {
		return err
	}
	logs.Info("Initialize serve repo with name: %s ...", repoName)

	repoURL := fmt.Sprintf("%s/%s/%s.git", gogitsSSHURL(), username, repoName)
	repoPath := ResolveRepoPath(repoName, username)

	_, err = InitRepo(repoURL, username, email, repoPath)
	if err != nil {
		logs.Error("Failed to initialize default user's repo: %+v", err)
		return err
	}
	gogsHandler := gogs.NewGogsHandler(username, accessToken)
	if gogsHandler == nil {
		return fmt.Errorf("failed to create Gogs handler")
	}
	err = gogsHandler.CreateRepo(repoName)
	if err != nil {
		logs.Error("Failed to create repo: %s, error %+v", repoName, err)
		return err
	}
	err = gogsHandler.CreateHook(username, repoName)
	if err != nil {
		logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
	}

	CreateFile("readme.md", "Repo created by Board.", repoPath)

	repoHandler, err := OpenRepo(repoPath, username, email)
	if err != nil {
		logs.Error("Failed to open the repo: %s, error: %+v.", repoPath, err)
		return err
	}

	repoHandler.SimplePush("Add some struts.", "readme.md")
	if err != nil {
		logs.Error("Failed to push readme.md file to the repo: %+v", err)
		return err
	}

	jenkinsHandler := jenkins.NewJenkinsHandler()
	err = jenkinsHandler.CreateJobWithParameter(repoName, username, email)
	if err != nil {
		logs.Error("Failed to create Jenkins' job with repo name: %s, error: %+v", repoName, err)
		return err
	}
	return nil
}

func CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, message string) error {
	gogsHandler := gogs.NewGogsHandler(username, repoToken)
	prInfo, err := gogsHandler.CreatePullRequest(ownerName, repoName, title, message, compareInfo)
	if err != nil {
		logs.Error("Failed to create pull request to the repo: %s with username: %s", repoName, username)
		return err
	}
	if prInfo != nil && prInfo.HasCreated {
		err = gogsHandler.CreateIssueComment(ownerName, repoName, prInfo.Index, message)
		if err != nil {
			logs.Error("Failed to comment issue to the pull request ID: %d, error: %+v", prInfo.IssueID, err)
			return err
		}
	}
	return nil
}

func ResolveRepoName(projectName, username string) (repoName string, err error) {
	project, err := GetProjectByName(projectName)
	if err != nil {
		return
	}
	if project == nil {
		err = errors.New("invalid project name")
		return
	}
	repoName = project.Name
	if project.OwnerName != username {
		repoName = username + "_" + project.Name
	}
	return
}

func ResolveRepoPath(repoName, username string) (repoPath string) {
	repoPath = filepath.Join(baseRepoPath(), username, "contents", repoName)
	logs.Debug("Set repo path at file upload: %s", repoPath)
	return
}

func ResolveDockerfileName(imageName, tag string) string {
	if imageName == "" && tag == "" {
		return "Dockerfile"
	}
	return fmt.Sprintf("Dockerfile.%s_%s", imageName, tag)
}
