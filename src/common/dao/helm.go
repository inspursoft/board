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

func AddHelmRelease(release model.ReleaseModel) (int64, error) {
	o := orm.NewOrm()
	releaseId, err := o.Insert(&release)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	// add the release and service relations
	if release.Services != nil {
		ptmt, err := o.Raw(`insert into release_service(releaseid, serviceid) values (?, ?)`).Prepare()
		if err != nil {
			return 0, err
		}
		defer ptmt.Close()
		for i := range release.Services {
			_, err := ptmt.Exec(releaseId, release.Services[i].ServiceId)
			if err != nil {
				return 0, err
			}
		}
	}
	return releaseId, err
}

func GetHelmRelease(release model.ReleaseModel) (*model.Release, error) {
	o := orm.NewOrm()
	releaseRs := o.Raw(` select board.release.id as i_d, board.release.name as name, board.release.project_id as project_id, 
								board.release.repoid as repository_id, board.release.chart as chart, board.release.chartversion as chart_version, 
								board.release.value as value, board.release.workload as workloads 
								from board.release 
								where board.release.id=?`, release.ID)
	ret := new(model.Release)
	err := releaseRs.QueryRow(ret)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	relations := make([]model.ReleaseService, 0, 0)
	rs := o.Raw(`select release_service.releaseid as release_id, release_service.serviceid as service_id, service_status.name as service_name 
						from release_service left join service_status on release_service.serviceid=service_status.id
						where release_service.releaseid=?`, release.ID)
	_, err = rs.QueryRows(&relations)
	if err != nil {
		if err != orm.ErrNoRows {
			return nil, err
		}
	}

	ret.Services = relations
	return ret, err
}

func DeleteHelmRelease(release model.ReleaseModel) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&release)
	// the release and service relation was deleted by cascade.
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func GetHelmReleasesByRepositoryId(repoid int64) ([]model.Release, error) {
	o := orm.NewOrm()
	releaseRs := o.Raw(` select board.release.id as i_d, board.release.name as name, board.release.project_id as project_id, 
								board.release.repoid as repository_id, board.release.chart as chart, board.release.chartversion as chart_version, 
								board.release.value as value, board.release.workload as workloads 
								from board.release
								where board.release.repoid=?`, repoid)
	ret := make([]model.Release, 0, 0)
	_, err := releaseRs.QueryRows(&ret)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	relations := make([]model.ReleaseService, 0, 0)
	rs := o.Raw(`select release_service.releaseid as release_id, release_service.serviceid as service_id, service_status.name as service_name 
						from release_service left join service_status on release_service.serviceid=service_status.id
						left join board.release on release_service.releaseid = board.release.id
						where board.release.repoid=?`, repoid)
	_, err = rs.QueryRows(&relations)
	if err != nil {
		if err != orm.ErrNoRows {
			return nil, err
		}
	}

	//make a map for release
	releaseMap := make(map[int64]*model.Release, len(ret))
	for i := range ret {
		releaseMap[ret[i].ID] = &ret[i]
	}
	for i := range relations {
		releaseMap[relations[i].ReleaseId].Services = append(releaseMap[relations[i].ReleaseId].Services, relations[i])
	}
	return ret, err
}

func GetHelmReleases() ([]model.Release, error) {
	o := orm.NewOrm()
	releaseRs := o.Raw(` select board.release.id as i_d, board.release.name as name, board.release.project_id as project_id, 
								board.release.repoid as repository_id, board.release.chart as chart, board.release.chartversion as chart_version, 
								board.release.value as value, board.release.workload as workloads 
								from board.release`)
	ret := make([]model.Release, 0, 0)
	_, err := releaseRs.QueryRows(&ret)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	relations := make([]model.ReleaseService, 0, 0)
	rs := o.Raw(`select release_service.releaseid as release_id, release_service.serviceid as service_id, service_status.name as service_name 
						from release_service left join service_status on release_service.serviceid=service_status.id
						left join board.release on release_service.releaseid = board.release.id`)
	_, err = rs.QueryRows(&relations)
	if err != nil {
		if err != orm.ErrNoRows {
			return nil, err
		}
	}

	//make a map for release
	releaseMap := make(map[int64]*model.Release, len(ret))
	for i := range ret {
		releaseMap[ret[i].ID] = &ret[i]
	}
	for i := range relations {
		releaseMap[relations[i].ReleaseId].Services = append(releaseMap[relations[i].ReleaseId].Services, relations[i])
	}
	return ret, err
}
