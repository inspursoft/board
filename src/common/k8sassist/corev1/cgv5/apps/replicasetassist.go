// a temp file for building and guiding
package apps

import (
	"git/inspursoft/board/src/common/model"

	"k8s.io/client-go/kubernetes/typed/apps/v1beta2"
)

type replicasets struct {
	namespace  string
	replicaset v1beta2.ReplicaSetInterface
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

func NewReplicaSets(namespace string, replicaset v1beta2.ReplicaSetInterface) *replicasets {
	return &replicasets{
		namespace:  namespace,
		replicaset: replicaset,
	}
}
