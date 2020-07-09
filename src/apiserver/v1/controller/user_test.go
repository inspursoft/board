package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/assert"
)

func isExternalAuth() bool {
	AuthMode := utils.GetStringValue("AUTH_MODE")
	if AuthMode != "" && AuthMode != "db_auth" {
		return true
	}
	return false
}

func TestUserAction(t *testing.T) {
	// external auth ignore the test
	if isExternalAuth() {
		return
	}

	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	user := model.User{
		Username:    "testuser",
		Password:    "MTIzNDU2YT8=",
		Email:       "testuser@test.com",
		Realname:    "testuser",
		Comment:     "this is just a test account",
		SystemAdmin: 0,
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("user marshal error: %v", err)
	}
	// init one assert
	assert := assert.New(t)
	// add user
	t.Log("adding user")
	r, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/adduser?token=%s", token), bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Add User fail.") {
		t.FailNow()
	}
	defer cleanUp(user.Username)
	// get users
	t.Log("getting users")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users?username=%s&token=%s", user.Username, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get Users fail.") {
		t.FailNow()
	}
	readUsers := make([]model.User, 0)
	err = json.Unmarshal(w.Body.Bytes(), &readUsers)
	if err != nil {
		t.Fatalf("user unmarshal error: %v", err)
	}
	if len(readUsers) == 0 {
		t.Fatalf("can't find the user which just created")
	}
	readUser := readUsers[0]

	t.Log("getting one user")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get User fail.") {
		t.FailNow()
	}

	// udpate user
	t.Log("updating user")
	readUser.Comment = "update comment"
	body, err = json.Marshal(readUser)
	if err != nil {
		t.Fatalf("read user marshal error: %v", err)
	}
	r, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Update Users fail.") {
		t.FailNow()
	}

	t.Log("getting one user")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get User fail.") {
		t.FailNow()
	}

	updatedUser := new(model.User)
	err = json.Unmarshal(w.Body.Bytes(), updatedUser)
	if err != nil {
		t.Fatalf("updateduser unmarshal error: %v", err)
	}

	if !assert.Equal(readUser.Comment, updatedUser.Comment, "Upate User Comment fail.") {
		t.FailNow()
	}

	// toggle system admin
	t.Log("toggling user")
	updatedUser.SystemAdmin = 1
	body, err = json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("updated user marshal error: %v", err)
	}
	r, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d/systemadmin?token=%s", readUser.ID, token), bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Update Users fail.") {
		t.FailNow()
	}

	t.Log("getting one user")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get User fail.") {
		t.FailNow()
	}

	toggledUser := new(model.User)
	err = json.Unmarshal(w.Body.Bytes(), toggledUser)
	if err != nil {
		t.Fatalf("toggledUser unmarshal error: %v", err)
	}

	if !assert.Equal(updatedUser.SystemAdmin, toggledUser.SystemAdmin, "Toggle User SystemAdmin fail.") {
		t.FailNow()
	}

	t.Log("deleteing user")
	r, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), nil)
	w = httptest.NewRecorder()
	beego.BConfig.WebConfig.EnableXSRF = false
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code, "Delete User fail.")
}

func TestChangeUserAccount(t *testing.T) {
	// external auth ignore the test
	if isExternalAuth() {
		return
	}

	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	user := model.User{
		Username:    "testuseraccount",
		Password:    "dGVzdHVzZXJwYXNzd3Jk",
		Email:       "testuseraccount@test.com",
		Realname:    "testuseraccount",
		Comment:     "this is just a test account",
		SystemAdmin: 0,
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("user marshal error: %v", err)
	}
	// init one assert
	assert := assert.New(t)

	// add user
	t.Log("adding user")
	r, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/adduser?token=%s", token), bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Add User fail.") {
		t.FailNow()
	}
	defer cleanUp(user.Username)

	// get users
	t.Log("getting users")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users?username=%s&token=%s", user.Username, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get Users fail.") {
		t.FailNow()
	}
	readUsers := make([]model.User, 0)
	err = json.Unmarshal(w.Body.Bytes(), &readUsers)
	if err != nil {
		t.Fatalf("user unmarshal error: %v", err)
	}
	if len(readUsers) == 0 {
		t.Fatalf("can't find the user which just created")
	}
	readUser := readUsers[0]

	t.Log("getting one user")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get User fail.") {
		t.FailNow()
	}

	// login using testuseraccount
	userToken := loginTest(t, user.Username, user.Password)
	defer logoutTest(t, user.Username)
	// change user account
	t.Log("updating user")
	readUser.Comment = "update comment"
	body, err = json.Marshal(readUser)
	if err != nil {
		t.Fatalf("read user marshal error: %v", err)
	}
	r, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/changeaccount?token=%s", userToken), bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Change User account fail.") {
		t.FailNow()
	}

	t.Log("getting one user")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get User fail.") {
		t.FailNow()
	}

	updatedUser := new(model.User)
	err = json.Unmarshal(w.Body.Bytes(), updatedUser)
	if err != nil {
		t.Fatalf("updateduser unmarshal error: %v", err)
	}
	if !assert.Equal(readUser.Comment, updatedUser.Comment, "Upate User Comment fail.") {
		t.FailNow()
	}

}

func TestChangePasswordAction(t *testing.T) {
	// external auth ignore the test
	if isExternalAuth() {
		return
	}

	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	user := model.User{
		Username:    "testpassword",
		Password:    "dGVzdHVzZXJwYXNzd3Jk",
		Email:       "testpasswordpwd@test.com",
		Realname:    "testpasswordpwd",
		Comment:     "this is just a test account",
		SystemAdmin: 0,
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("user marshal error: %v", err)
	}
	// init one assert
	assert := assert.New(t)

	// add user
	t.Log("adding user")
	r, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/adduser?token=%s", token), bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Add User fail.") {
		t.FailNow()
	}
	defer cleanUp(user.Username)

	// get users
	t.Log("getting users")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users?username=%s&token=%s", user.Username, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get Users fail.") {
		t.FailNow()
	}
	readUsers := make([]model.User, 0)
	err = json.Unmarshal(w.Body.Bytes(), &readUsers)
	if err != nil {
		t.Fatalf("user unmarshal error: %v", err)
	}
	if len(readUsers) == 0 {
		t.Fatalf("can't find the user which just created")
	}
	readUser := readUsers[0]

	t.Log("getting one user")
	r, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d?token=%s", readUser.ID, token), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Get User fail.") {
		t.FailNow()
	}

	// login using testuseraccount
	userToken := loginTest(t, user.Username, user.Password)
	defer logoutTest(t, user.Username)

	// admin change user password has bug. so we use use itself
	t.Log("change user password")
	changePwd := new(model.ChangePassword)
	changePwd.OldPassword = user.Password
	changePwd.NewPassword = "MTIzNDU2YT8="
	body, err = json.Marshal(changePwd)
	if err != nil {
		t.Fatalf("read user marshal error: %v", err)
	}
	r, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d/password?token=%s", readUser.ID, userToken), bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if !assert.Equal(http.StatusOK, w.Code, "Change User password fail.") {
		t.FailNow()
	}

	// using new password to login
	loginTest(t, user.Username, changePwd.NewPassword)

}

func cleanUp(username string) {
	o := orm.NewOrm()
	rs := o.Raw("delete from user where username = ?", username)
	r, err := rs.Exec()
	if err != nil {
		logs.Error("Error occurred while deleting user: %+v", err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		logs.Error("Error occurred while deleting user: %+v", err)
	}
}
