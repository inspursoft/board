package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:CommonController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:CommonController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:job_id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:CommonController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:CommonController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/deploy`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"],
		beego.ControllerComments{
			Method:           "PodLogs",
			Router:           `/:job_id/logs/:pod_name`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"],
		beego.ControllerComments{
			Method:           "Pods",
			Router:           `/:job_id/pods`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"],
		beego.ControllerComments{
			Method:           "Status",
			Router:           `/:job_id/status`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"],
		beego.ControllerComments{
			Method:           "Existing",
			Router:           `/existing`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"] = append(beego.GlobalControllerRouter["git/inspursoft/board/src/apiserver/controllers/jobs:SupplementController"],
		beego.ControllerComments{
			Method:           "Selectable",
			Router:           `/selectable`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
