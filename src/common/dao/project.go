package dao

import (
	"git/inspursoft/board/src/common/model"

	"time"

	"github.com/astaxie/beego/orm"
)

func AddProject(project model.Project) (int64, error) {
	o := orm.NewOrm()

	project.CreationTime = time.Now()
	project.UpdateTime = project.CreationTime

	projectID, err := o.Insert(&project)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return projectID, err
}

func GetProject(project model.Project, fieldNames ...string) (*model.Project, error) {
	o := orm.NewOrm()
	err := o.Read(&project, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &project, err
}

func GetProjectsByUser(query model.Project, userID int64) ([]*model.Project, error) {
	projectByUserSQL := `select  distinct p.id, p.name, p.comment, p.creation_time, 
	p.update_time, p.owner_id, p.owner_name, 
	p.public, p.toggleable, p.current_user_role_id, 
	p.service_count
from project p 
	left join project_member pm on p.id = pm.project_id
	left join user u on u.id = pm.user_id
where p.deleted = 0 
and p.name like ?
and (p.public = 1
	or p.id in (select p.id from project p left join project_member pm on p.id = pm.project_id  left join user u on u.id = pm.user_id where p.deleted = 0 and u.deleted = 0 and u.id=?)
	or exists (select * from user u where u.deleted = 0 and u.system_admin = 1 and u.id=?));`

	params := make([]interface{}, 0)
	params = append(params, "%"+query.Name+"%", userID, userID)

	if query.ServiceCount != 0 {
		projectByUserSQL += ` and p.service_count = ?`
		params = append(params, query.ServiceCount)
	}
	if query.CurrentUserRoleID != 0 {
		projectByUserSQL += ` and p.current_user_role_id = ?`
		params = append(params, query.CurrentUserRoleID)
	}

	var projects []*model.Project
	o := orm.NewOrm()
	_, err := o.Raw(projectByUserSQL, params).QueryRows(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func UpdateProject(project model.Project, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	project.UpdateTime = time.Now()
	projectID, err := o.Update(&project, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return projectID, err
}
