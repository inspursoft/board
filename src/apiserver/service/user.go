package service

import (
	"errors"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
)

func SignUp(user model.User) (bool, error) {
	user.Salt = utils.GenerateRandomString()
	user.Password = utils.Encrypt(user.Password, user.Salt)
	userID, err := dao.AddUser(user)
	if err != nil {
		return false, err
	}
	return (userID != 0), nil
}

func GetUserByID(userID int64) (*model.User, error) {
	query := model.User{ID: userID, Deleted: 0}
	user, err := dao.GetUser(query, "id", "deleted")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUsers(field string, value interface{}, selectedFields ...string) ([]*model.User, error) {
	return dao.GetUsers(field, value, selectedFields...)
}

func GetPaginatedUsers(field string, value interface{}, pageIndex int, pageSize int, selectedField ...string) (*model.PaginatedUsers, error) {
	return dao.GetPaginatedUsers(field, value, pageIndex, pageSize, selectedField...)
}

func UpdateUser(user model.User, selectedFields ...string) (bool, error) {
	if user.ID == 0 {
		return false, errors.New("No user ID provided.")
	}
	_, err := dao.UpdateUser(user, selectedFields...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteUser(userID int64) (bool, error) {
	user := model.User{ID: userID, Deleted: 1}
	_, err := dao.UpdateUser(user, "deleted")
	if err != nil {
		return false, err
	}
	return true, nil
}

func UserExists(fieldName string, value string, userID int64) (bool, error) {
	query := model.User{ID: userID, Username: value, Email: value}
	user, err := dao.GetUser(query, fieldName)
	if err != nil {
		return false, err
	}
	if userID == 0 {
		return (user != nil && user.ID != 0), nil
	}
	return (user != nil && user.ID != userID), nil
}

func IsSysAdmin(userID int64) (bool, error) {
	query := model.User{ID: userID}
	user, err := dao.GetUser(query, "id")
	if err != nil {
		return false, err
	}
	return (user != nil && user.ID != 0 && user.SystemAdmin == 1), nil
}
