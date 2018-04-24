package service

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"testing"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/orm"

	"github.com/stretchr/testify/assert"
)

const (
	projectID   = 1
	projectName = "testProject"
	userID      = 1
	userName    = "admin"
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

func cleanUpProject() {
	o := orm.NewOrm()
	affectedCount, err := o.Delete(&project)
	if err != nil {
		logs.Error("Failed to clean up project: %+v", err)
	}
	logs.Info("Deleted in project %d row(s) affected.", affectedCount)
}

func TestCreateProject(t *testing.T) {
	exists, err := ProjectExists(projectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while checking project name existing.")
	if !exists {
		isSuccess, err := CreateProject(project)
		assert.Nil(err, "Error occurred while creating project.")
		assert.Equal(true, isSuccess, "Failed to create project")
	}
}

func TestGetProject(t *testing.T) {
	query := model.Project{Name: projectName}
	project, err := GetProject(query, "name")
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while getting project by name.")
	assert.NotNilf(project, "Failed to get project by name: %s", projectName)
	assert.Equal(projectName, project.Name, "Project name is not equal to expected value.")
}

func TestProjectExists(t *testing.T) {
	exists, err := ProjectExists(projectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while checking project name existing.")
	assert.Equal(true, exists, "Project name does not exist.")
}

func TestProjectExistsByID(t *testing.T) {
	exists, err := ProjectExistsByID(projectID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while checking project ID existing.")
	assert.Equal(true, exists, "Project ID does not exist.")
}

func TestUpdateProject(t *testing.T) {
	isSuccess, err := UpdateProject(updatedProject, "public")
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while updating project public.")
	assert.Equal(true, isSuccess, "Failed to update project public.")
}

func TestGetProjectsByUser(t *testing.T) {
	query := model.Project{Name: "library"}
	projectList, err := GetProjectsByUser(query, userID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get projects by user.")
	assert.NotNilf(projectList, "Failed to get projects by name: %s", query.Name)
	assert.Lenf(projectList, 1, "Failed to get specific project by name with user ID: %d", userID)
}

func TestGetProjectsByMember(t *testing.T) {
	query := model.Project{}
	projectList, err := GetProjectsByMember(query, userID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get projects by member.")
	assert.NotNil(projectList, fmt.Sprintf("Failed to get projects by userID: %d", userID))
	assert.Lenf(projectList, 2, "Failed to get projects by member with userID: %d", userID)
}

func TestDeleteProject(t *testing.T) {
	isSuccess, err := DeleteProject(projectID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while deleting project.")
	assert.Equalf(true, isSuccess, "Failed to delete project by ID: %d", projectID)

	isSuccess, err = DeleteNamespace(projectName)
	assert.Nil(err, "Error occurred while deleting namespace.")
	assert.Equalf(true, isSuccess, "Failed to delete namespace by name: %s", projectName)

	project, err := GetProject(model.Project{ID: projectID}, "id")
	assert.Nilf(err, "Error occurred while getting project by ID: %d", projectID)
	assert.Nilf(project, "Project with ID: %d is not nil.", projectID)
}
