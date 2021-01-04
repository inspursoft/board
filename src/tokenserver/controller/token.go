package controller

import (
	"fmt"
	"github.com/inspursoft/board/src/tokenserver/service"

	"net/http"

	"github.com/inspursoft/board/src/common/model"

	"io/ioutil"

	"encoding/json"

	"github.com/astaxie/beego"
)

type TokenServiceController struct {
	beego.Controller
}

func (t *TokenServiceController) Render() error {
	return nil
}

func (t *TokenServiceController) ServeStatus(status int, message string) {
	ms := make(map[string]interface{})
	ms["status"] = status
	ms["message"] = message
	t.Data["json"] = ms
	t.Ctx.ResponseWriter.WriteHeader(status)
	t.ServeJSON()
}

func (t *TokenServiceController) Post() {
	if t.Ctx.Request.Body == nil {
		t.ServeStatus(http.StatusBadRequest, "Posting token request body is nil.")
		return
	}
	var err error
	reqData, err := ioutil.ReadAll(t.Ctx.Request.Body)
	if err != nil {
		t.ServeStatus(http.StatusInternalServerError, "Failed to get data from request.")
		return
	}
	var tokenPayload map[string]interface{}
	err = json.Unmarshal(reqData, &tokenPayload)
	if err != nil {
		t.ServeStatus(http.StatusInternalServerError, "Failed to unmarshal JSON.")
		return
	}
	tokenString, err := service.Sign(tokenPayload)
	if err != nil {
		t.ServeStatus(http.StatusInternalServerError, "Failed to sign token.")
		return
	}
	t.Data["json"] = model.Token{TokenString: tokenString}
	t.ServeJSON()
}

func (t *TokenServiceController) Get() {
	token := t.GetString("token")
	payload, err := service.Verify(token)
	if err != nil {
		t.ServeStatus(http.StatusUnauthorized, fmt.Sprintf("Failed to verify token: %+v", err))
		return
	}
	t.Data["json"] = payload
	t.ServeJSON()
}

func init() {
	ns := beego.NewNamespace("/tokenservice/", beego.NSRouter("/token", &TokenServiceController{}))
	beego.AddNamespace(ns)
}
