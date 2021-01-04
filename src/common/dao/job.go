package dao

import (
	"github.com/inspursoft/board/src/common/model"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddJob(job model.JobStatusMO) (int64, error) {
	o := orm.NewOrm()

	job.CreationTime = time.Now()
	job.UpdateTime = job.CreationTime

	jobID, err := o.Insert(&job)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return jobID, err
}

func GetJob(job model.JobStatusMO, fieldNames ...string) (*model.JobStatusMO, error) {
	o := orm.NewOrm()
	err := o.Read(&job, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &job, err
}

func UpdateJob(job model.JobStatusMO, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	job.UpdateTime = time.Now()
	jobID, err := o.Update(&job, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return jobID, err
}

func DeleteJob(job model.JobStatusMO) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&job)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func DeleteJobByNames(jobs []model.JobStatusMO) (int64, error) {
	if jobs == nil || len(jobs) == 0 {
		return 0, nil
	}
	sql, params := generateDeleteJobByNamesSQL(jobs)
	result, err := orm.NewOrm().Raw(sql, params).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func generateDeleteJobByNamesSQL(jobs []model.JobStatusMO) (string, []interface{}) {
	values := ""
	params := make([]interface{}, 0, 2*len(jobs))
	for i, job := range jobs {
		params = append(params, job.Name, job.ProjectName)
		if i != 0 {
			values += ","
		}
		values += ` (?,?)`
	}
	sql := `delete from job_status where (name, project_name) in ( ` + values + ` )`
	return sql, params
}

func generateJobStatusSQL(query model.JobStatusMO, userID int64) (string, []interface{}) {
	sql := `select distinct s.id, s.name, s.project_id, s.project_name, u.username as owner_name, s.owner_id, s.creation_time, s.update_time, s.status, s.source
	from job_status s 
		left join project p on p.id = s.project_id
		left join user u on u.id = s.owner_id
	where s.deleted = 0 and s.status >= 1
	and u.deleted = 0 and ( u.id = ? 
		or exists (select u.id from user u where u.deleted = 0 and u.system_admin = 1 and u.id = ?))`

	params := make([]interface{}, 0)
	params = append(params, userID, userID)

	if query.ProjectID != 0 {
		params = append(params, query.ProjectID)
		sql += ` and s.project_id = ?`
	}
	if query.Name != "" {
		params = append(params, "%"+query.Name+"%")
		sql += ` and s.name like ? `
	}
	return sql, params
}

func queryJobStatus(sql string, params []interface{}) ([]*model.JobStatusMO, error) {
	jobList := make([]*model.JobStatusMO, 0)
	_, err := orm.NewOrm().Raw(sql, params).QueryRows(&jobList)
	if err != nil {
		return nil, err
	}
	return jobList, nil
}

func GetJobData(query model.JobStatusMO, userID int64) ([]*model.JobStatusMO, error) {
	sql, params := generateJobStatusSQL(query, userID)
	return queryJobStatus(sql, params)
}

func GetPaginatedJobData(query model.JobStatusMO, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedJobStatus, error) {
	sql, params := generateJobStatusSQL(query, userID)
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

	jobList, err := queryJobStatus(sql, params)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedJobStatus{
		JobStatusList: jobList,
		Pagination:    pagination,
	}, nil
}

func SyncJobData(job model.JobStatusMO) (int64, error) {
	var jobquery model.JobStatusMO
	jobquery.Name = job.Name
	o := orm.NewOrm()
	err := o.Read(&jobquery, "name")
	if err != orm.ErrNoRows {
		return 0, nil
	}

	jobID, err := o.Insert(&job)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return jobID, err
}

func GetJobs(field string, value interface{}, selectedFields ...string) ([]model.JobStatusMO, error) {
	jobs := make([]model.JobStatusMO, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable("job_status").Filter(field, value).Filter("status", 1).Filter("deleted", 0).All(&jobs, selectedFields...)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
