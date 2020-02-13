package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
        beego.ControllerComments{
            Method: "Applycfg",
            Router: `/applycfg`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
        beego.ControllerComments{
            Method: "Initialize",
            Router: `/initialize`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
        beego.ControllerComments{
            Method: "Restart",
            Router: `/restart`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
        beego.ControllerComments{
            Method: "Shutdown",
            Router: `/shutdown`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:AccController"],
        beego.ControllerComments{
            Method: "Verify",
            Router: `/verify`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:CfgController"],
        beego.ControllerComments{
            Method: "GetKey",
            Router: `/pubkey`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:MoniController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/adminserver/controllers:MoniController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
