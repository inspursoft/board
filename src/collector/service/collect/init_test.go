package collect_test

// import (
// 	"fmt"
// 	"git/inspursoft/board/src/collector/service/collect"
// 	"testing"

// 	"github.com/astaxie/beego/orm"
// 	_ "github.com/go-sql-driver/mysql"
// )

// func TestRunOneCycle(t *testing.T) {
// 	collect.SetInitVar("10.110.18.26", "8080")
// 	fmt.Println("Initializing DB registration.")
// 	orm.RegisterDriver("mysql", orm.DRMySQL)
// 	err := orm.RegisterDataBase("default", "mysql",
// 		"root:root123@tcp(localhost:3306)/board?charset=utf8")
// 	orm.RunSyncdb("default", false, true)
// 	if err != nil {
// 		fmt.Printf("Error occurred on registering DB: %+v\n", err)
// 	}
// 	collect.RunOneCycle()
// }
