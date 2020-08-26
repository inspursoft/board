package gitlabci_test

import (
	"git/inspursoft/board/src/apiserver/service/devops/gitlabci"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGitlabCI(t *testing.T) {
	job1 := gitlabci.Job{
		Stage:  "test1",
		Script: []string{"echo hello"},
		Tags:   []string{"board-test-vm"},
	}

	job2 := gitlabci.Job{
		Stage:  "test2",
		Script: []string{"echo world"},
		Tags:   []string{"board-test-vm"},
	}
	ci := make(map[string]gitlabci.Job)
	ci["job1"] = job1
	ci["job2"] = job2
	var gc gitlabci.GitlabCI
	err := gc.GenerateGitlabCI(ci, ".")
	assert := assert.New(t)
	assert.Nil(err, "Failed to create Gitlab CI yaml file")
	os.Remove("output.yaml")
}
