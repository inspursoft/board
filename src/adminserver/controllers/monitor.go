package controllers

import (
	"git/inspursoft/board/src/adminserver/service"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// MonitorController includes operations about monitoring.
type MonitorController struct {
	BaseController
}

// @Title Get
// @Description monitor Board module containers
// @Param	token	query 	string	true	"token"
// @Success 200 {object} []models.Boardinfo	success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router / [get]
func (m *MonitorController) Get() {
	containers, err := service.GetMonitor()
	if err != nil {
		logs.Error(err)
		m.CustomAbort(http.StatusBadRequest, err.Error())
	}
	//apply struct to JSON value.
	m.Data["json"] = containers
	m.ServeJSON()
}
