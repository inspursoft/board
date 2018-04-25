package service

import (
	"fmt"
	"os"
	"path/filepath"

	"git/inspursoft/board/src/apiserver/service/devops/gogs"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var (
	repoName       = "testproject"
	forkedRepoName = repoName
	repoPath       = utils.GetConfig("BASE_REPO_PATH")

	user1 = model.User{
		Username: "testuser1",
		Password: "123456a?",
		Email:    "testuser1@inspur.com",
	}
	user2 = model.User{
		Username: "testuser2",
		Password: "123456a?",
		Email:    "testuser2@inspur.com",
	}

	token1 *gogs.AccessToken
	token2 *gogs.AccessToken

	defaultBranch = "master"

	gogitsBaseURL = utils.GetConfig("GOGITS_BASE_URL")
	gogitsRepoURL = utils.GetConfig("GOGITS_REPO_URL")

	prInfo *gogs.PullRequestInfo
)

func prepareGogitsTest() {
	os.RemoveAll(sshKeyPath())
	os.RemoveAll(repoPath())
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
	err = ConfigSSHAccess(user1.Username, token1.Sha1)
	assert.Nilf(err, "Error occurred while config SSH access: %+v to %s", err, user1.Username)
	_, err = InitRepo(fmt.Sprintf("%s/%s/%s.git", gogitsRepoURL(), user1.Username, repoName), user1.Username, filepath.Join(repoPath(), user1.Username))
	assert.Nilf(err, "Failed to initialize repo: %+v", err)

	err = ConfigSSHAccess(user2.Username, token2.Sha1)
	assert.Nilf(err, "Error occurred while config SSH access: %+v to %s", err, user2.Username)
}

func TestGogitsCreateRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gogs.NewGogsHandler(user1.Username, token1.Sha1).CreateRepo(repoName)
	assert.Nilf(err, "Error occurred while creating repo: %+v", err)

	CreateFile("test.txt", "This is test file.", filepath.Join(repoPath(), user1.Username))
	SimplePush(filepath.Join(repoPath(), user1.Username), user1.Username, user1.Email, "Initial Commit", "test.txt")
}

func TestGogitsCreateHook(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gogs.NewGogsHandler(user1.Username, token1.Sha1).CreateHook(user1.Username, repoName)
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

	repoHandler, err := InitRepo(fmt.Sprintf("%s/%s/%s.git", gogitsRepoURL(), user2.Username, forkedRepoName), user2.Username, filepath.Join(repoPath(), user2.Username))
	assert.Nilf(err, "Failed to initialize repo: %+v", err)

	err = repoHandler.Pull()
	assert.Nilf(err, "Failed to pull updates from repo: %+v", err)

	CreateFile("test.txt", "This is another test file.", filepath.Join(repoPath(), user2.Username))
	SimplePush(filepath.Join(repoPath(), user2.Username), user2.Username, user2.Email, "Update test.txt file.", "test.txt")

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
	err = gogs.NewGogsHandler(user1.Username, token1.Sha1).DeleteRepo(repoName)
	assert.Nilf(err, "Error occurred while deleting repo: %+v to %s", err, user1.Username)

	err = gogs.NewGogsHandler(user2.Username, token2.Sha1).DeleteRepo(forkedRepoName)
	assert.Nil(err, "Error occurred while deleting repo: %+v to %s", err, user2.Username)
}
