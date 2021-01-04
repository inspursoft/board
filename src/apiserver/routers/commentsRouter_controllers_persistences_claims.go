package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/persistences/claims:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/persistences/claims:CommonController"],
        beego.ControllerComments{
            Method: "Add",
            Router: `/:persistence_id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/persistences/claims:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/persistences/claims:CommonController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/:persistence_id/:claim_id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/persistences/claims:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/persistences/claims:CommonController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:persistence_id/:claim_id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
