package controllers

import (
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// BoardController controlls Board up and down.
type BoardController struct {
	beego.Controller
}

// @Title Restart
// @Description restart Board
// @Param	token	query 	string	true		"token"
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router /restart [post]
func (b *BoardController) Restart() {
	var host models.Account
	token := b.GetString("token")
	result := service.VerifyToken(token)
	if !result {
		b.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		b.ServeJSON()	
	} else {
		utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
		if err := service.Restart(&host); err != nil {
			logs.Error(err)
			b.CustomAbort(http.StatusBadRequest, err.Error())
		}
		b.ServeJSON()	
	}
}

// @Title Applycfg
// @Description apply cfg and restart Board
// @Param	token	query 	string	true	"token"
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router /applycfg [post]
func (b *BoardController) Applycfg() {
	var host models.Account
	token := b.GetString("token")
	result := service.VerifyToken(token)
	if !result {
		b.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		b.ServeJSON()	
	} else {
		utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
		if err := service.Applycfg(&host); err != nil {
			logs.Error(err)
			b.CustomAbort(http.StatusBadRequest, err.Error())
		}
		b.ServeJSON()	
	}
}

// @Title Shutdown
// @Description shutdown board
// @Param	token	query 	string	true	"token"
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router /shutdown [post]
func (b *BoardController) Shutdown() {
	var host models.Account
	token := b.GetString("token")
	result := service.VerifyToken(token)
	if !result {
		b.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		b.ServeJSON()	
	} else {
		utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
		if err := service.Shutdown(&host); err != nil {
			logs.Error(err)
			b.CustomAbort(http.StatusBadRequest, err.Error())
		}
		b.ServeJSON()	
	}
}
