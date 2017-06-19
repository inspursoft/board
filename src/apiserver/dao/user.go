package dao

import (
	"apiserver/model"

	"github.com/astaxie/beego/orm"
)

func AddUser(user model.User) (int64, error) {
	o := orm.NewOrm()
	userID, err := o.Insert(&user)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

func UpdateUser(user model.User, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	userID, err := o.Update(&user, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

func GetUser(user model.User, fieldNames ...string) (*model.User, error) {
	o := orm.NewOrm()
	err := o.Read(&user, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return &user, nil
		}
		return nil, err
	}
	return &user, err
}

func GetUsers(field string, value interface{}, selectedFields ...string) ([]*model.User, error) {
	o := orm.NewOrm()
	var users []*model.User
	var err error
	qs := o.QueryTable("user")
	if value == nil {
		_, err = qs.All(&users, selectedFields...)
	} else {
		_, err = qs.Filter(field+"__contains", value).
			All(&users, selectedFields...)
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(user model.User) (int64, error) {
	o := orm.NewOrm()
	return o.Delete(&user)
}
