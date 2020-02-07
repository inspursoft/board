package controllers

import (
	"board-adminserver/src/backend/models"
	"board-adminserver/src/backend/service"
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
// @Success 200 success
// @Failure 400 bad request
// @router / [post]
func (u *CfgController) Post() {
	var cfg models.Configuration
	var statusCode int = http.StatusOK
	//transferring JSON to struct.
	json.Unmarshal(u.Ctx.Input.RequestBody, &cfg)
	statusMessage := service.UpdateCfg(&cfg)
	if statusMessage == "BadRequest" {
		statusCode = http.StatusBadRequest
	}
	u.Ctx.ResponseWriter.WriteHeader(statusCode)
	u.ServeJSON()
}

// @Title GetAll
// @Description return all cfg parameters
// @Param	which	query 	string	false	"which file to get"
// @Success 200 {object} models.Configuration	success
// @Failure 400 bad request
// @router / [get]
func (u *CfgController) GetAll() {
	var statusCode int = http.StatusOK
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

// @Title GetKey
// @Description return public key
// @Success 200 {object} string	success
// @Failure 400 bad request
// @router /pubkey [get]
func (u *CfgController) GetKey() {
	pubkey := service.GetKey()
	u.Data["json"] = pubkey
	u.ServeJSON()
}
