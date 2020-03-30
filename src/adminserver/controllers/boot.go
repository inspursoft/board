package controllers

import (
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"fmt"
)

// BootController includes operations about booting config.
type BootController struct {
	beego.Controller
}

// @Title Initdb
// @Description init db password and max number of connections.
// @Param	body	body 	models.DBconf	true	"body for db conf"
// @Success 200 success
// @Failure 400 bad request
// @router /initdb [post]
func (b *BootController) Initdb() {
	var db models.DBconf
	utils.UnmarshalToJSON(b.Ctx.Request.Body, &db)
	if err := service.InitDB(&db); err != nil {
		b.CustomAbort(http.StatusBadRequest, err.Error())
		logs.Error(err)
		return
	}
	b.ServeJSON()	
	return
}

// @Title Startdb
// @Description ssh to host and docker-compose up the db
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 400 bad request
// @router /startdb [post]
func (b *BootController) Startdb() {
	var host models.Account
	utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
	if err := service.StartDB(&host); err != nil {
		b.CustomAbort(http.StatusBadRequest, err.Error())
		logs.Error(err)
		return
	}
	b.ServeJSON()	
	return
}

// @Title StartBoard
// @Description ssh to host and docker-compose up the Board
// @Param	token	query 	string	true	"token"
// @Param	body	body 	models.Account	true	"body for host acc info"
// @Success 200 success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router /startboard [post]
func (b *BootController) Start() {
	var host models.Account
	token := b.GetString("token")
	result := service.VerifyToken(token)
	if result == false {
		b.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		b.ServeJSON()	
		return
	} else {
		utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
		if err := service.StartBoard(&host); err != nil {
			b.CustomAbort(http.StatusBadRequest, err.Error())
			logs.Error(err)
			return
		}
		b.ServeJSON()	
		return
	}
}

// @Title CheckDB
// @Description Check db status
// @Success 200 success
// @Failure 400 bad request
// @router /checkdb [get]
func (b *BootController) CheckDB() {
	var statusCode int
	if service.CheckDB() == true {
		statusCode = http.StatusOK
	} else {
		b.CustomAbort(http.StatusBadRequest, fmt.Sprintf("DB is down."))
		return
	}
	b.Ctx.ResponseWriter.WriteHeader(statusCode)
	b.ServeJSON()	
	return
}

