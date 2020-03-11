package collect

import (
	"git/inspursoft/board/src/common/model"
	"sync"
)

var ThreadCountGet sync.WaitGroup
var ThreadMap sync.WaitGroup

func RunOneCycle() error {
	var newSource SourceMap
	newSource.gainResource()
	ThreadCountGet.Wait()
	newSource.MapRun()
	PodList = model.PodList{}
	NodeList = model.NodeList{}
	ServiceList = model.ServiceList{}
	podItem = []model.Pod{}
	return nil
}

func (c *SourceMap) gainResource() {
	ThreadCountGet.Add(3)
	timeList()
	go c.GainPods()
	go c.GainNodes()
	go c.GainServices()
}

func (c *SourceMap) MapRun() {
	ThreadMap.Add(1)
	c.maps.PreMap()
	c.maps.dashaboardCollect5s()
	ThreadMap.Done()
}
