package controllers

import (
	"encoding/json"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"net/http"

	"github.com/astaxie/beego"
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
	var statusCode int = http.StatusOK
	var passwd models.Password
	json.Unmarshal(a.Ctx.Input.RequestBody, &passwd)

	v, statusMessage := service.VerifyPassword(&passwd)
	if v == true {
		a.Data["json"] = "success"
	} else {
		a.Data["json"] = "wrong"
	}

	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
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
	var statusCode int = http.StatusOK
	//transferring JSON to struct.
	json.Unmarshal(a.Ctx.Input.RequestBody, &acc)
	statusMessage := service.Initialize(&acc)
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
	a.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	body	body 	models.Account	true	"body for user account"
// @Success 200 {object} string success
// @Failure 400 bad request
// @router /login [post]
func (a *AccController) Login() {
	var acc models.Account
	var statusCode int = http.StatusOK
	//transferring JSON to struct.
	json.Unmarshal(a.Ctx.Input.RequestBody, &acc)
	permission, statusMessage, token := service.Login(&acc)
	if permission == true {
		a.Data["json"] = token
	} else {
		a.Data["json"] = ""
		statusCode = http.StatusBadRequest
	}
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
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
	var statusCode int = http.StatusOK
	statusMessage := service.CreateUUID()
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
	a.ServeJSON()
}

// @Title ValidateUUID
// @Description validate the UUID
// @Param	body	body 	models.UUID	true	"UUID"
// @Success 200 {object} string success
// @Failure 400 bad request
// @router /ValidateUUID [post]
func (a *AccController) ValidateUUID() {
	var statusCode int = http.StatusOK
	var uuid models.UUID
	json.Unmarshal(a.Ctx.Input.RequestBody, &uuid)
	result, statusMessage := service.ValidateUUID(uuid.UUID)
	if result == true {
		a.Data["json"] = "validate success"
	} else {
		a.Data["json"] = "validate failure"
		statusCode = http.StatusBadRequest
	}
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
	a.ServeJSON()
}
