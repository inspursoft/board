package dao

import (
	"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego/orm"
)

type filterCondition struct {
	name   string
	values []interface{}
}

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
	_, err := o.QueryTable("repository").OrderBy("name").All(&repos, selectedFields...)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func AddHelmRelease(release model.ReleaseModel) (int64, error) {
	o := orm.NewOrm()
	releaseId, err := o.Insert(&release)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return releaseId, err
}

func GetHelmRelease(release model.ReleaseModel) (*model.ReleaseModel, error) {
	o := orm.NewOrm()
	err := o.Read(&release)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &release, err
}

func DeleteHelmRelease(release model.ReleaseModel) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&release)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func GetHelmReleasesByRepositoryAndUser(repoid, userid int64) ([]model.ReleaseModel, error) {
	filters := []filterCondition{}
	if repoid > 0 {
		filters = append(filters, filterCondition{
			name:   "repoid",
			values: []interface{}{repoid},
		})
	}
	if userid > 0 {
		filters = append(filters, filterCondition{
			name:   "owner_id",
			values: []interface{}{userid},
		})
	}
	return GetHelmReleases(filters...)
}

func GetHelmReleases(filters ...filterCondition) ([]model.ReleaseModel, error) {
	releases := make([]model.ReleaseModel, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("release").OrderBy("-creation_time")
	for _, filter := range filters {
		if filter.name != "" {
			qs = qs.Filter(filter.name, filter.values...)
		}
	}
	_, err := qs.All(&releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}
