package systems

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
)

//Operations about system info
type SupplementController struct {
	c.BaseController
}

func (i *SupplementController) Prepare() {
	i.EnableXSRF = false
}

// @Title Get system information.
// @Description Get system information.
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /info [get]
func (i *SupplementController) Info() {
	systemInfo, err := service.GetSystemInfo()
	if err != nil {
		i.InternalError(err)
		return
	}
	i.RenderJSON(systemInfo)
}

// @Title Get system resources information.
// @Description Get system resources information.
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /resources [get]
func (i *SupplementController) Resources() {
	systemResources, err := service.GetSystemResourcesInfo()
	if err != nil {
		i.InternalError(err)
		return
	}
	i.RenderJSON(systemResources)
}

// @Title Get Kubernetes information.
// @Description Get Kubernetes information.
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /kubernetes-info [get]
func (i *SupplementController) KubernetesInfo() {
	kubernetesInfo, err := service.GetKubernetesInfo()
	if err != nil {
		i.InternalError(err)
		return
	}
	i.RenderJSON(kubernetesInfo)
}
