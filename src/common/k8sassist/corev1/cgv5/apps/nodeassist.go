package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

// nodes implements NodeInterface
type nodes struct {
	node v1.NodeInterface
}

// Support a string as label selector parameter, default null
func (n *nodes) List(args ...string) (*model.NodeList, error) {
	var listOption = metav1.ListOptions{}
	if len(args) > 0 {
		listOption.LabelSelector = args[0]
	}
	k8sNodeList, err := n.node.List(listOption)
	if err != nil {
		logs.Error("List Nodes failed. Err:%+v", err)
		return nil, err
	}

	modelNodeList := types.FromK8sNodeList(k8sNodeList)
	return modelNodeList, nil
}

func (n *nodes) Get(name string) (*model.Node, error) {
	k8sNode, err := n.node.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get Node of %s failed. Err:%+v", name, err)
		return nil, err
	}

	modelNode := types.FromK8sNode(k8sNode)
	return modelNode, nil
}

func (n *nodes) Create(node *model.Node) (*model.Node, error) {
	k8sNode := types.ToK8sNode(node)
	k8sNode, err := n.node.Create(k8sNode)
	if err != nil {
		logs.Error("Create Node of %s failed. Err:%+v", node.Name, err)
		return nil, err
	}

	modelNode := types.FromK8sNode(k8sNode)
	return modelNode, nil
}

func (n *nodes) Update(node *model.Node) (*model.Node, error) {
	k8sNode, err := n.node.Get(node.Name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get Node of %s failed when updating node. Err:%+v", node.Name, err)
		return nil, err
	}
	types.UpdateK8sNode(k8sNode, node)
	k8sNode, err = n.node.Update(k8sNode)
	if err != nil {
		logs.Error("Update Node of %s failed. Err:%+v", node.Name, err)
		return nil, err
	}

	modelNode := types.FromK8sNode(k8sNode)
	return modelNode, nil
}

func (n *nodes) UpdateStatus(node *model.Node) (*model.Node, error) {
	k8sNode, err := n.node.Get(node.Name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get Node of %s failed when updating node status. Err:%+v", node.Name, err)
		return nil, err
	}
	types.UpdateK8sNode(k8sNode, node)
	k8sNode, err = n.node.UpdateStatus(k8sNode)
	if err != nil {
		logs.Error("UpdateStatus Node of %s failed. Err:%+v", node.Name, err)
		return nil, err
	}

	modelNode := types.FromK8sNode(k8sNode)
	return modelNode, nil
}

func (n *nodes) Delete(name string) error {
	err := n.node.Delete(name, nil)
	if err != nil {
		logs.Error("Delete Node of %s failed. Err:%+v", name, err)
		return err
	}
	return nil
}

func NewNodes(node v1.NodeInterface) *nodes {
	return &nodes{
		node: node,
	}
}
