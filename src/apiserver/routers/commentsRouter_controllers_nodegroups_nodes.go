package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:CommonController"],
        beego.ControllerComments{
            Method: "Add",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:CommonController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/:nodegroup_id/:node_id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:CommonController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:nodegroup_id/:node_id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/nodegroups/nodes:SupplementController"],
        beego.ControllerComments{
            Method: "Toggle",
            Router: `/:nodegroup_id/:node_id/toggle`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
