package jenkins_test

import (
	"git/inspursoft/board/src/apiserver/service/devops/jenkins"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	os.Exit(0)
}

var user = model.User{
	Username: "tester",
	Password: "123456a?",
	Email:    "tester@inspur.com",
}

func TestCreateJob(t *testing.T) {
	err := jenkins.NewJenkinsHandler().CreateJobWithParameter("testproject")
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while creating job: %+v", err)
}

func TestDeleteJob(t *testing.T) {
	err := jenkins.NewJenkinsHandler().DeleteJob("testproject")
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while deleting job: %+v", err)
}
