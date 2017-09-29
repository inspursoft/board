package collect

import model "git/inspursoft/board/src/collector/model/collect"

var minuteCounterI *int
var hourCounterI *int
var dayCounterI *int
var nodeID int64
var podID int64
var serviceDashboardID *[12]int64
var minuteServiceDashboardID *[60]int64
var hourServiceDashboardID *[24]int64

type KvMap struct {
	PodMap            []model.PodKvMap
	ServiceMap        []model.ServiceKvMap
	PodContainerCount map[string]int64
	ServiceTemp       []ServiceLog
	ServiceCount      map[string]ServiceLog
	serviceLog        ServiceLog
}

type ServiceLog struct {
	ServiceName     string
	PodName         string
	PodNumber       int64
	ContainerNumber int64
	TimeListId      int64
}
type SourceMap struct {
	nodes    model.Node
	pods     model.Pod
	services model.Service
	maps     KvMap
}
type GainKubernetes interface {
	GainPods() error
	GainNodes() error
	GainServices() error
	MapRun()
}
