package service

import (
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

	logs.Info("Initialize serve repo with name: %s ...", projectName)

	repoURL := fmt.Sprintf("%s/%s/%s.git", gogitsSSHURL(), username, projectName)
	repoPath := fmt.Sprintf("%s/%s/%s", baseRepoPath(), username, projectName)
	_, err = InitRepo(repoURL, username, repoPath)
	if err != nil {
		logs.Error("Failed to initialize default user's repo: %+v", err)
		return err
	}
	gogsHandler := gogs.NewGogsHandler(username, accessToken)
	if gogsHandler == nil {
		return fmt.Errorf("failed to create Gogs handler")
	}
	err = gogsHandler.CreateRepo(projectName)
	if err != nil {
		logs.Error("Failed to create repo: %s, error %+v", projectName, err)
		return err
	}
	err = gogsHandler.CreateHook(username, projectName)
	if err != nil {
		logs.Error("Failed to create hook to repo: %s, error: %+v", projectName, err)
	}
	err = CopyFile("parser.py", filepath.Join(repoPath, "parser.py"))
	if err != nil {
		logs.Error("Failed to copy parser.py file to repo: %+v", err)
		return err
	}
	CreateFile("readme.md", "Repo created by Board.", repoPath)
	err = SimplePush(repoPath, username, email, "Add some struts.", "readme.md", "parser.py")
	if err != nil {
		logs.Error("Failed to push readme.md file to the repo.")
		return err
	}

	jenkinsHandler := jenkins.NewJenkinsHandler()
	err = jenkinsHandler.CreateJobWithParameter(projectName)
	if err != nil {
		logs.Error("Failed to create Jenkins' job with project name: %s, error: %+v", projectName, err)
		return err
	}
	return nil
}
