// @APIVersion 1.0.0
// @Title Admin server API
// @Description Admin server API
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"git/inspursoft/board/src/adminserver/controllers"
	"git/inspursoft/board/src/adminserver/controllers/node"
	"github.com/astaxie/beego"
)

func init() {

	//registering a router address to the indicated controller that will handle requests for this url.
	ns := beego.NewNamespace("/v1/admin",
		beego.NSNamespace("/configuration",
			beego.NSInclude(
				&controllers.CfgController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&controllers.AccController{},
			),
		),
		beego.NSNamespace("/monitor",
			beego.NSInclude(
				&controllers.MoniController{},
			),
		),
		beego.NSNamespace("/node",
			beego.NSRouter("/add", &node.Controller{}, "get:AddNodeAction"),
			beego.NSRouter("/delete", &node.Controller{}, "delete:DeleteNodeAction"),
			beego.NSRouter("/list", &node.Controller{}, "get:GetNodeListAction"),
		),
	)
	beego.AddNamespace(ns)
}
