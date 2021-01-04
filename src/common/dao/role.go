package dao

import (
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

func GetRole(role model.Role, selectedFields ...string) (*model.Role, error) {
	o := orm.NewOrm()
	err := o.Read(&role, selectedFields...)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func GetRoles(role model.Role) ([]*model.Role, error) {
	o := orm.NewOrm()
	var roles []*model.Role
	_, err := o.QueryTable("role").All(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
