// an adapter file from board to k8s for auto-scale
package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/autoscaling/v1"
)

type autoscales struct {
	namespace string
	autoscale v1.HorizontalPodAutoscalerInterface
}

func (as *autoscales) Create(autoscale *model.AutoScale) (*model.AutoScale, error) {
	k8sHPA := types.ToK8sAutoScale(autoscale)
	k8sHPA, err := as.autoscale.Create(k8sHPA)
	if err != nil {
		logs.Error("Create auto scale of %s/%s failed. Err:%+v", autoscale.Name, as.namespace, err)
		return nil, err
	}

	return types.FromK8sAutoScale(k8sHPA), nil
}

func (as *autoscales) Update(autoscale *model.AutoScale) (*model.AutoScale, error) {
	k8sHPA, err := as.autoscale.Get(autoscale.Name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get auto scale of %s failed when updating node. Err:%+v", autoscale.Name, err)
		return nil, err
	}
	types.UpdateK8sAutoScale(k8sHPA, autoscale)

	k8sHPA, err = as.autoscale.Update(k8sHPA)
	if err != nil {
		logs.Error("Update auto scale of %s/%s failed. Err:%+v", autoscale.Name, as.namespace, err)
		return nil, err
	}
	return types.FromK8sAutoScale(k8sHPA), nil
}

func (as *autoscales) UpdateStatus(autoscale *model.AutoScale) (*model.AutoScale, error) {
	k8sHPA := types.ToK8sAutoScale(autoscale)
	k8sHPA, err := as.autoscale.UpdateStatus(k8sHPA)
	if err != nil {
		logs.Error("Create auto scale status of %s/%s failed. Err:%+v", autoscale.Name, as.namespace, err)
		return nil, err
	}

	return types.FromK8sAutoScale(k8sHPA), nil
}

func (as *autoscales) Delete(name string) error {
	err := as.autoscale.Delete(name, nil)
	if err != nil {
		logs.Error("delete auto scale of %s/%s failed. Err:%+v", name, as.namespace, err)
	}
	return err
}

func (as *autoscales) List() (*model.AutoScaleList, error) {
	asList, err := as.autoscale.List(meta_v1.ListOptions{})
	if err != nil {
		logs.Error("list auto scale failed. Err:%+v", err)
		return nil, err
	}

	return types.FromK8sAutoScaleList(asList), nil
}

func (as *autoscales) Get(name string) (*model.AutoScale, error) {
	autoscaleinstance, err := as.autoscale.Get(name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get auto scale of %s failed. Err:%+v", name, err)
		return nil, err
	}
	return types.FromK8sAutoScale(autoscaleinstance), nil
}

func NewAutoScales(namespace string, autoscale v1.HorizontalPodAutoscalerInterface) *autoscales {
	return &autoscales{
		namespace: namespace,
		autoscale: autoscale,
	}
}
