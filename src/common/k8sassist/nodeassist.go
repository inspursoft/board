// a temp file for building and guiding
package k8sassist

import (
	//api "k8s.io/client-go/pkg/api"
	//v1 "k8s.io/client-go/pkg/api/v1"
	//watch "k8s.io/client-go/pkg/watch"
	//rest "k8s.io/client-go/rest"
	"git/inspursoft/board/src/common/model"
)

// nodes implements NodeInterface
type nodes struct {
	client NodeCliInterface
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

// newNodes returns a Nodes
func NewNodes() (*nodes, error) {
	return &nodes{}, nil
}
