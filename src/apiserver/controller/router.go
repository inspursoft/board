package controller

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const root = "/"
const apiprefix = "/api/v1/"

func ConfigRouters() *mux.Router {
	var r *mux.Router = mux.NewRouter()
	s := r.PathPrefix(apiprefix).Subrouter()
	s.HandleFunc("/sign-in", SignInAction).Methods("POST")
	s.HandleFunc("/sign-up", SignUpAction).Methods("POST")
	s.HandleFunc("/users", GetUsersAction).Methods("GET")
	s.HandleFunc("/users/{id:[0-9]+}", OperateUserAction).Methods("PUT", "DELETE")
	if _, err := os.Stat("./swagger"); err == nil {
		r.PathPrefix(root).Handler(http.FileServer(http.Dir("./swagger")))
	}
	return r
}
