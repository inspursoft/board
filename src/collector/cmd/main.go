package main

import (
	"runtime"
	"time"
	"math/rand"
	"git/inspursoft/board/src/collector/cmd/app"
	_"git/inspursoft/board/src/collector/model/collect"
	_"git/inspursoft/board/src/collector/dao/registration"
	"net/http"
	"git/inspursoft/board/src/collector/control"
	"git/inspursoft/board/src/collector/util"
)

func main() {
	util.Logger.SetInfo("The cpu core is", runtime.NumCPU(), ",The app would use all of cores")
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	go app.Run()
	han:=func() http.Handler {
		s, _ := control.CollectRouters()
		return s
	}()
	err:=http.ListenAndServe(":8080", han)
	util.Logger.SetFatal(err)
}
