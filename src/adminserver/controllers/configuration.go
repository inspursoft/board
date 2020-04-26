package controllers

import (
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// CfgController includes operations about cfg
type CfgController struct {
	BaseController
}

// @Title Post
// @Description update cfg
// @Param	body	body	models.Configuration	true	"parameters"
// @Param	token	query 	string	true	"token"
// @Success 200 success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router / [post]
func (c *CfgController) Post() {
	var cfg models.Configuration

	//transferring JSON to struct.
	err := utils.UnmarshalToJSON(c.Ctx.Request.Body, &cfg)
	if err != nil {
		logs.Error("Failed to unmarshal data: %+v", err)
		c.CustomAbort(http.StatusBadRequest, err.Error())
	}
	err = service.UpdateCfg(&cfg)
	if err != nil {
		logs.Error(err)
		c.CustomAbort(http.StatusBadRequest, err.Error())
	}
	c.ServeJSON()

}

// @Title GetAll
// @Description return all cfg parameters
// @Param	which	query 	string	false	"which file to get"
// @Param	token	query 	string	true	"token"
// @Success 200 {object} models.Configuration	success
// @Failure 400 bad request
// @Failure 401 unauthorized
// @router / [get]
func (c *CfgController) GetAll() {

	which := c.GetString("which")
	cfg, err := service.GetAllCfg(which)
	if err != nil {
		logs.Error(err)
		c.CustomAbort(http.StatusBadRequest, err.Error())
	}
	//apply struct to JSON value.
	c.Data["json"] = cfg
	c.ServeJSON()

}

func (c *CfgController) GetKey() {
	pubkey := service.GetKey()
	c.Data["json"] = models.Key{Key: pubkey}
	c.ServeJSON()
}
