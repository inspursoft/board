// a temp file for building and guiding
package k8sassist

import (
	//api "k8s.io/client-go/pkg/api"
	//v1 "k8s.io/client-go/pkg/api/v1"
	//watch "k8s.io/client-go/pkg/watch"
	//rest "k8s.io/client-go/rest"
	"git/inspursoft/board/src/common/model"
)

type replicasets struct {
	ns string
}

func (d *replicasets) Create(*model.ReplicaSet) (*model.ReplicaSet, error) {
	return nil, nil
}

func (d *replicasets) Update(*model.ReplicaSet) (*model.ReplicaSet, error) {
	return nil, nil
}

func (d *replicasets) UpdateStatus(*model.ReplicaSet) (*model.ReplicaSet, error) {
	return nil, nil
}

func (d *replicasets) Delete(name string) error {
	return nil
}

func (d *replicasets) Get(name string) (*model.ReplicaSet, error) {
	return nil, nil
}

func (d *replicasets) List(opts model.ListOptions) (*model.ReplicaSetList, error) {
	return nil, nil
}

var _ ReplicaSetCliInterface = &replicasets{}

// newNodes returns a Nodes
func NewReplicaSets(namespace string) (*replicasets, error) {
	return &replicasets{ns: namespace}, nil
}

// ReplicaSetInterface has methods to work with ReplicaSet resources.
type ReplicaSetCliInterface interface {
	Create(*model.ReplicaSet) (*model.ReplicaSet, error)
	Update(*model.ReplicaSet) (*model.ReplicaSet, error)
	UpdateStatus(*model.ReplicaSet) (*model.ReplicaSet, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.ReplicaSet, error)
	List(opts model.ListOptions) (*model.ReplicaSetList, error)
	//Watch(opts v1.ListOptions) (watch.Interface, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1beta1.ReplicaSet, err error)
	//ReplicaSetExpansion
}
