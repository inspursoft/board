package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/configmaps:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/configmaps:CommonController"],
		beego.ControllerComments{
			Method:           "Add",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/configmaps:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/configmaps:CommonController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/:configmap_id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/configmaps:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/configmaps:CommonController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:configmap_id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
