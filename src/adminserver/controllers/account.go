package controllers

import (
	"git/inspursoft/board/src/common/utils"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"net/http"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// AccController includes operations about account.
type AccController struct {
	beego.Controller
}

// @Title Verify
// @Description verify input password
// @Param	body	body 	models.Password	true	"The password"
// @Success 200 {object} string success
// @Failure 400 bad request
// @router /verify [post]
func (a *AccController) Verify() {
	var passwd models.Password
	utils.UnmarshalToJSON(a.Ctx.Request.Body, &passwd)
	v, err := service.VerifyPassword(&passwd)
	if err != nil {
		logs.Error(err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if v {
		a.Data["json"] = "success"
	} else {
		a.Data["json"] = "wrong"
	}
	a.ServeJSON()
}

// @Title Initialize
// @Description initialize username and password
// @Param	body	body 	models.Account	true	"body for user account"
// @Success 200 success
// @Failure 400 bad request
// @router /initialize [post]
func (a *AccController) Initialize() {
	var acc models.Account
	//transferring JSON to struct.
	utils.UnmarshalToJSON(a.Ctx.Request.Body, &acc)
	err := service.Initialize(&acc)
	if err != nil {
		logs.Error(err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	a.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	body	body 	models.Account	true	"body for user account"
// @Success 200 {object} string success
// @Failure 400 bad request
// @Failure 403 forbidden
// @router /login [post]
func (a *AccController) Login() {
	var acc models.Account
	//transferring JSON to struct.
	utils.UnmarshalToJSON(a.Ctx.Request.Body, &acc)
	permission, err, token := service.Login(&acc)
	if err != nil {
		logs.Error(err)
		if err.Error() == "Forbidden" {
			a.CustomAbort(http.StatusForbidden, err.Error())
		}
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if permission {
		a.Data["json"] = token
	} else {
		a.CustomAbort(http.StatusBadRequest, "login failed")
	}
	a.ServeJSON()
}


// @Title Install
// @Description judge if it's the first time open admin server.
// @Success 200 {object} string success
// @Failure 400 bad request
// @router /install [get]
func (a *AccController) Install() {
	install := service.Install()
	if install == models.InitStatusTrue {
		a.Data["json"] = "yes"
	} else if install == models.InitStatusFirst{
		a.Data["json"] = "step1"
	} else if install == models.InitStatusSecond{
		a.Data["json"] = "step2"
	} else if install == models.InitStatusThird{
		a.Data["json"] = "step3"
	} else {
		a.Data["json"] = "no"
	}
	a.ServeJSON()
}

// @Title CreateUUID
// @Description create UUID
// @Success 200 success
// @Failure 400 bad request
// @router /createUUID [post]
func (a *AccController) CreateUUID() {
	err := service.CreateUUID()
	if err != nil {
		logs.Error(err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	a.ServeJSON()
}

// @Title ValidateUUID
// @Description validate the UUID
// @Param	body	body 	models.UUID	true	"UUID"
// @Success 200 {object} string success
// @Failure 400 bad request
// @router /ValidateUUID [post]
func (a *AccController) ValidateUUID() {
	var uuid models.UUID
	utils.UnmarshalToJSON(a.Ctx.Request.Body, &uuid)
	result, err := service.ValidateUUID(uuid.UUID)
	if err != nil {
		logs.Error(err)
		a.CustomAbort(http.StatusBadRequest, err.Error())
	}
	if result {
		a.Data["json"] = "validate success"
	} else {
		a.Data["json"] = "validate failure"
		a.CustomAbort(http.StatusBadRequest, "validate failure")
	}
	a.ServeJSON()
}
