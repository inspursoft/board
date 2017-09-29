package collect

import (
	"git/inspursoft/board/src/collector/model/collect/dashboard"
	"git/inspursoft/board/src/collector/dao"
	"git/inspursoft/board/src/collector/util"
	"git/inspursoft/board/src/collector/model/collect"
)

func init() {
	minuteCounterI = new(int)
	hourCounterI = new(int)
	dayCounterI = new(int)
	serviceDashboardID = new([12]int64)
	minuteServiceDashboardID = new([60]int64)
	hourServiceDashboardID = new([24]int64)
}
func (c *KvMap) dashaboardCollect5s() {
	c.serviceLog.TimeListId = (*serviceDashboardID)[*minuteCounterI]
	for k, v := range c.ServiceCount {
		s := setServerSecond(v.PodNumber, k, v.ContainerNumber, (*serviceDashboardID)[*minuteCounterI])
		dao.InsertDb(&s)
	}
	c.startCollectMinute()()
}
func setServerSecond(podNumber int64, serviceName string,
	containerNumber int64, timeListId int64) dashboard.ServiceDashboardSecond {
	return dashboard.ServiceDashboardSecond{
		PodNumber:       podNumber,
		ServiceName:     serviceName,
		ContainerNumber: containerNumber,
		TimeListId:      timeListId, }
}

func setNode(timeUnit string)  () {
	switch timeUnit {
	case "minute":
		var s dashboard.NodeDashboardMinute
		temp:=dao.CalcNode(timeUnit).([]collect.Node)
		for _,v:=range temp{
			s=dashboard.NodeDashboardMinute(v)
			dao.InsertDb(&(s))
		}
	case "hour":
		var s dashboard.NodeDashboardHour
		temp:=dao.CalcNode(timeUnit).([]dashboard.NodeDashboardMinute)
		for _,v:=range temp{
			s=dashboard.NodeDashboardHour(v)
			dao.InsertDb(&(s))
		}
	case "day":
		var s dashboard.NodeDashboardDay
		temp:=dao.CalcNode(timeUnit).([]dashboard.NodeDashboardHour)
		for _,v:=range temp{
			s=dashboard.NodeDashboardDay(v)
			dao.InsertDb(&(s))
		}

	}


}

func (c *KvMap) startCollectMinute() func() {
	*minuteCounterI = *minuteCounterI + 1
	if *minuteCounterI == 12 {
		*minuteCounterI=0
		util.Logger.SetInfo("start collect minute")
		setNode("minute")
		return c.DashaboardCollectMinute
	} else {
		return func() {
			util.Logger.SetInfo("minuteCounterI", *minuteCounterI)
		}
	}
}

func (c *KvMap) DashaboardCollectMinute() {
	a := minuteCalc()
	*minuteCounterI = 0
	for _, v := range a {
		var s dashboard.ServiceDashboardMinute
		s.PodNumber = v.PodNumber
		s.ContainerNumber = v.ContainerNumber
		s.ServiceName = v.ServiceName
		s.TimeListId = (*serviceDashboardID)[*minuteCounterI]
		dao.InsertDb(&s)
		(*minuteServiceDashboardID)[*hourCounterI] = s.TimeListId
		util.Logger.SetDebug("*minuteServiceDashboardID", *minuteServiceDashboardID)
	}
	*hourCounterI = *hourCounterI + 1
	if *hourCounterI == 60 {
		c.DashaboardCollectHour()
		setNode("hour")
	}
	util.Logger.SetInfo("*hourCounterI", *hourCounterI)
}
func (c *KvMap) DashaboardCollectHour() {
	a := hourCalc()
	*hourCounterI = 0
	for _, v := range a {
		var s dashboard.ServiceDashboardHour
		s.PodNumber = v.PodNumber
		s.ContainerNumber = v.ContainerNumber
		s.ServiceName = v.ServiceName
		s.TimeListId = (*serviceDashboardID)[*minuteCounterI]
		dao.InsertDb(&s)
		(*hourServiceDashboardID)[*dayCounterI] = s.TimeListId
		util.Logger.SetDebug("*hourServiceDashboardID", *hourServiceDashboardID)
	}

	*dayCounterI = *dayCounterI + 1
	if *dayCounterI == 24 {
		c.DashaboardCollectDay()
		setNode("day")
	}
}
func (c *KvMap) DashaboardCollectDay() {
	a := dayCalc()
	*dayCounterI = 0
	for _, v := range a {
		var s dashboard.ServiceDashboardDay
		s.PodNumber = v.PodNumber
		s.ContainerNumber = v.ContainerNumber
		s.ServiceName = v.ServiceName
		s.TimeListId = (*serviceDashboardID)[*minuteCounterI]
		dao.InsertDb(&s)
	}

	util.Logger.SetInfo("*hourCounterI", *dayCounterI)
}
