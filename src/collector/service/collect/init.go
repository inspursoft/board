package collect

import (
	"sync"
)

var ThreadCountGet sync.WaitGroup
var ThreadMap sync.WaitGroup

func RunOneCycle() error {
	var newSource GainKubernetes = new(SourceMap)
	gainResource(newSource)
	ThreadCountGet.Wait()
	newSource.MapRun()
	return nil
}

func gainResource(newSource GainKubernetes) {
	ThreadCountGet.Add(3)
	timeList()
	go newSource.GainPods()
	go newSource.GainNodes()
	go newSource.GainServices()
}

func (c *SourceMap) MapRun() {
	ThreadMap.Add(1)
	c.maps.PreMap()
	c.maps.dashaboardCollect5s()
	ThreadMap.Done()
}
