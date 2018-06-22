package dao

import (
	"git/inspursoft/board/src/common/model"

	"time"

	//"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddOperation(operation model.Operation) (int64, error) {
	o := orm.NewOrm()

	operation.CreationTime = time.Now()
	operation.UpdateTime = operation.CreationTime

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

func queryOperations(field string, value interface{}) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable("operation")
	if value == nil {
		return qs
	}
	return qs.Filter(field+"__contains", value)
}

func GetOperations(query model.Operation, fromtime string, totime string) ([]*model.Operation, error) {

	var operationSQL = "select * from operation "

	opt := "where"

	// TODO: Use params
	if query.ObjectType != "" {
		operationSQL = operationSQL + opt + "type = " + `"` + query.ObjectType + `"`
		if opt == "where" {
			opt = "and"
		}
	}
	if query.Action != "" {
		operationSQL = operationSQL + opt + "action = " + `"` + query.Action + `"`
		if opt == "where" {
			opt = "and"
		}
	}
	if query.UserName != "" {
		operationSQL = operationSQL + opt + "user_name = " + `"` + query.UserName + `"`
		if opt == "where" {
			opt = "and"
		}
	}
	if query.Status != "" {
		operationSQL = operationSQL + opt + "status = " + `"` + query.Status + `"`
		if opt == "where" {
			opt = "and"
		}
	}

	// fix me
	if fromtime != "" {
		operationSQL = operationSQL + opt + "creation_time >= " + `"` + query.CreationTime.String() + `"`
		if opt == "where" {
			opt = "and"
		}
	}

	if totime != "" {
		operationSQL = operationSQL + opt + "creation_time <= " + `"` + query.CreationTime.String() + `"`
		if opt == "where" {
			opt = "and"
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
