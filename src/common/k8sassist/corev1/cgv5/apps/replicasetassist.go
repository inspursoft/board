// a temp file for building and guiding
package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
)

type replicasets struct {
	namespace  string
	replicaset v1.ReplicaSetInterface
}

func (r *replicasets) Create(rs *model.ReplicaSet) (*model.ReplicaSet, error) {
	k8sRS := types.ToK8sReplicaSet(rs)
	k8sRS, err := r.replicaset.Create(k8sRS)
	if err != nil {
		logs.Error("Create ReplicaSet of %s/%s failed. Err:%+v", rs.Name, r.namespace, err)
		return nil, err
	}

	modelRS := types.FromK8sReplicaSet(k8sRS)
	return modelRS, nil
}

func (r *replicasets) Update(rs *model.ReplicaSet) (*model.ReplicaSet, error) {
	k8sRS := types.ToK8sReplicaSet(rs)
	k8sRS, err := r.replicaset.Update(k8sRS)
	if err != nil {
		logs.Error("Update ReplicaSet of %s/%s failed. Err:%+v", rs.Name, r.namespace, err)
		return nil, err
	}

	modelRS := types.FromK8sReplicaSet(k8sRS)
	return modelRS, nil
}

func (r *replicasets) UpdateStatus(rs *model.ReplicaSet) (*model.ReplicaSet, error) {
	k8sRS := types.ToK8sReplicaSet(rs)
	k8sRS, err := r.replicaset.UpdateStatus(k8sRS)
	if err != nil {
		logs.Error("UpdateStatus ReplicaSet of %s/%s failed. Err:%+v", rs.Name, r.namespace, err)
		return nil, err
	}

	modelRS := types.FromK8sReplicaSet(k8sRS)
	return modelRS, nil
}

func (r *replicasets) Delete(name string) error {
	err := r.replicaset.Delete(name, nil)
	if err != nil {
		logs.Error("Delete ReplicaSet of %s/%s failed. Err:%+v", name, r.namespace, err)
	}
	return err
}

func (r *replicasets) Get(name string) (*model.ReplicaSet, error) {
	rs, err := r.replicaset.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get ReplicaSet of %s/%s failed. Err:%+v", name, r.namespace, err)
		return nil, err
	}

	modelRS := types.FromK8sReplicaSet(rs)
	return modelRS, nil
}

func (r *replicasets) List(opts model.ListOptions) (*model.ReplicaSetList, error) {
	rsList, err := r.replicaset.List(metav1.ListOptions{})
	if err != nil {
		logs.Error("List ReplicaSets failed. Err:%+v", err)
		return nil, err
	}

	modelRSList := types.FromK8sReplicaSetList(rsList)
	return modelRSList, nil
}

func NewReplicaSets(namespace string, replicaset v1.ReplicaSetInterface) *replicasets {
	return &replicasets{
		namespace:  namespace,
		replicaset: replicaset,
	}
}
