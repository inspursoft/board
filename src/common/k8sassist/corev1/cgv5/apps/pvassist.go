// an adapter file from board to k8s for persistent volumes
package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

type persistentvolume struct {
	namespace string
	pv        v1.PersistentVolumeInterface
}

func (p *persistentvolume) Create(pv *model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error) {
	k8sPV := types.ToK8sPV(pv)
	k8sPV, err := p.pv.Create(k8sPV)
	if err != nil {
		logs.Error("Create pv of %s failed. Err:%+v", pv.Name, err)
		return nil, err
	}

	return types.FromK8sPV(k8sPV), nil
}

//TODO support update later
func (p *persistentvolume) Update(pv *model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error) {
	k8sPV, err := p.pv.Get(pv.Name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get PV of %s failed when updating node. Err:%+v", pv.Name, err)
		return nil, err
	}
	types.UpdateK8sPV(k8sPV, pv)

	k8sPV, err = p.pv.Update(k8sPV)
	if err != nil {
		logs.Error("Update PV of %s failed. Err:%+v", pv.Name, err)
		return nil, err
	}
	return types.FromK8sPV(k8sPV), nil
}

func (p *persistentvolume) UpdateStatus(pv *model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error) {
	k8sPV := types.ToK8sPV(pv)
	k8sPV, err := p.pv.UpdateStatus(k8sPV)
	if err != nil {
		logs.Error("Create PV status of %s failed. Err:%+v", pv.Name, err)
		return nil, err
	}

	return types.FromK8sPV(k8sPV), nil
}

func (p *persistentvolume) Delete(name string) error {
	err := p.pv.Delete(name, nil)
	if err != nil {
		logs.Error("delete PV of %s failed. Err:%+v", name, err)
	}
	return err
}

func (p *persistentvolume) List() (*model.PersistentVolumeList, error) {
	pvList, err := p.pv.List(meta_v1.ListOptions{})
	if err != nil {
		logs.Error("list PV failed. Err:%+v", err)
		return nil, err
	}

	return types.FromK8sPVList(pvList), nil
}

func (p *persistentvolume) Get(name string) (*model.PersistentVolumeK8scli, error) {
	pvinstance, err := p.pv.Get(name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get PV of %s failed. Err:%+v", name, err)
		return nil, err
	}
	return types.FromK8sPV(pvinstance), nil
}

func NewPersistentVolume(pv v1.PersistentVolumeInterface) *persistentvolume {
	return &persistentvolume{
		pv: pv,
	}
}
