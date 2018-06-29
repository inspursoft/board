package dao

import (
	"git/inspursoft/board/src/common/model"

	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddOperation(operation model.Operation) (int64, error) {
	o := orm.NewOrm()

	operation.CreationTime = time.Now()

	operationID, err := o.Insert(&operation)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return operationID, err
}

func GetOperation(operation model.Operation, fieldNames ...string) (*model.Operation, error) {
	o := orm.NewOrm()
	err := o.Read(&operation, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &operation, err
}

func UpdateOperation(operation model.Operation, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	operation.UpdateTime = time.Now()
	fieldNames = append(fieldNames, "update_time")
	operationID, err := o.Update(&operation, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return operationID, err
}

func DeleteOperation(operation model.Operation) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&operation)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func GetOperations(query model.Operation, fromtime string, totime string) ([]*model.Operation, error) {

	var operationSQL = "select * from operation"

	opt := " where "

	// TODO: Use params
	if query.ObjectType != "" {
		operationSQL = operationSQL + opt + fmt.Sprintf("type = '%s'", query.ObjectType)
		if opt == " where " {
			opt = " and "
		}
	}
	if query.Action != "" {
		operationSQL = operationSQL + opt + fmt.Sprintf("action = '%s'", query.Action)
		if opt == " where " {
			opt = " and "
		}
	}
	if query.UserName != "" {
		operationSQL = operationSQL + opt + fmt.Sprintf("user_name = '%s'", query.UserName)
		if opt == " where " {
			opt = " and "
		}
	}
	if query.Status != "" {
		operationSQL = operationSQL + opt + fmt.Sprintf("status = '%s'", query.Status)
		if opt == " where " {
			opt = " and "
		}
	}

	// fix me
	if fromtime != "" {
		operationSQL = operationSQL + opt + fmt.Sprintf("creation_time >= '%s'", query.CreationTime.String())
		if opt == " where " {
			opt = " and "
		}
	}

	if totime != "" {
		operationSQL = operationSQL + opt + fmt.Sprintf("creation_time <= '%s'", query.CreationTime.String())
		if opt == " where " {
			opt = " and "
		}
	}

	operations := make([]*model.Operation, 0)
	//_, err := orm.NewOrm().Raw(operationSQL, params).QueryRows(&operations)
	_, err := orm.NewOrm().Raw(operationSQL).QueryRows(&operations)
	if err != nil {
		return nil, err
	}
	return operations, nil
}

func GetPaginatedOperations(query model.OperationParam, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedOperations, error) {
	sql, params := generateOperationSQL(query)
	var err error

	pagination := &model.Pagination{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	pagination.TotalCount, err = getTotalRecordCount(sql, params)
	if err != nil {
		return nil, err
	}
	sql += getOrderSQL(operationTable, orderField, orderAsc) + ` limit ?, ?`
	params = append(params, pagination.GetPageOffset(), pagination.PageSize)
	
	logs.Debug("GetPaginatedOperations: +%+v", pagination.String())

	operationList, err := queryOperations(sql, params)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedOperations{
		OperationList: operationList,
		Pagination:        pagination,
	}, nil
}

func queryOperations(sql string, params []interface{}) ([]*model.Operation, error) {
	o := orm.NewOrm()
	operationList := make([]*model.Operation, 0)
	_, err := o.Raw(sql, params).QueryRows(&operationList)
	
	if err != nil {
		return nil, err
	}
	return operationList, nil
}

func generateOperationSQL(query model.OperationParam) (string, []interface{}) {
	sql := `select * from operation o where 1 = 1 `

	params := make([]interface{}, 0)
	params = append(params)

	if query.Operation_object != "" {
		params = append(params, "%"+query.Operation_object+"%")
		sql += ` and o.object_type like ? `
	}
	if(query.Operation_action != ""){
		params = append(params, "%"+query.Operation_action+"%")
		sql += ` and o.action like ? `
	}
	if(query.Operation_user != ""){
		params = append(params, "%"+query.Operation_user+"%")
		sql += ` and o.user_name like ? `
	}
	if(query.Operation_status != ""){
		params = append(params, "%"+query.Operation_status+"%")
		sql += ` and o.status like ? `
	}
	if(query.Operation_fromdate != ""){
		params = append(params, query.Operation_fromdate)
		sql += ` and o.creation_time >= ? `
	}
	if(query.Operation_todate != ""){
		params = append(params, query.Operation_todate)
		sql += ` and o.creation_time <= ? `
	}
	
	return sql, params
}