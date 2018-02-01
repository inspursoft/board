package service

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockInitRepoURL = "ssh://git@localhost:10022"
var mockInitRepoPath = "/repos"
var mockProjectName = "library"
var mockUsername = "admin"
var mockEmail = "admin@inspur.com"
var mockRepoPath = filepath.Join(mockInitRepoPath, mockUsername, mockProjectName)
var mockFileName = "temp.md"
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
	configurations := make(map[string]string)
	configurations["job_name"] = "process_image"
	configurations["file_name"] = "Dockerfile"
	err := CreateBaseDirectory(configurations, mockRepoPath)
	assert := assert.New(t)
	assert.Nilf(err, "Failed to create base directory: %+v", err)
	_, err = handler.Add("META.cfg")
	assert.Nilf(err, "Failed to add files to repo: %+v", err)
	_, err = handler.Add("process-image/.placehold.tmp")
	assert.Nilf(err, "Failed to add files to repo: %+v", err)
	_, err = handler.Add("process-service/.placehold.tmp")
	assert.Nilf(err, "Failed to add files to repo: %+v", err)
	_, err = handler.Add("rolling-update/.placehold.tmp")
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
