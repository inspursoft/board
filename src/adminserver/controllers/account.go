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
	permission, statusMessage := service.Login(&acc)
	if permission == true {
		a.Data["json"] = "login success"
	} else {
		a.Data["json"] = "login failure"
	}
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
	a.ServeJSON()
}

// @Title Restart
// @Description restart Board
// @Param	token	query 	string	true		"token"
// @Success 200 success
// @Failure 400 bad request
// @router /restart [get]
func (a *AccController) Restart() {
	var statusCode int = http.StatusOK
	statusMessage := service.Restart("/root/BOARD/Deploy")
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
	a.ServeJSON()
}

// @Title Applycfg
// @Description apply cfg and restart Board
// @Param	token	query 	string	true	"token"
// @Success 200 success
// @Failure 400 bad request
// @router /applycfg [get]
func (a *AccController) Applycfg() {
	var statusCode int = http.StatusOK
	statusMessage := service.Applycfg("/root/BOARD/Deploy")
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
	a.ServeJSON()
}

// @Title Shutdown
// @Description shutdown board
// @Param	token	query 	string	true	"token"
// @Success 200 success
// @Failure 400 bad request
// @router /shutdown [get]
func (a *AccController) Shutdown() {
	var statusCode int = http.StatusOK
	statusMessage := service.Shutdown("/root/BOARD/Deploy")
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	a.Ctx.ResponseWriter.WriteHeader(statusCode)
	a.ServeJSON()
}
