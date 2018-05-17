package apps

import (
	"git/inspursoft/board/src/common/model"

	"k8s.io/client-go/kubernetes/typed/core/v1"
)

// nodes implements NodeInterface
type nodes struct {
	node v1.NodeInterface
}

func (n *nodes) List() (*model.NodeList, error) {
	return &model.NodeList{}, nil
}

func (n *nodes) Get(nodename string) (*model.Node, error) {
	return &model.Node{}, nil
}

func (n *nodes) Create(newnode *model.Node) (*model.Node, error) {
	return newnode, nil
}

func (n *nodes) Update(newnode *model.Node) (*model.Node, error) {
	return newnode, nil
}

func (n *nodes) UpdateStatus(*model.Node) (*model.Node, error) {
	return nil, nil
}

func (n *nodes) Delete(s string) error {
	return nil
}

func NewNodes(node v1.NodeInterface) *nodes {
	return &nodes{
		node: node,
	}
}
