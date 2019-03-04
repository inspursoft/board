package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	helmpkg "git/inspursoft/board/src/apiserver/service/helm"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

var (
	defaultRepoID int64 = 1

	NotExistError = fmt.Errorf("does not exist")
)

type matchedServiceAndDeployment struct {
	Service    *model.K8sInfo
	Deployment *model.K8sInfo
}

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

func toEntry(repo *model.Repository) *helmpkg.Entry {
	entry := &helmpkg.Entry{
		Name: repo.Name,
		URL:  repo.URL,
		Type: repo.Type,
	}
	return entry
}

func GetRepoDetail(repo *model.Repository, nameRegex string) (*model.RepositoryDetail, error) {
	var err error
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
			if nameRegex == "" || strings.Contains(chartname, nameRegex) {
				allCharts = append(allCharts, &model.ChartVersions{
					Name:     chartname,
					Versions: versions,
				})
			}
		}
		//sort the chart by chartname
		sort.Sort(helmpkg.SortChartVersionsByName(allCharts))
		detail.ChartVersionsList = allCharts
	}

	return &detail, nil
}

func GetPaginatedRepoDetail(repo *model.Repository, nameRegex string, pageIndex, pageSize int) (*model.PaginatedRepositoryDetail, error) {
	detail, err := GetRepoDetail(repo, nameRegex)
	if err != nil {
		return nil, err
	}

	return paginateRepositoryDetail(detail, pageIndex, pageSize), nil
}

func paginateRepositoryDetail(detail *model.RepositoryDetail, pageIndex, pageSize int) *model.PaginatedRepositoryDetail {
	var pagedDetail model.PaginatedRepositoryDetail
	pagedDetail.Repository = detail.Repository

	page := &model.Pagination{
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: int64(len(detail.ChartVersionsList)),
	}
	pageCount := page.GetPageCount()
	if pageIndex < pageCount {
		pagedDetail.PaginatedChartVersions = model.PaginatedChartVersions{
			Pagination:        page,
			ChartVersionsList: detail.ChartVersionsList[(pageIndex-1)*pageSize : pageIndex*pageSize],
		}
	} else if pageIndex == pageCount {
		pagedDetail.PaginatedChartVersions = model.PaginatedChartVersions{
			Pagination:        page,
			ChartVersionsList: detail.ChartVersionsList[(pageIndex-1)*pageSize:],
		}
	} //else do nothing when the pageIndex out of band.
	return &pagedDetail
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
	workloads, err := chartrepo.InstallChart(chartname, chartversion, name, projectname, values, helmhost)
	if err != nil {
		return err
	}

	err = initBoardReleaseFromKubernetesHelmRelease(name, workloads, projectid, projectname, repo.ID, repo.Name, ownerid, ownername)
	if err != nil {
		// remove the release in the helm, ignore the err
		removeerr := helmpkg.DeleteRelease(name, helmhost)
		if removeerr != nil {
			logs.Warning("remove the release %s error: %+v", name, removeerr)
		}
	}
	return nil
}

func initBoardReleaseFromKubernetesHelmRelease(name, workloads string, projectid int64, projectname string, repoid int64, reponame string, ownerid int64, ownername string) error {
	r, err := insertReleaseIntoDatabase(name, workloads, projectid, projectname, repoid, reponame, ownerid, ownername)
	if err != nil {
		return err
	}

	err = addHelmReleaseToBoardService(r)
	if err != nil {
		// remove the release from database
		dao.DeleteHelmRelease(model.ReleaseModel{ID: r.ID})
	}
	return nil
}

func insertReleaseIntoDatabase(name, workloads string, projectid int64, projectname string, repoid int64, reponame string, ownerid int64, ownername string) (*model.ReleaseModel, error) {
	// retrieve the release detail info
	r, err := helmpkg.GetRelease(name, utils.GetStringValue("HELM_HOST"))
	if err != nil {
		logs.Error("Get release %s info from helm error:%+v", name, err)
		return nil, err
	}
	var update time.Time
	if r.Updated != "-" {
		update, err = time.Parse(time.ANSIC, r.Updated)
		if err != nil {
			logs.Warning("Parse the release %s time error: %+v", name, r.Updated)
			err = nil //ignore this err
		}
	}
	m := model.ReleaseModel{
		Name:           name,
		ProjectID:      projectid,
		ProjectName:    projectname,
		Workloads:      workloads,
		RepositoryID:   repoid,
		RepostiroyName: reponame,
		OwnerID:        ownerid,
		OwnerName:      ownername,
		UpdateTime:     update,
		CreateTime:     update,
	}
	id, err := dao.AddHelmRelease(m)
	if err != nil {
		return nil, err
	}
	m.ID = id
	return &m, nil
}

func addHelmReleaseToBoardService(r *model.ReleaseModel) error {
	// add the kubernetes resources to board
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	//resolve the templateInfo into kubernetes service and deployments.....
	return model.NewK8sHelper().Visit(r.Workloads, func(infos []*model.K8sInfo, err error) error {
		if err != nil {
			return err
		}
		matched, err := getBoardMatchedServiceAndDeployments(r.ProjectName, infos)
		if err != nil {
			return err
		}

		// add the service into board
		for i := range matched {
			// check service existents
			exist, err := ServiceExists(matched[i].Service.Name, matched[i].Service.Namespace)
			if err != nil {
				return err
			}
			if exist {
				return fmt.Errorf("The service %s has already exist in project %s", matched[i].Service.Name, matched[i].Service.Namespace)
			}
			svcType, err := getServiceTypeFromYaml(k8sclient, matched[i].Service.Namespace, matched[i].Service.Source)
			if err != nil {
				return err
			}

			_, err = CreateServiceConfig(model.ServiceStatus{
				Name:           matched[i].Service.Name,
				ProjectID:      r.ProjectID,
				ProjectName:    r.ProjectName,
				Comment:        "service created by helm",
				OwnerID:        r.OwnerID,
				OwnerName:      r.OwnerName,
				Status:         defaultStatus,
				Type:           svcType,
				Public:         defaultPublic,
				CreationTime:   r.CreateTime,
				UpdateTime:     r.UpdateTime,
				Source:         helm,
				ServiceYaml:    matched[i].Service.Source,
				DeploymentYaml: matched[i].Deployment.Source,
			})
			if err != nil {
				return err
			}
			break
		}

		return nil
	})
}

func getBoardMatchedServiceAndDeployments(namespace string, infos []*model.K8sInfo) ([]matchedServiceAndDeployment, error) {
	svcs, err := rivers.FromSlice(infos).Filter(func(t stream.T) bool {
		info := t.(*model.K8sInfo)
		if info.Kind == "Service" {
			info.Namespace = namespace
			return true
		}
		return false
	}).GroupBy(func(t stream.T) stream.T {
		info := t.(*model.K8sInfo)
		return model.K8sInfo{Name: info.Name, Namespace: info.Namespace}
	})
	if err != nil {
		return nil, err
	}

	matched := []matchedServiceAndDeployment{}
	err = rivers.FromSlice(infos).Filter(func(t stream.T) bool {
		info := t.(*model.K8sInfo)
		if info.Kind == "Deployment" {
			info.Namespace = namespace
			return true
		}
		return false
	}).OnData(func(t stream.T, emitter stream.Emitter) {
		info := t.(*model.K8sInfo)
		if matchedSvcs, ok := svcs[model.K8sInfo{Name: info.Name, Namespace: info.Namespace}]; ok {
			svc := matchedSvcs[0].(*model.K8sInfo)
			svc.Source, err = setYamlNamespace(svc.Source, svc.Namespace)
			if err != nil {
				logs.Warning("set the service namespace error:%+v", err)
				return
			}
			info.Source, err = setYamlNamespace(info.Source, info.Namespace)
			if err != nil {
				logs.Warning("set the service namespace error:%+v", err)
				return
			}
			emitter.Emit(matchedServiceAndDeployment{Service: svc, Deployment: info})
		}
	}).CollectAs(&matched)
	if err != nil {
		return nil, err
	}
	return matched, nil
}

func getServiceTypeFromYaml(k8sclient *k8sassist.K8sAssistClient, namespace, yaml string) (int, error) {
	svcType := model.ServiceTypeUnknown
	//check service type
	svcModel, err := k8sclient.AppV1().Service(namespace).CheckYaml(strings.NewReader(yaml))
	if err != nil {
		return svcType, err
	}

	switch svcModel.Type {
	case "ClusterIP":
		svcType = model.ServiceTypeClusterIP
	case "NodePort":
		svcType = model.ServiceTypeNormalNodePort
	}
	return svcType, nil
}

func ListReleases(repo *model.Repository, userid int64) ([]model.Release, error) {
	var models []model.ReleaseModel
	var err error
	if repo == nil {
		models, err = dao.GetHelmReleases(&dao.ReleaseFilter{OwnerID: userid})
	} else {
		models, err = dao.GetHelmReleases(&dao.ReleaseFilter{RepositoryID: repo.ID, OwnerID: userid})
	}
	// get the releases from helm cmd
	list, err := helmpkg.ListAllReleases(utils.GetStringValue("HELM_HOST"))
	if err != nil {
		return nil, err
	}
	releases := map[string]helmpkg.Release{}
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
	err = helmpkg.DeleteRelease(release.Name, utils.GetStringValue("HELM_HOST"))
	if err != nil {
		return err
	}
	_, err = dao.DeleteHelmRelease(model.ReleaseModel{
		ID: releaseid,
	})
	if err != nil {
		return err
	}
	// remove the service entry from database
	// add the kubernetes resources to board
	//resolve the templateInfo into kubernetes service and deployments.....
	err = model.NewK8sHelper().Visit(release.Workloads, func(infos []*model.K8sInfo, err error) error {
		if err != nil {
			return err
		}
		matched, err := getBoardMatchedServiceAndDeployments(release.ProjectName, infos)
		if err != nil {
			return err
		}

		// delete the service from board
		for i := range matched {
			dao.DeleteServiceByNameAndProjectName(model.ServiceStatus{
				Name:        matched[i].Service.Name,
				ProjectName: matched[i].Service.Namespace,
			})
		}

		return nil
	})
	return err
}

func GetReleaseDetail(releaseid int64) (*model.ReleaseDetail, error) {
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

	helmhost := utils.GetStringValue("HELM_HOST")
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
		ProjectID:      m.ProjectID,
		ProjectName:    m.ProjectName,
		RepositoryID:   m.RepositoryID,
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

func CheckReleaseNames(name string) (bool, error) {
	models, err := dao.GetHelmReleases(&dao.ReleaseFilter{Name: name})
	if err != nil {
		return false, err
	}

	if len(models) > 0 {
		return true, nil
	}

	// get the releases from helm cmd
	list, err := helmpkg.ListAllReleases(utils.GetStringValue("HELM_HOST"))
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

func setYamlNamespace(yamlstr, namespace string) (string, error) {
	return model.NewK8sHelper().Transform(yamlstr, func(in interface{}) (interface{}, error) {
		switch s := in.(type) {
		case map[string]interface{}:
			setNamespace(s, namespace)
			return s, nil
		default:
			return nil, fmt.Errorf("The type of object %+v is unknown", in)
		}
	})
}

func setNamespace(obj map[string]interface{}, value interface{}) {
	utils.SetNestedField(obj, value, "metadata", "namespace")
}

func SyncHelmReleaseWithK8s(projectname string) error {
	project, err := GetProjectByName(projectname)
	if err != nil {
		return err
	}
	// get repo
	repo, err := GetRepository(defaultRepoID)
	if err != nil {
		return err
	}
	// get releases from kubernetes
	// get the releases from helm cmd
	list, err := helmpkg.ListDeployedReleasesByNamespace(utils.GetStringValue("HELM_HOST"), projectname)
	if err != nil {
		return err
	}

	// get release from dao
	helmhost := utils.GetStringValue("HELM_HOST")
	models, err := dao.GetHelmReleases(nil)
	for i := range list.Releases {
		exist := false
		for j := range models {
			if list.Releases[i].Name == models[j].Name {
				exist = true
				break
			}
		}
		// sync the release which does not exist in board release table.
		if !exist {
			// get the release manifest
			load, err := helmpkg.GetReleaseManifest(list.Releases[i].Name, helmhost)
			if err != nil {
				logs.Warning("Get release %s workloads when synchronizing error:%+v", list.Releases[i].Name, err)
				continue
			}
			// add the release into database.
			err = initBoardReleaseFromKubernetesHelmRelease(list.Releases[i].Name, load, project.ID, project.Name, repo.ID, repo.Name, int64(project.OwnerID), project.OwnerName)
			if err != nil {
				// ignore the initReleaseError
				logs.Warning("sync release from kubernetes release %s error: %+v", list.Releases[i].Name, err)
			}
		}
	}
	return nil
}
