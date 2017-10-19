package dao

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() {
	var err error
	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		logs.Error("Faild to load app.conf: %+v", err)
	}
	dbPassword := conf.String("dbPassword")
	logs.Info("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", "root:"+dbPassword+"@tcp(mysql:3306)/board?charset=utf8")
	if err != nil {
		logs.Error("error occurred on registering DB: %+v", err)
	}
}
