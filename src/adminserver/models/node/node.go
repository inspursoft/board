package node

import (
	"time"
)

const AddDeleteNodeJsonFileName = "/root/ansible_k8s/addNode.json"
const AddNodeYamlFileName = "/root/ansible_k8s/addNode.yml"
const DeleteNodeYamlFileName = "/root/ansible_k8s/uninstallnode.yml"
const AddDeleteNodeFileName = "/root/ansible_k8s/addNode"

type WsNodeResponseStatus int
type ActionType int

const (
	WsNodeResponseUnKnown WsNodeResponseStatus = 0
	WsNodeResponseStart   WsNodeResponseStatus = 1
	WsNodeResponseNormal  WsNodeResponseStatus = 2
	WsNodeResponseError   WsNodeResponseStatus = 3
	WsNodeResponseWarning WsNodeResponseStatus = 4
	WsNodeResponseSuccess WsNodeResponseStatus = 5
	WsNodeResponseFailed  WsNodeResponseStatus = 6
)

const (
	ActionTypeAddNode    ActionType = 0
	ActionTypeDeleteNode ActionType = 1
)

type WsNodeResponse struct {
	Message string               `json:"message"`
	Status  WsNodeResponseStatus `json:"status"`
}

type NodeListType struct {
	Ip           string    `json:"ip"`
	CreationTime time.Time `json:"creation_time"`
}
