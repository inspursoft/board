package collect

import (
	"git/inspursoft/board/src/collector/util"
)
func (c *KvMap) PreMap() {
	for _, service := range c.ServiceMap {
		for _, pod := range c.PodMap {
			if service.Name == pod.Name && service.Value == pod.Value {
				c.serviceLog.ServiceName = service.Belong
				c.serviceLog.PodName = pod.Belong
				c.serviceLog.ContainerNumber = c.PodContainerCount[pod.Belong]
				c.ServiceTemp = append(c.ServiceTemp, c.serviceLog)
			}
		}
	}
	for _, v := range c.ServiceTemp {
		if v.PodNumber == 0 {
			v.PodNumber = 1
		}
		if v.ServiceName == c.ServiceCount[v.ServiceName].ServiceName {
			v.PodNumber = c.ServiceCount[v.ServiceName].PodNumber
			v.ContainerNumber = v.ContainerNumber + c.ServiceCount[v.ServiceName].ContainerNumber
			v.PodNumber = v.PodNumber + 1
			c.ServiceCount[v.ServiceName] = v
		} else {
			c.ServiceCount[v.ServiceName] = v
		}
	}
	util.Logger.SetInfo("c.ServiceCount", c.ServiceCount)
}
