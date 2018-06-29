package dao

import (
	"fmt"
	"git/inspursoft/board/src/common/utils"
	"strings"

	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	projectTable                = "PROJECT"
	userTable			= "USER"
	operationTable			= "OPERATION"
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
	operationTableCreateTimeField= "OPERATION_CREATION_TIME"
)

var orderFields = map[string]string{
	projectTableNameField:       "name",
	userTableNameField:          "username",
	serviceTableNameField:       "name",
	projectTableCreateTimeField: "creation_time",
	userTableCreateTimeField:    "creation_time",
	serviceTableCreateTimeField: "creation_time",
	operationTableCreateTimeField:"creation_time",
}

func InitDB() {

	dbIP := utils.GetStringValue("DB_IP")
	dbPort := utils.GetIntValue("DB_PORT")
	dbPassword := utils.GetStringValue("DB_PASSWORD")

	//init models
	model.InitModelDB()
	logs.Info("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:%s@tcp(%s:%d)/board?charset=utf8", dbPassword, dbIP, dbPort))
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
		logs.Info(key)
		return ""
	}
	if orderAsc != 0 {
		return fmt.Sprintf(` order by %s`, orderFields[key])
	}
	return fmt.Sprintf(` order by %s desc`, orderFields[key])
}
