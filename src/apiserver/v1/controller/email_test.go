package controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/inspursoft/board/src/apiserver/v1/controller"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var emailReq = controller.EmailPingParam{
	Username: "admin@myserver.com",
	Password: AdminPassword,
	Hostname: "smtp.myserver.com",
	Port:     12225,
}

func TestEmailPing(t *testing.T) {
	assert := assert.New(t)
	token := signIn(AdminUsername, AdminPassword)
	defer signOut(AdminUsername)
	assert.NotEmpty(token, "signIn error")

	data, err := json.Marshal(emailReq)
	assert.Nilf(err, "Failed to marshal SMTP server request: %+v", err)

	reqURL := "/api/v1/email/ping?token=" + token
	r, _ := http.NewRequest("POST", reqURL, bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code, "Get registry fail.")
	logs.Info("Tested GetImageRegistry %s pass", w.Body.String())
}
