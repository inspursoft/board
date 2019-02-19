package service

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	helmpkg "git/inspursoft/board/src/apiserver/service/helm"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
)

var NotExistError = fmt.Errorf("does not exist")

func ListRepositories() ([]model.Repository, error) {
	return dao.GetHelmRepositories()
}

func AddRepository(repo model.Repository) (int64, error) {
	return dao.AddHelmRepository(repo)
}

func GetRepository(id int64) (*model.Repository, error) {
	repo := model.Repository{
		ID: id,
	}

	r, err := dao.GetHelmRepository(repo)
	if err != nil {
		return nil, err
	}
	if r == nil {
		logs.Error("the repository %d does not exist", id)
		return nil, NotExistError
	}
	return r, nil
}

func UpdateRepository(repo model.Repository) error {
	id, err := dao.UpdateHelmRepository(repo)
	if err != nil {
		return err
	}
	if id == 0 {
		logs.Error("the repository %d does not exist", repo.ID)
		return NotExistError
	}
	return nil
}

func DeleteRepository(id int64) error {
	repo, err := dao.DeleteHelmRepository(model.Repository{
		ID: id,
	})
	if err != nil {
		return err
	}
	if repo == 0 {
		logs.Error("the repository %d does not exist", id)
		return NotExistError
	}
	return nil
}

func CheckRepoNameNotExist(name string) (bool, error) {
	// get the hpaname from storage
	repos, err := dao.GetHelmRepositories()
	if err != nil {
		return false, err
	}
	for i := range repos {
		if repos[i].Name == name {
			return true, nil
		}
	}
	return false, nil
}

func toEntry(repo *model.Repository) *helmpkg.Entry {
	entry := &helmpkg.Entry{
		Name: repo.Name,
		URL:  repo.URL,
		Type: repo.Type,
	}
	return entry
}

func GetRepoDetail(repo *model.Repository, nameRegex string, pageIndex, pageSize int) (*model.RepositoryDetail, error) {
	var err error
	var filter *regexp.Regexp
	if nameRegex != "" {
		filter, err = regexp.Compile(nameRegex)
		if err != nil {
			return nil, err
		}
	}
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return nil, err
	}
	var detail model.RepositoryDetail
	detail.Repository = *repo

	if chartrepo.IndexFile != nil {
		var allCharts []*model.ChartVersions
		//filter the results.
		for chartname, versions := range chartrepo.IndexFile.Entries {
			if filter == nil || filter.FindString(chartname) != "" {
				allCharts = append(allCharts, &model.ChartVersions{
					Name:     chartname,
					Versions: versions,
				})
			}
		}
		//sort the chart by chartname
		sort.Sort(helmpkg.SortChartVersionsByName(allCharts))
		page := &model.Pagination{
			PageIndex:  pageIndex,
			PageSize:   pageSize,
			TotalCount: int64(len(allCharts)),
		}
		page.GetPageCount()
		if pageIndex == 0 {
			detail.PaginatedChartVersions = model.PaginatedChartVersions{
				Pagination:        page,
				ChartVersionsList: allCharts,
			}
		} else if (pageIndex-1)*pageSize < len(allCharts) {
			if pageSize == 0 || pageIndex*pageSize >= len(allCharts) {
				detail.PaginatedChartVersions = model.PaginatedChartVersions{
					Pagination:        page,
					ChartVersionsList: allCharts[(pageIndex-1)*pageSize:],
				}
			} else {
				detail.PaginatedChartVersions = model.PaginatedChartVersions{
					Pagination:        page,
					ChartVersionsList: allCharts[(pageIndex-1)*pageSize : pageIndex*pageSize],
				}
			}
		}
	}

	return &detail, nil
}

func GetChartDetail(repo *model.Repository, chartname, chartversion string) (*model.Chart, error) {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return nil, err
	}
	chart, err := chartrepo.FetchChart(chartname, chartversion)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

func UploadChart(repo *model.Repository, chartfile string) error {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return err
	}
	return chartrepo.UploadChart(chartfile)
}

func DeleteChart(repo *model.Repository, chartname, chartversion string) error {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return err
	}
	return chartrepo.DeleteChart(chartname, chartversion)
}

func InstallChart(repo *model.Repository, chartname, chartversion, name string, projectid int64, projectname, values string, ownerid int64, ownername string) error {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return err
	}
	helmhost := utils.GetStringValue("HELM_HOST")
	if helmhost == "" {
		return fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
	workloads, err := chartrepo.InstallChart(chartname, chartversion, name, projectname, values, helmhost)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// remove the release in the helm, ignore the err
			removeerr := helmpkg.DeleteReleaseFromRepository(name, helmhost)
			if removeerr != nil {
				logs.Warning("remove the release %s error: %+v", name, removeerr)
			}
		}
	}()
	// retrieve the release detail info
	r, err := helmpkg.GetRelease(name, helmhost)
	if err != nil {
		logs.Warning("Get release %s info from helm error:%+v", name, err)
		return err
	}
	var update time.Time
	if r.Updated != "-" {
		update, err = time.Parse(time.ANSIC, r.Updated)
		if err != nil {
			logs.Warning("Parse the release %s time error: %+v", name, r.Updated)
			err = nil //ignore this err
		}
	}
	model := model.ReleaseModel{
		Name:           name,
		ProjectId:      projectid,
		ProjectName:    projectname,
		Workloads:      workloads,
		RepositoryId:   repo.ID,
		RepostiroyName: repo.Name,
		OwnerID:        ownerid,
		OwnerName:      ownername,
		UpdateTime:     update,
		CreateTime:     update,
	}
	_, err = dao.AddHelmRelease(model)
	return err
}

func ListReleases(repo *model.Repository, userid int64) ([]model.Release, error) {
	var models []model.ReleaseModel
	var err error
	if repo == nil {
		models, err = dao.GetHelmReleasesByRepositoryAndUser(-1, userid)
	} else {
		models, err = dao.GetHelmReleasesByRepositoryAndUser(repo.ID, userid)
	}
	// get the releases from helm cmd
	helmhost := utils.GetStringValue("HELM_HOST")
	if helmhost == "" {
		return nil, fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
	list, err := helmpkg.ListReleases(helmhost)
	if err != nil {
		return nil, err
	}
	releases := make(map[string]helmpkg.Release)
	for i := range list.Releases {
		releases[list.Releases[i].Name] = list.Releases[i]
	}
	// list the release base on the database.
	var ret []model.Release
	for _, m := range models {
		if r, ok := releases[m.Name]; ok {
			ret = append(ret, generateModelRelease(&m, &r))
		} else {
			logs.Warning("the release %s does not exist in kubernetes", m.Name)
		}
	}
	return ret, err
}

func GetReleaseFromDB(releaseid int64) (*model.ReleaseModel, error) {
	return dao.GetHelmRelease(model.ReleaseModel{
		ID: releaseid,
	})
}

func DeleteRelease(releaseid int64) error {
	helmhost := utils.GetStringValue("HELM_HOST")
	if helmhost == "" {
		return fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
	release, err := dao.GetHelmRelease(model.ReleaseModel{
		ID: releaseid,
	})
	if err != nil {
		return err
	}
	if release == nil {
		logs.Error("the release with id %d does not exist", releaseid)
		return NotExistError
	}
	err = helmpkg.DeleteReleaseFromRepository(release.Name, helmhost)
	if err != nil {
		return err
	}
	_, err = dao.DeleteHelmRelease(model.ReleaseModel{
		ID: releaseid,
	})
	return err
}

func GetReleaseDetail(releaseid int64) (*model.ReleaseDetail, error) {
	helmhost := utils.GetStringValue("HELM_HOST")
	if helmhost == "" {
		return nil, fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
	release, err := dao.GetHelmRelease(model.ReleaseModel{
		ID: releaseid,
	})
	if err != nil {
		return nil, err
	}
	if release == nil {
		logs.Error("the release with id %d does not exist", releaseid)
		return nil, NotExistError
	}
	releaseChan := make(chan *helmpkg.Release)
	go func() {
		r, err := helmpkg.GetRelease(release.Name, helmhost)
		if err != nil {
			logs.Warning("Get release %s from helm error:%+v", release.Name, err)
		}
		releaseChan <- r
	}()

	notesChan := make(chan string)
	go func() {
		note, err := helmpkg.GetReleaseNotes(release.Name, helmhost)
		if err != nil {
			logs.Warning("Get release %s notes from helm error:%+v", release.Name, err)
		}
		notesChan <- note
	}()

	statusChan := make(chan string)
	go func() {
		status, err := helmpkg.GetReleaseStatus(release.Name, helmhost)
		if err != nil {
			logs.Warning("Get release %s status from helm error:%+v", release.Name, err)
		}
		statusChan <- status
	}()

	loadChan := make(chan string)
	go func() {
		load, err := helmpkg.GetReleaseManifest(release.Name, helmhost)
		if err != nil {
			logs.Warning("Get release %s workloads from helm error:%+v", release.Name, err)
		}
		loadChan <- load
	}()

	//get the result
	helmrelease := <-releaseChan

	notes := <-notesChan
	workloads := <-loadChan
	status := <-statusChan

	detail := model.ReleaseDetail{
		Release:        generateModelRelease(release, helmrelease),
		Workloads:      workloads,
		Notes:          notes,
		WorkloadStatus: status,
	}
	return &detail, err
}

func generateModelRelease(m *model.ReleaseModel, r *helmpkg.Release) model.Release {
	var chart, version, status string
	if r != nil {
		index := strings.LastIndex(r.Chart, "-")
		chart = r.Chart
		if index != -1 {
			chart = r.Chart[:index]
			if index != len(r.Chart)-1 {
				version = r.Chart[index+1:]
			}
		}
		status = r.Status
	}

	return model.Release{
		ID:             m.ID,
		Name:           m.Name,
		ProjectId:      m.ProjectId,
		ProjectName:    m.ProjectName,
		RepositoryId:   m.RepositoryId,
		RepositoryName: m.RepostiroyName,
		Chart:          chart,
		ChartVersion:   version,
		OwnerID:        m.OwnerID,
		OwnerName:      m.OwnerName,
		Status:         status,
		UpdateTime:     m.UpdateTime,
		CreateTime:     m.CreateTime,
	}
}

func GetHelmReleaseResources() ([]*model.Info, error) {
	models, err := dao.GetHelmReleasesByRepositoryAndUser(-1, -1)
	if err != nil {
		return nil, err
	}

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	//resolve the templateInfo into kubernetes service and deployments.....
	mapper := k8sclient.AppV1().Mapper()

	ret := []*model.Info{}
	for _, m := range models {
		mapper.Visit(m.Workloads, func(infos []*model.Info, err error) error {
			if err != nil {
				logs.Warning("ignore analysis the workload error: %s", err.Error())
			}
			for i := range infos {
				if infos[i].Namespace == "" {
					infos[i].Namespace = m.ProjectName
				}
			}
			ret = append(ret, infos...)
			return nil
		})
	}
	return ret, nil
}

func CheckReleaseNames(name string) (bool, error) {
	models, err := dao.GetHelmReleasesByRepositoryAndUser(-1, -1)
	if err != nil {
		return false, err
	}
	for i := range models {
		if models[i].Name == name {
			return true, nil
		}
	}
	// get the releases from helm cmd
	helmhost := utils.GetStringValue("HELM_HOST")
	if helmhost == "" {
		return false, fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
	list, err := helmpkg.ListReleases(helmhost)
	if err != nil {
		return false, err
	}

	for i := range list.Releases {
		if list.Releases[i].Name == name {
			return true, nil
		}
	}
	return false, nil
}
