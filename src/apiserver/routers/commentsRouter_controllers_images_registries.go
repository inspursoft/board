package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images/registries:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images/registries:CommonController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/:project_id/:image_name`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images/registries:CommonController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images/registries:CommonController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:project_id/:image_name`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
