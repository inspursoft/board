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

func DeleteServiceByNames(services []model.ServiceStatus) (int64, error) {
	if services == nil || len(services) == 0 {
		return 0, nil
	}
	sql, params := generateDeleteServiceByNamesSQL(services)
	result, err := orm.NewOrm().Raw(sql, params).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func generateDeleteServiceByNamesSQL(services []model.ServiceStatus) (string, []interface{}) {
	values := ""
	params := make([]interface{}, 0, 2*len(services))
	for i, svc := range services {
		params = append(params, svc.Name, svc.ProjectName)
		if i != 0 {
			values += ","
		}
		values += ` (?,?)`
	}
	sql := `delete from service_status where (name, project_name) in ( ` + values + ` )`
	return sql, params
}

func generateServiceStatusSQL(query model.ServiceStatusFilter, userID int64) (string, []interface{}) {
	sql := `select distinct s.id, s.name, s.project_id, s.project_name, u.username as owner_name, s.owner_id, s.creation_time, s.update_time, s.status, s.type, s.public, s.source, s.source_id,
	(select if(count(s0.id), 1, 0) from service_status s0 where s0.deleted = 0 and s0.id = s.id and s0.project_id in (
		select p0.id
		from project p0
		left join project_member pm0 on p0.id = pm0.project_id
		left join user u0 on u0.id = pm0.user_id
			where p0.deleted = 0 and u0.deleted = 0 and u0.id = ?) or exists (
				select u0.id
			  from user u0
			  where u0.deleted = 0 and u0.system_admin = 1 and u0.id = ?)) as is_member
	from service_status s 
		left join project_member pm on s.project_id = pm.project_id
		left join project p on p.id = pm.project_id
		left join user u on u.id = s.owner_id
	where s.deleted = 0 and s.status >= 1
	and (s.public = 1 
		or s.project_id in (select p.id from project p left join project_member pm on p.id = pm.project_id  left join user u on u.id = pm.user_id where p.deleted = 0 and u.deleted = 0 and u.id = ?)
		or exists (select u.id from user u where u.deleted = 0 and u.system_admin = 1 and u.id = ?))`

	params := make([]interface{}, 0)
	params = append(params, userID, userID, userID, userID)

	if query.ProjectID != 0 {
		params = append(params, query.ProjectID)
		sql += ` and s.project_id = ?`
	}
	if query.Name != "" {
		params = append(params, "%"+query.Name+"%")
		sql += ` and s.name like ? `
	}
	if query.Source != nil {
		params = append(params, query.Source)
		sql += ` and s.source = ? `
	}
	if query.SourceID != nil {
		params = append(params, query.SourceID)
		sql += ` and s.source_id = ? `
	}
	if query.ProjectID != 0 {
		params = append(params, query.ProjectID)
		sql += ` and s.project_id = ? `
	}
	return sql, params
}

func queryServiceStatus(sql string, params []interface{}) ([]*model.ServiceStatusMO, error) {
	serviceList := make([]*model.ServiceStatusMO, 0)
	_, err := orm.NewOrm().Raw(sql, params).QueryRows(&serviceList)
	if err != nil {
		return nil, err
	}
	return serviceList, nil
}

func GetServiceData(query model.ServiceStatusFilter, userID int64) ([]*model.ServiceStatusMO, error) {
	sql, params := generateServiceStatusSQL(query, userID)
	return queryServiceStatus(sql, params)
}

func GetPaginatedServiceData(query model.ServiceStatusFilter, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedServiceStatus, error) {
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
	sql += getOrderSQL(orderField, orderAsc) + ` limit ?, ?`
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

func GetServices(field string, value interface{}, selectedFields ...string) ([]model.ServiceStatus, error) {
	services := make([]model.ServiceStatus, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable("service_status").Filter(field, value).Filter("deleted", 0).All(&services, selectedFields...)
	if err != nil {
		return nil, err
	}
	return services, nil
}
