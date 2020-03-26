package controllers

import (
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
)

// CfgController includes operations about cfg
type CfgController struct {
	beego.Controller
}

// @Title Post
// @Description update cfg
// @Param	body	body	models.Configuration	true	"parameters"
// @Param	token	query 	string	true	"token"
// @Success 200 success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router / [post]
func (u *CfgController) Post() {
	var cfg models.Configuration
	var statusCode int = http.StatusOK
	token := u.GetString("token")
	result := service.VerifyToken(token)
	if result == false {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.ServeJSON()	
		return
	} else {
		//transferring JSON to struct.
		json.Unmarshal(u.Ctx.Input.RequestBody, &cfg)
		statusMessage := service.UpdateCfg(&cfg)
		if statusMessage == "BadRequest" {
			statusCode = http.StatusBadRequest
		}
		u.Ctx.ResponseWriter.WriteHeader(statusCode)
		u.ServeJSON()
	}
	
}

// @Title GetAll
// @Description return all cfg parameters
// @Param	which	query 	string	false	"which file to get"
// @Param	token	query 	string	true	"token"
// @Success 200 {object} models.Configuration	success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router / [get]
func (u *CfgController) GetAll() {
	var statusCode int = http.StatusOK
	token := u.GetString("token")
	result := service.VerifyToken(token)
	if result == false {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.ServeJSON()	
		return
	} else {
		which := u.GetString("which")
		cfg, statusMessage := service.GetAllCfg(which)
		if statusMessage == "BadRequest" {
			statusCode = http.StatusBadRequest
		}
		u.Ctx.ResponseWriter.WriteHeader(statusCode)
		//apply struct to JSON value.
		u.Data["json"] = cfg
		u.ServeJSON()
	}
}

// @Title GetKey
// @Description return public key
// @Param	token	query 	string	true	"token"
// @Success 200 {object} string	success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router /pubkey [get]
func (u *CfgController) GetKey() {
	token := u.GetString("token")
	result := service.VerifyToken(token)
	if result == false {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.ServeJSON()	
		return
	} else {
		pubkey := service.GetKey()
		u.Data["json"] = pubkey
		u.ServeJSON()
	}
	
}
