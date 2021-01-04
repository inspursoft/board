package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:ImageController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:ImageController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/:project_id/:service_id/images`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:ImageController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:ImageController"],
        beego.ControllerComments{
            Method: "Patch",
            Router: `/:project_id/:service_id/images`,
            AllowHTTPMethods: []string{"patch"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:NodeGroupController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:NodeGroupController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/:project_id/:service_id/nodegroups`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:NodeGroupController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/services/rollingupdates:NodeGroupController"],
        beego.ControllerComments{
            Method: "Patch",
            Router: `/:project_id/:service_id/nodegroups`,
            AllowHTTPMethods: []string{"patch"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
