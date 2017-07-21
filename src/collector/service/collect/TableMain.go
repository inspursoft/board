package collect

import (
	"git/inspursoft/board/src/collector/dao"
	"git/inspursoft/board/src/collector/util"
	"strconv"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	modelK8s "k8s.io/client-go/pkg/api/v1"
	"log"
	"git/inspursoft/board/src/collector/model/collect"
	"git/inspursoft/board/src/collector/model/collect/dashboard"
	"strings"
	"github.com/google/cadvisor/info/v2"
	"os"
)

var PodList modelK8s.PodList
var NodeList modelK8s.NodeList
var ServiceList modelK8s.ServiceList
var KuberMasterIp string
var KuberMasterStatus bool
var podItem []modelK8s.Pod
//var nodeCollect []collect.Node

func init() {
	KuberMasterIp = os.Getenv("KUBEIP")
	//KuberMasterIp = "http://10.110.18.107:8080"
	K8sApiLinkTest()
}
func K8sApiLinkTest() {
	_, err := http.Get(KuberMasterIp + "/version")
	if err != nil {
		util.Logger.SetFatal(err)
		KuberMasterStatus = false
	} else {
		KuberMasterStatus = true
	}
	log.Printf("%s\t%s\t%s\t", "KuberMasterStatus status is ", strconv.FormatBool(KuberMasterStatus), time.Now())
}

//get resource form k8s api-server
func k8sGet(resource interface{}, urls string) {
	if body, err2 := ioutil.ReadAll(func() *http.Response {
		resp, err1 := http.Get(KuberMasterIp + urls)
		if err1 != nil {
			util.Logger.SetFatal(err1)
		}
		return resp
	}().Body); err2 != nil {
		util.Logger.SetFatal(err2)
	} else {
		err3 := json.Unmarshal(body, &resource)
		if err3 != nil {
			util.Logger.SetFatal(err2)
		}
	}
}

// insert time list table
func timeList() {
	var t dashboard.TimeListLog
	t.RecordTime = time.Now().Unix()
	(*serviceDashboardID)[*minuteCounterI], _ = dao.InsertDb(&t)
	util.Logger.SetInfo((*serviceDashboardID)[*minuteCounterI], *serviceDashboardID, *minuteCounterI)
}

//get nodes info from k8s apiserver
func (this *SourceMap) GainPods() error {
	k8sGet(&PodList, "/api/v1/pods")
	this.maps.ServiceCount = make(map[string]ServiceLog)
	this.maps.PodContainerCount = make(map[string]int64)
	getPods(this, PodList.Items)
	util.Logger.SetInfo("pods is insert")
	ThreadCountGet.Done()
	return nil
}

//get insert data for nodes k8s info from the func

func (resource SourceMap) GainNodes() error {
	var nodeCollect []collect.Node
	k8sGet(&NodeList, "/api/v1/nodes")
	for _, v := range NodeList.Items {
		var nodes = resource.nodes
		nodes.NodeName = v.Name
		nodes.CreateTime = v.CreationTimestamp.Format("2006-01-02 15:04:05")
		for k, v := range v.Status.Capacity {
			switch k {
			case "cpu":
				nodes.NumbersCpuCore = v.String()
			case "memory":
				nodes.MemorySize = v.String()
			case "alpha.kubernetes.io/nvidia-gpu":
				nodes.NumbersGpuCore = v.String()
			case "pods":
				nodes.PodLimit = v.String()
			}
		}

		nodes.TimeListId = (*serviceDashboardID)[*minuteCounterI]
		nodes.InternalIp = v.Status.Addresses[1].Address
		cpu, mem := getNodePs(v.Status.Addresses[1].Address)
		s := GetNodeMachine(v.Status.Addresses[1].Address)
		a, _ := s.(struct {
			outCapacity int64
			outUse      int64
		})
		nodes.StorageTotal = a.outCapacity
		nodes.StorageUse = a.outUse
		nodes.CpuUsage = float32(cpu)
		nodes.MemUsage = float32(mem)
		nodeCollect = append(nodeCollect, nodes)
		dao.InsertDb(&nodes)

	}
	util.Logger.SetInfo("nodes is insert")
	ThreadCountGet.Done()
	return nil

}

func getNodeJson(pre interface{}, method string, url string) {
	var r http.Request
	r.ParseForm()
	bodyStr := strings.TrimSpace(r.Form.Encode())
	request, _ := http.NewRequest(method, url, strings.NewReader(bodyStr))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, pre)
}
func GetNodeMachine(ip string) (interface{}) {
	if ip == "127.0.0.1" {
		return nil
	}
	url := "http://" + ip + ":4194/api/v2.0/storage"
	var storage []v2.MachineFsStats
	getNodeJson(&storage, "GET", url)
	var outCapacity uint64
	var outUse uint64
	for _, v := range storage {
		outCapacity = *v.Capacity + outCapacity
		outUse = *v.Usage + outUse
	}
	return struct {
		outCapacity int64
		outUse      int64
	}{outCapacity: int64(outCapacity),
		outUse:    int64(outUse)}
}

//get nodes ps info
func getNodePs(ip string) (cpu float32, mem float32) {
	var y []v2.ProcessInfo
	var r http.Request
	r.ParseForm()

	bodyStr := strings.TrimSpace(r.Form.Encode())
	if ip == "127.0.0.1" {
		return
	}
	request, _ := http.NewRequest("GET", "http://"+ip+":4194/api/v2.0/ps/", strings.NewReader(bodyStr))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	util.Logger.SetFatal(ip)
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &y)
	var c, m float32

	for _, v := range y {
		c = c + v.PercentCpu
		m = m + v.PercentMemory
	}
	cpu = c
	mem = m
	return
}

//get server info
func (resource *SourceMap) GainServices() error {
	k8sGet(&ServiceList, "/api/v1/services")
	for _, v := range ServiceList.Items {
		var service = resource.services
		service.CreateTime = v.CreationTimestamp.Format("2006-01-02 15:04:05")
		service.ServiceName = v.Name
		service.TimeListId = (*serviceDashboardID)[*minuteCounterI]
		for k, v := range v.Spec.Selector {
			var kvMap collect.ServiceKvMap
			kvMap.Name = k
			kvMap.Value = v
			kvMap.Belong = service.ServiceName
			kvMap.TimeListId = (*serviceDashboardID)[*minuteCounterI]
			resource.maps.ServiceMap = append(resource.maps.ServiceMap, kvMap)
			dao.InsertDb(&kvMap)
		}
		dao.InsertDb(&service)
	}
	ThreadCountGet.Done()
	return nil
}
