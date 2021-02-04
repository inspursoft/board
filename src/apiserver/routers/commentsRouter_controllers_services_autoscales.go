package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/:project_id/:service_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"],
		beego.ControllerComments{
			Method:           "Add",
			Router:           `/:project_id/:service_id`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           `/:project_id/:service_id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/services/autoscales:CommonController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:project_id/:service_id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
