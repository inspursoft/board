package nodeModel

import "bytes"

const AddNodeYamlFile = "/root/ansible_k8s/addNode.yml"
const RemoveNodeYamlFile = "/root/ansible_k8s/uninstallnode.yml"
const AddRemoveNodeFile = "/root/ansible_k8s/addNode"
const AddRemoveShellFile = "/root/ansible_k8s/addNode.sh"

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

type Pagination struct {
	PageIndex  int   `json:"page_index"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	PageCount  int   `json:"page_count"`
}

type AddNodePostData struct {
	NodeIp string `json:"node_ip"`
}

type NodeLogDetailArray = []NodeLogDetail;

// database table's name: node-status
type NodeStatus struct {
	Id           int    `json:"id"`
	Ip           string `json:"ip"`
	CreationTime int64  `json:"creation_time"`
}

// database table's name: node-log
type NodeLog struct {
	Id           int        `json:"id"`
	Ip           string     `json:"ip"`
	LogType      ActionType `json:"log_type"`
	Success      bool       `json:"success"`
	Pid          int        `json:"pid"`
	CreationTime int64      `json:"creation_time"`
}

// database table's name: node-log-detail-info
type NodeLogDetailInfo struct {
	Id           int    `json:"id"`
	CreationTime int64  `json:"creation_time"`
	Detail       string `json:"detail"`
}

type NodeLogCache struct {
	DetailBuffer bytes.Buffer
	NodeLogPtr *NodeLog
}

type PaginatedNodeLogList struct {
	Pagination *Pagination `json:"pagination"`
	LogList    *[]NodeLog  `json:"log_list"`
}

type NodeLogDetail struct {
	Message string                `json:"message"`
	Status  NodeLogResponseStatus `json:"status"`
}

type NodeListType struct {
	Ip           string `json:"ip"`
	CreationTime int64  `json:"creation_time"`
}
