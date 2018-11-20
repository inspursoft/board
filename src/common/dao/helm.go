package dao

import (
	"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego/orm"
)

func AddHelmRepository(repo model.Repository) (int64, error) {
	o := orm.NewOrm()
	repoId, err := o.Insert(&repo)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return repoId, err
}

func GetHelmRepository(repo model.Repository, fieldNames ...string) (*model.Repository, error) {
	o := orm.NewOrm()
	err := o.Read(&repo, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &repo, err
}

func UpdateHelmRepository(repo model.Repository, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	repoId, err := o.Update(&repo, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return repoId, err
}

func DeleteHelmRepository(repo model.Repository) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&repo)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func GetHelmRepositories(selectedFields ...string) ([]model.Repository, error) {
	repos := make([]model.Repository, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable("repository").All(&repos, selectedFields...)
	if err != nil {
		return nil, err
	}
	return repos, nil
}
