package auth

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/assert"
)

const (
	adminUserID     = 1
	initialPassword = "123456a?"
)

func updateAdminPassword() {
	salt := utils.GenerateRandomString()
	encryptedPassword := utils.Encrypt(initialPassword, salt)
	user := model.User{ID: adminUserID, Password: encryptedPassword, Salt: salt}
	isSuccess, err := service.UpdateUser(user, "password", "salt")
	if err != nil {
		logs.Error("Failed to update user password: %+v", err)
	}
	if isSuccess {
		logs.Info("Admin password has been updated successfully.")
	} else {
		logs.Info("Failed to update admin initial password.")
	}
}

func connectToDB() {
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(localhost:3306)/board?charset=utf8")
	if err != nil {
		logs.Error("Failed to connect to DB.")
	}
}

func TestMain(m *testing.M) {
	connectToDB()
	updateAdminPassword()
	os.Exit(m.Run())
}

func TestSignIn(t *testing.T) {
	assert := assert.New(t)
	currentAuth, err := GetAuth("db_auth")
	u, err := (*currentAuth).DoAuth("admin", "123456a?")
	assert.Nil(err, "Error occurred while calling SignIn method.")
	assert.NotNil(u, "User is nil.")
	assert.Equal("admin", u.Username, "Signed in failed.")
}
