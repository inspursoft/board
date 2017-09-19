package dao

import (
	"git/inspursoft/board/src/common/model"

	"time"

	"github.com/astaxie/beego/orm"
)

func AddService(service model.ServiceStatus) (int64, error) {
	o := orm.NewOrm()

	service.CreationTime = time.Now()
	service.UpdateTime = service.CreationTime

	serviceID, err := o.Insert(&service)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return serviceID, err
}

func GetService(service model.ServiceStatus, fieldNames ...string) (*model.ServiceStatus, error) {
	o := orm.NewOrm()
	err := o.Read(&service, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &service, err
}

func UpdateService(service model.ServiceStatus, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	service.UpdateTime = time.Now()
	serviceID, err := o.Update(&service, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return serviceID, err
}

func GetServiceData() ([]model.ServiceStatus, error) {
	var serviceList []model.ServiceStatus

	o := orm.NewOrm()
	serviceModel := new(model.ServiceStatus)
	qs := o.QueryTable(serviceModel)
	_, err := qs.All(&serviceList)
	if err != nil {
		return nil, err
	}

	return serviceList, err
}
