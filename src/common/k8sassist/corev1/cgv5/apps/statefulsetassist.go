package apps

import (
	"encoding/json"
	//	"errors"
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"io"
	"io/ioutil"

	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type statefulsets struct {
	namespace   string
	statefulset v1.StatefulSetInterface
}

//var (
//	namespacesErr = errors.New("Namespace value isn't consistent with project name")
//)

func (d *statefulsets) processStatefulSetHandler(statefuleset *model.StatefulSet, handler func(*appsv1.StatefulSet) (*appsv1.StatefulSet, error)) (customModel *model.StatefulSet, primitiveData []byte, err error) {
	k8sStatefulSet := types.ToK8sStatefulSet(statefuleset)
	handledSta, err := handler(k8sStatefulSet)
	if err != nil {
		logs.Error("Failed to handle StatefulSet of %s/%s failed. Err:%+v", handledSta.Name, handledSta.Namespace, err)
		return nil, nil, err
	}
	customModel = types.FromK8sStatefulSet(handledSta)
	primitiveData, err = yaml.Marshal(types.GenerateStatefulSetConfig(handledSta))
	if err != nil {
		logs.Error("Failed to marshal primitive from statefulset config, error: %+v", err)
		return
	}
	return
}

func (d *statefulsets) Create(statefuleset *model.StatefulSet) (*model.StatefulSet, []byte, error) {
	return d.processStatefulSetHandler(statefuleset, d.statefulset.Create)
}

func (d *statefulsets) Update(statefuleset *model.StatefulSet) (*model.StatefulSet, []byte, error) {
	return d.processStatefulSetHandler(statefuleset, d.statefulset.Update)
}

func (d *statefulsets) UpdateStatus(statefuleset *model.StatefulSet) (*model.StatefulSet, []byte, error) {
	return d.processStatefulSetHandler(statefuleset, d.statefulset.UpdateStatus)
}

func (d *statefulsets) Delete(name string) error {
	deletePolicy := types.DeletePropagationForeground
	err := d.statefulset.Delete(name, &types.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logs.Error("Delete statefulset of %s/%s failed. Err:%+v", name, d.namespace, err)
	}
	return err
}

func (d *statefulsets) Get(name string) (*model.StatefulSet, []byte, error) {
	statefulset, err := d.statefulset.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get statefulset of %s/%s failed. Err:%+v", name, d.namespace, err)
		return nil, nil, err
	}
	statefulsetConfig := types.GenerateStatefulSetConfig(statefulset)
	statefulsetfileInfo, err := yaml.Marshal(statefulsetConfig)
	if err != nil {
		logs.Error("Marshal statefulset failed, error: %v", err)
		return nil, nil, err
	}
	modelSta := types.FromK8sStatefulSet(statefulset)
	return modelSta, statefulsetfileInfo, nil
}

func (d *statefulsets) List() (*model.StatefulSetList, error) {
	statefulsetList, err := d.statefulset.List(metav1.ListOptions{})
	if err != nil {
		logs.Error("List statefulsets failed. Err:%+v", err)
		return nil, err
	}
	modelStaList := types.FromK8sStatefulSetList(statefulsetList)
	return modelStaList, nil
}

func (d *statefulsets) Patch(name string, pt model.PatchType, data []byte, subresources ...string) (result *model.StatefulSet, err error) {
	k8sSta, err := d.statefulset.Patch(name, k8stypes.PatchType(string(pt)), data, subresources...)
	if err != nil {
		logs.Error("Patch statefulset of %s/%s failed. Err:%+v", name, d.namespace, err)
		return nil, err
	}
	modelSta := types.FromK8sStatefulSet(k8sSta)
	return modelSta, nil
}

func (d *statefulsets) PatchToK8s(name string, pt model.PatchType, statefulset *model.StatefulSet) (*model.StatefulSet, []byte, error) {
	k8sStatefulSet := types.ToK8sStatefulSet(statefulset)
	rollConfig, err := json.Marshal(k8sStatefulSet)
	if err != nil {
		logs.Debug("Marshal rollingUpdateConfig failed %+v\n", k8sStatefulSet)
		return nil, nil, err
	}

	k8sSta, err := d.statefulset.Patch(name, k8stypes.PatchType(pt), rollConfig)
	if err != nil {
		logs.Error("PatchK8s statefulset of %s/%s failed. Err:%+v", statefulset.Name, d.namespace, err)
		return nil, nil, err
	}

	statefulsetConfig := types.GenerateStatefulSetConfig(k8sSta)
	statefulsetfileInfo, err := yaml.Marshal(statefulsetConfig)
	if err != nil {
		logs.Error("Marshal statefulset failed, error: %v", err)
		return nil, nil, err
	}
	modelSta := types.FromK8sStatefulSet(k8sSta)
	return modelSta, statefulsetfileInfo, nil
}

func (d *statefulsets) CreateByYaml(r io.Reader) (*model.StatefulSet, error) {
	context, err := ioutil.ReadAll(r)
	if err != nil {
		logs.Error("Read file failed, error: %v", err)
		return nil, err
	}

	var statefulset appsv1.StatefulSet
	err = yaml.Unmarshal(context, &statefulset)
	if err != nil {
		logs.Error("Unmarshal statefulset failed, error: %v", err)
		return nil, err
	}

	if statefulset.ObjectMeta.Namespace != d.namespace {
		logs.Error(namespacesErr.Error())
		return nil, namespacesErr
	}

	statefulsetInfo, err := d.statefulset.Create(&statefulset)
	if err != nil {
		logs.Error("Create statefulset failed, error: %v", err)
		return nil, err
	}

	return types.FromK8sStatefulSet(statefulsetInfo), nil
}

func (d *statefulsets) CheckYaml(r io.Reader) (*model.StatefulSet, error) {
	context, err := ioutil.ReadAll(r)
	if err != nil {
		logs.Error("Read file failed, error: %v", err)
		return nil, err
	}

	var statefulset appsv1.StatefulSet
	err = yaml.Unmarshal(context, &statefulset)
	if err != nil {
		logs.Error("Unmarshal statefulset failed, error: %v", err)
		return nil, err
	}

	if statefulset.ObjectMeta.Namespace != d.namespace {
		logs.Error(namespacesErr.Error())
		return nil, namespacesErr
	}

	return types.FromK8sStatefulSet(&statefulset), nil
}

// NewStatefulSets is to create a statefulset adapter
func NewStatefulSets(namespace string, statefulset v1.StatefulSetInterface) *statefulsets {
	return &statefulsets{
		namespace:   namespace,
		statefulset: statefulset,
	}
}
