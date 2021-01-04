package dao

import (
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

type SearchProjectResult struct {
	OwnerName   string `json:"owner_name" orm:"column(owner_name)"`
	ProjectName string `json:"project_name" orm:"column(project_name)"`
	IsPublic    bool   `json:"is_public"  orm:"column(is_public)"`
}
type SearchUserResult struct {
	UserName  string `json:"user_name" orm:"column(user_name)"`
	RoleName  string `json:"role_name" orm:"column(role_name)"`
	UserEmail string `json:"user_email" orm:"column(user_email)"`
}

func SearchPrivateProject(projectName string, userName string) ([]SearchProjectResult, error) {
	var searchRes []SearchProjectResult
	sql := `select distinct p.name as project_name, 
	p.owner_name, 
	p.public as is_public
from project p 
left join project_member pm on p.id = pm.project_id
left join user u on u.id = pm.user_id
where p.deleted = 0 
and p.name like ?
and (p.public = 1
or p.id in (select p.id from project p left join project_member pm on p.id = pm.project_id  left join user u on u.id = pm.user_id where p.deleted = 0 and u.deleted = 0 and u.username = ?)
or exists (select * from user u where u.deleted = 0 and u.system_admin = 1 and u.username = ?));`
	o := orm.NewOrm()
	_, err := o.Raw(sql, "%"+projectName+"%", userName, userName).QueryRows(&searchRes)
	return searchRes, err
}
func SearchPublicProject(projectName string) ([]SearchProjectResult, error) {
	var searchRes []SearchProjectResult
	sql := `select
	  owner_name as owner_name,
	  project.name as project_name,
	  project.public as is_public
	from project
	where deleted = 0 
	and public = 1
	and project.name like ?;`
	o := orm.NewOrm()
	_, err := o.Raw(sql, "%"+projectName+"%").QueryRows(&searchRes)
	return searchRes, err
}
func SearchUser(activeUser string, searchName string) ([]SearchUserResult, error) {
	var searchRes []SearchUserResult
	sql := ` select distinct
  u.username as user_name,
  r.name     as role_name,
  u.email    as user_email
from user u 
  left join  project_member pm on pm.user_id = u.id
  left  join role r on pm.role_id = r.id
  where u.deleted = 0
			and exists (select * from user where deleted = 0 and system_admin = 1 and username = ? ) 
			and u.username like ?;`
	o := orm.NewOrm()
	_, err := o.Raw(sql, activeUser, "%"+searchName+"%").QueryRows(&searchRes)
	return searchRes, err
}

func SearchPublicSvr(para string) ([]model.ServiceStatus, error) {
	var svr []model.ServiceStatus
	o := orm.NewOrm()
	qs := o.QueryTable("service_status")
	_, err := qs.Filter("deleted", 0).Filter("public", 1).Filter("status__gte", 1).Filter("name__contains", para).All(&svr)
	return svr, err
}
