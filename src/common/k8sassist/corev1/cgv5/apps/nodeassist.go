package apps

import (
	"git/inspursoft/board/src/common/model"
)

// nodes implements NodeInterface
type nodes struct {
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

func (n *nodes) Delete(s string) error {
	return nil
}
