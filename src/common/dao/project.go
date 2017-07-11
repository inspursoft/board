package dao

import (
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

func AddProject(project model.Project) (int64, error) {
	o := orm.NewOrm()
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
			return &project, nil
		}
		return nil, err
	}
	return &project, err
}

func GetProjects(field string, value interface{}, fieldNames ...string) ([]*model.Project, error) {
	o := orm.NewOrm()
	var projects []*model.Project
	var err error
	qs := o.QueryTable("project")
	if value == nil {
		_, err = qs.All(&projects, fieldNames...)
	} else {
		_, err = qs.Filter(field+"__contains", value).
			All(&projects, fieldNames...)
	}
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func UpdateProject(project model.Project, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	projectID, err := o.Update(&project, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return projectID, err
}

func DeleteProject(project model.Project) (int64, error) {
	o := orm.NewOrm()
	return o.Delete(&project, "id")
}
