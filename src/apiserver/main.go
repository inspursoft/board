package main

import (
	"apiserver/controller"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", controller.ConfigRouters())
}
