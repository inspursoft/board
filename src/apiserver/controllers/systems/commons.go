package systems

import "github.com/astaxie/beego"

//Operations about system info
type CommonController struct {
	beego.Controller
}

// @Title Get system information.
// @Description Get system information.
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /info [get]
func (i *CommonController) Info() {

}

// @Title Get system resources information.
// @Description Get system resources information.
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /resources [get]
func (i *CommonController) Resources() {

}
