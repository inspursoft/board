package controllers

import (
	"git/inspursoft/board/src/adminserver/service"
	"net/http"

	"github.com/astaxie/beego"
)

// BoardController controlls Board up and down.
type BoardController struct {
	beego.Controller
}

// @Title Restart
// @Description restart Board
// @Param	token	query 	string	true		"token"
// @Success 200 success
// @Failure 400 bad request
// @router /restart [get]
func (b *BoardController) Restart() {
	var statusCode int = http.StatusOK
	statusMessage := service.Restart("/root/BOARD/Deploy")
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	b.Ctx.ResponseWriter.WriteHeader(statusCode)
	b.ServeJSON()
}

// @Title Applycfg
// @Description apply cfg and restart Board
// @Param	token	query 	string	true	"token"
// @Success 200 success
// @Failure 400 bad request
// @router /applycfg [get]
func (b *BoardController) Applycfg() {
	var statusCode int = http.StatusOK
	statusMessage := service.Applycfg("/root/BOARD/Deploy")
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	b.Ctx.ResponseWriter.WriteHeader(statusCode)
	b.ServeJSON()
}

// @Title Shutdown
// @Description shutdown board
// @Param	token	query 	string	true	"token"
// @Success 200 success
// @Failure 400 bad request
// @router /shutdown [get]
func (b *BoardController) Shutdown() {
	var statusCode int = http.StatusOK
	statusMessage := service.Shutdown("/root/BOARD/Deploy")
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	b.Ctx.ResponseWriter.WriteHeader(statusCode)
	b.ServeJSON()
}