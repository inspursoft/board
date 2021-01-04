// an adapter file from board to k8s for persistent volume Claims
package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

type persistentvolumeclaim struct {
	namespace string
	pvc       v1.PersistentVolumeClaimInterface
}

func (p *persistentvolumeclaim) Create(pvc *model.PersistentVolumeClaimK8scli) (*model.PersistentVolumeClaimK8scli, error) {
	k8sPVC := types.ToK8sPVC(pvc)
	k8sPVC, err := p.pvc.Create(k8sPVC)
	if err != nil {
		logs.Error("Create pvc of failed. %v Err:%+v", k8sPVC, err)
		return nil, err
	}

	return types.FromK8sPVC(k8sPVC), nil
}

//TODO support update later
func (p *persistentvolumeclaim) Update(pvc *model.PersistentVolumeClaimK8scli) (*model.PersistentVolumeClaimK8scli, error) {
	k8sPVC, err := p.pvc.Get(pvc.Name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get PVC of %s failed when updating node. Err:%+v", pvc.Name, err)
		return nil, err
	}
	types.UpdateK8sPVC(k8sPVC, pvc)

	k8sPVC, err = p.pvc.Update(k8sPVC)
	if err != nil {
		logs.Error("Update PV of %s failed. Err:%+v", pvc.Name, err)
		return nil, err
	}
	return types.FromK8sPVC(k8sPVC), nil
}

func (p *persistentvolumeclaim) UpdateStatus(pvc *model.PersistentVolumeClaimK8scli) (*model.PersistentVolumeClaimK8scli, error) {
	k8sPVC := types.ToK8sPVC(pvc)
	k8sPVC, err := p.pvc.UpdateStatus(k8sPVC)
	if err != nil {
		logs.Error("Create PV status of %s failed. Err:%+v", pvc.Name, err)
		return nil, err
	}

	return types.FromK8sPVC(k8sPVC), nil
}

func (p *persistentvolumeclaim) Delete(name string) error {
	err := p.pvc.Delete(name, nil)
	if err != nil {
		logs.Error("delete PVC of %s failed. Err:%+v", name, err)
	}
	return err
}

func (p *persistentvolumeclaim) List() (*model.PersistentVolumeClaimList, error) {
	pvcList, err := p.pvc.List(meta_v1.ListOptions{})
	if err != nil {
		logs.Error("list PVC failed. Err:%+v", err)
		return nil, err
	}

	return types.FromK8sPVCList(pvcList), nil
}

func (p *persistentvolumeclaim) Get(name string) (*model.PersistentVolumeClaimK8scli, error) {
	pvcinstance, err := p.pvc.Get(name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get PVC of %s failed. Err:%+v", name, err)
		return nil, err
	}
	return types.FromK8sPVC(pvcinstance), nil
}

func NewPersistentVolumeClaim(namespace string, pvc v1.PersistentVolumeClaimInterface) *persistentvolumeclaim {
	return &persistentvolumeclaim{
		namespace: namespace,
		pvc:       pvc,
	}
}
