package app

import (
	"git/inspursoft/board/src/collector/service/collect"
	"sync"
	"time"
	"git/inspursoft/board/src/collector/util"
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
	for range ticker.C {
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
