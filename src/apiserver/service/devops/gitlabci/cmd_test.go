package gitlabci_test

import (
	"git/inspursoft/board/src/apiserver/service/devops/gitlabci"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGitlabCI(t *testing.T) {
	var gc gitlabci.GitlabCI
	tagName := "abc:v1.0"
	image := gitlabci.Image{
		Name:       "kaniko",
		Entrypoint: []string{"bash", "-c", "ls", "-al"},
	}
	job1 := gitlabci.Job{
		Image:  image,
		Stage:  "test1",
		Script: []string{"echo hello"},
		Tags:   []string{"board-test-vm"},
	}
	job2 := gitlabci.Job{
		Image: image,
		Stage: "test2",
		Script: []string{
			"echo world",
			gc.WriteMultiLine("docker build  -f toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong -t %s .", tagName),
			gc.WriteMultiLine("docker build  -f toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong -t %s .", tagName),
			gc.WriteMultiLine("docker build  -f toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong -t %s .", tagName),
		},
		Tags: []string{"board-test-vm"},
	}
	ci := make(map[string]gitlabci.Job)
	ci["job1"] = job1
	ci["job2"] = job2

	err := gc.GenerateGitlabCI(ci, ".")
	assert := assert.New(t)
	assert.Nil(err, "Failed to create Gitlab CI yaml file")
	os.Remove(gitlabci.GitlabCIFilename)
}
