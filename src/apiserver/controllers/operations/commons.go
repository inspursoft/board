package operations

import "github.com/astaxie/beego"

// Operations about recording operations
type CommonController struct {
	beego.Controller
}

// @Title List all operations
// @Description List all operations.
// @Param	search	query	string	false	"Query item for operations"
// @Success 200 Successful listed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (c *CommonController) List() {

}

// @Title Record operation
// @Description Record operation
// @Param	body	body 	models.operation.vm.Operation	true	"View model for operation."
// @Success 200 Successful added.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [post]
func (c *CommonController) Add() {

}
