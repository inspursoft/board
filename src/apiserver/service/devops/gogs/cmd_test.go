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
	Username: "user12",
	Password: "123456a?",
	Email:    "user12@inspur.com",
}

var mockPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzaDsh+RgEO+VdXnKKFfH0a2GLfomldSrUCS0wfvXBXETmhUJ+r5pvyZBXlIoUd4D3kMPKnKk1oqYa4qks31BYEajfHYpMVve5MhBNKZM5wS+MlL1Aa6vxMwCJcjp0X6vpzOjtD3TEdkQtqxyPsYm11fK0XeWILZBinOR9L6vBIOwjaz891VgNmM6RBZtbCKy8RV8ejevsFkUWcYh71+85HqHPp0DiB0CefZTpz8G3HM+941E9K0FWY82slgBKtUEjvxShSVUmMPbY3i/hjLCaqS5+UQqpzosuZlMtpgzyKEDF0iIXU5+sOAOYpHOnBvxzZ+XpKOJ845WLPeSzgDjv wangkun@wangkuns-MacBook-Pro.local`

var token *accessToken

func TestSignUp(t *testing.T) {
	err := SignUp(user)
	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("Error occurred while signing up: %+v", err))
}

func TestCreateAccessToken(t *testing.T) {
	var err error
	token, err = CreateAccessToken(user.Username, user.Password)
	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("Error occurred while creating token: %+v", err))
	assert.NotNil(token, "Failed to initialize Gogs access token.")
	logs.Info("Access Token: %+v", token)
}

// func TestCreatePublicKey(t *testing.T) {
// 	err := NewGogsHandler(user.Username, token.Sha1).CreatePublicKey("userPublicKey", mockPublicKey)
// 	assert := assert.New(t)
// 	assert.Nilf(err, fmt.Sprintf("Error occurred while creating public key: %+v", err))
// }

// func TestDeletePublicKey(t *testing.T) {
// 	err:=NewGogsHandler(user.Username, token.Sha1).DeletePublicKey(2)
// }

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
