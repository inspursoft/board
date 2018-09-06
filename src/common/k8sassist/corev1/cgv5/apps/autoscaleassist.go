// an adapter file from board to k8s for auto-scale
package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/autoscaling/v1"
)

type autoscales struct {
	namespace string
	autoscale v1.AutoscalingV1Interface
}

func (as *autoscales) Get(name string) (*model.AutoScale, error) {
	autoscaleinstance, err := as.autoscale.HorizontalPodAutoscalers(as.namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		logs.Error("Get auto scale of %s failed. Err:%+v", name, err)
		return nil, err
	}
	return types.FromK8sAutoScale(autoscaleinstance), nil
}

func NewAutoScales(namespace string, autoscale v1.AutoscalingV1Interface) *autoscales {
	return &autoscales{
		namespace: namespace,
		autoscale: autoscale,
	}
}
