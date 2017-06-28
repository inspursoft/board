package main

import (
	"git/inspursoft/board/src/apiserver/controller"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", controller.ConfigRouters())
}
