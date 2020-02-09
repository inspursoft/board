package controllers

import (
	"git/inspursoft/board/src/adminserver/service"
	"net/http"

	"github.com/astaxie/beego"
)

// MoniController includes operations about monitoring.
type MoniController struct {
	beego.Controller
}

// @Title Get
// @Description monitor Board module containers
// @Param	token	query 	string	true	"token"
// @Success 200 {object} []models.Boardinfo	success
// @Failure 400 bad request
// @router / [get]
func (m *MoniController) Get() {
	var statusCode int = http.StatusOK
	containers, statusMessage := service.GetMonitor()
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	m.Ctx.ResponseWriter.WriteHeader(statusCode)
	//apply struct to JSON value.
	m.Data["json"] = containers
	m.ServeJSON()
}
