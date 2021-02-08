package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "GetNodeListAction",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "RemoveNodeAction",
			Router:           `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "AddNodeAction",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "ControlStatusAction",
			Router:           `/:node_name`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "CallBackAction",
			Router:           `/callback`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "DeleteNodeLog",
			Router:           `/log`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "GetNodeLogDetail",
			Router:           `/log`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "GetNodeLogList",
			Router:           `/logs`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers/nodeController:NodeController"],
		beego.ControllerComments{
			Method:           "PreparationAction",
			Router:           `/preparation`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
