package dao_test

import (
	"fmt"
	"github.com/inspursoft/board/src/common/dao"
	"testing"
)

func TestGetDashboardServiceList(t *testing.T) {
	s := dao.DashboardNodeDao{}
	s.TimeCount = 499
	s.TimeUnit = "second"
	s.TimeStamp = 1501586374
	fmt.Println(s.GetTotalNodeData())
	fmt.Println(s.GetNodeListData())
	fmt.Println(s.GetNodeData())
	se := dao.DashboardServiceDao{}
	se.TimeCount = 499
	se.TimeStamp = 1500372187
	se.DuraTime = 100000
	se.TimeUnit = "minute"
	se.Name = "demoshow"
	fmt.Println(se.GetTotalServiceData())
	s0, s1, s2 := se.GetServiceData()
	fmt.Println(s0, s1, s2)
	fmt.Println(se.GetLimitTime())
	fmt.Println(s.GetLimitTime())
}
