package service_test

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

var user = model.User{
	Username: "Tester",
	Email:    "tester@inspur.com",
	Password: "123456a?",
}

func TestGetUserByID(t *testing.T) {
	assert := assert.New(t)
	u, err := service.GetUserByID(1)
	assert.Nil(err, "Error occurred while calling GetUserByID method.")
	assert.NotNil(u, "User does not exists.")
	assert.Equal("boardadmin", u.Username, "Username is not equal to be expected.")
}

func TestUsernameExists(t *testing.T) {
	assert := assert.New(t)
	exists, err := service.UserExists("username", "", 0)
	assert.Nil(err, "Error occurred while checking username exists.")
	assert.False(exists, "Username exists.")
}
