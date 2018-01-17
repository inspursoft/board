package gogs

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

func prepare() {

}

func cleanUp() {

}

func TestMain(m *testing.M) {
	os.Exit(func() int {
		prepare()
		m.Run()
		cleanUp()
		return 0
	}())
}

var user = model.User{
	Username: "tester",
	Password: "123456a?",
	Email:    "tester@inspur.com",
}

var token *accessToken

func TestCreateAccessToken(t *testing.T) {
	var err error
	token, err = CreateAccessToken(user.Username, user.Password)
	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("Error occurred while creating token: %+v", err))
	assert.NotNil(token, "Failed to initialize Gogs access token.")
	logs.Info("Access Token: %+v", token)
}

func TestSignUp(t *testing.T) {
	err := SignUp(user)
	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("Error occurred while signing up: %+v", err))
}

func TestCreateRepo(t *testing.T) {
	err := NewGogsHandler(user.Username, token.Sha1).CreateRepo("myrepo")
	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("Error occurred while creating repo: %+v", err))
}

func TestDeleteRepo(t *testing.T) {
	err := NewGogsHandler(user.Username, token.Sha1).DeleteRepo("myrepo")
	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("Error occurred while deleting repo: %+v", err))
}
