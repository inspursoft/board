package registration
import (
	"fmt"

	"github.com/astaxie/beego/orm"
	 _"github.com/go-sql-driver/mysql"
	"git/inspursoft/board/src/collector/cmd/app"
	"git/inspursoft/board/src/collector/util"
)

func init() {
	fmt.Println("Initializing DB registration.")
	sqlip:=app.RunFlag.ServerDbIp
	sqlpass:=app.RunFlag.ServerDbPassword
	sqlpo:=app.RunFlag.ServerDbPort
	util.Logger.SetInfo("root:"+sqlpass+"@tcp("+sqlip+":"+sqlpo+")/k8s?charset=utf8")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql",
		"root:"+sqlpass+"@tcp("+sqlip+":"+sqlpo+")/k8s?charset=utf8")
	orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Error occurred on registering DB: %+v\n", err)
	}
}

