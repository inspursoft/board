package configmaps

import (
	"github.com/astaxie/beego"
)

// Operations about config maps
type CommonController struct {
	beego.Controller
}

// @Title List all config maps
// @Description List all for config maps.
// @Param	configmap_id	path	int	false	"ID of persistence"
// @Param	search	query	string	false	"Query item for config maps"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:configmap_id [get]
func (c *CommonController) List() {

}

// @Title Add config map
// @Description Add config map.
// @Param	body	body 	models.configmaps.vm.ConfigMap	true	"View model for config map."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}

// @Title Delete config map by ID
// @Description Delete config map by ID.
// @Param	configmap_id	path	int	true	"ID of config map"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:configmap_id [delete]
func (c *CommonController) Delete() {

}
