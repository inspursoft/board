package dao

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(10.165.33.196:3306)/board?charset=utf8")
	if err != nil {
		fmt.Errorf("error occurred on registering DB: %+v\n", err)
	}
}
