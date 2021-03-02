package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/emails:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/emails:SupplementController"],
		beego.ControllerComments{
			Method:           "Forgot",
			Router:           `/forgot`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/emails:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/emails:SupplementController"],
		beego.ControllerComments{
			Method:           "Notification",
			Router:           `/notification`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/emails:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/emails:SupplementController"],
		beego.ControllerComments{
			Method:           "Ping",
			Router:           `/ping`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
