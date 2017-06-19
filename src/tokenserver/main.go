package main

import (
	"database/sql"

	"encoding/json"
	//"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"common/models"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//"golang.org/x/crypto/bcrypt"
)

// Declare a global variable to store the Redis connection pool.
var db *sql.DB
var authPort = "4000"
var secretKey = []byte("secret")
var mysqlPassword = ""

//type User struct {
//	Username string `json:"user_name"`
//	Password string `json:"user_password"`
//}

//Init just loads the default variables
func init() {
	authPort = os.Getenv("AUTH_PORT")
	secretKey = []byte(os.Getenv("SERECT_KEY"))
	mysqlPassword = os.Getenv("MYSQL_ROOT_PASSWORD")
}

func prepareDatabase(host, port, user, password, database string) (conn *sql.DB, err error) {

	// db, err := sql.Open("mysql", "root:root123@tcp(localhost:3306)/board?charset=utf8")
	connString := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8"
	log.Println(connString)

	conn, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Open database error:", err)
		return nil, err
	}
	err = conn.Ping()

	if err != nil {
		log.Fatal("Ping database error:", err)
		return nil, err
	}
	return conn, err
}

// gets you a token if you pass the right credentials
func login(w http.ResponseWriter, r *http.Request) {

	var err error
	var password string

	user := new(models.User)

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unmarshal error!!%v\n", err)
	}
	defer r.Body.Close()

	//log.Printf("user: %s, %s\n", user.Username, user.Password)

	db, err := prepareDatabase("mysql", "3306", "root", mysqlPassword, "board")
	if err != nil {
		log.Fatal("Database can't connect:", err)
	}
	defer db.Close()

	err = db.QueryRow("select password from user where username = ?", user.Username).Scan(&password)
	if err != nil {
		log.Fatalf("query error!!%v\n", err)
		return
	}

	if err == nil {
		// compare passwords
		// err = bcrypt.CompareHashAndPassword(password, []byte(r.FormValue("password")))
		// if it doesn't match

		if user.Password != password {
			http.Error(w, "Wrong password", 401)
			return
		}
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		//claims["admin"] = false
		claims["username"] = user.Username
		// 24 hour token
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString(secretKey)
		w.Write([]byte(tokenString))
	} else {
		// User not found
		http.Error(w, "User not found", 200)
		return
	}
}

// register a new user, gives you a token if the user -> password
// is not registered already
func register(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "POST" {
		http.Error(w, "Forbidden", 403)
		return
	}

	user := new(models.User)

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unmarshal error!!%v\n", err)
	}
	defer r.Body.Close()

	db, err := prepareDatabase("mysql", "3306", "root", mysqlPassword, "board")
	if err != nil {
		log.Fatal("Database can't connect:", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into user(username,password)values(?,?)")
	if err != nil {
		log.Fatalf("User insert prepare error!!%v\n", err)
		return
	}
	defer stmt.Close()

	//log.Printf("user: %s, %s\n", user.Username, user.Password)
	if result, err := stmt.Exec(user.Username, user.Password); err == nil {
		if id, err := result.LastInsertId(); err == nil {
			log.Println("insert id : ", id)
		} else {
			log.Fatalf("insert id err:%v\n", err)
		}
	} else {
		log.Fatalf("insert err:%v\n", err)
	}

	// check if the user is already registered
	//	exists, err := redis.Bool(conn.Do("EXISTS", email))
	//	if exists {
	//		w.Write([]byte("Email taken"))
	//		return
	//	}
	//	// get password from the post request form value
	//	password := r.FormValue("password")
	//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),
	//		bcrypt.DefaultCost)
	//	if err != nil {
	//		panic(err)
	//	}
	//	// Set user -> password in redis
	//	_, err = conn.Do("SET", email, string(hashedPassword[:]))
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	token := jwt.New(jwt.SigningMethodHS256)
	//	claims := token.Claims.(jwt.MapClaims)
	//	claims["admin"] = false
	//	claims["email"] = email
	//	// 24 hour token
	//	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	//	tokenString, _ := token.SignedString(SIGN_KEY)
	//	w.Write([]byte(tokenString))
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/login", login).Methods("POST")
	r.HandleFunc("/api/v1/register", register).Methods("POST")
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
