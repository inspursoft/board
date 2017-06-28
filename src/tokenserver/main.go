package main

import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"git/inspursoft/board/src/common/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var authPort = "4000"
var secretKey = []byte("secret")

//Init just loads the default variables
func init() {
	authPort = os.Getenv("AUTH_PORT")
	secretKey = []byte(os.Getenv("SERECT_KEY"))
}

// gets you a token
func token(w http.ResponseWriter, r *http.Request) {

	var err error

	user := new(model.User)

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unmarshal error!!%v\n", err)
	}
	defer r.Body.Close()

	//log.Printf("user: %s, %s\n", user.Username, user.Password)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	//claims["admin"] = false
	claims["username"] = user.Username
	// 24 hour token
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(secretKey)
	w.Write([]byte(tokenString))

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/token", token).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + authPort,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
