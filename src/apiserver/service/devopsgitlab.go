package service

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service/devops/gitlab"
	"git/inspursoft/board/src/apiserver/service/devops/jenkins"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"time"

	"github.com/astaxie/beego/logs"
)

var gitlabAdminToken = utils.GetConfig("GITLAB_ADMIN_TOKEN")

type GitlabDevOps struct{}

func (g GitlabDevOps) SignUp(user model.User) error {
	userCreation, err := gitlab.NewGitlabHandler(gitlabAdminToken()).CreateUser(user)
	if err != nil {
		logs.Error("Failed to sign up via Gitlab API, error: %+v", err)
		return err
	}
	logs.Debug("Successful signed up user: %+v", userCreation)
	return nil
}

func (g GitlabDevOps) CreateAccessToken(username string, password string) (string, error) {
	userCreation := gitlab.UserCreation{Name: username, Username: username}
	token, err := gitlab.NewGitlabHandler(gitlabAdminToken()).ImpersonationToken(userCreation)
	if err != nil {
		logs.Error("Failed to create access token via Gitlab API, error %+v", err)
		return "", err
	}
	return token.Token, nil
}

func (g GitlabDevOps) ConfigSSHAccess(username string, token string, publicKey string) error {
	addSSHKeyResponse, err := gitlab.NewGitlabHandler(token).AddSSHKey(fmt.Sprintf("%s's SSH access.", username), publicKey)
	if err != nil {
		logs.Error("Failed to config SSH access via Gitlab API, error: %+v", err)
		return err
	}
	logs.Debug("Successful configured SSH access: %+v", addSSHKeyResponse)
	return nil
}

func (g GitlabDevOps) CreateRepoAndJob(userID int64, projectName string) error {
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

	repoURL := fmt.Sprintf("%s/%s/%s.git", GogitsSSHURL(), username, repoName)
	repoPath := ResolveRepoPath(repoName, username)

	_, err = InitRepo(repoURL, username, email, repoPath)
	if err != nil {
		logs.Error("Failed to initialize default user's repo: %+v", err)
		return err
	}
	gitlabHandler := gitlab.NewGitlabHandler(accessToken)
	if gitlabHandler == nil {
		return fmt.Errorf("failed to create Gitlab handler")
	}
	projectCreation, err := gitlabHandler.CreateRepo(*user, model.Project{Name: projectName})
	if err != nil {
		logs.Error("Failed to create repo via Gitlab API, error %+v", err)
		return err
	}
	logs.Debug("Successful created Gitlab project: %+v", projectCreation)
	// err = gogsHandler.CreateHook(username, repoName)
	// if err != nil {
	// 	logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
	// }

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
	err = jenkinsHandler.CreateJobWithParameter(repoName)
	if err != nil {
		logs.Error("Failed to create Jenkins' job with repo name: %s, error: %+v", repoName, err)
		return err
	}
	logs.Debug("Waiting for services to be created...")
	time.Sleep(time.Second * 12)
	return nil
}

func (g GitlabDevOps) ForkRepo(forkedUser model.User, baseRepoName string) error {
	return nil
}

func (g GitlabDevOps) CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, message string) error {
	return nil
}

func (g GitlabDevOps) DeleteRepo(username string, repoName string) error {
	// user, err := GetUserByName(username)
	// if err != nil {
	// 	logs.Error("Failed to get user by name: %s, with error: %+v", username, err)
	// 	return err
	// }

	// gitlab.NewGitlabHandler(user.RepoToken).DeleteProject()
	return nil
}
