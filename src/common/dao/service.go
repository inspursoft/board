package dao

import (
	"git/inspursoft/board/src/common/model"

	"time"

	"github.com/astaxie/beego/logs"
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

func DeleteService(service model.ServiceStatus) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&service)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func generateServiceStatusSQL(query model.ServiceStatus, userID int64) (string, []interface{}) {
	sql := `select distinct s.id, s.name, s.project_name, u.username as owner_name, s.owner_id, s.creation_time, s.status, s.public
	from service_status s 
		left join project_member pm on s.project_id = pm.project_id
		left join project p on p.id = pm.project_id
		left join user u on u.id = s.owner_id
	where s.deleted = 0 and s.status >= 1
	and (p.public = 1 
		or s.public = 1
		or s.project_id in (select p.id from project p left join project_member pm on p.id = pm.project_id  left join user u on u.id = pm.user_id where p.deleted = 0 and u.deleted = 0 and u.id = ?)
		or exists (select * from user u where u.deleted = 0 and u.system_admin = 1 and u.id = ?))`

	params := make([]interface{}, 0)
	params = append(params, userID, userID)

	if query.Name != "" {
		params = append(params, "%"+query.Name+"%")
		sql += ` and s.name like ? `
	}
	return sql, params
}

func queryServiceStatus(sql string, params []interface{}) ([]*model.ServiceStatus, error) {
	serviceList := make([]*model.ServiceStatus, 0)
	_, err := orm.NewOrm().Raw(sql, params).QueryRows(&serviceList)
	if err != nil {
		return nil, err
	}
	return serviceList, nil
}

func GetServiceData(query model.ServiceStatus, userID int64) ([]*model.ServiceStatus, error) {
	sql, params := generateServiceStatusSQL(query, userID)
	return queryServiceStatus(sql, params)
}

func GetPaginatedServiceData(query model.ServiceStatus, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedServiceStatus, error) {
	sql, params := generateServiceStatusSQL(query, userID)
	var err error

	pagination := &model.Pagination{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	pagination.TotalCount, err = getTotalRecordCount(sql, params)
	if err != nil {
		return nil, err
	}
	sql += getOrderSQL(serviceTable, orderField, orderAsc) + ` limit ?, ?`
	params = append(params, pagination.GetPageOffset(), pagination.PageSize)
	logs.Debug("%+v", pagination.String())

	serviceList, err := queryServiceStatus(sql, params)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedServiceStatus{
		ServiceStatusList: serviceList,
		Pagination:        pagination,
	}, nil
}

//GetService(servicequery, "id")
func SyncServiceData(service model.ServiceStatus) (int64, error) {

	var servicequery model.ServiceStatus
	servicequery.Name = service.Name
	o := orm.NewOrm()
	err := o.Read(&servicequery, "name")
	if err != orm.ErrNoRows {
		return 0, nil
	}

	serviceID, err := o.Insert(&service)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return serviceID, err
}

func GetSelectableServices(pName string, sName string) ([]string, error) {
	o := orm.NewOrm()
	sql := `select s.name
	from service_status s 
	where s.deleted = 0 and s.status >= 1
	and s.project_name = ? and s.name != ?`

	params := make([]interface{}, 0)
	params = append(params, pName, sName)

	var serviceList []string
	_, err := o.Raw(sql, params).QueryRows(&serviceList)
	return serviceList, err
}
