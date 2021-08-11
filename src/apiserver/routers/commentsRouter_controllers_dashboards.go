package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"],
		beego.ControllerComments{
			Method:           "Data",
			Router:           `/data`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"],
		beego.ControllerComments{
			Method:           "Node",
			Router:           `/nodes`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"],
		beego.ControllerComments{
			Method:           "ServerTime",
			Router:           `/server_time`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/dashboards:DataController"],
		beego.ControllerComments{
			Method:           "Service",
			Router:           `/services`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
