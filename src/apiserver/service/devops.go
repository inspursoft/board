package service

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/apiserver/service/devops/gogs"
	"git/inspursoft/board/src/apiserver/service/devops/jenkins"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

var BaseRepoPath = utils.GetConfig("BASE_REPO_PATH")
var GogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var JenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")
var jenkinsNodeIP = utils.GetConfig("JENKINS_NODE_IP")
var jenkinsNodeSSHPort = utils.GetConfig("JENKINS_NODE_SSH_PORT")
var jenkinsNodeUsername = utils.GetConfig("JENKINS_NODE_USERNAME")
var jenkinsNodePassword = utils.GetConfig("JENKINS_NODE_PASSWORD")
var jenkinsNodeVolume = utils.GetConfig("JENKINS_NODE_VOLUME")
var kvmToolsPath = utils.GetConfig("KVM_TOOLS_PATH")
var kvmRegistryPath = utils.GetConfig("KVM_REGISTRY_PATH")
var kvmRegistrySize = utils.GetConfig("KVM_REGISTRY_SIZE")
var kvmRegistryPort = utils.GetConfig("KVM_REGISTRY_PORT")
var kvmToolkitsPath = utils.GetConfig("KVM_TOOLKITS_PATH")
var apiServerURL = utils.GetConfig("BOARD_API_BASE_URL")

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

	repoURL := fmt.Sprintf("%s/%s/%s.git", GogitsSSHURL(), username, repoName)
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
	jenkinsHookURL := fmt.Sprintf("%s/generic-webhook-trigger/invoke", JenkinsBaseURL())
	err = gogsHandler.CreateHook(username, repoName, jenkinsHookURL, "push")
	if err != nil {
		logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
	}

	project, err := GetProjectByName(projectName)
	if err != nil {
		return fmt.Errorf("failed to get project: %+v", err)
	}
	if project == nil {
		logs.Error("Failed to get project by name: %s in DevOps procedure.", projectName)
	}
	pullRequestHookURL := fmt.Sprintf("%s/projects/%d/pull-requests/reset?skip=PRD", apiServerURL(), project.ID)
	err = gogsHandler.CreateHook(username, repoName, pullRequestHookURL, "pull_request")
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
	err = jenkinsHandler.CreateJobWithParameter(repoName)
	if err != nil {
		logs.Error("Failed to create Jenkins' job with repo name: %s, error: %+v", repoName, err)
		return err
	}
	logs.Debug("Waiting for services to be created...")
	time.Sleep(time.Second * 12)
	return nil
}

func ForkRepo(forkedUser *model.User, baseRepoName string) error {
	if forkedUser == nil {
		return errors.New("forked user is nil")
	}
	username := forkedUser.Username
	email := forkedUser.Email
	repoToken := forkedUser.RepoToken

	repoName, err := ResolveRepoName(baseRepoName, username)
	if err != nil {
		logs.Error("Failed to resolve repo name with base name: %s and username: %s.", baseRepoName, username)
		return err
	}

	project, err := GetProjectByName(baseRepoName)
	if err != nil {
		logs.Error("Failed to get project by name: %s, error: %+v", baseRepoName, err)
		return err
	}
	if project == nil {
		return errors.New("project name doesn't exist")
	}

	gogsHandler := gogs.NewGogsHandler(username, repoToken)
	err = gogsHandler.ForkRepo(project.OwnerName, baseRepoName, repoName, "Forked repo.")
	if err != nil {
		return err
	}
	jenkinsHookURL := fmt.Sprintf("%s/generic-webhook-trigger/invoke", JenkinsBaseURL())
	gogsHandler.CreateHook(username, repoName, jenkinsHookURL, "push")
	if err != nil {
		logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
		return err
	}
	pullRequestHookURL := fmt.Sprintf("%s/projects/%d/pull-requests/reset?skip=PRD", apiServerURL(), project.ID)
	err = gogsHandler.CreateHook(username, repoName, pullRequestHookURL, "pull_request")
	if err != nil {
		logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
	}

	repoURL := fmt.Sprintf("%s/%s/%s.git", GogitsSSHURL(), username, repoName)
	repoPath := ResolveRepoPath(repoName, username)
	_, err = InitRepo(repoURL, username, email, repoPath)
	if err != nil {
		logs.Error("Failed to initialize project repo: %+v", err)
		return err
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
	err = jenkinsHandler.CreateJobWithParameter(repoName)
	if err != nil {
		logs.Error("Failed to create Jenkins' job with project name: %s, error: %+v", repoName, err)
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
	members, err := GetProjectMembers(project.ID)
	if err != nil {
		return
	}
	isMember := false
	for _, m := range members {
		if m.Username == username {
			isMember = true
		}
	}
	repoName = project.Name
	if isMember && project.OwnerName != username {
		repoName = username + "_" + project.Name
	}
	logs.Debug("Resolved repo name as: %s.", repoName)
	return
}

func ResolveRepoPath(repoName, username string) (repoPath string) {
	repoPath = filepath.Join(BaseRepoPath(), username, "contents", repoName)
	logs.Debug("Set repo path at file upload: %s", repoPath)
	return
}

func ResolveDockerfileName(imageName, tag string) string {
	imageName = imageName[strings.LastIndex(imageName, "/")+1:]
	return fmt.Sprintf("Dockerfile.%s_%s", imageName, tag)
}

func PrepareKVMHost() error {
	sshPort, _ := strconv.Atoi(jenkinsNodeSSHPort())
	sshHandler, err := NewSecureShell(jenkinsNodeIP(), sshPort, jenkinsNodeUsername(), jenkinsNodePassword())
	if err != nil {
		return err
	}
	kvmToolsNodePath := filepath.Join(kvmToolkitsPath(), "kvm")
	kvmRegistryNodePath := filepath.Join(kvmToolkitsPath(), "kvmregistry")
	err = sshHandler.ExecuteCommand(fmt.Sprintf("mkdir -p %s %s", kvmToolsNodePath, kvmRegistryNodePath))
	if err != nil {
		return err
	}
	err = sshHandler.SecureCopy(kvmToolsPath(), kvmToolsNodePath)
	if err != nil {
		return err
	}
	err = sshHandler.SecureCopy(kvmRegistryPath(), kvmRegistryNodePath)
	if err != nil {
		return err
	}
	return sshHandler.ExecuteCommand(fmt.Sprintf(`
		cd %s && chmod +x kvmregistry && nohup ./kvmregistry -size %s -port %s > kvmregistry.out 2>&1 &`,
		kvmRegistryNodePath, kvmRegistrySize(), kvmRegistryPort()))
}

func ReleaseKVMRegistryByJobName(jobName string) error {
	return utils.SimpleGetRequestHandle(fmt.Sprintf("http://%s:%s/release-node?job_name=%s", jenkinsNodeIP(), kvmRegistryPort(), jobName))
}
