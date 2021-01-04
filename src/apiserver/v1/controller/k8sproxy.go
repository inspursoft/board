package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"github.com/astaxie/beego/logs"
	"net/http"
)

type K8SProxyController struct {
	c.BaseController
}

func (u *K8SProxyController) Prepare() {
	u.EnableXSRF = false
	u.ResolveSignedInUser()
}

func (u *K8SProxyController) GetK8SProxyConfig() {
	config, err := service.GetK8SProxyConfig()
	if err != nil {
		logs.Debug("Failed to get k8s proxy config: %+v", err)
		u.InternalError(err)
		return
	}
	u.RenderJSON(config)
}

func (u *K8SProxyController) SetK8SProxyConfig() {
	if !u.IsSysAdmin {
		u.CustomAbortAudit(http.StatusForbidden, "User should be admin")
		return
	}
	// only audit the modify operation
	u.RecordOperationAudit()
	var config model.K8SProxyConfig
	var err error
	err = u.ResolveBody(&config)
	if err != nil {
		u.InternalError(err)
		return
	}
	err = service.SetK8SProxyConfig(config)
	if err != nil {
		logs.Debug("Failed to set k8s proxy: %+v", err)
		u.InternalError(err)
	}
}

func (u *K8SProxyController) ProxyAction() {
	config, err := service.GetK8SProxyConfig()
	if err != nil {
		u.CustomAbortAudit(http.StatusServiceUnavailable, fmt.Sprintf("Proxy to kubernetes error: %+v.", err))
		return
	}
	if !config.Enable {
		u.CustomAbortAudit(http.StatusForbidden, "You should turn on the kubernetes proxy setting.")
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
