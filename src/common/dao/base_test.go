package dao

import (
	"fmt"
	"git/inspursoft/board/src/common/utils"
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
	os.Exit(m.Run())
}
