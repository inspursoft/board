package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
	v2 "github.com/google/cadvisor/info/v2"
	//modelK8s "k8s.io/client-go/pkg/api/v1"
	//"golang.org/x/crypto/ssh"
)

type NodeStatus int

const (
	_ NodeStatus = iota
	Running
	Unschedulable
	Unknown
	AutonomousOffline
)

const (
	K8sLabel         = "kubernetes.io"
	K8sNamespaces    = "kube-system cadvisor"
	K8sMasterLabel   = "node-role.kubernetes.io/master"
	K8sEdgeNodeLabel = "node-role.kubernetes.io/edge"
	NodeTypeMaster   = "master"
	NodeTypeEdge     = "edge"
	NodeTypeNode     = "node"
)

type NodeListResult struct {
	NodeName   string            `json:"node_name"`
	NodeIP     string            `json:"node_ip"`
	Status     NodeStatus        `json:"status"`
	CreateTime int64             `json:"create_time"`
	Labels     map[string]string `json:"labels"`
	NodeType   string            `json:"node_type"`
}

type NodeInfo struct {
	NodeName      string  `json:"node_name" orm:"column(node_name)"`
	NodeIP        string  `json:"node_ip" orm:"column(node_ip)"`
	CreateTime    int64   `json:"create_time" orm:"column(create_time)"`
	CPUUsage      float32 `json:"cpu_usage" orm:"column(cpu_usage)"`
	NumberCPUCore int     `json:"numbers_cpu_core" orm:"column(numbers_cpu_core)"`
	MemoryUsage   float32 `json:"memory_usage" orm:"column(memory_usage)"`
	MemorySize    int     `json:"memory_size" orm:"column(memory_size)"`
	StorageTotal  uint64  `json:"storage_total" orm:"column(storage_total)"`
	StorageUse    uint64  `json:"storage_use" orm:"column(storage_usage)"`
}

func GetNodes() (nodes []NodeInfo, err error) {
	defer func() { recover() }()
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nodecli := k8sclient.AppV1().Node()

	Node, err := nodecli.List()
	if err != nil {
		logs.Error("Failed to get Node List")
		return
	}
	for _, v := range Node.Items {
		var mlimit int
		var CPUCore int
		mlimit, err = strconv.Atoi(fmt.Sprintf("%s", v.Status.Capacity["memory"]))
		if err != nil {
			logs.Error("Failed to get the number of memory of %s Node.", v.Status.Addresses[1].Address)
		}
		CPUCore, err = strconv.Atoi(fmt.Sprintf("%s", v.Status.Capacity["cpu"]))
		if err != nil {
			logs.Error("Failed to get the number of CPU core of %s Node.", v.Status.Addresses[1].Address)
		}
		time := v.CreationTimestamp.Unix()
		var ps []v2.ProcessInfo
		nodeIP := getNodeAddress(v, "InternalIP")
		getFromRequest("http://"+nodeIP+":4194/api/v2.0/ps/", &ps)
		var c, m float32
		for _, v := range ps {
			c += v.PercentCpu
			m += v.PercentMemory
		}
		cpu := c
		mem := m
		var fs []v2.MachineFsStats
		getFromRequest("http://"+nodeIP+":4194/api/v2.0/storage", &fs)
		var capacity uint64
		var use uint64
		for _, v := range fs {
			capacity += *v.Capacity
			use += *v.Usage
		}
		nodes = append(nodes, NodeInfo{
			NodeName:      getNodeAddress(v, "Hostname"),
			NodeIP:        nodeIP,
			CreateTime:    time,
			CPUUsage:      cpu,
			MemoryUsage:   mem,
			MemorySize:    mlimit,
			NumberCPUCore: CPUCore,
			StorageTotal:  capacity,
			StorageUse:    use,
		})
	}
	return
}

// Get a node struct by the nodename
func GetNode(nodeName string) (node NodeInfo, err error) {
	nodes, err := GetNodes()
	if err != nil {
		logs.Error("Failed to get Node information.")
		return
	}
	for _, node = range nodes {
		if strings.EqualFold(node.NodeName, nodeName) {
			return
		}
	}
	return NodeInfo{}, nil
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
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nodecli := k8sclient.AppV1().Node()

	Node, err := nodecli.List()
	if err != nil {
		logs.Error("Failed to get Node List")
		return
	}

	for _, v := range Node.Items {
		nodetype := getNodeType(v)
		if nodetype != NodeTypeEdge {
			res = append(res, NodeListResult{
				NodeName:   getNodeAddress(v, "Hostname"),
				NodeIP:     getNodeAddress(v, "InternalIP"),
				CreateTime: v.CreationTimestamp.Unix(),
				Labels:     v.ObjectMeta.Labels,
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
				}(),
				NodeType: nodetype})
		} else {
			var nodeitem NodeListResult
			if name, ok := v.Labels["kubernetes.io/hostname"]; ok {
				nodeitem.NodeName = name
			} else {
				nodeitem.NodeName = v.NodeIP
			}
			nodeitem.NodeIP = getNodeAddress(v, "InternalIP")
			nodeitem.CreateTime = v.CreationTimestamp.Unix()
			nodeitem.Labels = v.ObjectMeta.Labels
			nodeitem.Status = Unknown
			nodeitem.NodeType = nodetype

			// update status
			for _, cond := range v.Status.Conditions {
				if strings.EqualFold(string(cond.Type), "Ready") && cond.Status == model.ConditionTrue {
					nodeitem.Status = Running
				}
			}
			if nodeitem.Status != Running {
				//TODO Ping the edgenode is not the only condition for AutonomousOffline
				status, err := utils.PingIPAddr(v.NodeIP)
				if err != nil {
					logs.Error("Failed to ping IPAddr: %s, error: %+v", v.NodeIP, err)
				} else if !status {
					logs.Debug("The edge node %s is in AutonomousOffline", v.NodeIP)
					nodeitem.Status = AutonomousOffline
				}
			}
			res = append(res, nodeitem)
		}
	}
	return
}

//Get the node type based on labels
func getNodeType(node model.Node) string {
	if _, ok := node.ObjectMeta.Labels[K8sMasterLabel]; ok {
		return NodeTypeMaster
	}
	if _, ok := node.ObjectMeta.Labels[K8sEdgeNodeLabel]; ok {
		return NodeTypeEdge
	}
	return NodeTypeNode
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

// Check Nodegroup in DB, existing true
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
	config.KubeConfigPath = kubeConfigPath()
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

// Get the groups that this node belong to
func GetGroupOfNode(nodeName string) ([]string, error) {
	var groups []string
	//nInterface, err := k8sassist.NewNodes()
	//if err != nil {
	//	logs.Error("Failed to get node client interface")
	//	return nil, err
	//}

	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nInterface := k8sclient.AppV1().Node()

	nNode, err := nInterface.Get(nodeName)
	if err != nil {
		logs.Error("Failed to get K8s node")
		return nil, err
	}
	for key, _ := range nNode.ObjectMeta.Labels {
		if !strings.Contains(key, K8sLabel) {
			exist, err := NodeGroupExists(key)
			if err != nil {
				logs.Error("Failed to get nodegroup")
				return nil, err
			}
			if exist {
				groups = append(groups, key)
			}
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
	config.KubeConfigPath = kubeConfigPath()
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
				KubeConfigPath: kubeConfigPath(),
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
		KubeConfigPath: kubeConfigPath(),
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

// Check the node name existing in cluster, existing return ture
func NodeExists(nodeName string) (bool, error) {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Node()

	nodeList, err := n.List()
	if err != nil {
		logs.Error("Failed to check node list in cluster", nodeName)
		return false, err
	}

	for _, nd := range (*nodeList).Items {
		if nodeName == nd.Name {
			logs.Info("Nodename existing %+v", nodeName)
			return true, nil
		}
	}
	return false, nil
}

// Get a node in kubernetes cluster
func GetNodebyName(nodeName string) (*model.Node, error) {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Node()

	node, err := n.Get(nodeName)
	if err != nil {
		logs.Error("Failed to get node: %s, error: %+v", nodeName, err)
		return nil, err
	}
	logs.Info("Node in K8s %+v", node)
	return node, nil

}

// Get a node in kubernetes cluster
func GetNodebyIP(nodeIP string) (*model.Node, error) {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Node()
	nodeList, err := n.List()
	if err != nil {
		logs.Error("Failed to get nodelist: %s, error: %+v", nodeIP, err)
		return nil, err
	}
	for _, nodeitem := range nodeList.Items {

		if nodeitem.NodeIP == nodeIP {
			logs.Info("Node in K8s %+v", nodeitem)
			return &nodeitem, nil
		}
	}
	return nil, nil
}

// Create a node in kubernetes cluster
func CreateNode(node model.NodeCli) (*model.Node, error) {

	nExists, err := NodeExists(node.NodeName)
	if err != nil {
		return nil, err
	}
	if nExists {
		logs.Info("Node name %s already exists in cluster.", node.NodeName)
		return nil, nil
	}

	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Node()

	var nodek8s model.Node
	nodek8s.ObjectMeta.Name = node.NodeName
	nodek8s.ObjectMeta.Labels = node.Labels
	nodek8s.Taints = node.Taints

	newnode, err := n.Create(&nodek8s)
	if err != nil {
		logs.Error("Failed to create node: %s, error: %+v", node.NodeName, err)
		return nil, err
	}
	logs.Info("New Node in K8s %+v", newnode)
	return newnode, nil

}

// Delete a node from kubernetes cluster, should do clean work first before this func
func DeleteNode(nodeName string) (bool, error) {
	nExists, err := NodeExists(nodeName)
	if err != nil {
		return false, err
	}
	if !nExists {
		logs.Info("Name %s not exists in cluster.", nodeName)
		return false, nil
	}

	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Node()

	err = n.Delete(nodeName)
	if err != nil {
		logs.Error("Failed to delete node %s", nodeName)
		return false, err
	}
	return true, nil
}

func getNodeAddress(v model.Node, t string) string {
	for _, addr := range v.Status.Addresses {
		if string(addr.Type) == t {
			return addr.Address
		}
	}

	logs.Warning("The value is null when get the field of %s in node", t)
	return ""
}

// Get a node control status
func GetNodeControlStatus(nodeName string) (*model.NodeControlStatus, error) {
	var nodecontrol model.NodeControlStatus
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nInterface := k8sclient.AppV1().Node()

	nNode, err := nInterface.Get(nodeName)
	if err != nil {
		logs.Error("Failed to get K8s node")
		return nil, err
	}
	nodecontrol.NodeName = nNode.Name
	nodecontrol.NodeIP = nNode.NodeIP
	nodecontrol.NodePhase = string(nNode.Status.Phase)
	nodecontrol.NodeUnschedule = nNode.Unschedulable
	nodecontrol.NodeDeletable = true

	// Phase is deprecated, if null, use Status, fix me
	if nodecontrol.NodePhase == "" {
		nodecontrol.NodePhase = func() string {
			for _, cond := range nNode.Status.Conditions {
				if strings.EqualFold(string(cond.Type), "Ready") && cond.Status == model.ConditionTrue {
					return "Running"
				}
			}
			return "Unknown"
		}()
	}

	//Check master, add a check for master, undeletable
	if _, ok := nNode.ObjectMeta.Labels[K8sMasterLabel]; ok {
		nodecontrol.NodeType = NodeTypeMaster
		nodecontrol.NodeDeletable = false
		logs.Debug("Master Node %s", nodecontrol.NodeName)
		return &nodecontrol, nil
	}

	// Get service instances
	// si, err := GetNodeServiceInstances(nodeName)
	// if err != nil {
	// 	logs.Error("Failed to get K8s node service instances")
	// 	return nil, err
	// }
	// if si != nil {
	// 	nodecontrol.Service_Instances = *si
	// 	logs.Debug("Node service instances: %v", nodecontrol.Service_Instances)
	// }
	pInterface := k8sclient.AppV1().Pod(model.NamespaceAll)
	podList, err := pInterface.List(model.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName)})
	if err != nil {
		logs.Error("Failed to get K8s pods")
		return nil, err
	}
	for _, podinstance := range podList.Items {
		var instance model.ServiceInstance
		instance.ProjectName = podinstance.Namespace
		instance.ServiceInstanceName = podinstance.Name
		nodecontrol.Service_Instances = append(nodecontrol.Service_Instances, instance)

		//TODO Need check the deletable by pod list information, owner reference
		if !strings.Contains(K8sNamespaces, podinstance.Namespace) {
			nodecontrol.NodeDeletable = false
		}
	}
	return &nodecontrol, nil
}

// Get a node service instances
func GetNodeServiceInstances(nodeName string) (*[]model.ServiceInstance, error) {
	var instances []model.ServiceInstance
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	pInterface := k8sclient.AppV1().Pod(model.NamespaceAll)
	podList, err := pInterface.List(model.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName)})
	if err != nil {
		logs.Error("Failed to get K8s pods")
		return nil, err
	}
	for _, podinstance := range podList.Items {
		var instance model.ServiceInstance
		instance.ProjectName = podinstance.Namespace
		instance.ServiceInstanceName = podinstance.Name
		instances = append(instances, instance)
	}
	return &instances, err
}

//Drain node serivce instances by adminserver
func DrainNodeServiceInstanceByAdminServer(nodeName string) error {
	//TODO call adminserver do kubectl drain
	return nil
}

//Drain node serivce instances
func DrainNodeServiceInstance(nodeName string) error {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	pInterface := k8sclient.AppV1().Pod(model.NamespaceAll)
	podList, err := pInterface.List(model.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName)})
	if err != nil {
		logs.Error("Failed to get K8s pods")
		return err
	}
	for _, podinstance := range podList.Items {

		//TODO Need to delete the pod based on its owner reference
		logs.Debug("pod %s, kind %v", podinstance.Name, podinstance.ObjectMeta.Labels)

		//TODO Need to support pod evict

		//If not support evict, use pod delete simply
		if strings.Contains(K8sNamespaces, podinstance.Namespace) {
			//Skip the pods of k8s self
			continue
		}
		podcli := k8sclient.AppV1().Pod(podinstance.Namespace)
		err = podcli.Delete(podinstance.Name)
		if err != nil {
			logs.Error("Failed to delete pod %s %v", podinstance.Name, err)
			//TODO fix me, whether continue to delete the rest
			//return err
		} else {
			logs.Debug("pod %s deleted", podinstance.Name)
		}
	}
	return nil
}

// Create an edge node in kubernetes cluster
func CreateEdgeNode(edgenode model.EdgeNodeCli) (*model.Node, error) {
	// This is to control the adding of edgenode manually

	// TODO run ansible docker script to add an edge node by admin api
	logs.Debug("To install edgenode by ansible: %v", edgenode)

	// Add in k8s
	var node model.NodeCli
	node.NodeName = edgenode.NodeName
	node.Labels = make(map[string]string)
	node.Labels[K8sEdgeNodeLabel] = ""
	node.Labels["name"] = edgenode.NodeName
	node.Labels["kubernetes.io/hostname"] = edgenode.NodeName
	node.Taints = append(node.Taints, model.Taint{Key: "edge", Value: node.NodeName, Effect: model.TaintEffectNoExecute})
	if edgenode.RegistryMode == "auto" {
		node.Labels["edge"] = "true"
	}
	return CreateNode(node)
}

// check the edge node hostname config
func CheckEdgeHostname(edgenode model.EdgeNodeCli) (bool, error) {
	// var sshUser = "root"
	// var sshPort = 22

	// sshHandler, err := NewSecureShell(edgenode.NodeIP, sshPort, sshUser, edgenode.Password)
	// if err != nil {
	// 	logs.Debug("Failed to dail edgenode %s %v", edgenode.NodeIP, err)
	// 	return false, err
	// }
	// defer sshHandler.client.Close()

	// session, err := sshHandler.client.NewSession()
	// if err != nil {
	// 	logs.Debug("Failed to get session edgenode %s %v", edgenode.NodeIP, err)
	// 	return false, err
	// }
	// defer session.Close()

	// combo, err := session.CombinedOutput("hostname")
	// if err != nil {
	// 	logs.Debug("Failed to get hostname edgenode %s %v", edgenode.NodeIP, err)
	// 	return false, err
	// }
	// sshhostname := strings.Replace(string(combo), "\n", "", -1)
	// logs.Debug("Edge hostname:", sshhostname)

	sshhostname, err := GetEdgeHostname(edgenode.NodeIP, edgenode.Password)
	if err != nil {
		logs.Debug("Failed to get Edge hostname %s", edgenode.NodeIP)
		return false, err
	}

	//TODO Check the hostname config in edge yaml

	if edgenode.NodeName != sshhostname {
		logs.Debug("Failed config %s edgenode %s", edgenode.NodeName, sshhostname)
		return false, errors.New("edge node hostname mismatched")
	}
	return true, nil
}

//Get the edge hostname by IP
func GetEdgeHostname(edgeIP string, edgePassword string) (string, error) {
	var sshUser = "root"
	var sshPort = 22

	if edgeIP == "" || edgePassword == "" {
		logs.Debug("IP address or Password invalid")
		return "", errors.New("IP address or Password invalid")
	}

	sshHandler, err := NewSecureShell(edgeIP, sshPort, sshUser, edgePassword)
	if err != nil {
		logs.Debug("Failed to dail edgenode %s %v", edgeIP, err)
		return "", err
	}
	defer sshHandler.client.Close()

	session, err := sshHandler.client.NewSession()
	if err != nil {
		logs.Debug("Failed to get session edgenode %s %v", edgeIP, err)
		return "", err
	}
	defer session.Close()

	combo, err := session.CombinedOutput("hostname")
	if err != nil {
		logs.Debug("Failed to get hostname edgenode %s %v", edgeIP, err)
		return "", err
	}
	sshhostname := strings.Replace(string(combo), "\n", "", -1)
	logs.Debug("Edge node %s hostname: %s", edgeIP, sshhostname)
	return sshhostname, nil
}

//Get Node List by labelselector
func GetNodeListbyLabel(labelselector string) (*model.NodeList, error) {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Node()
	nodeList, err := n.List(labelselector)
	if err != nil {
		logs.Error("Failed to get nodelist: %s, error: %+v", labelselector, err)
		return nil, err
	}
	return nodeList, nil
}
