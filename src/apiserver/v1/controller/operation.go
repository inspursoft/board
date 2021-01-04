package controller

import (
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type OperationController struct {
	c.BaseController
}

func (o *OperationController) Prepare() {
	o.EnableXSRF = false
	o.ResolveSignedInUser()
}

func (o *OperationController) OperationList() {
	if !o.IsSysAdmin {
		o.CustomAbortAudit(http.StatusForbidden, "Insufficient permissions.")
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

	orderFieldValue, err := service.ParseOrderField("operation", orderField)
	if err != nil {
		o.CustomAbortAudit(http.StatusBadRequest, err.Error())
		return
	}

	paginatedoperations, err := service.GetPaginatedOperationList(optparam, pageIndex, pageSize, orderFieldValue, orderAsc)
	if err != nil {
		o.InternalError(err)
		return
	}
	o.RenderJSON(paginatedoperations)
}

func (o *OperationController) CreateOperation() {
	operation := model.Operation{}
	err := o.ResolveBody(&operation)
	if err != nil {
		o.InternalError(err)
		return
	}
	err = service.CreateOperationAudit(&operation)
	if err != nil {
		logs.Error("Failed to create operation Audit. Error:%+v", err)
		o.InternalError(err)
		return
	}
	o.RenderJSON(operation)
}
