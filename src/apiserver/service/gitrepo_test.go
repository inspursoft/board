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
	configurations["flag"] = "image"
	configurations["extras"] = "10.110.13.134:5000/library/myimage20180201:v1.0"
	configurations["file_name"] = "Dockerfile"
	configurations["docker_registry"] = "10.110.13.134:5000"
	configurations["apiserver"] = "10.165.14.97:8089"
	configurations["value"] = "Dockerfile"

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
