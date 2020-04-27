package controllers

import (
	"git/inspursoft/board/src/adminserver/service"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type BaseController struct {
	beego.Controller
}

func (b *BaseController) Prepare() {
	token := b.Ctx.Request.Header.Get("token")
	if token == "" {
		token = b.GetString("token")
	}

	if err := service.CheckBoard(); err != nil {
		result, err := service.VerifyUUIDToken(token)
		if err != nil {
			logs.Error(err)
			b.CustomAbort(http.StatusBadRequest, err.Error())
		}
		if !result {
			b.CustomAbort(http.StatusUnauthorized, "Unauthorized")
		}
	} else {
		if user := service.GetCurrentUser(token); user == nil {
			b.CustomAbort(http.StatusUnauthorized, "Unauthorized")
		}
	}
}
