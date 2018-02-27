package controller

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

const (
	AdminUsername = "admin"
	AdminPassword = "123456a?"
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

func TestGetNode(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	nodeIP := os.Getenv("NODE_IP")
	r, _ := http.NewRequest("GET", "/api/v1/node?node_name="+nodeIP+"&token="+token, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "Get Node fail.")
}

func TestNodeList(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	r, _ := http.NewRequest("GET", "/api/v1/nodes?token="+token, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "Get Nodes fail.")
}

func TestNodeToggle(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	nodeIP := os.Getenv("NODE_IP")

	r, _ := http.NewRequest("GET", "/api/v1/node/toggle?node_name="+nodeIP+"&node_status=false&token="+token, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert := assert.New(t)
	if !assert.Equal(http.StatusOK, w.Code, "Toggle Node false fail.") {
		t.FailNow()
	}

	r, _ = http.NewRequest("GET", "/api/v1/node/toggle?node_name="+nodeIP+"&node_status=true&token="+token, nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code, "Toggle Node true fail.")
}
