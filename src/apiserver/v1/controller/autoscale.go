package controller

import (
	"fmt"
	"net/http"
	"strconv"

	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
)

type AutoScaleController struct {
	c.BaseController
}

func (as *AutoScaleController) resolveServiceInfo() (s *model.ServiceStatus, err error) {
	serviceID, err := strconv.Atoi(as.Ctx.Input.Param(":id"))
	if err != nil {
		as.InternalError(err)
		return
	}
	// Get the project info of this service
	s, err = service.GetServiceByID(int64(serviceID))
	if err != nil {
		as.InternalError(err)
		return
	}
	if s == nil {
		as.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	return
}

func (as *AutoScaleController) CreateAutoScaleAction() {
	// make sure the service exist
	svc, err := as.resolveServiceInfo()
	if err != nil {
		return
	}

	// resolve the hpa
	hpa := new(model.ServiceAutoScale)
	err = as.ResolveBody(hpa)
	if err != nil {
		return
	}
	logs.Info("Added autoscale %s: %+v", hpa.HPAName, hpa)

	// override the fields
	hpa.ServiceID = svc.ID

	// do some check
	exist, err := service.CheckAutoScaleExist(svc, hpa.HPAName)
	if err != nil {
		as.InternalError(err)
		return
	} else if exist {
		as.CustomAbortAudit(http.StatusConflict, fmt.Sprintf("AutoScale %s already exists in cluster.", hpa.HPAName))
		return
	}

	// add the hpa to k8s
	hpa, err = service.CreateAutoScale(svc, hpa)
	if err != nil {
		as.InternalError(err)
		return
	}
	as.RenderJSON(hpa)
}

func (as *AutoScaleController) ListAutoScaleAction() {
	// make sure the service exist
	svc, err := as.resolveServiceInfo()
	if err != nil {
		return
	}
	logs.Info("list all autoscales of service %s", svc.Name)

	// list the hpas from storage
	hpas, err := service.ListAutoScales(svc)
	if err != nil {
		as.InternalError(err)
		return
	}
	for _, hpa := range hpas {
		_, exist, err := service.GetAutoScaleK8s(svc.ProjectName, hpa.HPAName)
		if err != nil {
			as.InternalError(err)
			return
		} else if exist {
			hpa.HPAStatus = 1
		} else {
			logs.Debug("Not found hpa %s in system", hpa.HPAName)
			hpa.HPAStatus = 0
		}
	}

	as.RenderJSON(hpas)
}

func (as *AutoScaleController) UpdateAutoScaleAction() {
	// make sure the service exist
	svc, err := as.resolveServiceInfo()
	if err != nil {
		return
	}

	// get the hpa id
	hpaid, err := strconv.Atoi(as.Ctx.Input.Param(":hpaid"))
	if err != nil {
		as.InternalError(err)
		return
	}

	// resolve the hpa
	hpa := new(model.ServiceAutoScale)
	err = as.ResolveBody(hpa)
	if err != nil {
		return
	}

	logs.Info("update autoscale %d to %+v", hpaid, hpa)
	// override the fields
	hpa.ID = int64(hpaid)
	hpa.ServiceID = svc.ID

	// do some check
	autoscale, err := service.GetAutoScale(svc, hpa.ID)
	if err != nil {
		as.InternalError(err)
		return
	} else if autoscale == nil {
		as.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Autoscale %d does not exists.", hpa.ID))
		return
	} else if autoscale.HPAName != hpa.HPAName {
		as.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("can't change Autoscale %s's name to %s", autoscale.HPAName, hpa.HPAName))
		return
	} else if autoscale.ServiceID != hpa.ServiceID {
		as.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("can't change Autoscale's service id %d to %d", autoscale.ServiceID, hpa.ServiceID))
		return
	}

	hpa, err = service.UpdateAutoScale(svc, hpa)
	if err != nil {
		as.InternalError(err)
		return
	}
	as.RenderJSON(hpa)
}

func (as *AutoScaleController) DeleteAutoScaleAction() {
	// make sure the service exist
	svc, err := as.resolveServiceInfo()
	if err != nil {
		return
	}

	// get the hpa id
	hpaid, err := strconv.Atoi(as.Ctx.Input.Param(":hpaid"))
	if err != nil {
		as.InternalError(err)
		return
	}
	logs.Info("delete autoscale %d", hpaid)

	// do some check
	autoscale, err := service.GetAutoScale(svc, int64(hpaid))
	if err != nil {
		as.InternalError(err)
		return
	} else if autoscale == nil {
		as.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Autoscale %d does not exists.", int64(hpaid)))
		return
	}

	// delete the autoscale
	err = service.DeleteAutoScale(svc, int64(hpaid))
	if err != nil {
		as.InternalError(err)
		return
	}
}
