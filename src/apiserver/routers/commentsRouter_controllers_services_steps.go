package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/steps:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/steps:CommonController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:project_id/:service_id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/steps:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/steps:CommonController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/:project_id/:service_id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/steps:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/steps:CommonController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:project_id/:service_id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
