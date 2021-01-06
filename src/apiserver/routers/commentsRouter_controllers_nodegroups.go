package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:CommonController"],
		beego.ControllerComments{
			Method:           "Add",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:CommonController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/:nodegroup_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:CommonController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:nodegroup_id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/nodegroups:SupplementController"],
		beego.ControllerComments{
			Method:           "Toggle",
			Router:           `/:nodegroup_id/existing`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
