package controllers

import (
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// BootController includes operations about booting config.
type BootController struct {
	BaseController
}

func (b *BootController) Prepare() {}

func (b *BootController) Initdb() {
	var db models.DBconf
	err := utils.UnmarshalToJSON(b.Ctx.Request.Body, &db)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		b.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if err = service.InitDB(&db); err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusBadRequest, err.Error())
	}
	b.ServeJSON()
}

func (b *BootController) Startdb() {
	var host models.Account
	err := utils.UnmarshalToJSON(b.Ctx.Request.Body, &host)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		b.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if err = service.StartDB(&host); err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusBadRequest, err.Error())
	}
	b.ServeJSON()
}

func (b *BootController) CheckDB() {
	if err := dao.CheckDB(); err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusBadRequest, fmt.Sprintf("DB is down."))
	}
	b.ServeJSON()
}

// @Title CheckSysStatus
// @Description return the current system status.
// @Success 200 {object} models.InitSysStatus success
// @Failure 400 bad request
// @router /checksysstatus [get]
func (b *BootController) CheckSysStatus() {
	this, err := service.CheckSysStatus()
	if err != nil {
		logs.Error(err)
		b.CustomAbort(http.StatusBadRequest, err.Error())
	}
	b.Data["json"] = models.InitSysStatus{Status: this}
	b.ServeJSON()
}
