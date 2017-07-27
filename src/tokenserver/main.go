package main

import (
	_ "git/inspursoft/board/src/tokenserver/controller"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run(":4000")
}
