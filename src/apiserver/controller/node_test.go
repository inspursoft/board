package controller_test

// import (
// 	"git/inspursoft/board/src/common/utils"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/astaxie/beego"
// 	"github.com/stretchr/testify/assert"
// )

// var nodeIP = utils.GetConfig("NODE_IP")

// func TestGetNode(t *testing.T) {
// 	token := adminLoginTest(t)
// 	defer adminLogoutTest(t)

// 	r, _ := http.NewRequest("GET", "/api/v1/node?node_name="+nodeIP()+"&token="+token, nil)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	assert := assert.New(t)
// 	assert.Equal(http.StatusOK, w.Code, "Get Node fail.")
// }

// func TestNodeList(t *testing.T) {
// 	token := adminLoginTest(t)
// 	defer adminLogoutTest(t)

// 	r, _ := http.NewRequest("GET", "/api/v1/nodes?token="+token, nil)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	assert := assert.New(t)
// 	assert.Equal(http.StatusOK, w.Code, "Get Nodes fail.")
// }

// func TestNodeToggle(t *testing.T) {
// 	token := adminLoginTest(t)
// 	defer adminLogoutTest(t)

// 	r, _ := http.NewRequest("GET", "/api/v1/node/toggle?node_name="+nodeIP()+"&node_status=false&token="+token, nil)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	assert := assert.New(t)
// 	if !assert.Equal(http.StatusOK, w.Code, "Toggle Node false fail.") {
// 		t.FailNow()
// 	}

// 	r, _ = http.NewRequest("GET", "/api/v1/node/toggle?node_name="+nodeIP()+"&node_status=true&token="+token, nil)
// 	w = httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	assert.Equal(http.StatusOK, w.Code, "Toggle Node true fail.")
// }
