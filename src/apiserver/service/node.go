package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	"github.com/google/cadvisor/info/v2"
	//modelK8s "k8s.io/client-go/pkg/api/v1"
)

type NodeStatus int

const (
	_ NodeStatus = iota
	Running
	Unschedulable
	Unknown
)

const (
	K8sLabel = "kubernetes.io"
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

func GetNode(nodeName string) (node NodeInfo, err error) {
	defer func() { recover() }()
	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nodecli := k8sclient.AppV1().Node()

	Node, err := nodecli.List()
	if err != nil {
		logs.Error("Failed to get Node List")
		return
	}
	for _, v := range Node.Items {
		var mlimit string
		if strings.EqualFold(v.Status.Addresses[1].Address, nodeName) {
			//for k, v := range v.Status.Capacity {
			//	switch k {
			//	case "memory":
			//		mlimit = v.String()
			//	}
			//}
			mlimit = fmt.Sprintf("%s", v.Status.Capacity["memory"])
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

	//var nodecli model.NodeCli
	defer func() { recover() }()

	//nodecli, err := k8sassist.NewNodes()

	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nodecli := k8sclient.AppV1().Node()

	Node, err := nodecli.List()
	if err != nil {
		logs.Error("Failed to get Node List")
		return
	}

	for _, v := range Node.Items {
		res = append(res, NodeListResult{
			NodeName: v.Status.Addresses[1].Address,
			NodeIP:   v.Status.Addresses[1].Address,
			Status: func() NodeStatus {
				if v.Unschedulable {
					return Unschedulable
				}
				for _, cond := range v.Status.Conditions {
					if strings.EqualFold(string(cond.Type), "Ready") && cond.Status == model.ConditionTrue {
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

func GetNodeGroup(nodeGroup model.NodeGroup, selectedFields ...string) (*model.NodeGroup, error) {
	n, err := dao.GetNodeGroup(nodeGroup, selectedFields...)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func GetNodeGroupList() ([]model.NodeGroup, error) {
	return dao.GetNodeGroups()
}

func NodeGroupExists(nodeGroupName string) (bool, error) {
	query := model.NodeGroup{GroupName: nodeGroupName}
	nodegroup, err := dao.GetNodeGroup(query, "name")
	if err != nil {
		return false, err
	}
	return (nodegroup != nil && nodegroup.ID != 0), nil
}

func AddNodeToGroup(nodeName string, groupName string) error {
	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nInterface := k8sclient.AppV1().Node()

	nNode, err := nInterface.Get(nodeName)
	if err != nil {
		logs.Error("Failed to get K8s node")
		return err
	}
	//logs.Info(nNode)

	logs.Debug(nNode.ObjectMeta.Labels)
	nNode.ObjectMeta.Labels[groupName] = "true"

	newNode, err := nInterface.Update(nNode)
	if err != nil {
		logs.Error("Failed to update K8s node")
		return err
	}
	logs.Debug(newNode)
	return nil
}

func GetGroupOfNode(nodeName string) ([]string, error) {
	var groups []string
	//nInterface, err := k8sassist.NewNodes()
	//if err != nil {
	//	logs.Error("Failed to get node client interface")
	//	return nil, err
	//}

	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nInterface := k8sclient.AppV1().Node()

	nNode, err := nInterface.Get(nodeName)
	if err != nil {
		logs.Error("Failed to get K8s node")
		return nil, err
	}
	for key, _ := range nNode.ObjectMeta.Labels {
		if !strings.Contains(key, K8sLabel) {
			groups = append(groups, key)
		}
	}
	return groups, nil
}

func NodeOrNodeGroupExists(nodeOrNodeGroupName string) (bool, error) {
	nodeGroupExists, err := NodeGroupExists(nodeOrNodeGroupName)
	if err != nil {
		return false, err
	}
	if !nodeGroupExists {
		res, err := GetNode(nodeOrNodeGroupName)
		if err != nil {
			return false, err
		}
		if res.NodeName == "" {
			return false, nil
		}
	}
	return true, nil
}

func RemoveNodeFromGroup(nodeName string, groupName string) error {
	//nInterface, err := k8sassist.NewNodes()
	//if err != nil {
	//	logs.Error("Failed to get node client interface")
	//	return err
	//}
	var config k8sassist.K8sAssistConfig
	config.K8sMasterURL = kubeMasterURL()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nInterface := k8sclient.AppV1().Node()

	nNode, err := nInterface.Get(nodeName)
	if err != nil {
		logs.Error("Failed to get K8s node")
		return err
	}
	//logs.Debug(nNode.ObjectMeta.Labels)
	delete(nNode.ObjectMeta.Labels, groupName)

	newNode, err := nInterface.Update(nNode)
	if err != nil {
		logs.Error("Failed to update K8s node")
		return err
	}
	logs.Debug(newNode.ObjectMeta.Labels)
	return nil
}

func RemoveNodeGroup(groupName string) error {
	// Check nodegroup in DB
	ngQuery, err := GetNodeGroup(model.NodeGroup{GroupName: groupName}, "name")
	if err != nil {
		logs.Error("Failed to get group %s in DB", groupName)
		return err
	}
	if ngQuery == nil {
		logs.Info("%s not in system DB", groupName)
		return nil
	}
	if ngQuery.Deleted == 1 {
		logs.Info("%s deleted in system DB", groupName)
		return nil
	}

	//TODO：Need to change it, do not traverse all nodes in huge cluster
	nodeList := GetNodeList()
	for _, nodeinfo := range nodeList {
		groupList, err := GetGroupOfNode(nodeinfo.NodeName)
		if err != nil {
			logs.Error("Failed to check node %s group", nodeinfo.NodeName)
			return err
		}
		for _, g := range groupList {
			if groupName == g {
				// Remove this groupname from node
				err = RemoveNodeFromGroup(nodeinfo.NodeName, groupName)
				if err != nil {
					logs.Error("Failed to remove %s from node %s", g, nodeinfo.NodeName)
					return err
				}
				break
			}
		}
	}
	// Remove it in group DB
	_, err = dao.DeleteNodeGroup(*ngQuery)
	if err != nil {
		logs.Error("Failed to delete %s in DB", ngQuery.GroupName)
		return err
	}
	return nil
}

func RemovePodByNode(node string) error {
	podList, err := GetPods()
	if err != nil {
		logs.Info("Failed to get pods from system", err)
		return err
	}
	for _, v := range podList.Items {
		if v.Status.HostIP == node {
			logs.Info("Gracefully remove the pod %s from node %s", v.Name, node)
			k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
				K8sMasterURL: kubeMasterURL(),
			})
			//TODO need evict in released version
			err = k8sclient.AppV1().Pod(v.Namespace).Delete(v.Name)
			if err != nil {
				logs.Info("Failed to Delete pod", v.Name, err)
				return err
			}
		}
	}
	return nil
}

func GetNodesAvailableResources() ([]model.NodeAvailableResources, error) {
	var resources []model.NodeAvailableResources
	c := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: kubeMasterURL(),
	})
	l, err := c.AppV1().Node().List()
	if err != nil {
		logs.Debug("Failed to get node list %v", c)
		return nil, err
	}
	logs.Debug("Node List: %v", l)
	for _, node := range l.Items {
		// TODO: check the status of node
		var noderesource model.NodeAvailableResources
		noderesource.NodeName = node.Name
		noderesource.CPUAvail = string(node.Status.Allocatable[model.ResourceCPU])
		noderesource.MemAvail = string(node.Status.Allocatable[model.ResourceMemory])
		resources = append(resources, noderesource)
	}
	return resources, nil
}
