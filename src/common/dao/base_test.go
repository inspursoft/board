package dao

import (
	"os"
	"testing"
	"fmt"

	"github.com/astaxie/beego/orm"
)

func TestMain(m *testing.M) {
	hostIP:=os.Getenv("HOST_IP")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:root123@tcp(%s:3306)/board?charset=utf8", hostIP))
	os.Exit(m.Run())
}
