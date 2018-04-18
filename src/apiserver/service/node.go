package service

import (
	"encoding/json"
	"errors"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"

	"strings"

	//"github.com/astaxie/beego/logs"
	"github.com/google/cadvisor/info/v2"
	modelK8s "k8s.io/client-go/pkg/api/v1"
)

type NodeStatus int

const (
	_ NodeStatus = iota
	Running
	Unschedulable
	Unknown
)

type NodeListResult struct {
	NodeName string     `json:"node_name"`
	NodeIP   string     `json:"node_ip"`
	Status   NodeStatus `json:"status"`
}
type NodeInfo struct {
	NodeName     string  `json:"node_name" orm:"column(node_name)"`
	NodeIP       string  `json:"node_ip" orm:"column(node_ip)"`
	CreateTime   int64   `json:"create_time" orm:"column(create_time)"`
	CPUUsage     float32 `json:"cpu_usage" orm:"column(cpu_usage)"`
	MemoryUsage  float32 `json:"memory_usage" orm:"column(memory_usage)"`
	MemorySize   string  `json:"memory_size" orm:"column(memory_size)"`
	StorageTotal uint64  `json:"storage_total" orm:"column(storage_total)"`
	StorageUse   uint64  `json:"storage_use" orm:"column(storage_usage)"`
}

var kubeNodeURL = utils.GetConfig("KUBE_NODE_URL")

func GetNode(nodeName string) (node NodeInfo, err error) {
	var Node modelK8s.NodeList
	defer func() { recover() }()
	err = getFromRequest(kubeNodeURL(), &Node)
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
				CPUUsage:     cpu,
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
func GetNodeList() (res []NodeListResult) {

	var Node modelK8s.NodeList
	defer func() { recover() }()
	err := getFromRequest(kubeNodeURL(), &Node)
	if err != nil {
		return
	}
	for _, v := range Node.Items {
		res = append(res, NodeListResult{
			NodeName: v.Status.Addresses[1].Address,
			NodeIP:   v.Status.Addresses[1].Address,
			Status: func() NodeStatus {
				if v.Spec.Unschedulable {
					return Unschedulable
				}
				for _, cond := range v.Status.Conditions {
					if strings.EqualFold(string(cond.Type), "Ready") && cond.Status == modelK8s.ConditionTrue {
						return Running
					}
				}
				return Unknown
			}()})
	}
	return
}

func CreateNodeGroup(nodeGroup model.NodeGroup) (*model.NodeGroup, error) {
	nodeGroupID, err := dao.AddNodeGroup(nodeGroup)
	if err != nil {
		return nil, err
	}
	nodeGroup.ID = nodeGroupID
	return &nodeGroup, err
}

func UpdateNodeGroup(n model.NodeGroup, fieldNames ...string) (bool, error) {
	if n.ID == 0 {
		return false, errors.New("no Node group ID provided")
	}
	_, err := dao.UpdateNodeGroup(n, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteNodeGroupByID(n model.NodeGroup) (int64, error) {
	if n.ID == 0 {
		return 0, errors.New("no Node Group ID provided")
	}
	num, err := dao.DeleteNodeGroup(n)
	if err != nil {
		return 0, err
	}
	return num, nil
}
