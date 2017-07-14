package service

import (
	"errors"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"log"
)

func SignUp(user model.User) (bool, error) {
	userID, err := dao.AddUser(user)
	if err != nil {
		return false, err
	}
	return (userID != 0), nil
}

func SignIn(principal string, password string) (bool, error) {
	query := model.User{Username: principal, Password: password}
	user, err := dao.GetUser(query, "username", "password")
	if err != nil {
		log.Printf("Failed to get user in SignIn: %+v\n", err)
		return false, err
	}
	return (user != nil && user.ID != 0), nil
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

func UsernameExists(username string) (bool, error) {
	query := model.User{Username: username}
	user, err := dao.GetUser(query, "username")
	if err != nil {
		return false, err
	}
	return (user != nil && user.ID != 0), nil
}

func EmailExists(email string) (bool, error) {
	query := model.User{Email: email}
	user, err := dao.GetUser(query, "email")
	if err != nil {
		return false, err
	}
	return (user != nil && user.ID != 0), nil
}

func IsSysAdmin(userID int64) (bool, error) {
	query := model.User{ID: userID}
	user, err := dao.GetUser(query, "id")
	if err != nil {
		return false, err
	}
	return (user != nil && user.ID != 0 && user.SystemAdmin == 1), nil
}
