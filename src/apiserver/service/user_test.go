package service

import (
	"git/inspursoft/board/src/common/model"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

var user = model.User{
	Username: "Tester",
	Password: "123456a?",
}

func cleanUpUser() {
	o := orm.NewOrm()
	affectedCount, err := o.Delete(&user)
	if err != nil {
		logs.Error("Failed to clean up user: %+v", err)
	}
	logs.Info("Deleted  in user %d row(s) affected.", affectedCount)
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

func TestUsernameExists(t *testing.T) {
	assert := assert.New(t)
	exists, err := UserExists("username", "", 0)
	assert.Nil(err, "Error occurred while checking username exists.")
	assert.False(exists, "Username exists.")
}
