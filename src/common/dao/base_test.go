package dao_test

import (
	"fmt"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/orm"
)

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	dbIP := utils.GetStringValue("DB_IP")
	dbPort := utils.GetStringValue("DB_PORT")
	dbPassword := utils.GetStringValue("DB_PASSWORD")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:%s@tcp(%s:%s)/board?charset=utf8", dbPassword, dbIP, dbPort))
	model.InitModelDB()
	os.Exit(m.Run())
}
