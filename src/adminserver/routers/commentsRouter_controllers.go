package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
		beego.ControllerComments{
			Method:           "CreateUUID",
			Router:           `/createUUID`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BoardController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BoardController"],
		beego.ControllerComments{
			Method:           "Applycfg",
			Router:           `/applycfg`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BoardController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BoardController"],
		beego.ControllerComments{
			Method:           "Shutdown",
			Router:           `/shutdown`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BoardController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BoardController"],
		beego.ControllerComments{
			Method:           "Start",
			Router:           `/start`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BootController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:BootController"],
		beego.ControllerComments{
			Method:           "CheckSysStatus",
			Router:           `/checksysstatus`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:MonitorController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:MonitorController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
