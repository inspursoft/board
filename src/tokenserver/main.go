package main

import (
	_ "github.com/inspursoft/board/src/tokenserver/controller"
	"github.com/inspursoft/board/src/tokenserver/service"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(4)

	err := service.InitService()
	if err != nil {
		logs.Error("Init token server config error: %+v", err)
		panic(err)
	}
	beego.Run(":4000")
}
