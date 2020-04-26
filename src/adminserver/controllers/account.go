package controllers

import (
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// AccController includes operations about account.
type AccController struct {
	BaseController
}

func (a *AccController) Prepare() {}

func (a *AccController) Verify() {
	var passwd models.Password
	err := utils.UnmarshalToJSON(a.Ctx.Request.Body, &passwd)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	v, err := service.VerifyPassword(&passwd)
	if err != nil {
		logs.Error(err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if !v {
		a.CustomAbort(http.StatusBadRequest, "wrong password")
	}
	a.ServeJSON()
}

func (a *AccController) Initialize() {
	var acc models.Account
	err := utils.UnmarshalToJSON(a.Ctx.Request.Body, &acc)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	err = service.Initialize(&acc)
	if err != nil {
		logs.Error(err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	a.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	body	body 	models.Account	true	"body for user account"
// @Success 200 {object} models.TokenString success
// @Failure 400 bad request
// @Failure 403 forbidden
// @router /login [post]
func (a *AccController) Login() {
	var acc models.Account
	err := utils.UnmarshalToJSON(a.Ctx.Request.Body, &acc)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	permission, err, token := service.Login(&acc)
	if err != nil {
		logs.Error(err)
		if err.Error() == "Forbidden" {
			a.CustomAbort(http.StatusForbidden, err.Error())
		}
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if permission {
		a.Data["json"] = models.TokenString{TokenString: token}
	} else {
		a.CustomAbort(http.StatusBadRequest, "login failed")
	}
	a.ServeJSON()
}

// @Title CreateUUID
// @Description create UUID
// @Success 200 success
// @Success 202 accepted
// @Failure 400 bad request
// @router /createUUID [post]
func (a *AccController) CreateUUID() {
	err := service.CreateUUID()
	if err != nil {
		logs.Error(err)
		if err.Error() == "another admin user has signed in other place" {
			a.Ctx.ResponseWriter.WriteHeader(http.StatusAccepted)
		}
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	a.ServeJSON()
}

// @Title ValidateUUID
// @Description validate the UUID
// @Param	body	body 	models.UUID	true	"UUID"
// @Success 200 success
// @Failure 400 bad request
// @router /ValidateUUID [post]
func (a *AccController) ValidateUUID() {
	var uuid models.UUID
	err := utils.UnmarshalToJSON(a.Ctx.Request.Body, &uuid)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	result, err := service.ValidateUUID(uuid.UUIDstring)
	if err != nil {
		logs.Error(err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if !result {
		a.CustomAbort(http.StatusBadRequest, "validate failure")
	}
	a.ServeJSON()
}
