package app

import (
	"git/inspursoft/board/src/collector/dao"
	"git/inspursoft/board/src/collector/service/collect"
	"git/inspursoft/board/src/collector/util"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

var ThreadCount sync.WaitGroup
var statusSwitchOn chan bool
var statusSwitchOff chan bool
var statusSwitchLast *bool

func init() {
	statusSwitchLast = new(bool)
	*statusSwitchLast = true
}

func Run() (err error) {
	statusSwitchOn = make(chan bool)
	statusSwitchOff = make(chan bool)
	collectMainInOnCycle()
	return nil
}

func TurnStatus(status bool) {
	switch status {
	case true:
		util.Logger.SetInfo("turn the SwitchOn")
		statusSwitchOn <- true
	case false:
		util.Logger.SetInfo("turn the SwitchOff")
		statusSwitchOff <- false
	}
}

func runOneCycle() {
	collect.RunOneCycle()
	collect.ThreadCountGet.Wait()
	collect.ThreadMap.Wait()
	ThreadCount.Done()
}
func collectMainInOnCycle() {
	ticker := time.NewTicker(time.Millisecond * 5000)
	util.Logger.SetInfo("main routine is run")
	count := 0
	for range ticker.C {
		count = count + 1
		if count%20000 == 0 {
			cleanStales(20)
			count = 0
		}
		ThreadCount.Add(1)
		util.Logger.SetInfo("run with the state", "statusSwitchLast is", *statusSwitchLast)
		select {
		case i := <-statusSwitchOn:
			if i != *statusSwitchLast {
				util.Logger.SetInfo("into the select thread in statusSwitchOn")
				go runOneCycle()
				*statusSwitchLast = true
				ThreadCount.Wait()
			}
		case i := <-statusSwitchOff:
			if i != *statusSwitchLast {
				util.Logger.SetInfo("into the select thread in statusSwitchOff")
				*statusSwitchLast = false
				ThreadCount.Done()
			}
		default:
			switch *statusSwitchLast {
			case true:
				util.Logger.SetInfo("into the select thread in default on")
				go runOneCycle()
				ThreadCount.Wait()
			case false:
				util.Logger.SetInfo("into the select thread in default off")
				ThreadCount.Done()
			}
		}
	}
}

func cleanStales(dateSpan int) {
	ThreadCount.Add(1)
	for _, tableName := range []string{
		"node", "node_dashboard_day", "node_dashboard_hour", "node_dashboard_minute",
		"pod", "pod_kv_map", "service", "service_dashboard_day",
		"service_dashboard_hour", "service_dashboard_minute",
		"service_dashboard_second", "service_kv_map",
	} {
		affected, err := dao.DeleteStaleData(tableName, dateSpan)
		if err != nil {
			logs.Error("Failed to delete table: %s, error: %+v", tableName, err)
		}
		logs.Debug("Affected delete table: %s rows: %d", tableName, affected)
		time.Sleep(time.Second * 1)
	}

	affected, err := dao.DeleteStaleTimeList(dateSpan)
	if err != nil {
		logs.Error("Failed to delete table time_list_log, error: %+v", err)
	}
	logs.Debug("Affected delete table time_list_log rows: %d", affected)
	ThreadCount.Done()
}
