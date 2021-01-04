package service_test

import (
	"fmt"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	projectID   = 1
	projectName = "testProject"
	userID      = 1
	userName    = "boardadmin"
)

var project = model.Project{
	Name:      "testProject",
	OwnerID:   userID,
	OwnerName: userName,
}

var updatedProject = model.Project{
	ID:     1,
	Public: 1,
}

func TestCreateProject(t *testing.T) {
	exists, err := service.ProjectExists(projectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while checking project name existing.")
	if !exists {
		isSuccess, err := service.CreateProject(project)
		assert.Nil(err, "Error occurred while creating project.")
		assert.Equal(true, isSuccess, "Failed to create project")
	}
}

func TestGetProject(t *testing.T) {
	project, err := service.GetProjectByName(projectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while getting project by name.")
	assert.NotNilf(project, "Failed to get project by name: %s", projectName)
	assert.Equal(projectName, project.Name, "Project name is not equal to expected value.")
}

func TestProjectExists(t *testing.T) {
	exists, err := service.ProjectExists(projectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while checking project name existing.")
	assert.Equal(true, exists, "Project name does not exist.")
}

func TestProjectExistsByID(t *testing.T) {
	exists, err := service.ProjectExistsByID(projectID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while checking project ID existing.")
	assert.Equal(true, exists, "Project ID does not exist.")
}

func TestUpdateProject(t *testing.T) {
	isSuccess, err := service.UpdateProject(updatedProject, "public")
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while updating project public.")
	assert.Equal(true, isSuccess, "Failed to update project public.")
}

func TestGetProjectsByUser(t *testing.T) {
	query := model.Project{Name: "library"}
	projectList, err := service.GetProjectsByUser(query, userID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get projects by user.")
	assert.NotNilf(projectList, "Failed to get projects by name: %s", query.Name)
	assert.Lenf(projectList, 1, "Failed to get specific project by name with user ID: %d", userID)
}

func TestGetProjectsByMember(t *testing.T) {
	query := model.Project{}
	projectList, err := service.GetProjectsByMember(query, userID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get projects by member.")
	assert.NotNil(projectList, fmt.Sprintf("Failed to get projects by userID: %d", userID))
}

// func TestDeleteProject(t *testing.T) {
// 	isSuccess, err := service.DeleteProject(userID, projectID)
// 	assert := assert.New(t)
// 	assert.Nil(err, "Error occurred while deleting project.")
// 	assert.Equalf(true, isSuccess, "Failed to delete project by ID: %d", projectID)

// 	isSuccess, err = service.DeleteNamespace(projectName)
// 	assert.Nil(err, "Error occurred while deleting namespace.")
// 	assert.Equalf(true, isSuccess, "Failed to delete namespace by name: %s", projectName)

// 	project, err := service.GetProjectByID(projectID)
// 	assert.Nilf(err, "Error occurred while getting project by ID: %d", projectID)
// 	assert.Nilf(project, "Project with ID: %d is not nil.", projectID)
// }
