package controllers

import (
	"git/inspursoft/board/src/adminserver/service"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// MoniController includes operations about monitoring.
type MoniController struct {
	beego.Controller
}

func (m *MoniController) Prepare() {
	token := m.Ctx.Request.Header.Get("token")
	if token == "" {
		token = m.GetString("token")
	}
	result, err := service.VerifyToken(token)
	if err != nil {
		logs.Error(err)
		m.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if !result {
		m.CustomAbort(http.StatusUnauthorized, "Unauthorized")	
	} 
}

// @Title Get
// @Description monitor Board module containers
// @Param	token	query 	string	true	"token"
// @Success 200 {object} []models.Boardinfo	success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router / [get]
func (m *MoniController) Get() {
	containers, err := service.GetMonitor()
	if err != nil {
		logs.Error(err)
		m.CustomAbort(http.StatusBadRequest, err.Error())
	}
	//apply struct to JSON value.
	m.Data["json"] = containers
	m.ServeJSON()
}
