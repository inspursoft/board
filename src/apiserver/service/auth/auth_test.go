package auth

import (
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/assert"
)

func connectToDB() {
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(localhost:3306)/board?charset=utf8")
	if err != nil {
		logs.Error("Failed to connect to DB.")
	}
}

func TestMain(m *testing.M) {
	connectToDB()
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
