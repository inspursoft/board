package dao

import (
	"os"
	"testing"

	"github.com/astaxie/beego/orm"
)

func TestMain(m *testing.M) {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root123@tcp(localhost:3306)/board?charset=utf8")
	os.Exit(m.Run())
}
