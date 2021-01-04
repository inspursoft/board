package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type daemonsets struct {
	namespace string
	daemonset v1.DaemonSetInterface
}

func (d *daemonsets) processDaemonSetHandler(daemonset *model.DaemonSet, handler func(*appsv1.DaemonSet) (*appsv1.DaemonSet, error)) (customModel *model.DaemonSet, primitiveData []byte, err error) {
	k8sDaemonSet := types.ToK8sDaemonSet(daemonset)
	handledDs, err := handler(k8sDaemonSet)
	if err != nil {
		logs.Error("Failed to handle DaemonSet of %s/%s failed. Err:%+v", handledDs.Name, handledDs.Namespace, err)
		return nil, nil, err
	}
	customModel = types.FromK8sDaemonSet(handledDs)
	primitiveData, err = yaml.Marshal(types.GenerateDaemonSetConfig(handledDs))
	if err != nil {
		logs.Error("Failed to marshal primitive from daemonset config, error: %+v", err)
		return
	}
	return
}

func (d *daemonsets) Create(daemonset *model.DaemonSet) (*model.DaemonSet, []byte, error) {
	return d.processDaemonSetHandler(daemonset, d.daemonset.Create)
}

func (d *daemonsets) Update(daemonset *model.DaemonSet) (*model.DaemonSet, []byte, error) {
	return d.processDaemonSetHandler(daemonset, d.daemonset.Update)
}

func (d *daemonsets) UpdateStatus(daemonset *model.DaemonSet) (*model.DaemonSet, []byte, error) {
	return d.processDaemonSetHandler(daemonset, d.daemonset.UpdateStatus)
}

func (d *daemonsets) Delete(name string) error {
	deletePolicy := types.DeletePropagationForeground
	err := d.daemonset.Delete(name, &types.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logs.Error("Delete daemonset of %s/%s failed. Err:%+v", name, d.namespace, err)
	}
	return err
}

func (d *daemonsets) Get(name string) (*model.DaemonSet, []byte, error) {
	daemonset, err := d.daemonset.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get daemonset of %s/%s failed. Err:%+v", name, d.namespace, err)
		return nil, nil, err
	}
	daemonsetConfig := types.GenerateDaemonSetConfig(daemonset)
	daemonsetfileInfo, err := yaml.Marshal(daemonsetConfig)
	if err != nil {
		logs.Error("Marshal daemonset failed, error: %v", err)
		return nil, nil, err
	}
	modelDs := types.FromK8sDaemonSet(daemonset)
	return modelDs, daemonsetfileInfo, nil
}

func (d *daemonsets) List(opts model.ListOptions) (*model.DaemonSetList, error) {
	daemonsetList, err := d.daemonset.List(types.ToK8sListOptions(opts))
	if err != nil {
		logs.Error("List daemonsets failed. Err:%+v", err)
		return nil, err
	}
	modelDSList := types.FromK8sDaemonSetList(daemonsetList)
	return modelDSList, nil
}

// NewDaemonSets is to create a daemonset adapter
func NewDaemonSets(namespace string, daemonset v1.DaemonSetInterface) *daemonsets {
	return &daemonsets{
		namespace: namespace,
		daemonset: daemonset,
	}
}
