package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"],
		beego.ControllerComments{
			Method:           "Add",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/:project_id/:service_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           `/:project_id/:service_id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:CommonController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:project_id/:service_id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"],
		beego.ControllerComments{
			Method:           "Info",
			Router:           `/:project_id/:service_id/info`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"],
		beego.ControllerComments{
			Method:           "Publicity",
			Router:           `/:project_id/:service_id/publicity`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"],
		beego.ControllerComments{
			Method:           "Route",
			Router:           `/:project_id/:service_id/route`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"],
		beego.ControllerComments{
			Method:           "Seletable",
			Router:           `/:project_id/:service_id/selectable`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"],
		beego.ControllerComments{
			Method:           "Status",
			Router:           `/:project_id/:service_id/status`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"],
		beego.ControllerComments{
			Method:           "Test",
			Router:           `/:project_id/:service_id/test`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services:SupplementController"],
		beego.ControllerComments{
			Method:           "Toggle",
			Router:           `/:project_id/:service_id/toggle`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
