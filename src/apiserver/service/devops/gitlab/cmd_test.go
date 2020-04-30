package gitlab_test

import (
	"git/inspursoft/board/src/apiserver/service/devops/gitlab"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var user = model.User{
	Username: "testuser18",
	Email:    "testuser18@inspur.com",
	Password: "123456a?",
}
var project = model.Project{
	Name: "myrepo",
}
var createdUser gitlab.UserCreation
var token gitlab.ImpersonationToken
var createdProject gitlab.ProjectCreation

var adminAccessToken = "si1Z1eUZUVui7XFarUyW"

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	os.Exit(m.Run())
}

func TestUserCreation(t *testing.T) {
	var err error
	createdUser, err = gitlab.NewGitlabHandler(adminAccessToken).CreateUser(user)
	assert.New(t).Nilf(err, "Error occurred while creating user via Gitlab API: %+v", err)
}

func TestImpersonateToken(t *testing.T) {
	var err error
	token, err = gitlab.NewGitlabHandler(adminAccessToken).ImpersonationToken(createdUser)
	logs.Debug("Impersonated token: %+v", token)
	assert.New(t).Nilf(err, "Error occurred while impersonating token via Gitlab API: %+v", err)
}

func TestCreateRepo(t *testing.T) {
	var err error
	createdProject, err = gitlab.NewGitlabHandler(token.Token).CreateRepo(user, project)
	assert.New(t).Nilf(err, "Error occurred while creating repo via Gitlab API: %+v", err)
}

func TestDeleteRepo(t *testing.T) {
	err := gitlab.NewGitlabHandler(token.Token).DeleteProject(createdProject.ID)
	assert.New(t).Nilf(err, "Error occurred while deleting project via Gitlab API: %+v", err)
}

func TestUserDeletion(t *testing.T) {
	err := gitlab.NewGitlabHandler(adminAccessToken).DeleteUser(createdUser.ID)
	assert.New(t).Nilf(err, "Error occurred while deleting user via Gitlab API: %+v", err)
}
