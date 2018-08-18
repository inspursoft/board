package dao

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	projectTable                = "PROJECT"
	userTable                   = "USER"
	serviceTable                = "SERVICE"
	nameField                   = "NAME"
	createTimeField             = "CREATE_TIME"
	defaultField                = "CREATE_TIME"
	projectTableNameField       = "PROJECT_NAME"
	userTableNameField          = "USER_NAME"
	serviceTableNameField       = "SERVICE_NAME"
	projectTableCreateTimeField = "PROJECT_CREATE_TIME"
	userTableCreateTimeField    = "USER_CREATE_TIME"
	serviceTableCreateTimeField = "SERVICE_CREATE_TIME"
)

var orderFields = map[string]string{
	projectTableNameField:       "name",
	userTableNameField:          "username",
	serviceTableNameField:       "name",
	projectTableCreateTimeField: "creation_time",
	userTableCreateTimeField:    "creation_time",
	serviceTableCreateTimeField: "creation_time",
}

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

func getOrderSQL(orderTable string, orderField string, orderAsc int) string {
	key := fmt.Sprintf("%s_%s", strings.ToUpper(orderTable), strings.ToUpper(orderField))
	if orderFields[key] == "" {
		return fmt.Sprintf(` order by %s desc`, orderFields[defaultField])
	}
	if orderAsc != 0 {
		return fmt.Sprintf(` order by %s`, orderFields[key])
	}
	return fmt.Sprintf(` order by %s desc`, orderFields[key])
}
