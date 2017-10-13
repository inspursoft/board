package service

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var user = model.User{
	Username: "Tester",
	Password: "123456a?",
}

func connectToDB() {
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(mysql:3306)/board?charset=utf8")
	if err != nil {
		logs.Error("Failed to connect to DB.")
	}
}

func cleanUp(username string) {
	o := orm.NewOrm()
	rs := o.Raw("delete from user where username = ?", username)
	r, err := rs.Exec()
	if err != nil {
		logs.Error("Error occurred while deleting user: %+v", err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		logs.Error("Error occurred while deleting user: %+v", err)
	}
	if affected == 0 {
		logs.Error("Failed to delete user")
	} else {
		logs.Error("Successful cleared up.")
	}
}

func TestMain(m *testing.M) {
	connectToDB()
	cleanUp(user.Username)
	os.Exit(m.Run())
}

func TestSignIn(t *testing.T) {
	u, err := SignIn("admin", "123456a?")
	if err != nil {
		t.Errorf("Error occurred while sign in: %+v\n", err)
	}
	if u == nil {
		t.Error("Sign in failed.")
	}
	if u.Username == "admin" {
		t.Log("Signed in successfully.")
	}
}

func TestGetUserByID(t *testing.T) {
	u, err := GetUserByID(1)
	if err != nil {
		t.Errorf("Error occurred while sign in: %+v\n", err)
	}
	if u == nil {
		t.Fatal("Sign in failed.")
	}
	if u.Username == "admin" {
		t.Log("Signed in successfully.")
	}
}

func TestSignUp(t *testing.T) {
	status, err := SignUp(user)
	if err != nil {
		t.Fatalf("Failed to sign up: %+v\n", err)
	}
	if status {
		t.Log("Successful signed up.")
	}
}
