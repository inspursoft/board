package controllers

import (
	"github.com/inspursoft/board/src/adminserver/service"
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
		logs.Info(err)
		result, err := service.VerifyUUIDToken(token)
		if err != nil {
			logs.Error(err)
			b.CustomAbort(http.StatusInternalServerError, err.Error())
		}
		if !result {
			b.CustomAbort(http.StatusUnauthorized, "UUID invalid or timeout")
		}
		b.Ctx.ResponseWriter.Header().Set("token", token)
	} else {
		user, newtoken := service.GetCurrentUser(token)
		if user == nil {
			b.CustomAbort(http.StatusUnauthorized, "Unauthorized")
		}
		b.Ctx.ResponseWriter.Header().Set("token", newtoken)
	}
}
