package service_test

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockProjectName = "testgitproject"
	mockUsername    = "boardadmin"
	mockEmail       = "boardadmin@inspur.com"
	mockRepoPath    string
)

func TestGitInitRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	mockRepoPath = filepath.Join(repoPath(), mockUsername, mockProjectName)
	user, err := service.GetUserByName(mockUsername)
	assert.Nilf(err, "Failed to get user: %+v", err)
	service.ConfigSSHAccess(mockUsername, user.RepoToken)

	repoHandler, err := service.InitRepo(fmt.Sprintf("%s/%s/%s.git", gogitsRepoURL(), mockUsername, mockProjectName), mockUsername, mockEmail, mockRepoPath)
	assert.Nilf(err, "Failed to initialize repo: %+v", err)
	assert.NotNilf(repoHandler, "Error occurred while creating repo handler: %+v", err)
}

func TestGitOpenRepo(t *testing.T) {
	repoHandler, err := service.OpenRepo(mockRepoPath, mockUsername, mockEmail)
	assert := assert.New(t)
	assert.Nilf(err, "Failed to open repo: %+v", err)
	assert.NotNilf(repoHandler, "Error occurred while openning repo handler: %+v", err)
}

func TestGitAddFileToRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	service.CreateFile("target.txt", "Add target.txt file.", mockRepoPath)
	repoHandler, err := service.OpenRepo(mockRepoPath, mockUsername, mockEmail)
	_, err = repoHandler.Add("target.txt")
	assert.Nilf(err, "Failed to add items to repo: %+v", err)
}

func TestGitCommitToRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	repoHandler, err := service.OpenRepo(mockRepoPath, mockUsername, mockEmail)
	_, err = repoHandler.Commit("Initial commit.")
	assert.Nilf(err, "Failed to commit to repo: %+v", err)
}
