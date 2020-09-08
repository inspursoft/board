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
	Username: "testuser24",
	Email:    "testuser24@inspur.com",
	Password: "123456a?",
}

var forkUser = model.User{
	Username: "testuser25",
	Email:    "testuser25@inspur.com",
	Password: "123456a?",
}

var project = model.Project{
	Name: "myrepo11",
}

var createdUser gitlab.UserInfo
var createdForkUser gitlab.UserInfo

var createdUserToken gitlab.ImpersonationToken
var createdForkUserToken gitlab.ImpersonationToken

var addSSHKeyResponse gitlab.AddSSHKeyResponse
var createdProject gitlab.ProjectCreation
var createdForkProject gitlab.ProjectCreation

var createdMR gitlab.MRCreation

var branch = "master"
var sourceBranch = branch
var targetBranch = sourceBranch

var fileInfo = gitlab.FileInfo{
	Name:    "README.md",
	Path:    "content/README.md",
	Content: "# myrepo",
}

var adminAccessToken = utils.GetConfig("GITLAB_ADMIN_TOKEN")
var sshPubKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDaa046MKqllR1bE0pfPYwcVYHmBx291OzeWj5VHS6FCsVeLnky99pJigp3uwDz68uTDOx1I+zUU3XE39o4591isCCbM9ba5l2hKvGnHUoTRdG6Pkc9gy+OdKIJMGFca58Bt1hhPCa5FT8cQadsSnr7rGmg1O5tfG6a9mjzKFjn3nNNlYi5U6BsJxD3ReV5mVkFea5wH2yMzrHCSxTQiyLM8owB9Dem7Mrqz799sfB9MjC6ryVGwJd8oZOxGCB7hNz/Eenb+EUjdevxLFAVZgakTk4vDm/ubVfQjrdxGg4MaAbD4+kYNezEfh9c5W2uC0QlZHQhItEoMqytmWmjeZF7 root@10.110.25.227"

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	os.Exit(m.Run())
}

func TestUserCreation(t *testing.T) {
	var err error
	createdUser, err = gitlab.NewGitlabHandler(adminAccessToken()).CreateUser(user)
	assert.New(t).Nilf(err, "Error occurred while creating user via Gitlab API: %+v", err)
}

func TestImpersonateToken(t *testing.T) {
	var err error
	createdUserToken, err = gitlab.NewGitlabHandler(adminAccessToken()).ImpersonationToken(createdUser)
	logs.Debug("Impersonated token: %+v", createdUserToken)
	assert.New(t).Nilf(err, "Error occurred while impersonating token via Gitlab API: %+v", err)
}

// func TestAddSSHKey(t *testing.T) {
// 	var err error
// 	addSSHKeyResponse, err = gitlab.NewGitlabHandler(createdUserToken.Token).AddSSHKey("user-ssh-key", sshPubKey)
// 	assert := assert.New(t)
// 	assert.Nilf(err, "Error occurred while adding SSH key via gitlab API: %+v", err)
// 	assert.NotNilf(addSSHKeyResponse, "Failed to get response after adding SSH key.", nil)

// }

func TestCreateRepo(t *testing.T) {
	var err error
	createdProject, err = gitlab.NewGitlabHandler(createdUserToken.Token).CreateRepo(user, project)
	assert.New(t).Nilf(err, "Error occurred while creating repo via Gitlab API: %+v", err)
	project.ID = int64(createdProject.ID)
}

func TestGetRepo(t *testing.T) {
	foundProjectList, err := gitlab.NewGitlabHandler(createdUserToken.Token).GetRepoInfo(project)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while get repo via Gitlab API: %+v", err)
	assert.NotNilf(foundProjectList, "Failed to get repo after creating repo.", nil)
}

func TestCreateFileToRepo(t *testing.T) {
	targetProject := model.Project{
		ID:   int64(createdProject.ID),
		Name: createdProject.Name,
	}
	fileCreation, err := gitlab.NewGitlabHandler(createdUserToken.Token).ManipulateFile("create", user, targetProject, branch, fileInfo)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while creating file to repo via Gitlab API: %+v", err)
	assert.NotNilf(fileCreation, "Failed to create file: %+v", fileInfo)
}

func TestForkRepo(t *testing.T) {
	var err error
	createdForkUser, err = gitlab.NewGitlabHandler(adminAccessToken()).CreateUser(forkUser)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while creating fork user via Gitlab API: %+v", err)
	assert.NotNilf(createdForkUser, "Failed to create fork user with detail: %+v", forkUser)
	forkUser.ID = int64(createdForkUser.ID)

	createdForkUserToken, err = gitlab.NewGitlabHandler(adminAccessToken()).ImpersonationToken(createdForkUser)
	assert.Nilf(err, "Error occurred while impersonating token to fork user via Gitlab API: %+v", err)
	assert.NotNilf(createdForkUser, "Failed to impersonate token with forked user detail: %+v", forkUser)

	memberUser, err := gitlab.NewGitlabHandler(createdUserToken.Token).AddMemberToRepo(forkUser, project)
	assert.Nilf(err, "Error occurred while adding member: %s to the project ID: %d, error: %+v", forkUser.Username, project.ID, err)
	assert.NotNilf(memberUser, "Failed to add member to the project with detail: %+v", forkUser)

	forkRepoName := forkUser.Username + "_" + createdProject.Name
	assert.Nilf(err, "Error occurred while resolving repo name with username %+v", err)

	createdForkProject, err = gitlab.NewGitlabHandler(createdForkUserToken.Token).ForkRepo(createdProject.ID, forkRepoName)
	assert.Nilf(err, "Error occurred while forking project via Gitlab API: %+v", err)
	assert.NotNilf(createdForkProject, "Failed to fork project with detail: %+v", createdForkProject)
}

func TestCreateMR(t *testing.T) {
	fileInfo.Content = "# myrepo with updated"
	forkProject := model.Project{ID: int64(createdForkProject.ID)}
	gitlabHandler := gitlab.NewGitlabHandler(createdForkUserToken.Token)
	updatedFileCreation, err := gitlabHandler.ManipulateFile("update", forkUser, forkProject, branch, fileInfo)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while updating file via Gitlab API: %+v", err)
	assert.NotNilf(updatedFileCreation, "Failed to update file with forked user detail: %+v", forkUser)
	user.ID = int64(createdUser.ID)
	project.ID = int64(createdProject.ID)
	mrCreation, err := gitlabHandler.CreateMR(forkUser, forkProject, project, sourceBranch, targetBranch, "Update README.md", "Update README.md file.")
	assert.Nilf(err, "Error occurred while creating MR via Gitlab API: %+v", err)
	assert.NotNilf(mrCreation, "Failed to create MR with detail: %+v", mrCreation)
}

func TestMergeMR(t *testing.T) {
	gitlabHandler := gitlab.NewGitlabHandler(createdUserToken.Token)
	mrList, err := gitlabHandler.ListMR(project)
	assert := assert.New(t)
	assert.Lenf(mrList, 1, "No MR found for repo: %s", project.Name)
	assert.Nilf(err, "Error occurred while list MR via Gitlab API: %+v", err)
	createdMR = mrList[0]
	mrAcceptance, err := gitlabHandler.AcceptMR(project, createdMR.IID)
	assert.Nilf(err, "Error occurred while merging MR via Gitlab API: %+v", err)
	assert.NotNilf(mrAcceptance, "Failed to merge MR with detail: %+v", mrAcceptance)
}

func TestDeleteRepo(t *testing.T) {
	assert := assert.New(t)
	err := gitlab.NewGitlabHandler(createdForkUserToken.Token).DeleteProject(createdForkProject.ID)
	assert.Nilf(err, "Error occurred while deleting fork project via Gitlab API: %+v", err)
	err = gitlab.NewGitlabHandler(createdUserToken.Token).DeleteProject(createdProject.ID)
	assert.Nilf(err, "Error occurred while deleting project via Gitlab API: %+v", err)
}

func TestUserDeletion(t *testing.T) {
	err := gitlab.NewGitlabHandler(adminAccessToken()).DeleteUser(createdForkUser.ID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while deleting fork user via Gitlab API: %+v", err)
	err = gitlab.NewGitlabHandler(adminAccessToken()).DeleteUser(createdUser.ID)
	assert.Nilf(err, "Error occurred while deleting user via Gitlab API: %+v", err)
}
