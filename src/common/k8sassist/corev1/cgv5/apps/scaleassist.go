// a temp file for building and guiding
package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

type scales struct {
	namespace string
	scale     v1beta1.ScaleInterface
}

func (s *scales) Update(kind string, scale *model.Scale) (*model.Scale, error) {
	k8sScale := types.ToK8sScale(scale)
	newk8sScale, err := s.scale.Update(kind, k8sScale)
	if err != nil {
		logs.Error("Update Scale of %s/%s failed. Err:%+v", scale.Name, s.namespace, err)
		return nil, err
	}

	modelScale := types.FromK8sScale(newk8sScale)
	return modelScale, nil
}

func (s *scales) Get(kind string, name string) (*model.Scale, error) {
	scaleinstance, err := s.scale.Get(kind, name)
	if err != nil {
		logs.Error("Get scale of %s failed. Err:%+v", name, err)
		return nil, err
	}

	return types.FromK8sScale(scaleinstance), nil
}

func NewScales(namespace string, scale v1beta1.ScaleInterface) *scales {
	return &scales{
		namespace: namespace,
		scale:     scale,
	}
}
