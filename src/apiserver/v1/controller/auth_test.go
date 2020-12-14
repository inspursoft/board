package controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

const (
	AdminUsername = "boardadmin"
	AdminPassword = "MTIzNDU2YT8="
)

func loginTest(t *testing.T, username, password string) string {
	token := signIn(username, password)
	assert := assert.New(t)
	if !assert.NotEmpty(token, "signIn error") {
		// logs and failNow
		t.Fatalf("%s Failed to login\n", username)
	}
	return token
}

func logoutTest(t *testing.T, username string) {
	err := signOut(username)
	if err != nil {
		t.Fatalf("%s Failed to logout", username)
	}
}

func adminLoginTest(t *testing.T) string {
	return loginTest(t, AdminUsername, AdminPassword)
}

func adminLogoutTest(t *testing.T) {
	logoutTest(t, AdminUsername)
}

func TestGetSystemInfo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/v1/systeminfo", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "Get Systeminfo fail.")
}

func signIn(name, password string) string {
	var reqUser model.User
	var token model.Token
	reqUser.Username = name
	reqUser.Password = password

	req, err := json.Marshal(reqUser)
	if err != nil {
		return ""
	}
	body := ioutil.NopCloser(strings.NewReader(string(req)))
	r, _ := http.NewRequest("POST", "/api/v1/sign-in", body)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		return ""
	}
	err = json.Unmarshal(w.Body.Bytes(), &token)
	if err != nil {
		return ""
	}
	return token.TokenString
}

func signOut(name string) error {
	reqURL := "/api/v1/log-out?username=" + name
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		errstr := fmt.Sprintf("Logout error: %+v", w.Body.String())
		return errors.New(errstr)
	}
	return nil
}

func TestSignInOutAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn(AdminUsername, AdminPassword)
	assert.NotEmpty(token, "signIn error")

	err := signOut(AdminUsername)
	assert.Nil(err, "signOut error")
}

func TestCurrentUserAction(t *testing.T) {
	var user model.User

	assert := assert.New(t)
	token := signIn(AdminUsername, AdminPassword)
	defer signOut(AdminUsername)
	assert.NotEmpty(token, "signIn error")

	reqURL := "/api/v1/users/current?token=" + token
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code, "Get current user fail.")

	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(err, "Unmarshal user error.")
	assert.Equal(AdminUsername, user.Username, "Get current user error.")
}

func TestUserExists(t *testing.T) {
	assert := assert.New(t)
	token := signIn(AdminUsername, AdminPassword)
	defer signOut(AdminUsername)
	assert.NotEmpty(token, "signIn error")

	reqURL := "/api/v1/user-exists?token=" + token + "&target=username&value=boardadmin"
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusConflict, w.Code, "User exist test fail.")

	reqURL = "/api/v1/user-exists?token=" + token + "&target=username&value=admin1"
	r, _ = http.NewRequest("GET", reqURL, nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code, "User exist test fail.")
}

func TestSignUpAction(t *testing.T) {
	var user model.User
	user.Username = "testuser1"
	user.Email = "testuser1@inspur.com"
	user.Password = `MTIjJHF3RVI=`
// 	assert := assert.New(t)
// 	req, err := json.Marshal(user)
// 	assert.Nil(err, "user marshal fail")
// 	body := ioutil.NopCloser(strings.NewReader(string(req)))

// 	fmt.Println(body)
// 	reqURL := "/api/v1/sign-up"
// 	r, _ := http.NewRequest("POST", reqURL, body)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	assert.Equal(http.StatusOK, w.Code, "User sign up test fail.")
}
