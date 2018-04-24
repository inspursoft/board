package service

import (
	"fmt"
	"git/inspursoft/board/src/common/utils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockInitRepoURL  = utils.GetConfig("GOGITS_REPO_URL")
	mockInitRepoPath = utils.GetConfig("BASE_REPO_PATH")
	mockProjectName  = "testgitproject"
	mockUsername     = "admin"
	mockEmail        = "admin@inspur.com"
	mockRepoPath     string
)

func TestGitInitRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	mockRepoPath = filepath.Join(mockInitRepoPath(), mockUsername, mockProjectName)
	user, err := GetUserByName(mockUsername)
	assert.Nilf(err, "Failed to get user: %+v", err)
	ConfigSSHAccess(mockUsername, user.RepoToken)

	repoHandler, err := InitRepo(fmt.Sprintf("%s/%s/%s.git", mockInitRepoURL(), mockUsername, mockProjectName), mockUsername, mockRepoPath)
	assert.Nilf(err, "Failed to initialize repo: %+v", err)
	assert.NotNilf(repoHandler, "Error occurred while creating repo handler: %+v", err)
}

func TestGitOpenRepo(t *testing.T) {
	repoHandler, err := OpenRepo(mockRepoPath, mockUsername)
	assert := assert.New(t)
	assert.Nilf(err, "Failed to open repo: %+v", err)
	assert.NotNilf(repoHandler, "Error occurred while openning repo handler: %+v", err)
}

func TestGitAddFileToRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	CreateFile("target.txt", "Add target.txt file.", mockRepoPath)
	repoHandler, err := OpenRepo(mockRepoPath, mockUsername)
	_, err = repoHandler.Add("target.txt")
	assert.Nilf(err, "Failed to add items to repo: %+v", err)
}

func TestGitCommitToRepo(t *testing.T) {
	var err error
	assert := assert.New(t)
	repoHandler, err := OpenRepo(mockRepoPath, mockUsername)
	_, err = repoHandler.Commit("Initial commit.", mockUsername, mockEmail)
	assert.Nilf(err, "Failed to commit to repo: %+v", err)
}
