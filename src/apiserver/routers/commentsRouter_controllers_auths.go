package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"],
		beego.ControllerComments{
			Method:           "SignIn",
			Router:           `/sign-in`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"],
		beego.ControllerComments{
			Method:           "SignOut",
			Router:           `/sign-out`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"],
		beego.ControllerComments{
			Method:           "SignUp",
			Router:           `/sign-up`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:AuthController"],
		beego.ControllerComments{
			Method:           "ThirdParty",
			Router:           `/third-party`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:PasswordController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/auths:PasswordController"],
		beego.ControllerComments{
			Method:           "Reset",
			Router:           `/reset`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
