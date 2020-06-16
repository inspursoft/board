package controller

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/utils"
	"github.com/astaxie/beego/logs"
	"net/http"
)

type K8SProxyController struct {
	c.BaseController
}

func (u *K8SProxyController) Prepare() {
	u.EnableXSRF = false
	u.ResolveSignedInUser()
	//	u.RecordOperationAudit()
}

func (u *K8SProxyController) ProxyAction() {
	info, err := service.GetSystemInfo()
	if err != nil {
		u.CustomAbortAudit(http.StatusServiceUnavailable, fmt.Sprintf("Proxy to kubernetes error: %+v.", err))
		return
	}
	if !info.K8SProxyEnabled {
		u.CustomAbortAudit(http.StatusForbidden, "You should turn on the kubernetes proxy setting.")
		return
	}
	if !u.IsSysAdmin {
		u.CustomAbortAudit(http.StatusForbidden, "User should be admin")
		return
	}
	handler, err := k8sProxyHandler()
	if err != nil {
		msg := fmt.Sprintf("Proxy to kubernetes error: %+v.", err)
		logs.Info(msg)
		u.CustomAbortAudit(http.StatusServiceUnavailable, msg)
		return
	}
	handler.ServeHTTP(u.Ctx.ResponseWriter, u.Ctx.Request)
}

func k8sProxyHandler() (http.Handler, error) {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: utils.GetConfig("KUBE_CONFIG_PATH")(),
	})
	return k8sclient.AppV1().Proxy().ProxyAPI("/kubernetes/")
}
