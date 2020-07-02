package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"git/inspursoft/board/src/apiserver/service/devops/gogs"
	"git/inspursoft/board/src/apiserver/service/devops/jenkins"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/astaxie/beego/logs"
	"golang.org/x/net/context"
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

type gogsJenkinsPushRepositoryPayload struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
}

type gogsJenkinsPushPayload struct {
	Repository   gogsJenkinsPushRepositoryPayload `json:"repository"`
	NodeSelector string                           `json:"node_selector"`
}

type LegacyDevOps struct{}

func (l LegacyDevOps) SignUp(user model.User) error {
	return gogs.SignUp(user)
}

func (l LegacyDevOps) CreateAccessToken(username string, password string) (string, error) {
	accessToken, err := gogs.CreateAccessToken(username, password)
	if err != nil {
		return "", err
	}
	return accessToken.Sha1, nil
}

func (l LegacyDevOps) CommitAndPush(repoName string, isRemoved bool, username string, email string, items ...CommitItem) error {
	repoPath := ResolveRepoPath(repoName, username)
	repoHandler, err := OpenRepo(repoPath, username, email)
	if err != nil {
		logs.Error("Failed to open repo: %+v", err)
		return err
	}
	if isRemoved {
		repoHandler.ToRemove()
	}
	itemNames := []string{}
	for _, commitItem := range items {
		itemNames = append(itemNames, commitItem.PathWithName)
	}
	return repoHandler.SimplePush(itemNames...)
}

func (l LegacyDevOps) ConfigSSHAccess(username string, token string, publicKey string) error {
	return gogs.NewGogsHandler(username, token).CreatePublicKey(fmt.Sprintf("%s's access public key", username), publicKey)
}

func (l LegacyDevOps) CreateRepoAndJob(userID int64, projectName string) error {

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logs.Error("Failed to execute Project and repo creation: %+v", err)
			}
			logs.Info("Successful create repo and Jenkins job.")
		}
		gogsHandler := gogs.NewGogsHandler(username, accessToken)
		if gogsHandler == nil {
			logs.Error("failed to create Gogs handler")
			cancel()
		}
		err = gogsHandler.CreateRepo(repoName)
		if err != nil {
			logs.Error("Failed to create repo: %s, error %+v", repoName, err)
			cancel()
		}
		hookURL := fmt.Sprintf("%s/jenkins-job/invoke", boardAPIBaseURL())
		err = gogsHandler.CreateHook(username, repoName, hookURL)
		if err != nil {
			logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
		}

		CreateFile("readme.md", "Repo created by Board.", repoPath)

		repoHandler, err := OpenRepo(repoPath, username, email)
		if err != nil {
			logs.Error("Failed to open the repo: %s, error: %+v.", repoPath, err)
		}

		repoHandler.SimplePush("Add some struts.", "readme.md")
		if err != nil {
			logs.Error("Failed to push readme.md file to the repo: %+v", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logs.Error("Failed to execute Project and repo creation: %+v", err)
			}
			logs.Info("Successful create repo and Jenkins job.")
		}
		jenkinsHandler := jenkins.NewJenkinsHandler()
		err = jenkinsHandler.CreateJobWithParameter(repoName)
		if err != nil {
			logs.Error("Failed to create Jenkins' job with repo name: %s, error: %+v", repoName, err)
			cancel()
		}
	}()

	wg.Wait()
	return nil
}

func (l LegacyDevOps) ForkRepo(forkedUser model.User, baseRepoName string) error {
	if forkedUser.ID == 0 {
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

	hookURL := fmt.Sprintf("%s/jenkins-job/invoke", boardAPIBaseURL())
	gogsHandler.CreateHook(username, repoName, hookURL)
	if err != nil {
		logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
		return err
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

func (l LegacyDevOps) CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, message string) error {
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

func (l LegacyDevOps) MergePullRequest(repoName, repoToken string) error {
	return fmt.Errorf("unsupport merge pull request feature with the Gogits repo service")
}

func (l LegacyDevOps) DeleteRepo(username string, repoName string) error {
	user, err := GetUserByName(username)
	if err != nil {
		logs.Error("Failed to get user by name: %s, error: %+v", username, err)
		return err
	}
	return gogs.NewGogsHandler(user.Username, user.RepoToken).DeleteRepo(user.Username, repoName)
}

func (l LegacyDevOps) CustomHookPushPayload(rawPayload []byte, nodeSelection string) error {
	var cp gogsJenkinsPushPayload
	err := json.Unmarshal(rawPayload, &cp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON custom push payload: %+v", err)
	}
	cp.NodeSelector = nodeSelection
	logs.Debug("Resolve for push event hook payload: %+v", cp)
	header := http.Header{
		"content-type":   []string{"json"},
		"X-Gitlab-Event": []string{"Push Hook"},
	}
	return utils.SimplePostRequestHandle(fmt.Sprintf("%s/generic-webhook-trigger/invoke", JenkinsBaseURL()), header, cp)
}

func (l LegacyDevOps) GetRepoFile(username string, repoName string, branch string, filePath string) ([]byte, error) {
	return nil, fmt.Errorf("unimplement get repo files feature with the Gogits repo service")
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
