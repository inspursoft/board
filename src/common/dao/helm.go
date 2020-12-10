package dao

import (
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

func AddHelmRepository(repo model.HelmRepository) (int64, error) {
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

func GetHelmRepository(repo model.HelmRepository, fieldNames ...string) (*model.HelmRepository, error) {
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

func UpdateHelmRepository(repo model.HelmRepository, fieldNames ...string) (int64, error) {
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

func DeleteHelmRepository(repo model.HelmRepository) (int64, error) {
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

func GetHelmRepositories(selectedFields ...string) ([]model.HelmRepository, error) {
	repos := make([]model.HelmRepository, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable("helm_repository").OrderBy("name").All(&repos, selectedFields...)
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

func GetHelmReleaseByName(name string) (*model.ReleaseModel, error) {
	release := &model.ReleaseModel{Name: name}
	o := orm.NewOrm()
	err := o.Read(release, "name")
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return release, err
}

func GetAllHelmReleases() ([]model.ReleaseModel, error) {
	releases := make([]model.ReleaseModel, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable("helm_release").OrderBy("-creation_time").All(&releases)

	if err != nil {
		return nil, err
	}
	return releases, nil
}

func GetAllHelmReleasesByProjectName(projectName string) ([]model.ReleaseModel, error) {
	releases := make([]model.ReleaseModel, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable("helm_release").Filter("project_name", projectName).OrderBy("-creation_time").All(&releases)

	if err != nil {
		return nil, err
	}
	return releases, nil
}

func GetHelmReleasesByUserID(userID int64) ([]model.ReleaseModel, error) {
	sql, params := generateHelmReleasesSQL(userID, "")
	return queryHelmReleases(sql, params)
}

func GetHelmReleasesByUserIDAndProjectName(userID int64, projectName string) ([]model.ReleaseModel, error) {
	sql, params := generateHelmReleasesSQL(userID, projectName)
	return queryHelmReleases(sql, params)
}

func generateHelmReleasesSQL(userID int64, projectName string) (string, []interface{}) {
	sql := `select distinct hr.id, hr.name, hr.project_id, hr.project_name, hr.repository_id, hr.repository, hr.workloads, hr.owner_id, hr.owner_name, hr.creation_time, hr.update_time
	from helm_release hr 
		left join project_member pm on hr.project_id = pm.project_id
		left join project p on p.id = pm.project_id
		left join user u on u.id = hr.owner_id
	where hr.owner_id = ? and (hr.project_id in (select p.id from project p left join project_member pm on p.id = pm.project_id  left join user u on u.id = pm.user_id where p.deleted = 0 and u.deleted = 0 and u.id = ?)
		or exists (select * from user u where u.deleted = 0 and u.system_admin = 1 and u.id = ?))`
	values := []interface{}{userID, userID, userID}
	if projectName != "" {
		sql += " and hr.project_name= ?"
		values = append(values, projectName)
	}
	sql += " order by creation_time desc"
	return sql, values
}

func queryHelmReleases(sql string, params []interface{}) ([]model.ReleaseModel, error) {
	releaseList := make([]model.ReleaseModel, 0)
	_, err := orm.NewOrm().Raw(sql, params).QueryRows(&releaseList)
	if err != nil {
		return nil, err
	}
	return releaseList, nil
}
