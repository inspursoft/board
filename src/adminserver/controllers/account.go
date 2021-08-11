package controllers

import (
	"git/inspursoft/board/src/adminserver/common"
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

// @Title Login
// @Description Logs user into the system
// @Param	body	body 	models.Account	true	"body for user account"
// @Success 200 {object} models.TokenString success
// @Failure 500 Internal Server Error
// @Failure 403 forbidden
// @router /login [post]
func (a *AccController) Login() {
	var acc models.Account
	var permission bool
	var token string
	var err error

	err = utils.UnmarshalToJSON(a.Ctx.Request.Body, &acc)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		a.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	if acc.Username != "boardadmin" {
		a.CustomAbort(http.StatusForbidden, common.ErrForbidden.Error())
	}

	permission, token, err = service.Login(&acc)
	if err != nil {
		logs.Error(err)
		if err == common.ErrForbidden {
			a.CustomAbort(http.StatusForbidden, err.Error())
		}
		a.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	if !permission {
		a.CustomAbort(http.StatusInternalServerError, "login failed")
	} else {
		a.Data["json"] = models.TokenString{TokenString: token}
	}
	a.ServeJSON()
}

// @Title CreateUUID
// @Description create UUID
// @Success 200 success
// @Success 202 accepted
// @Failure 500 Internal Server Error
// @router /createUUID [post]
func (a *AccController) CreateUUID() {
	err := service.CreateUUID()
	if err != nil {
		logs.Error(err)
		if err == common.ErrAdminLogin {
			a.Ctx.ResponseWriter.WriteHeader(http.StatusAccepted)
			return
		}
		a.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	a.ServeJSON()
}
