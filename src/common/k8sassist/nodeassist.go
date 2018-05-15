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

// NodeCli Interface has methods to work with Node resources in k8s-assist.
// How to:  nodeCli, err := k8sassist.NewNodes()
//          nodeInstance, err := nodeCli.Get(nodename)
type NodeCliInterface interface {
	Create(*model.Node) (*model.Node, error)
	Update(*model.Node) (*model.Node, error)
	UpdateStatus(*model.Node) (*model.Node, error)
	Delete(name string) error
	Get(name string) (*model.Node, error)
	List() (*model.NodeList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Node, err error)
}
