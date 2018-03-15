package controller

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

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
