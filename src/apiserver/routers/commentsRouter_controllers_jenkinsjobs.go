package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/jenkinsjobs:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/jenkinsjobs:SupplementController"],
        beego.ControllerComments{
            Method: "JobNumber",
            Router: `/:user_id/:build_number`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/jenkinsjobs:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/jenkinsjobs:SupplementController"],
        beego.ControllerComments{
            Method: "Console",
            Router: `/console`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/jenkinsjobs:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/jenkinsjobs:SupplementController"],
        beego.ControllerComments{
            Method: "Stop",
            Router: `/stop`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
