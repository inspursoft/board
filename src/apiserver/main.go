package main

import (
	_ "git/inspursoft/board/src/apiserver/router"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run(":8088")
}
