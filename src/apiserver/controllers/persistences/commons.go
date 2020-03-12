package persistences

import (
	"github.com/astaxie/beego"
)

// Operations about persistent volumes
type CommonController struct {
	beego.Controller
}

// @Title List all persistent volumes
// @Description List all for persistent volumes.
// @Param	persistence_id	path	int	false	"ID of persistence"
// @Param	search	query	string	false	"Query item for persistences"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:persistence_id [get]
func (c *CommonController) List() {

}

// @Title Add persistent volume
// @Description Add persistent volume.
// @Param	body	body 	models.persistences.vm.PersistenceVolume	true	"View model for persistence volume."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}

// @Title Update persistent volume by ID
// @Description Update persistent volume by ID.
// @Param	persistence_id	path	int	true	"ID of persistent volume"
// @Param	body	body	models.persistences.vm.PersistenceVolume	true	"View model for persistence volume."
// @Param	action	query	string	true	"Option of update."
// @Success 200 Successful updated.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:persistence_id [put]
func (c *CommonController) Update() {

}

// @Title Delete persistence volume by ID
// @Description Delete persistent volume by ID.
// @Param	persistence_id	path	int	true	"ID of persistent volume"
// @Success 200 Successful deleted.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:persistence_id [delete]
func (c *CommonController) Delete() {

}
