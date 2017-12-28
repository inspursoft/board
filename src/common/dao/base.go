package dao

import (
	"fmt"

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
	dbHost := conf.String("dbHost")
	dbPassword := conf.String("dbPassword")

	logs.Info("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:%s@tcp(%s:3306)/board?charset=utf8", dbPassword, dbHost))
	if err != nil {
		logs.Error("error occurred on registering DB: %+v", err)
		panic(err)
	}
}

func getTotalRecordCount(baseSQL string, params []interface{}) (int64, error) {
	o := orm.NewOrm()
	var count int64
	err := o.Raw(`select count(*) from (`+baseSQL+`) t`, params).QueryRow(&count)
	if err != nil {
		logs.Error("failed to retrieve for total count: %+v", err)
		return 0, err
	}
	return count, nil
}
