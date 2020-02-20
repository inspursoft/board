package nodeModel

import (
	"time"
)

const AddNodeYamlFile = "/root/ansible_k8s/addNode.yml"
const RemoveNodeYamlFile = "/root/ansible_k8s/uninstallnode.yml"
const AddRemoveNodeFile = "/root/ansible_k8s/addNode"
const AddNodeListJson = "/root/ansible_k8s/addNodeInfo/addNodeList.json"
const AddNodeHistoryJson = "/root/ansible_k8s/addNodeInfo/addNodeHistory.json"
const AddNodeLogPath = "/root/ansible_k8s/addNodeInfo/Logs/"

//const AddNodeYamlFile = "/Users/liyanqing/ansible_k8s/addNode.yml"
//const RemoveNodeYamlFile = "/Users/liyanqing/ansible_k8s/uninstallnode.yml"
//const AddRemoveNodeFile = "/Users/liyanqing/ansible_k8s/addNode"
//const AddNodeListJson = "/Users/liyanqing/ansible_k8s/addNodeInfo/addNodeList.json"
//const AddNodeHistoryJson = "/Users/liyanqing/ansible_k8s/addNodeInfo/addNodeHistory.json"
//const AddNodeLogPath = "/Users/liyanqing/ansible_k8s/addNodeInfo/Logs/"

type NodeLogResponseStatus int
type ActionType int

const (
	NodeLogResponseUnKnown NodeLogResponseStatus = 0
	NodeLogResponseStart   NodeLogResponseStatus = 1
	NodeLogResponseNormal  NodeLogResponseStatus = 2
	NodeLogResponseError   NodeLogResponseStatus = 3
	NodeLogResponseWarning NodeLogResponseStatus = 4
	NodeLogResponseSuccess NodeLogResponseStatus = 5
	NodeLogResponseFailed  NodeLogResponseStatus = 6
)

const (
	ActionTypeAddNode    ActionType = 0
	ActionTypeDeleteNode ActionType = 1
)

type NodeLogDetailArray = []NodeLogDetail;

type LogHistory struct {
	Ip           string     `json:"ip"`
	Type         ActionType `json:"type"`
	Success      bool       `json:"success"`
	Pid          int        `json:"pid"`
	CreationTime int64      `json:"creation_time"`
	Completed    bool       `json:"completed"`
}

type NodeLogDetail struct {
	Message string                `json:"message"`
	Status  NodeLogResponseStatus `json:"status"`
}

type NodeListType struct {
	Ip           string    `json:"ip"`
	CreationTime time.Time `json:"creation_time"`
}
