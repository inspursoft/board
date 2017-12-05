package service

import (
	"git/inspursoft/board/src/common/model"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

var user = model.User{
	Username: "Tester",
	Password: "123456a?",
}

func TestGetUserByID(t *testing.T) {
	assert := assert.New(t)
	u, err := GetUserByID(1)
	assert.Nil(err, "Error occurred while calling GetUserByID method.")
	assert.NotNil(u, "User does not exists.")
	assert.Equal("admin", u.Username, "Username is not equal to be expected.")
}

func TestSignUp(t *testing.T) {
	assert := assert.New(t)
	status, err := SignUp(user)
	assert.Nil(err, "Error occurred while calling SignUp method.")
	assert.True(status, "Signed up failed.")
}
