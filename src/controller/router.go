package controller

import "github.com/gorilla/mux"

func ConfigRouters() *mux.Router {
	var r *mux.Router = mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/sign-in", SignInAction)
	s.HandleFunc("/sign-up", SignUpAction)

	s.HandleFunc("/users", GetUsersAction)
	s.HandleFunc("/users/{id:[0-9]+}", OperateUserAction).
		Methods("PUT", "DELETE")
	return r
}
