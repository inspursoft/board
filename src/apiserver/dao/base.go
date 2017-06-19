package dao

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:root123@tcp(mysql:3306)/board?charset=utf8")
	if err != nil {
		fmt.Errorf("Error occurred on registering DB: %+v\n", err)
	}
}
