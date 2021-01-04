package service

import (
	"fmt"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"

	"errors"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddOrUpdateProjectMember(projectID int64, userID int64, roleID int64) (bool, error) {
	projectMember := model.ProjectMember{ProjectID: projectID, UserID: userID, RoleID: roleID}
	_, err := dao.InsertOrUpdateProjectMember(projectMember)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetProjectMembers(projectID int64) ([]*model.ProjectMember, error) {
	return dao.GetProjectMembers(model.Project{ID: projectID})
}

func GetProjectAvailableMembers(projectID int64) ([]*model.User, error) {
	return dao.GetProjectAvailableMembers(model.Project{ID: projectID})
}

func DeleteProjectMember(projectID int64, userID int64) (bool, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("no user was found with provided user ID")
	}
	project, err := GetProjectByID(projectID)
	if err != nil {
		return false, fmt.Errorf("failed to get project by ID: %d, error: %+v", projectID, err)
	}

	repoName, err := ResolveRepoName(project.Name, user.Username)
	if err != nil {
		return false, fmt.Errorf("failed to resolve repo name with project: %s, username: %s, error: %+v", project.Name, user.Username, err)
	}
	err = CurrentDevOps().DeleteRepo(user.Username, repoName)
	if err != nil {
		logs.Warning("failed to delete repo with name: %s, error: %+v", repoName, err)
	}

	projectMember := model.ProjectMember{ID: projectID + userID}
	_, err = dao.DeleteProjectMember(projectMember)
	if err != nil {
		return false, err
	}
	return true, nil
}

func IsProjectMember(projectID int64, userID int64) (bool, error) {
	members, err := GetProjectMembers(projectID)
	if err != nil {
		return false, err
	}
	var isMember bool
	for _, m := range members {
		if m.UserID == userID {
			isMember = true
			break
		}
	}
	return isMember, nil
}

func GetRoleByID(roleID int64) (*model.Role, error) {
	role, err := dao.GetRole(model.Role{ID: roleID}, "id")
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return role, nil
}

func IsProjectMemberByName(projectName string, userID int64) (bool, error) {
	queryProject := model.Project{Name: projectName}
	project, err := GetProject(queryProject, "name")
	if err != nil {
		return false, err
	}
	if project == nil {
		return false, errors.New("invalid project ID")
	}
	return IsProjectMember(project.ID, userID)
}
