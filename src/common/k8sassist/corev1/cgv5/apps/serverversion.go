// a temp file for building and guiding
package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	discovery "k8s.io/client-go/discovery"
)

type serverversion struct {
	discvy discovery.DiscoveryInterface
}

func (s *serverversion) ServerVersion() (*model.KubernetesInfo, error) {
	info, err := s.discvy.ServerVersion()
	if err != nil {
		logs.Error("Get Kubernetes ServerInfo error. Err:%+v", err)
		return nil, err
	}

	modelInfo := types.FromK8sInfo(info)
	return modelInfo, nil
}

func NewServerVersion(discvy discovery.DiscoveryInterface) *serverversion {
	return &serverversion{
		discvy: discvy,
	}
}
