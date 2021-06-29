package nodeModel

import "bytes"

const BasePath = "/data/adminserver/ansible_k8s/"
const AddNodeYamlFile = "addnode"
const RemoveNodeYamlFile = "uninstallnode"
const NodeHostsFile = "addNode"
const LogFileDir = "log"
const HostFileDir = "hosts"
const PreEnvDir = "/data/pre-env"

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
	MasterIp       string `json:"master_ip"`
	NodeIp         string `json:"node_ip"`
	NodePassword   string `json:"node_password"`
	HostUsername   string `json:"host_username"`
	HostPassword   string `json:"host_password"`
	MasterPassword string `json:"master_password"`
}

type PreparationData struct {
	HostIp   string `json:"host_ip"`
	MasterIp string `json:"master_ip"`
}

type ContainerEnv struct {
	NodeIp         string `json:"node_ip"`
	NodePassword   string `json:"node_password"`
	HostIp         string `json:"host_ip"`
	HostUserName   string `json:"host_user_name"`
	HostPassword   string `json:"host_password"`
	MasterIp       string `json:"master_ip"`
	MasterPassword string `json:"master_password"`
	InstallFile    string `json:"install_file"`
	HostFile       string `json:"host_file"`
	LogId          int64  `json:"log_id"`
	LogTimestamp   int64  `json:"log_timestamp"`
}

type UpdateNodeLog struct {
	LogId       int    `json:"log_id"`
	Ip          string `json:"ip"`
	InstallFile string `json:"install_file"`
	LogFile     string `json:"log_file"`
	ExitCode    int    `json:"success"`
}

type NodeLogDetailArray = []NodeLogDetail

// database table's name: node-status
type NodeStatus struct {
	Id           int    `json:"id"`
	Ip           string `json:"ip"`
	CreationTime int64  `json:"creation_time"`
}

type ApiServerNodeListResult struct {
	NodeName   string            `json:"node_name"`
	NodeIP     string            `json:"node_ip"`
	NodeType   string            `json:"node_type"`
	Status     int               `json:"status"`
	CreateTime int64             `json:"create_time"`
	Labels     map[string]string `json:"labels"`
}

type NodeListResponse struct {
	Ip           string `json:"ip"`
	NodeName     string `json:"node_name"`
	CreationTime int64  `json:"creation_time"`
	LogTime      int64  `json:"log_time"`
	Status       int    `json:"status"`
	IsMaster     bool   `json:"is_master"`
	IsEdge       bool   `json:"is_edge"`
	Origin       int    `json:"origin"`
}

// database table's name: node-log
type NodeLog struct {
	Id           int        `json:"id"`
	Ip           string     `json:"ip"`
	LogType      ActionType `json:"log_type"`
	Success      bool       `json:"success"`
	Completed    bool       `json:"completed"`
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
	NodeLogPtr   *NodeLog
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
