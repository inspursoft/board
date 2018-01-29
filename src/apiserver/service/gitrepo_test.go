package service

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var mockInitRepoURL = "ssh://git@10.165.14.97:10022"
var mockInitRepoPath = "/Users/wangkun/repos"
var mockProjectName = "library"
var mockUsername = "admin"
var mockEmail = "admin@inspur.com"
var mockRepoPath = filepath.Join(mockInitRepoPath, mockUsername, mockProjectName)
var mockFileName = "readme.md"
var handler *repoHandler
var err error

func TestInitGitRepo(t *testing.T) {
	_, err := InitRepo(fmt.Sprintf("%s/%s/%s.git", mockInitRepoURL, mockUsername, mockProjectName), mockUsername, mockRepoPath)
	assert := assert.New(t)
	assert.Nilf(err, "Failed to initialize repo: %+v", err)
}

func TestOpenRepo(t *testing.T) {
	handler, err = OpenRepo(mockRepoPath, mockUsername)
	assert := assert.New(t)
	assert.Nilf(err, "Failed to open repo: %+v", err)
}

func TestAddFileToRepo(t *testing.T) {
	tempFilePath := filepath.Join(mockRepoPath, mockFileName)
	logs.Debug("temp file path: %s", tempFilePath)
	_, err := os.OpenFile(tempFilePath, os.O_CREATE, 0740)
	assert := assert.New(t)
	assert.Nilf(err, "Failed to create file: %+v", err)
	_, err = handler.Add(mockFileName)
	assert.Nilf(err, "Failed to add files to repo: %+v", err)
}

func TestCommitToRepo(t *testing.T) {
	handler, err = handler.Commit("Initial commit.", mockUsername, mockEmail)
	assert := assert.New(t)
	assert.Nilf(err, "Failed to commit to repo: %+v", err)
}
func TestPushFileToRepo(t *testing.T) {
	err = handler.Push()
	assert := assert.New(t)
	assert.Nilf(err, "Failed to push files to repo: %+v", err)
}
