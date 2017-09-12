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
	StorageTotal int64   `json:"storage_total" orm:"column(storage_total)"`
	StorageUse   int64   `json:"storage_use" orm:"column(storage_usage)"`
}

func GetNode(nodeName string) (node NodeInfo, err error) {
	var Node modelK8s.NodeList
	defer func() { recover() }()
	var url string
	url = fmt.Sprintf("%s:%s/api/v1/nodes", os.Getenv("KUBE_IP"), os.Getenv("KUBE_PORT"))
	err = getFromK8sApi(url, &Node)
	fmt.Println("dddd", err, url)
	for _, v := range Node.Items {
		var mlimit string
		if strings.Contains(v.Status.Addresses[1].Address, nodeName) {
			for k, v := range v.Status.Capacity {
				switch k {
				case "memory":
					mlimit = v.String()
				}
			}
			time := v.CreationTimestamp.Unix()
			var y []v2.ProcessInfo
			getFromK8sApi(nodeName+":4194/api/v2.0/ps/", &y)
			var c, m float32
			for _, v := range y {
				c = c + v.PercentCpu
				m = m + v.PercentMemory
			}
			cpu := c
			mem := m
			var fs []v2.MachineFsStats
			getFromK8sApi(nodeName+":4194/api/v2.0/storage", &fs)
			var Capacity uint64
			var Use uint64
			for _, v := range fs {
				Capacity = *v.Capacity + Capacity
				Use = *v.Usage + Use
			}
			node = NodeInfo{
				NodeName:     nodeName,
				NodeIP:       nodeName,
				CreateTime:   time,
				CpuUsage:     cpu,
				MemoryUsage:  mem,
				MemorySize:   mlimit,
				StorageTotal: int64(Capacity),
				StorageUse:   int64(Use),
			}
			break
		}
	}

	return
}
func getFromK8sApi(url string, source interface{}) (err error) {
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
