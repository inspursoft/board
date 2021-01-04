package dao

import (
	"fmt"
	"github.com/inspursoft/board/src/common/utils"

	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() {

	dbIP := utils.GetStringValue("DB_IP")
	dbPort := utils.GetIntValue("DB_PORT")
	dbPassword := utils.GetStringValue("DB_PASSWORD")

	//init models
	model.InitModelDB()
	logs.Info("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:%s@tcp(%s:%d)/board?charset=utf8&loc=Local", dbPassword, dbIP, dbPort))
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

func getOrderSQL(orderField string, orderAsc int) string {
	if orderAsc != 0 {
		return fmt.Sprintf(` order by %s`, orderField)
	}
	return fmt.Sprintf(` order by %s desc`, orderField)
}
