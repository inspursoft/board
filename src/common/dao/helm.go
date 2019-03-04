package dao

import (
	"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego/orm"
)

type ReleaseFilter struct {
	RepositoryID int64
	OwnerID      int64
	Name         string
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

func GetHelmReleases(filter *ReleaseFilter) ([]model.ReleaseModel, error) {
	releases := make([]model.ReleaseModel, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("release").OrderBy("-creation_time")

	if filter != nil {
		if filter.OwnerID <= 0 {
			qs = qs.Filter("owner_id", filter.OwnerID)
		}

		if filter.RepositoryID <= 0 {
			qs = qs.Filter("repository_id", filter.RepositoryID)
		}
		if filter.Name != "" {
			qs = qs.Filter("name", filter.Name)
		}
	}

	_, err := qs.All(&releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}
