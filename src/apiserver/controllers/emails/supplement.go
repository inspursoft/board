package emails

import (
	"github.com/astaxie/beego"
)

// Operations about Email service supplements
type SupplementController struct {
	beego.Controller
}

// @Title Ping target Email service
// @Description Ping target Email service
// @Param	email	query	string	true	"Email address"
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /ping [post]
func (c *SupplementController) Ping() {

}

// @Title Notify by Email service
// @Description Notify by Email service
// @Param	email	query	string	true	"Email address"
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /notification [post]
func (c *SupplementController) Notification() {

}
