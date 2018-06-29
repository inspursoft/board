package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/utils"
	"git/inspursoft/board/src/common/model"
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
	
	var optparam model.OperationParam;
	optparam.Operation_action = o.GetString("operation_action")
	optparam.Operation_user = o.GetString("operation_user")
	optparam.Operation_object = o.GetString("operation_object")
	optparam.Operation_status = o.GetString("operation_status")
	optparam.Operation_fromdate = o.GetString("operation_fromdate")
	optparam.Operation_todate = o.GetString("operation_todate")
	pageIndex, _ := o.GetInt("page_index", 1)
	pageSize, _ := o.GetInt("page_size", 10)
	
	orderField := o.GetString("order_field", "CREATION_TIME")  //默认以creation_time排序
	orderAsc, _ := o.GetInt("order_asc", 0)

	paginatedoperations, err := service.GetPaginatedOperationList(optparam, pageIndex, pageSize, orderField,orderAsc)
	if err != nil {
		o.internalError(err)
		return
	}
	o.renderJSON(paginatedoperations)
}