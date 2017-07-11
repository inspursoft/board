package service

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
)

func AddOrUpdateProjectMember(projectID int64, userID int64, roleID int64) (bool, error) {
	projectMember := model.ProjectMember{ProjectID: projectID, UserID: userID, RoleID: roleID}
	_, err := dao.InsertOrUpdateProjectMember(projectMember)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetProjectMembers(projectID int64, userID int64) ([]*model.User, error) {
	return dao.GetProjectMembers(model.Project{ID: projectID}, model.User{ID: userID})
}

func DeleteProjectMember(projectID int64, userID int64) (bool, error) {
	projectMember := model.ProjectMember{ID: projectID + userID}
	_, err := dao.DeleteProjectMember(projectMember)
	if err != nil {
		return false, err
	}
	return true, nil
}
