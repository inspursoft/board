package collect

import (
	"git/inspursoft/board/src/collector/model/collect"
	"git/inspursoft/board/src/collector/util"
	"k8s.io/apimachinery/pkg/labels"
)

func (c *KvMap) PreMap() {
	selectors := generateServiceSelectors(c.ServiceMap)
	labels := generatePodLabels(c.PodMap)
	for serviceName, selector := range selectors {
		for podName, label := range labels {
			if selector.Matches(label) {
				c.serviceLog.ServiceName = serviceName
				c.serviceLog.PodName = podName
				c.serviceLog.ContainerNumber = c.PodContainerCount[podName]
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

func generateServiceSelectors(serviceMaps []collect.ServiceKvMap) map[string]labels.Selector {
	svcSelectors := map[string]map[string]string{}
	for _, service := range serviceMaps {
		selectors := svcSelectors[service.Belong]
		if selectors == nil {
			selectors = map[string]string{}
			svcSelectors[service.Belong] = selectors
		}
		selectors[service.Name] = service.Value
	}
	labelSelectors := map[string]labels.Selector{}
	for serviceName, selectors := range svcSelectors {
		labelSelectors[serviceName] = labels.Set(selectors).AsSelector()
	}
	return labelSelectors
}

func generatePodLabels(podMaps []collect.PodKvMap) map[string]labels.Labels {
	podLabelsMap := map[string]map[string]string{}
	for _, pod := range podMaps {
		labels := podLabelsMap[pod.Belong]
		if labels == nil {
			labels = map[string]string{}
			podLabelsMap[pod.Belong] = labels
		}
		labels[pod.Name] = pod.Value
	}
	podLabels := map[string]labels.Labels{}
	for podName, labelsMap := range podLabelsMap {
		podLabels[podName] = labels.Set(labelsMap)
	}
	return podLabels
}
