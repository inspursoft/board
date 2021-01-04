package dao

import (
	"fmt"
	"github.com/inspursoft/board/src/common/model"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddUser(user model.User) (int64, error) {
	o := orm.NewOrm()
	user.CreationTime = time.Now()
	user.UpdateTime = user.CreationTime
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
	user.UpdateTime = time.Now()
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
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}

func queryUsers(field string, value interface{}) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable("user")
	if value == nil {
		return qs
	}
	return qs.Filter(field, value)
}

func GetUsers(field string, value interface{}, selectedFields ...string) ([]*model.User, error) {
	users := make([]*model.User, 0)
	_, err := queryUsers(field, value).All(&users, selectedFields...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetPaginatedUsers(field string, value interface{}, pageIndex int, pageSize int, orderField string, orderAsc int, selectedFields ...string) (*model.PaginatedUsers, error) {
	pagination := &model.Pagination{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}

	qs := queryUsers(field, value)

	var err error
	pagination.TotalCount, err = qs.Count()
	if err != nil {
		return nil, err
	}
	qs = qs.OrderBy(getOrderExprs(orderField, orderAsc)).Limit(pagination.PageSize).Offset(pagination.GetPageOffset())
	logs.Debug("%+v", pagination.String())

	users := make([]*model.User, 0)
	_, err = qs.All(&users, selectedFields...)
	if err != nil {
		return nil, err
	}

	return &model.PaginatedUsers{
		Pagination: pagination,
		UserList:   users,
	}, nil
}

func getOrderExprs(orderField string, orderAsc int) string {
	orderStr := strings.ToLower(orderField)
	if orderAsc != 0 {
		return orderStr
	}
	return fmt.Sprintf(`-%s`, orderStr)
}
