package dao

import (
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

func GetAllConfigs() ([]*model.Config, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("config")
	var configs []*model.Config
	_, err := qs.All(&configs)
	return configs, err
}
func GetConfig(name string) (*model.Config, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("config")
	var config model.Config
	_, err := qs.Filter("name__eq", name).All(&config)
	return &config, err
}

func AddOrUpdateConfig(config model.Config) (int64, error) {
	o := orm.NewOrm()
	ptmt, err := o.Raw(`insert into config
		 (name, value, comment) 
		 values (?, ?, ?)
			on duplicate key 
			update value = ?, comment = ?`).Prepare()
	if err != nil {
		return 0, err
	}
	defer ptmt.Close()
	r, err := ptmt.Exec(config.Name, config.Value, config.Comment, config.Value, config.Comment)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return r.RowsAffected()
}

func DeleteConfig(name string) (int64, error) {
	o := orm.NewOrm()
	ptmt, err := o.Raw(`delete from config where name = ?`).Prepare()
	if err != nil {
		return 0, err
	}
	defer ptmt.Close()
	r, err := ptmt.Exec(name)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return r.RowsAffected()
}
