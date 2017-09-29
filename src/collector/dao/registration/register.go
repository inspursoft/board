package registration

import (
	"fmt"

	"git/inspursoft/board/src/collector/cmd/app"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	fmt.Println("Initializing DB registration.")
	sqlip := app.RunFlag.ServerDbIp
	sqlpass := app.RunFlag.ServerDbPassword
	sqlpo := app.RunFlag.ServerDbPort
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//err := orm.RegisterDataBase("default", "mysql", "root:"+sqlpass+"@tcp(mysql:3306)/board?charset=utf8")
	connStr := fmt.Sprintf("root:%s@tcp(%s:%s)/board?charset=utf8", sqlpass, sqlip, sqlpo)
	err := orm.RegisterDataBase("default", "mysql",
		connStr)
	orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Error occurred on registering DB: %+v\n", err)
	}
}
