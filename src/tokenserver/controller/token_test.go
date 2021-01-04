package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

const (
	TOKEN_BAD = "abcd"
	PAYLOAD   = `{"id": "1", "username": "zhangsan", "email": "zhangsan@inspur.com", "realname": "zhangsan", "is_project_admin": 1, "is_system_admin": 0}`
)

type ErrorReader struct {
}

func (*ErrorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error reader")
}

func TestPost(t *testing.T) {
	assert := assert.New(t)

	// Post a nil body format reqeust.
	r, _ := http.NewRequest("POST", "/tokenservice/token", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.NotEqual(http.StatusOK, w.Code, "Request with nil body should failed.")

	// Post a not json body request.
	r, _ = http.NewRequest("POST", "/tokenservice/token", bytes.NewBuffer([]byte(TOKEN_BAD)))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.NotEqual(http.StatusOK, w.Code, "Request with not json body should failed.")

	// Post a error body request.
	r, _ = http.NewRequest("POST", "/tokenservice/token", nil)
	r.Body = ioutil.NopCloser(new(ErrorReader))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.NotEqual(http.StatusOK, w.Code, "Request with error body should failed.")
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	// Get without token param reqeust.
	r, _ := http.NewRequest("GET", "/tokenservice/token", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.NotEqual(http.StatusOK, w.Code, "Rrequest with no token parameter should failed.")

	// Get with a bad token parameter request.
	r, _ = http.NewRequest("GET", "/tokenservice/token?token="+TOKEN_BAD, nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.NotEqual(http.StatusOK, w.Code, "Request with bad token parameter should failed.")
}

func TestToken(t *testing.T) {
	assert := assert.New(t)

	// Retrive a token
	r, _ := http.NewRequest("POST", "/tokenservice/token", bytes.NewBuffer([]byte(PAYLOAD)))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code, "Retrive token fail.")

	respData := w.Body.String()
	token := model.Token{}
	err := json.Unmarshal([]byte(respData), &token)
	assert.Nil(err, fmt.Sprintf("Unmarshal json message %s to map error:%v", respData, err))

	// Decode the token
	r, _ = http.NewRequest("GET", "/tokenservice/token?token="+token.TokenString, nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code, "Decode token fail.")
	decode := w.Body.String()
	assert.JSONEq(PAYLOAD, decode, "The decoded token is not the same as origin message")
}
