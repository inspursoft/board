package collect

import (
	"sync"
	modelK8s "k8s.io/client-go/pkg/api/v1"
)

var ThreadCountGet sync.WaitGroup
var ThreadMap sync.WaitGroup

func RunOneCycle() error {
	var newSource SourceMap
	newSource.gainResource()
	ThreadCountGet.Wait()
	newSource.MapRun()
	PodList = modelK8s.PodList{}
	NodeList = modelK8s.NodeList{}
	ServiceList = modelK8s.ServiceList{}
	podItem = []modelK8s.Pod{}
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
