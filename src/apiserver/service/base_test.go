package service

import (
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func connectToDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(localhost:3306)/board?charset=utf8")
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
