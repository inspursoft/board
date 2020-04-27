package controllers

import (
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// BoardController controlls Board up and down.
type BoardController struct {
	BaseController
}

// @Title Start
// @Description start Board
// @Param	token	query 	string	true		"token"
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 500 Internal Server Error
// @Failure 401 unauthorized
// @router /start [post]
func (b *BoardController) Start() {
	var host models.Account
	err := utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	if err = service.Start(&host); err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	b.ServeJSON()
}

// @Title Applycfg
// @Description apply cfg and restart Board
// @Param	token	query 	string	true	"token"
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 500 Internal Server Error
// @Failure 401 unauthorized
// @router /applycfg [post]
func (b *BoardController) Applycfg() {
	var host models.Account
	err := utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	if err = service.Applycfg(&host); err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	b.ServeJSON()
}

// @Title Shutdown
// @Description shutdown board
// @Param	token	query 	string	true	"token"
// @Param	uninstall	query 	bool	true	"uninstall flag"
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 500 Internal Server Error
// @Failure 401 unauthorized
// @router /shutdown [post]
func (b *BoardController) Shutdown() {
	var host models.Account
	uninstall, err := b.GetBool("uninstall")
	if err != nil {
		logs.Error("Failed to get bool data: %+v", err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	err = utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	if err = service.Shutdown(&host, uninstall); err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	b.ServeJSON()
}
