package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"time"

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
	var base_format = "2006-01-02 15:04:05"
	var optparam model.OperationParam
	optparam.Action = o.GetString("operation_action")
	optparam.User = o.GetString("operation_user")
	optparam.Object = o.GetString("operation_object")
	optparam.Status = o.GetString("operation_status")
	starttime, _ := o.GetInt64("operation_fromdate", 0)
	totime, _ := o.GetInt64("operation_todate", 0)
	if starttime != 0 {
		optparam.Fromdate = time.Unix(starttime/1000, 0).Format(base_format)
	}
	if totime != 0 {
		optparam.Todate = time.Unix(totime/1000, 0).Format(base_format)
	}
	pageIndex, _ := o.GetInt("page_index", 1)
	pageSize, _ := o.GetInt("page_size", 10)
	orderField := o.GetString("order_field", "CREATION_TIME") //默认以creation_time排序
	orderAsc, _ := o.GetInt("order_asc", 0)

	paginatedoperations, err := service.GetPaginatedOperationList(optparam, pageIndex, pageSize, orderField, orderAsc)
	if err != nil {
		o.internalError(err)
		return
	}
	o.renderJSON(paginatedoperations)
}
