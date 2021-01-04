// a temp file for building and guiding
package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	scalepkg "k8s.io/client-go/scale"
)

type scales struct {
	namespace string
	scale     scalepkg.ScaleInterface
}

func (s *scales) Update(resource model.GroupResource, scale *model.Scale) (*model.Scale, error) {
	k8sScale := types.ToK8sScale(scale)
	newk8sScale, err := s.scale.Update(types.ToK8sGroupResource(resource), k8sScale)
	if err != nil {
		logs.Error("Update Scale of %s/%s failed. Err:%+v", scale.Name, s.namespace, err)
		return nil, err
	}

	modelScale := types.FromK8sScale(newk8sScale)
	return modelScale, nil
}

func (s *scales) Get(resource model.GroupResource, name string) (*model.Scale, error) {
	scaleinstance, err := s.scale.Get(types.ToK8sGroupResource(resource), name)
	if err != nil {
		logs.Error("Get scale of %s failed. Err:%+v", name, err)
		return nil, err
	}

	return types.FromK8sScale(scaleinstance), nil
}

func NewScales(namespace string, scale scalepkg.ScaleInterface) *scales {
	return &scales{
		namespace: namespace,
		scale:     scale,
	}
}
