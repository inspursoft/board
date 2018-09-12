package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
)

type AutoScaleController struct {
	BaseController
}

func (as *AutoScaleController) resolveServiceInfo() (s *model.ServiceStatus, err error) {
	serviceID, err := strconv.Atoi(as.Ctx.Input.Param(":id"))
	if err != nil {
		as.internalError(err)
		return
	}
	// Get the project info of this service
	s, err = service.GetServiceByID(int64(serviceID))
	if err != nil {
		as.internalError(err)
		return
	}
	if s == nil {
		as.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
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
	err = as.resolveBody(hpa)
	if err != nil {
		return
	}
	// override the fields
	hpa.ServiceID = svc.ID

	// add the hpa to k8s
	hpa, err = service.CreateAutoScale(svc, hpa)
	if err != nil {
		as.internalError(err)
		return
	}
	as.renderJSON(hpa)
}

func (as *AutoScaleController) ListAutoScaleAction() {
	// make sure the service exist
	svc, err := as.resolveServiceInfo()
	if err != nil {
		return
	}

	// list the hpas from storage
	hpas, err := service.ListAutoScales(svc)
	if err != nil {
		as.internalError(err)
		return
	}
	for _, hpa := range hpas {
		_, err = service.GetAutoScaleK8s(svc.ProjectName, hpa.HPAName)
		if err != nil {
			logs.Debug("Not found hpa %s in system", hpa.HPAName)
			hpa.HPAStatus = 0
		} else {
			hpa.HPAStatus = 1
		}
	}

	as.renderJSON(hpas)
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
		as.internalError(err)
		return
	}

	// resolve the hpa
	hpa := new(model.ServiceAutoScale)
	err = as.resolveBody(hpa)
	if err != nil {
		return
	}
	// override the fields
	hpa.ID = int64(hpaid)
	hpa.ServiceID = svc.ID

	hpa, err = service.UpdateAutoScale(svc, hpa)
	if err != nil {
		as.internalError(err)
		return
	}
	as.renderJSON(hpa)
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
		as.internalError(err)
		return
	}

	// delete the autoscale
	err = service.DeleteAutoScale(svc, int64(hpaid))
	if err != nil {
		as.internalError(err)
		return
	}
}
