package nfs

import (
	"github.com/astaxie/beego"
)

// Operations about storage NFS
type CommonController struct {
	beego.Controller
}

// @Title Add storage NFS
// @Description Add storage NFS.
// @Param	body	body 	models.storages.vm.NFS	true	"View model for storage NFS."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}
