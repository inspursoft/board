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

func UpdateOperation(operation model.Project, fieldNames ...string) (int64, error) {
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
