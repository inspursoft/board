package service

import (
	"errors"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
)

func CreateProject(project model.Project) (bool, error) {
	projectID, err := dao.AddProject(project)
	if err != nil {
		return false, err
	}
	userQuery := model.User{ID: 1}
	user, err := dao.GetUser(userQuery)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("none of user found")
	}
	projectMember := model.ProjectMember{ProjectID: projectID, UserID: user.ID}
	projectMemberID, err := dao.InsertOrUpdateProjectMember(projectMember)
	if err != nil {
		return false, errors.New("failed to create project member")
	}
	return (projectID != 0 && projectMemberID != 0), nil
}

func GetProject(project model.Project) (*model.Project, error) {
	p, err := dao.GetProject(project)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func ProjectExists(projectName string) (bool, error) {
	query := model.Project{Name: projectName}
	project, err := dao.GetProject(query, "name")
	if err != nil {
		return false, err
	}
	return (project.ID != 0), nil
}

func ProjectExistsByID(projectID int64) (bool, error) {
	query := model.Project{ID: projectID}
	project, err := dao.GetProject(query, "id")
	if err != nil {
		return false, err
	}
	return (project.Name != ""), nil
}

func UpdateProject(project model.Project, fieldNames ...string) (bool, error) {
	if project.ID == 0 {
		return false, errors.New("no Project ID provided")
	}
	_, err := dao.UpdateProject(project, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetProjects(fieldName string, value interface{}, selectedFields ...string) ([]*model.Project, error) {
	return dao.GetProjects(fieldName, value, selectedFields...)
}

func DeleteProject(projectID int64) (bool, error) {
	project := model.Project{ID: projectID}
	_, err := dao.DeleteProject(project)
	if err != nil {
		return false, err
	}
	return true, nil
}
