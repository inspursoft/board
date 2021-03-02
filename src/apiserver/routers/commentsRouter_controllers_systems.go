package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/systems:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/systems:SupplementController"],
		beego.ControllerComments{
			Method:           "Info",
			Router:           `/info`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/systems:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/systems:SupplementController"],
		beego.ControllerComments{
			Method:           "KubernetesInfo",
			Router:           `/kubernetes-info`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/systems:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/systems:SupplementController"],
		beego.ControllerComments{
			Method:           "Resources",
			Router:           `/resources`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
