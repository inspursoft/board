package controllers

import (
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"fmt"
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
	token := u.GetString("token")
	result := service.VerifyToken(token)
	if !result {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.ServeJSON()	
	} else {
		//transferring JSON to struct.
		utils.UnmarshalToJSON(u.Ctx.Request.Body, &cfg)
		err := service.UpdateCfg(&cfg)
		if err != nil {
			logs.Error(err)
			u.CustomAbort(http.StatusBadRequest, err.Error())
		}
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
	token := u.GetString("token")
	result := service.VerifyToken(token)
	if !result {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.ServeJSON()	
	} else {
		which := u.GetString("which")
		cfg, statusMessage := service.GetAllCfg(which)
		if statusMessage == "BadRequest" {
			u.CustomAbort(http.StatusBadRequest, fmt.Sprintf("Get config failed."))
		}
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
	if !result {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.ServeJSON()	
	} else {
		pubkey := service.GetKey()
		u.Data["json"] = pubkey
		u.ServeJSON()
	}
	
}
