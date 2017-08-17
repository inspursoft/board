package controller

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego/orm"
)

func init() {
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:123456@tcp(10.110.18.107:30000)/k8s?charset=utf8")
	orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Error occurred on registering DB: %+v\n", err)
	}
}
func TestDashboardNodeController_GetList(t *testing.T) {
	var s *DashboardNodeController
	s = new(DashboardNodeController)
	s.GetNodeList()
	s.GetService()
}
