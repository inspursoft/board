package dashboard

import (
	"fmt"
	_"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"testing"
)

func TestGetTotalDashboardServiceDao(t *testing.T) {
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:123456@tcp(10.110.18.107:30000)/k8s?charset=utf8")
	orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Error occurred on registering DB: %+v\n", err)
	}
	GetTotal("hour", "10", "1000000000000")
	GetTotal("day", "10", "10000000000000")
	GetTotal("minute", "10", "1000000000000")
	GetTotal("second", "10", "100000000000")

	a, _ := GetService("hour", "10", "1000000000000", "")
	b, _ := GetService("day", "10", "10000000000000", "")
	c, _ := GetService("minute", "10", "1000000000000", "")
	d, _ := GetService("second", "10", "100000000000", "")
	f, _ := GetService("second", "10", "100000000000", "redis-master")
	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Println(string(c))
	fmt.Println(string(d))
	fmt.Println(string(f))

}
