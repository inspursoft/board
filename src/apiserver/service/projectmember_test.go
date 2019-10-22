package service_test

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

const (
	testMemberProjectName = "testmemberproject"
	testMemberName        = "testmember"
	roleID                = 1
	roleName              = "projectAdmin"
	adminUserID           = 1
)

var testMember = model.User{
	Username: testMemberName,
	Email:    "testmember@inspur.com",
	Password: "123456a?",
}

var testMemberProject = model.Project{
	Name: testMemberProjectName,
}

func prepareForProjectMember() {
	var isSuccess bool
	var err error
	isSuccess, err = service.SignUp(testMember)
	if err != nil {
		logs.Error("Error occurred while signning up: %+v", err)
	}
	if !isSuccess {
		logs.Error("Failed to sign up user: %+v", testMember)
	}
	user, err := service.GetUserByName(testMemberName)
	if err != nil {
		logs.Error("Failed to get user by name: %s, error: %+v", testMemberName, err)
	}
	testMemberProject.OwnerID = int(user.ID)
	testMemberProject.OwnerName = user.Username
	isSuccess, err = service.CreateProject(testMemberProject)
	if err != nil {
		logs.Error("Error occurred while creating project: %+v", err)
	}
	if !isSuccess {
		logs.Error("Failed to create project for testing project member: %+v", testMemberProject)
	}
}

func TestAddOrUpdateProjectMember(t *testing.T) {
	prepareForProjectMember()
	project, err := service.GetProjectByName(testMemberProjectName)
	if err != nil {
		logs.Error("Failed to get project.")
	}
	assert := assert.New(t)
	isSuccess, err := service.AddOrUpdateProjectMember(project.ID, adminUserID, roleID)
	assert.Nilf(err, "Error occurred while adding project member by projectID: %d, userID: %d, roleID: %d", project.ID, adminUserID, roleID)
	assert.Equalf(true, isSuccess, "Failed to add project member by projectID: %d, userID: %d, roleID: %d", project.ID, adminUserID, roleID)
}

func TestGetProjectMembers(t *testing.T) {
	project, err := service.GetProjectByName(testMemberProjectName)
	if err != nil {
		logs.Error("Failed to get project.")
	}
	assert := assert.New(t)
	projectMemberList, err := service.GetProjectMembers(project.ID)
	assert.Nilf(err, "Error occurred while getting project member by projectID: %d.", project.ID)
	assert.Lenf(projectMemberList, 2, "Failed to get project members by projectID: %d.", project.ID)
}

func TestIsProjectMember(t *testing.T) {
	project, err := service.GetProjectByName(testMemberProjectName)
	if err != nil {
		logs.Error("Failed to get project.")
	}
	projectMember, err := service.GetUserByName(testMemberName)
	if err != nil {
		logs.Error("Failed to get project member.")
	}
	isSuccess, err := service.IsProjectMember(project.ID, projectMember.ID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while checking whether user is the project member by projectID: %d, userID: %d.", project.ID, projectMember.ID)
	assert.Equalf(true, isSuccess, "Failed to check whether user is the project member by projectID: %d, userID: %d.", project.ID, projectMember.ID)
}

func TestGetRoleByID(t *testing.T) {
	role, err := service.GetRoleByID(roleID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while getting role: %+v", err)
	assert.Equalf(roleName, role.Name, "Role name is not as expected: %s", roleName)
}

func TestIsProjectMemberByName(t *testing.T) {
	projectMember, err := service.GetUserByName(testMemberName)
	if err != nil {
		logs.Error("Failed to get project member.")
	}
	isSuccess, err := service.IsProjectMemberByName(testMemberProjectName, projectMember.ID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while checking whether the user is project member by project name: %s", testMemberProjectName)
	assert.Equalf(true, isSuccess, "User %s is project %s's member.", projectMember.Username, testMemberProjectName)
}

// func TestDeleteProjectMember(t *testing.T) {
// 	project, err := service.GetProjectByName(testMemberProjectName)
// 	if err != nil {
// 		logs.Error("Failed to get project.")
// 	}
// 	assert := assert.New(t)
// 	isSuccess, err := service.DeleteProjectMember(project.ID, adminUserID)
// 	assert.Nilf(err, "Error occurred while deleting project member by projectID: %d, userID: %d.", project.ID, adminUserID)
// 	assert.Equalf(true, isSuccess, "Failed to delete project member by projectID: %d, userID: %d.", project.ID, adminUserID)
// }
