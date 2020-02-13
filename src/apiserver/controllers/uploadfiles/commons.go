package uploadfiles

import (
	"github.com/astaxie/beego"
)

// Operations about files
type CommonController struct {
	beego.Controller
}

// @Title List all upload files
// @Description List all for upload files.
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (c *CommonController) List() {

}

// @Title Upload file
// @Description Upload file.
// @Param	nodegroup_id	path	int	true	"ID of node group"
// @Param	upload	formData	file	true	"File uploaded"
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}

// @Title Download uploaded file
// @Description Download uploaded file.
// @Success 200 Successful downloaded.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [head]
func (c *CommonController) Head() {

}

// @Title Delete uploaded file
// @Description Delete uploaded file.
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [delete]
func (c *CommonController) Delete() {

}
