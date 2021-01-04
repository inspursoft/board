package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:CommonController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:PasswordController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:PasswordController"],
        beego.ControllerComments{
            Method: "UpdatePassword",
            Router: `/password`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:ProbeController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:ProbeController"],
        beego.ControllerComments{
            Method: "Current",
            Router: `/current`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/users:SupplementController"],
        beego.ControllerComments{
            Method: "Exists",
            Router: `/existing`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
