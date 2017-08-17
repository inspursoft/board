package dao

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetDashboardServiceList(t *testing.T) {
	s:=DashboardNodeDao{}
	s.TimeCount=499
	s.TimeUnit="second"
	s.TimeStamp=1501586374
	s.Name="10.110.18.107"
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:123456@tcp(10.110.18.107:30000)/k8s?charset=utf8")
	orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Error occurred on registering DB: %+v\n", err)
	}
	fmt.Println(s.GetTotalNodeData())
	fmt.Println(s.GetNodeListData())
	fmt.Println(s.GetNodeData())
	se:=DashboardServiceDao{}
	se.TimeCount=499
	se.TimeStamp=1500372187
	se.DuraTime=100000
	se.TimeUnit="minute"
	se.Name="demoshow"
	fmt.Println(se.GetTotalServiceData())
	s0,s1,s2:=se.GetServiceData()
	fmt.Println(s0,s1,s2)
}
