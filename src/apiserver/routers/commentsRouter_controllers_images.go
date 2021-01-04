package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:BuildController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:BuildController"],
        beego.ControllerComments{
            Method: "ByDockerfile",
            Router: `/by_dockerfile`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:BuildController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:BuildController"],
        beego.ControllerComments{
            Method: "ByTemplate",
            Router: `/by_template`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:BuildController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:BuildController"],
        beego.ControllerComments{
            Method: "ByUploadedPackage",
            Router: `/by_uploaded_package`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:SupplementController"],
        beego.ControllerComments{
            Method: "CleanConfig",
            Router: `/:project_id/clean_config`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:SupplementController"],
        beego.ControllerComments{
            Method: "Exists",
            Router: `/:project_id/exists`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:SupplementController"] = append(beego.GlobalControllerRouter["github.com/inspursoft/board/src/apiserver/controllers/images:SupplementController"],
        beego.ControllerComments{
            Method: "Preview",
            Router: `/:project_id/preview`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
