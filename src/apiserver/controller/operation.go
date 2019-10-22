package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type OperationController struct {
	BaseController
}

func (o *OperationController) Prepare() {
	o.resolveSignedInUser()
}

func (o *OperationController) OperationList() {
	if !o.isSysAdmin {
		o.customAbort(http.StatusForbidden, "Insufficient permissions.")
		return
	}
	var optparam model.OperationParam
	optparam.Action = o.GetString("operation_action")
	optparam.User = o.GetString("operation_user")
	optparam.Object = o.GetString("operation_object")
	optparam.Status = o.GetString("operation_status")
	optparam.Fromdate, _ = o.GetInt64("operation_fromdate", 0)
	optparam.Todate, _ = o.GetInt64("operation_todate", 0)
	pageIndex, _ := o.GetInt("page_index", 1)
	pageSize, _ := o.GetInt("page_size", 10)
	orderField := o.GetString("order_field", "creation_time")
	orderAsc, _ := o.GetInt("order_asc", 0)

	paginatedoperations, err := service.GetPaginatedOperationList(optparam, pageIndex, pageSize, orderField, orderAsc)
	if err != nil {
		o.internalError(err)
		return
	}
	o.renderJSON(paginatedoperations)
}

func (o *OperationController) CreateOperation() {
	operation := model.Operation{}
	err := o.resolveBody(&operation)
	if err != nil {
		o.internalError(err)
		return
	}
	err = service.CreateOperationAudit(&operation)
	if err != nil {
		logs.Error("Failed to create operation Audit. Error:%+v", err)
		o.internalError(err)
		return
	}
	o.renderJSON(operation)
}
