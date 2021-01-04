package controllers

import (
	"github.com/inspursoft/board/src/adminserver/models"
	"github.com/inspursoft/board/src/adminserver/service"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// CfgController includes operations about cfg
type CfgController struct {
	BaseController
}

// @Title Put
// @Description update cfg
// @Param	body	body	models.Configuration	true	"parameters"
// @Param	token	query 	string	false	"token"
// @Success 200 success
// @Failure 500 Internal Server Error
// @Failure 401 unauthorized: token invalid/session timeout
// @router / [put]
func (c *CfgController) Put() {
	var cfg models.Configuration

	//transferring JSON to struct.
	err := utils.UnmarshalToJSON(c.Ctx.Request.Body, &cfg)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	err = service.UpdateCfg(&cfg)
	if err != nil {
		logs.Error(err)
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	c.ServeJSON()

}

// @Title GetAll
// @Description return all cfg parameters
// @Param	which	query 	string	false	"which file to get"
// @Param	token	query 	string	false	"token"
// @Success 200 {object} models.Configuration	success
// @Failure 500 Internal Server Error
// @Failure 401 unauthorized: token invalid/session timeout
// @router / [get]
func (c *CfgController) GetAll() {

	which := c.GetString("which")
	cfg, err := service.GetAllCfg(which, false)
	if err != nil {
		logs.Error(err)
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	//apply struct to JSON value.
	c.Data["json"] = cfg
	c.ServeJSON()

}
