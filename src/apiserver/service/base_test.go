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

func TestMain(m *testing.M) {
	connectToDB()
	os.Exit(m.Run())
}
