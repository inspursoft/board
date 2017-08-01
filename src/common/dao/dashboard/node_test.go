package dashboard

import (
	"fmt"
	_"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"regexp"
	"testing"
)

func TestQueryNode(t *testing.T) {
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:123456@tcp(10.110.18.107:30000)/k8s?charset=utf8")
	orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Error occurred on registering DB: %+v\n", err)
	}
	nodeName := "10.110.18.71"
	recordTime := "10000000000000"
	QueryNode(nodeName, recordTime, "10", "second")

}
func TestNodeList(t *testing.T) {
	ma:=`^(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[0-9])$`
	ip:=`10.110.18.1193`
	a,n:=regexp.MatchString(ma,ip)
	fmt.Println(a,n)
	NodeList()
}


