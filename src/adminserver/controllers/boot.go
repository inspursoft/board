package controllers

import (
	"encoding/json"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

// BootController includes operations about booting config.
type BootController struct {
	beego.Controller
}

// @Title Initdb
// @Description init db password and max number of connections.
// @Param	body	body 	models.DBconf	true	"body for db conf"
// @Success 200 {object} string success
// @Failure 400 bad request
// @router /initdb [post]
func (b *BootController) Initdb() {
	var db models.DBconf
	json.Unmarshal(b.Ctx.Input.RequestBody, &db)
	if err := service.InitDB(&db); err != nil {
		b.CustomAbort(http.StatusBadRequest, err.Error())
		logs.Error(err)
		return
	}
	b.ServeJSON()	
	return
}