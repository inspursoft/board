package dao

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

func init() {
	var err error
	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		logs.Error("Faild to load app.conf: %+v", err)
	}
	dbPassword := conf.String("dbPassword")
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", "root:"+dbPassword+"@tcp(mysql:3306)/board?charset=utf8")
	if err != nil {
		fmt.Printf("error occurred on registering DB: %+v", err)
	}
}
