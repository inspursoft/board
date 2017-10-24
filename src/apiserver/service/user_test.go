package service

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var user = model.User{
	Username: "Tester",
	Password: "123456a?",
}

func connectToDB() {
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(mysql:3306)/board?charset=utf8")
	if err != nil {
		logs.Error("Failed to connect to DB.")
	}
}

func cleanUp(username string) {
	o := orm.NewOrm()
	rs := o.Raw("delete from user where username = ?", username)
	r, err := rs.Exec()
	if err != nil {
		logs.Error("Error occurred while deleting user: %+v", err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		logs.Error("Error occurred while deleting user: %+v", err)
	}
	if affected == 0 {
		logs.Error("Failed to delete user")
	} else {
		logs.Error("Successful cleared up.")
	}
}

func TestMain(m *testing.M) {
	connectToDB()
	cleanUp(user.Username)
	os.Exit(m.Run())
}

func TestSignIn(t *testing.T) {
	assert := assert.New(t)
	u, err := SignIn("admin", "123456a?")
	assert.Nil(err, "Error occurred while calling SignIn method.")
	assert.NotNil(u, "User is nil.")
	assert.Equal("admin", u.Username, "Signed in failed.")
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
