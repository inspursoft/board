package service

import (
	"errors"
	"git/inspursoft/board/src/apiserver/dao"
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
	return (user.ID != 0), nil
}

func GetUsers(field string, value interface{}, selectedFields ...string) ([]*model.User, error) {
	return dao.GetUsers(field, value, selectedFields...)
}

func UpdateUser(user model.User) (bool, error) {
	if user.ID == 0 {
		return false, errors.New("No user ID provided.")
	}
	count, err := dao.UpdateUser(user)
	if err != nil {
		return false, err
	}
	return (count != 0), nil
}

func DeleteUser(userID int64) (bool, error) {
	user := model.User{ID: userID}
	affected, err := dao.DeleteUser(user)
	if err != nil {
		return false, err
	}
	return (affected != 0), nil
}

func UsernameExists(username string) (bool, error) {
	query := model.User{Username: username}
	user, err := dao.GetUser(query, "username")
	if err != nil {
		return false, err
	}
	return (user.ID != 0), nil
}

func EmailExists(email string) (bool, error) {
	query := model.User{Email: email}
	user, err := dao.GetUser(query, "email")
	if err != nil {
		return false, err
	}
	return (user.ID != 0), nil
}
