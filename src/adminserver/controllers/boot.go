package controllers

import (
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// BootController includes operations about booting config.
type BootController struct {
	BaseController
}

func (b *BootController) Prepare() {}

// @Title CheckSysStatus
// @Description return the current system status.
// @Success 200 {object} models.InitSysStatus success
// @Failure 500 Internal Server Error
// @router /checksysstatus [get]
func (b *BootController) CheckSysStatus() {
	this, err := service.CheckSysStatus()
	if err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	b.Data["json"] = models.InitSysStatus{Status: this, Log: b.buf.String()}
	b.ServeJSON()
}
