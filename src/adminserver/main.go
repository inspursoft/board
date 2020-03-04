package main

import (
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/models"
	_ "git/inspursoft/board/src/adminserver/routers"
	"github.com/astaxie/beego/logs"
	"os"
	"github.com/astaxie/beego"
)

func main() {
	//beego framework configuring and booting.
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	if _, err := os.Stat(dao.AdminServerDbFile); os.IsNotExist(err) {
		if errInitDb := dao.InitDatabase(); errInitDb != nil {
			logs.Error(errInitDb)
		}
	}
	if err := dao.RegisterDatabase(); err != nil {
		logs.Error(err)
	}
	models.RegisterModels();
	dao.InitGlobalCache()
	beego.Run()
}
