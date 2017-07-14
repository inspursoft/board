package controller

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const root = "/"
const apiprefix = "/api/v1"

func ConfigRouters() *mux.Router {
	var r *mux.Router = mux.NewRouter()
	s := r.PathPrefix(apiprefix).Subrouter()
	s.HandleFunc("/sign-in", SignInAction).Methods("POST")
	s.HandleFunc("/sign-up", SignUpAction).Methods("POST")
	s.HandleFunc("/adduser", AddUserAction).Methods("POST")
	s.HandleFunc("/users", GetUsersAction).Methods("GET")
	s.HandleFunc("/users/{id:[0-9]+}", OperateUserAction).Methods("GET", "PUT", "DELETE")
	s.HandleFunc("/users/{id:[0-9]+}/password", ChangePasswordAction).Methods("PUT")
	s.HandleFunc("/users/{id:[0-9]+}/systemadmin", ToggleSystemAdminAction).Methods("PUT")
	s.HandleFunc("/projects", ListAndCreateProjectAction).Methods("GET", "POST")
	s.HandleFunc("/projects/{id:[0-9]+}", GetAndDeleteProjectAction).Methods("GET", "DELETE")
	s.HandleFunc("/projects/{id:[0-9]+}/publicity", ToggleProjectPublicAction).Methods("PUT")
	s.HandleFunc("/projects/{id:[0-9]+}/members", OperateProjectMembersAction).Methods("GET", "POST", "DELETE")

	if _, err := os.Stat("./swagger"); err == nil {
		r.PathPrefix(root).Handler(http.FileServer(http.Dir("./swagger")))
	}
	return r
}
