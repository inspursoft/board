package controller

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
)

func TestGetSystemInfo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/v1/systeminfo", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Get Systeminfo fail: %+v", w.Body.String())
	} else {
		t.Log("Get Systeminfo successfully.")
	}
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
	} else {
		err = json.Unmarshal(w.Body.Bytes(), &token)
		if err != nil {
			return ""
		}
		return token.TokenString
	}
}

func signOut(name string) error {
	reqURL := "/api/v1/log-out?username=" + name
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		errstr := fmt.Sprintf("Logout error: %+v", w.Body.String())
		return errors.New(errstr)
	} else {
		return nil
	}
}

func TestSignInOutAction(t *testing.T) {
	token := signIn("admin", "123456a?")
	if token == "" {
		t.Errorf("signIn error.")
	} else {
		t.Log("signIn successfully.")
	}

	err := signOut("admin")
	if err != nil {
		t.Errorf("signOut error: %+v", err)
	} else {
		t.Log("signOut successfully.")
	}
}

func TestCurrentUserAction(t *testing.T) {
	token := signIn("admin", "123456a?")
	if token == "" {
		t.Errorf("signIn error")
		return
	}
	defer signOut("admin")

	reqURL := "/api/v1/users/current?token=" + token
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Get current user fail, %+v", w.Body.String())
	} else {
		var user model.User
		err := json.Unmarshal(w.Body.Bytes(), &user)
		if err != nil {
			t.Errorf("Get current user fail, %+v", err)
		} else {
			if user.Username != "admin" {
				t.Errorf("Get current user error, want:\"admin\", get:%+v.", user.Username)
			} else {
				t.Log("Get current user successfully.")
			}
		}
	}
}
