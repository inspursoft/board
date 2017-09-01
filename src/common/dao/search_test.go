package dao

import (
	"fmt"

	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func TestSearchPrivite(t *testing.T) {
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(10.110.18.107:3306)/board?charset=utf8")
	orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Error occurred on registering DB: %+v\n", err)
	}
	fmt.Println(SearchPrivite("l", "Admin"))
}
