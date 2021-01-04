package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/operations:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/operations:CommonController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/operations:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/operations:CommonController"],
        beego.ControllerComments{
            Method: "Add",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
