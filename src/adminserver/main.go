package main

import (
	_ "git/inspursoft/board/src/adminserver/routers"

	"github.com/astaxie/beego"
)

func main() {
	//beego framework configuring and booting.
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
