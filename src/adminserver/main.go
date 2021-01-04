package main

import (
	"github.com/inspursoft/board/src/adminserver/dao"
	"github.com/inspursoft/board/src/adminserver/models"
	_ "github.com/inspursoft/board/src/adminserver/routers"

	"github.com/astaxie/beego"
)

func main() {
	//beego framework configuring and booting.
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	dao.InitDatabase()
	models.RegisterModels()
	dao.InitGlobalCache()
	beego.Run()
}
