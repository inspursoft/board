package claims

import (
	"github.com/astaxie/beego"
)

// Operations about persistent volume claims
type CommonController struct {
	beego.Controller
}

// @Title List all persistent volume claims
// @Description List all for persistent volume claims.
// @Param	persistence_id	path	int	true	"ID of persistent volume"
// @Param	claim_id	path	int	true	"ID of persistent volume claim"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:persistence_id/:claim_id [get]
func (c *CommonController) List() {

}

// @Title Creat persistent volume claim
// @Description Create persistent volume claim.
// @Param	persistence_id	path	int	true	"ID of persistent volume"
// @Param	body	body 	models.persistences.vm.PersistenceVolumeClaim	true	"View model for persistence volume claim."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:persistence_id [post]
func (c *CommonController) Add() {

}

// @Title Delete persistence volume by ID
// @Description Delete persistent volume by ID.
// @Param	persistence_id	path	int	true	"ID of persistent volume"
// @Param	claim_id	path	int	true	"ID of persistent volume claim"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:persistence_id/:claim_id [delete]
func (c *CommonController) Delete() {

}
