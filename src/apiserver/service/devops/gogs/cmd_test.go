package gogs_test

import (
	"fmt"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/service/devops/gogs"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var (
	repoName       = "testproject51"
	forkedRepoName = repoName
	user1          = model.User{
		Username: "testuser51",
		Password: "123456a?",
		Email:    "testuser51@inspur.com",
	}
	user2 = model.User{
		Username: "testuser61",
		Password: "123456a?",
		Email:    "testuser61@inspur.com",
	}

	token1 *gogs.AccessToken
	token2 *gogs.AccessToken

	defaultBranch = "master"

	prInfo *gogs.PullRequestInfo

	user1RepoURL string

	gogitsBaseURL = utils.GetConfig("GOGITS_BASE_URL")
	gogitsRepoURL = utils.GetConfig("GOGITS_SSH_URL")
	sshKeyPath    = utils.GetConfig("SSH_KEY_PATH")
	repoPath      = utils.GetConfig("BASE_REPO_PATH")
)

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	utils.SetConfig("BASE_REPO_PATH", "/tmp/test-repos")
	utils.SetConfig("SSH_KEY_PATH", "/tmp/test-keys")
	os.Exit(0)
}

func TestGogitsSignUp(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = gogs.SignUp(user1)
	assert.Nilf(err, "Error occurred while signing up: %+v", err)
	err = gogs.SignUp(user2)
	assert.Nilf(err, "Error occurred while sign up: %+v", err)
}

func TestGogitsCreateAccessToken(t *testing.T) {
	var err error
	assert := assert.New(t)

	token1, err = gogs.CreateAccessToken(user1.Username, user1.Password)
	assert.Nilf(err, "Error occurred while creating token: %+v", err)
	assert.NotNilf(token1, "Failed to initialize Gogs access token to %s", user1.Username)
	logs.Info("Access Token: %+v", token1)

	token2, err = gogs.CreateAccessToken(user2.Username, user2.Password)
	assert.Nilf(err, "Error occurred while creating token: %+v", err)
	assert.NotNilf(token1, "Failed to initialize Gogs access token to %s", user2.Username)
	logs.Info("Access Token: %+v", token2)
}

func TestGogitsCreatePublicKey(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = service.ConfigSSHAccess(user1.Username, token1.Sha1)
	assert.Nilf(err, "Error occurred while config SSH access: %+v to %s", err, user1.Username)

	_, err = service.InitRepo(fmt.Sprintf("%s/%s/%s.git", gogitsRepoURL(), user1.Username, repoName), user1.Username, user1.Email, filepath.Join(repoPath(), user1.Username))
	assert.Nilf(err, "Failed to initialize repo: %+v", err)

	err = service.ConfigSSHAccess(user2.Username, token2.Sha1)
	assert.Nilf(err, "Error occurred while config SSH access: %+v to %s", err, user2.Username)
}

func TestGogitsCreateRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	gogsHandler := gogs.NewGogsHandler(user1.Username, token1.Sha1)
	err = gogsHandler.CreateRepo(repoName)
	assert.Nilf(err, "Error occurred while creating repo: %+v", err)

	service.CreateFile("test.txt", "This is test file.", filepath.Join(repoPath(), user1.Username))
	repoHandler, err := service.OpenRepo(filepath.Join(repoPath(), user1.Username), user1.Username, user1.Email)
	assert.Nilf(err, "Error occurred while openning Git repo handler: %+v", err)
	err = repoHandler.SimplePush("test.txt")
	assert.Nilf(err, "Failed to push files to repo: %+v", err)
}

func TestGogitsCreateHook(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gogs.NewGogsHandler(user1.Username, token1.Sha1).CreateHook(user1.Username, repoName, fmt.Sprintln("http://mock_url"))
	assert.Nilf(err, "Error occurred while creating hook to repo: %+s, error: %+v", repoName, err)
}

func TestGogitsFork(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = gogs.NewGogsHandler(user2.Username, token2.Sha1).ForkRepo(user1.Username, repoName, forkedRepoName, "Forked repo library from admin.")
	assert.Nilf(err, "Error occurred while forking in: %+v", err)
}

func TestGogitsCreatePullRequest(t *testing.T) {
	var err error
	assert := assert.New(t)

	repoHandler, err := service.InitRepo(fmt.Sprintf("%s/%s/%s.git", gogitsRepoURL(), user2.Username, forkedRepoName), user2.Username, user2.Email, filepath.Join(repoPath(), user2.Username))
	assert.Nilf(err, "Failed to initialize repo: %+v", err)

	err = repoHandler.Pull()
	assert.Nilf(err, "Failed to pull updates from repo: %+v", err)

	service.CreateFile("test.txt", "This is another test file.", filepath.Join(repoPath(), user2.Username))
	err = repoHandler.SimplePush("test.txt")
	assert.Nilf(err, "Failed to push files to repo: %+v", err)

	time.Sleep(time.Second * 3)
	prInfo, err = gogs.NewGogsHandler(user2.Username, token2.Sha1).
		CreatePullRequest(user1.Username, repoName, "Update readme.md", "Update readme.md111", fmt.Sprintf("%s...%s:%s", defaultBranch, user2.Username, defaultBranch))
	assert.Nilf(err, "Error occurred while creating pull request: %+v", err)
	assert.NotNilf(prInfo, "Error occurred while getting pull request info: %+v", prInfo)
}

func TestGogitsIssueComment(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gogs.NewGogsHandler(user1.Username, token1.Sha1).CreateIssueComment(user1.Username, repoName, prInfo.Index, "Some comments...")
	assert.Nilf(err, "Error occurred while issuing comments to the PR: %+v", err)
}

func TestDeleteRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gogs.NewGogsHandler(user2.Username, token2.Sha1).DeleteRepo(user2.Username, forkedRepoName)
	assert.Nil(err, "Error occurred while deleting repo: %+v to %s", err, user2.Username)
	err = gogs.NewGogsHandler(user1.Username, token1.Sha1).DeleteRepo(user1.Username, repoName)
	assert.Nilf(err, "Error occurred while deleting repo: %+v to %s", err, user1.Username)

	os.RemoveAll(sshKeyPath())
	os.RemoveAll(repoPath())
}
