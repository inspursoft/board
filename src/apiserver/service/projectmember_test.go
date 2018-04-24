package service

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"testing"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

const (
	testMemberProjectName = "testmemberproject"
	testMemberName        = "testmember"
	roleID                = 1 /* project admin */
	roleName              = "projectAdmin"
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
	isSuccess, err = SignUp(testMember)
	if err != nil {
		logs.Error("Error occurred while signning up: %+v", err)
	}
	if !isSuccess {
		logs.Error("Failed to sign up user: %+v", testMember)
	}
	isSuccess, err = CreateProject(testMemberProject)
	if err != nil {
		logs.Error("Error occurred while creating project: %+v", err)
	}
	if !isSuccess {
		logs.Error("Failed to create project for testing project member: %+v", testMemberProject)
	}
}

func cleanupProjectMember() {
	o := orm.NewOrm()
	var affectedCount int64
	var err error
	affectedCount, err = o.Delete(&testMemberProject)
	if err != nil {
		logs.Error("Failed to delete project: %+v", err)
	}
	logs.Info("Deleted in project %d row(s) affected.", affectedCount)
	affectedCount, err = o.Delete(&testMember)
	if err != nil {
		logs.Error("Failed to delete member: %+v", err)
	}
	logs.Info("Deleted in member %d row(s) affected.", affectedCount)
}

func getTestMemberProject() (*model.Project, error) {
	queryProject := model.Project{Name: testMemberProjectName}
	return GetProject(queryProject, "name")
}

func getTestMember() (*model.User, error) {
	userList, err := GetUsers("username", testMemberName, "id")
	if err != nil {
		logs.Error("Error occurred while getting users by name: %+v", err)
	}
	if len(userList) > 0 {
		return userList[0], nil
	}
	return nil, fmt.Errorf("failed to get user by name: %s", testMemberName)
}

func TestAddOrUpdateProjectMember(t *testing.T) {
	prepareForProjectMember()
	project, err := getTestMemberProject()
	if err != nil {
		logs.Error("Failed to get project.")
	}
	assert := assert.New(t)
	isSuccess, err := AddOrUpdateProjectMember(project.ID, adminUserID, roleID)
	assert.Nilf(err, "Error occurred while adding project member by projectID: %d, userID: %d, roleID: %d", project.ID, adminUserID, roleID)
	assert.Equalf(true, isSuccess, "Failed to add project member by projectID: %d, userID: %d, roleID: %d", project.ID, adminUserID, roleID)
}

func TestGetProjectMembers(t *testing.T) {
	project, err := getTestMemberProject()
	if err != nil {
		logs.Error("Failed to get project.")
	}
	assert := assert.New(t)
	projectMemberList, err := GetProjectMembers(project.ID)
	assert.Nilf(err, "Error occurred while getting project member by projectID: %d.", project.ID)
	assert.Lenf(projectMemberList, 1, "Failed to get project members by projectID: %d.", project.ID)
}

func TestHasProjectAdminRole(t *testing.T) {
	project, err := getTestMemberProject()
	if err != nil {
		logs.Error("Failed to get project.")
	}
	projectMember, err := getTestMember()
	if err != nil {
		logs.Error("Failed to get project member.")
	}
	isSuccess, err := HasProjectAdminRole(project.ID, projectMember.ID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while checking project admin role by projectID: %d, userID: %d.", project.ID, projectMember.ID)
	assert.Equalf(true, isSuccess, "Failed to check project admin role by projectID: %d, userID: %d.", project.ID, projectMember.ID)
}

func TestIsProjectMember(t *testing.T) {
	project, err := getTestMemberProject()
	if err != nil {
		logs.Error("Failed to get project.")
	}
	projectMember, err := getTestMember()
	if err != nil {
		logs.Error("Failed to get project member.")
	}
	isSuccess, err := IsProjectMember(project.ID, projectMember.ID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while checking whether user is the project member by projectID: %d, userID: %d.", project.ID, projectMember.ID)
	assert.Equalf(true, isSuccess, "Failed to check whether user is the project member by projectID: %d, userID: %d.", project.ID, projectMember.ID)
}

func TestGetRoleByID(t *testing.T) {
	role, err := GetRoleByID(roleID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while getting role: %+v", err)
	assert.Equalf(roleName, role.Name, "Role name is not as expected: %s", roleName)
}

func TestIsProjectMemberByName(t *testing.T) {
	projectMember, err := getTestMember()
	if err != nil {
		logs.Error("Failed to get project member.")
	}
	isSuccess, err := IsProjectMemberByName(projectName, projectMember.ID)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while checking whether the user is project member by project name: %s", testMemberProjectName)
	assert.Equalf(true, isSuccess, "User %s is project %s's member.", projectMember.Username, testMemberProjectName)
}

func TestDeleteProjectMember(t *testing.T) {
	project, err := getTestMemberProject()
	if err != nil {
		logs.Error("Failed to get project.")
	}
	assert := assert.New(t)
	isSuccess, err := DeleteProjectMember(project.ID, adminUserID)
	assert.Nilf(err, "Error occurred while deleting project member by projectID: %d, userID: %d.", project.ID, adminUserID)
	assert.Equalf(true, isSuccess, "Failed to delete project member by projectID: %d, userID: %d.", project.ID, adminUserID)
}
