package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"strings"

	"github.com/google/cadvisor/info/v2"
	modelK8s "k8s.io/client-go/pkg/api/v1"
)

type NodeInfo struct {
	NodeName     string  `json:"node_name" orm:"column(node_name)"`
	NodeIP       string  `json:"node_ip" orm:"column(node_ip)"`
	CreateTime   int64   `json:"create_time" orm:"column(create_time)"`
	CpuUsage     float32 `json:"cpu_usage" orm:"column(cpu_usage)"`
	MemoryUsage  float32 `json:"memory_usage" orm:"column(memory_usage)"`
	MemorySize   string  `json:"memory_size" orm:"column(memory_size)"`
	StorageTotal uint64   `json:"storage_total" orm:"column(storage_total)"`
	StorageUse   uint64   `json:"storage_use" orm:"column(storage_usage)"`
}

var NodeUrl = fmt.Sprintf("%s:%s/api/v1/nodes", os.Getenv("KUBE_IP"), os.Getenv("KUBE_PORT"))

func GetNode(nodeName string) (node NodeInfo, err error) {
	var Node modelK8s.NodeList
	defer func() { recover() }()
	var url string
	url = NodeUrl
	err = getFromRequest(url, &Node)
	if err != nil {
		return
	}
	for _, v := range Node.Items {
		var mlimit string
		if strings.EqualFold(v.Status.Addresses[1].Address, nodeName) {
			for k, v := range v.Status.Capacity {
				switch k {
				case "memory":
					mlimit = v.String()
				}
			}
			time := v.CreationTimestamp.Unix()
			var ps []v2.ProcessInfo
			getFromRequest("http://"+nodeName+":4194/api/v2.0/ps/", &ps)
			var c, m float32
			for _, v := range ps {
				c += v.PercentCpu
				m += v.PercentMemory
			}
			cpu := c
			mem := m
			var fs []v2.MachineFsStats
			getFromRequest("http://"+nodeName+":4194/api/v2.0/storage", &fs)
			var capacity uint64
			var use uint64
			for _, v := range fs {
				capacity += *v.Capacity
				use += *v.Usage
			}
			node = NodeInfo{
				NodeName:     nodeName,
				NodeIP:       nodeName,
				CreateTime:   time,
				CpuUsage:     cpu,
				MemoryUsage:  mem,
				MemorySize:   mlimit,
				StorageTotal: capacity,
				StorageUse:   use,
			}
			break
		}
	}

	return
}
func getFromRequest(url string, source interface{}) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, source)
	if err != nil {
		return
	}
	return nil
}
func SuspendNode(nodeName string) (bool, error) {
	return Suspend(nodeName)
}
func ResumeNode(nodeName string) (bool, error) {
	return Resume(nodeName)
}
