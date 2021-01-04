// an adapter file from board to k8s for persistent volume Claims
package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

type configMap struct {
	namespace string
	configmap v1.ConfigMapInterface
}

func (p *configMap) Create(cm *model.ConfigMap) (*model.ConfigMap, error) {
	k8sConfigMap := types.ToK8sConfigMap(cm)
	k8snewConfigMap, err := p.configmap.Create(k8sConfigMap)
	if err != nil {
		logs.Error("Create configmap of failed. %v Err:%+v", k8sConfigMap, err)
		return nil, err
	}

	return types.FromK8sConfigMap(k8snewConfigMap), nil
}

//TODO support update later
func (p *configMap) Update(cm *model.ConfigMap) (*model.ConfigMap, error) {
	k8sConfigMap, err := p.configmap.Get(cm.Name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get ConfigMap of %s failed when updating node. Err:%+v", cm.Name, err)
		return nil, err
	}
	types.UpdateK8sConfigMap(k8sConfigMap, cm)

	k8snewConfigMap, err := p.configmap.Update(k8sConfigMap)
	if err != nil {
		logs.Error("Update ConfigMap of %s failed. Err:%+v", cm.Name, err)
		return nil, err
	}
	return types.FromK8sConfigMap(k8snewConfigMap), nil
}

func (p *configMap) Delete(name string) error {
	err := p.configmap.Delete(name, nil)
	if err != nil {
		logs.Error("delete ConfigMap of %s failed. Err:%+v", name, err)
	}
	return err
}

func (p *configMap) List() (*model.ConfigMapList, error) {
	configmapList, err := p.configmap.List(meta_v1.ListOptions{})
	if err != nil {
		logs.Error("list ConfigMap failed. Err:%+v", err)
		return nil, err
	}

	return types.FromK8sConfigMapList(configmapList), nil
}

func (p *configMap) Get(name string) (*model.ConfigMap, error) {
	configmap, err := p.configmap.Get(name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get ConfigMap of %s failed. Err:%+v", name, err)
		return nil, err
	}
	return types.FromK8sConfigMap(configmap), nil
}

func NewConfigMap(namespace string, cm v1.ConfigMapInterface) *configMap {
	return &configMap{
		namespace: namespace,
		configmap: cm,
	}
}
