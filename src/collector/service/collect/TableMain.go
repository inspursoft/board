package collect

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/collector/dao"
	"git/inspursoft/board/src/collector/model/collect"
	"git/inspursoft/board/src/collector/model/collect/dashboard"
	"git/inspursoft/board/src/collector/util"
	"git/inspursoft/board/src/common/k8sassist"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git/inspursoft/board/src/common/model"

	"github.com/google/cadvisor/info/v2"
)

var ignoredNamespaces []string = []string{"kube-system", "istio-system"}
var PodList model.PodList
var NodeList model.NodeList
var ServiceList model.ServiceList
var KuberMasterIp string
var KuberMasterStatus bool
var podItem []model.Pod
var KuberMasterURL string
var kubeConfigPath string
var KuberPort string

func SetInitVar(ip string, port string) {
	KuberMasterIp = ip
	KuberPort = port
	KuberMasterURL = fmt.Sprintf("http://%s%s%s", KuberMasterIp, ":", KuberPort)
	kubeConfigPath = `/root/kubeconfig`
	//	pingK8sApiLink()
}

//func pingK8sApiLink() {
//	url := fmt.Sprintf("%s/version", KuberMasterURL)
//	cl := &http.Client{Timeout: time.Millisecond * 2000}
//	fmt.Println("url is ", url)
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	resp, _ := cl.Do(req)
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println("kubernetes version is ", string(body))
//}

// insert time list table
func timeList() {
	var t dashboard.TimeListLog
	t.RecordTime = time.Now().Unix()
	(*serviceDashboardID)[*minuteCounterI], _ = dao.InsertDb(&t)
	util.Logger.SetInfo((*serviceDashboardID)[*minuteCounterI], *serviceDashboardID, *minuteCounterI)
}

//get nodes info from k8s apiserver
func (this *SourceMap) GainPods() error {
	defer ThreadCountGet.Done()
	c := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath,
	})
	l, err := c.AppV1().Pod("").List(model.ListOptions{})
	if err != nil {
		return err
	}
	PodList = *l
	// filter the ignored pods.
	items := make([]model.Pod, 0, len(PodList.Items))
	for _, p := range PodList.Items {
		if !shouldIgnore(p.Namespace) {
			items = append(items, p)
		}
	}
	PodList.Items = items
	this.maps.ServiceCount = make(map[string]ServiceLog)
	this.maps.PodContainerCount = make(map[string]int64)
	getPods(this, PodList.Items)
	util.Logger.SetInfo("pods is insert")
	return nil
}

//get insert data for nodes k8s info from the func

func (resource SourceMap) GainNodes() error {
	defer ThreadCountGet.Done()
	var nodeCollect []collect.Node
	c := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath,
	})
	l, err := c.AppV1().Node().List()
	if err != nil {
		return err
	}
	NodeList = *l
	for _, v := range NodeList.Items {
		var nodes = resource.nodes
		nodes.NodeName = v.Name
		nodes.CreateTime = v.CreationTimestamp.Format("2006-01-02 15:04:05")
		var cpuCores int = 1
		for k, v := range v.Status.Capacity {
			switch k {
			case "cpu":
				nodes.NumbersCpuCore = fmt.Sprintf("%v", v)
				cpuCores, _ = strconv.Atoi(string(v))
			case "memory":
				nodes.MemorySize = fmt.Sprintf("%v", v)
			case "alpha.kubernetes.io/nvidia-gpu":
				nodes.NumbersGpuCore = fmt.Sprintf("%v", v)
			case "pods":
				nodes.PodLimit = fmt.Sprintf("%v", v)
			}
		}

		nodes.TimeListId = (*serviceDashboardID)[*minuteCounterI]
		nodes.InternalIp = v.Status.Addresses[1].Address
		if func(nodeCondition []model.NodeCondition) bool {
			for _, cond := range nodeCondition {
				if strings.EqualFold(string(cond.Type), "Ready") {
					if cond.Status != model.ConditionTrue {
						return false
					}
				}
			}
			return true
		}(v.Status.Conditions) {
			cpu, mem, err := getNodePs(v.Status.Addresses[1].Address, cpuCores)
			if err != nil {
				return err
			}
			s := GetNodeMachine(v.Status.Addresses[1].Address)
			a, _ := s.(struct {
				outCapacity int64
				outUse      int64
			})
			nodes.StorageTotal = a.outCapacity
			nodes.StorageUse = a.outUse
			nodes.CpuUsage = float32(cpu)
			nodes.MemUsage = float32(mem)
		} else {
			nodes.StorageTotal = 0
			nodes.StorageUse = 0
			nodes.CpuUsage = float32(0)
			nodes.MemUsage = float32(0)
			util.Logger.SetWarn("this node status is unkown", v.Status.Addresses[1].Address)
		}
		nodeCollect = append(nodeCollect, nodes)
		dao.InsertDb(&nodes)

	}
	util.Logger.SetInfo("nodes is insert")
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
func GetNodeMachine(ip string) interface{} {
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
		outUse: int64(outUse)}
}

//get nodes ps info
func getNodePs(ip string, cpuCores int) (cpu float32, mem float32, err error) {
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
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		util.Logger.SetError("Request node ps info error: %+v", err)
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.SetError("Read response body error: %+v", err)
		return
	}
	err = json.Unmarshal(body, &y)
	if err != nil {
		util.Logger.SetError("Unmarshal response body error: %+v", err)
		return
	}
	var c, m float32

	for _, v := range y {
		c = c + v.PercentCpu
		m = m + v.PercentMemory
	}
	cpu = c / float32(cpuCores)
	mem = m
	return
}

//get server info
func (resource *SourceMap) GainServices() error {
	defer ThreadCountGet.Done()
	c := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath,
	})
	l, err := c.AppV1().Service("").List()
	if err != nil {
		return err
	}
	ServiceList = *l
	// filter the ignored services.
	items := make([]model.Service, 0, len(ServiceList.Items))
	for _, v := range ServiceList.Items {
		if !shouldIgnore(v.Namespace) {
			items = append(items, v)
		}
	}
	ServiceList.Items = items
	for _, v := range ServiceList.Items {
		var service = resource.services
		service.CreateTime = v.CreationTimestamp.Format("2006-01-02 15:04:05")
		service.ServiceName = v.Name
		service.TimeListId = (*serviceDashboardID)[*minuteCounterI]
		for k, v := range v.Selector {
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
	return nil
}

func shouldIgnore(namespace string) bool {
	for _, ns := range ignoredNamespaces {
		if ns == namespace {
			return true
		}
	}
	return false
}
