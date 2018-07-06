package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type OperationController struct {
	BaseController
}

func (o *OperationController) Prepare() {
	o.resolveSignedInUser()
	o.isExternalAuth = utils.GetBoolValue("IS_EXTERNAL_AUTH")
}

func (o *OperationController) OperationList() {
	if o.isExternalAuth && o.currentUser.Username != "admin" {
		logs.Debug("Current AUTH_MODE is external auth.")
		o.customAbort(http.StatusPreconditionFailed, "Current AUTH_MODE is not available to the user.")
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
