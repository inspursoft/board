package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/inspursoft/board/src/apiserver/models/helms/repositories/vm"
	helmpkg "github.com/inspursoft/board/src/apiserver/service/helm"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
)

var (
	defaultRepoID int64 = 1

	NotExistError = fmt.Errorf("does not exist")
)

type WORKLOAD_TYPE string

const (
	WORKLOAD_TYPE_DEPLOYMENT  WORKLOAD_TYPE = "Deployment"
	WORKLOAD_TYPE_STATEFULSET WORKLOAD_TYPE = "StatefulSet"
)

type matchedServiceAndWorkload struct {
	Service     *model.K8sInfo
	ServiceType int
	Type        WORKLOAD_TYPE
	Workload    *model.K8sInfo
}

func vmRepositoryModel(repo model.HelmRepository) vm.HelmRepository {
	return vm.HelmRepository{
		ID:   repo.ID,
		Name: repo.Name,
		URL:  repo.URL,
		Type: repo.Type,
	}
}

func ListVMHelmRepositories() ([]vm.HelmRepository, error) {
	// list the repos from storage
	repos, err := dao.GetHelmRepositories()
	if err != nil {
		return nil, err
	}

	vmRepos := []vm.HelmRepository{}
	for _, r := range repos {
		vmRepos = append(vmRepos, vmRepositoryModel(r))
	}
	return vmRepos, nil
}

func ListHelmRepositories() ([]model.HelmRepository, error) {
	return dao.GetHelmRepositories()
}

func AddHelmRepository(repo model.HelmRepository) (int64, error) {
	return dao.AddHelmRepository(repo)
}

func GetHelmRepository(id int64) (*model.HelmRepository, error) {
	repo := model.HelmRepository{
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

func UpdateHelmRepository(repo model.HelmRepository) error {
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

func DeleteHelmRepository(id int64) error {
	repo, err := dao.DeleteHelmRepository(model.HelmRepository{
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

func toEntry(repo *model.HelmRepository) *helmpkg.Entry {
	entry := &helmpkg.Entry{
		Name: repo.Name,
		URL:  repo.URL,
		Type: repo.Type,
	}
	return entry
}

func GetRepoDetail(repo *model.HelmRepository, nameRegex string) (*model.HelmRepositoryDetail, error) {
	var err error
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return nil, err
	}
	var detail model.HelmRepositoryDetail
	detail.HelmRepository = *repo

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

func GetPaginatedRepoDetail(repo *model.HelmRepository, nameRegex string, pageIndex, pageSize int) (*model.PaginatedHelmRepositoryDetail, error) {
	detail, err := GetRepoDetail(repo, nameRegex)
	if err != nil {
		return nil, err
	}

	return paginateHelmRepositoryDetail(detail, pageIndex, pageSize), nil
}

func paginateHelmRepositoryDetail(detail *model.HelmRepositoryDetail, pageIndex, pageSize int) *model.PaginatedHelmRepositoryDetail {
	var pagedDetail model.PaginatedHelmRepositoryDetail
	pagedDetail.HelmRepository = detail.HelmRepository

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

func GetChartDetail(repo *model.HelmRepository, chartname, chartversion string) (*model.Chart, error) {
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

func UploadChart(repo *model.HelmRepository, chartfile string) error {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return err
	}
	return chartrepo.UploadChart(chartfile)
}

func DeleteChart(repo *model.HelmRepository, chartname, chartversion string) error {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return err
	}
	return chartrepo.DeleteChart(chartname, chartversion)
}

func InstallChart(repo *model.HelmRepository, target *model.Release) error {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return err
	}
	helmhost, err := getHelmHost()
	if err != nil {
		return err
	}
	err = chartrepo.InstallChart(target.Chart, target.ChartVersion, target.Name, target.ProjectName, target.Values, target.Answers, helmhost)
	if err != nil {
		return err
	}

	err = initBoardReleaseFromKubernetesHelmRelease(repo, target)
	if err != nil {
		// remove the release in the helm, ignore the err
		removeerr := helmpkg.DeleteRelease(target.Name, helmhost)
		if removeerr != nil {
			logs.Warning("remove the release %s error: %+v", target.Name, removeerr)
		}
	}
	return nil
}

func initBoardReleaseFromKubernetesHelmRelease(repo *model.HelmRepository, target *model.Release) error {
	r, err := insertReleaseIntoDatabase(repo, target)
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

func insertReleaseIntoDatabase(repo *model.HelmRepository, target *model.Release) (*model.ReleaseModel, error) {
	// retrieve the release detail info
	helmhost, err := getHelmHost()
	if err != nil {
		return nil, err
	}
	r, err := helmpkg.GetRelease(target.Name, helmhost)
	if err != nil {
		logs.Error("Get release %s info from helm error:%+v", target.Name, err)
		return nil, err
	}
	var update time.Time
	if r.Updated != "-" {
		update, err = time.Parse(time.ANSIC, r.Updated)
		if err != nil {
			logs.Warning("Parse the release %s time error: %+v", target.Name, r.Updated)
			err = nil //ignore this err
		}
	}
	// get the release workloads
	workloads, err := helmpkg.GetReleaseManifest(target.Name, helmhost)
	if err != nil {
		logs.Error("Get release %s workloads when synchronizing error:%+v", target.Name, err)
		return nil, err
	}
	m := model.ReleaseModel{
		Name:           target.Name,
		ProjectID:      target.ProjectID,
		ProjectName:    target.ProjectName,
		Workloads:      workloads,
		RepositoryID:   repo.ID,
		RepostiroyName: repo.Name,
		OwnerID:        target.OwnerID,
		OwnerName:      target.OwnerName,
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
	return processBoardMatchedServiceAndWorkloads(r.Workloads, r.ProjectName, func(matched *matchedServiceAndWorkload) (bool, error) {
		// process the statefulset
		if matched.Type == WORKLOAD_TYPE_STATEFULSET {
			matched.ServiceType = model.ServiceTypeStatefulSet
			return true, nil
		}
		// check service existents
		exist, err := ServiceExists(matched.Service.Name, matched.Service.Namespace)
		if err != nil {
			return false, err
		}
		if !exist {
			svcType, err := getServiceTypeFromYaml(k8sclient, matched.Service.Namespace, matched.Service.Source)
			if err != nil {
				return false, err
			}
			matched.ServiceType = svcType
			return true, nil
		}
		return false, nil
	}, func(matched []*matchedServiceAndWorkload) error {
		// add the service into board
		for i := range matched {
			_, err := CreateServiceConfig(model.ServiceStatus{
				Name:           matched[i].Service.Name,
				ProjectID:      r.ProjectID,
				ProjectName:    r.ProjectName,
				Comment:        "service created by helm",
				OwnerID:        r.OwnerID,
				OwnerName:      r.OwnerName,
				Status:         defaultStatus,
				Type:           matched[i].ServiceType,
				Public:         defaultPublic,
				CreationTime:   r.CreateTime,
				UpdateTime:     r.UpdateTime,
				Source:         helm,
				SourceID:       r.ID,
				ServiceYaml:    matched[i].Service.Source,
				DeploymentYaml: matched[i].Workload.Source,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func processBoardMatchedServiceAndWorkloads(workloads, projectname string, prepare func(matched *matchedServiceAndWorkload) (bool, error), process func(matched []*matchedServiceAndWorkload) error) error {
	return model.NewK8sHelper().Visit(workloads, func(infos []*model.K8sInfo, err error) error {
		if err != nil {
			return err
		}
		svcs, err := rivers.FromSlice(infos).Filter(func(t stream.T) bool {
			info := t.(*model.K8sInfo)
			if info.Kind == "Service" {
				info.Namespace = projectname
				return true
			}
			return false
		}).GroupBy(func(t stream.T) stream.T {
			info := t.(*model.K8sInfo)
			return model.K8sInfo{Name: info.Name, Namespace: info.Namespace}
		})
		if err != nil {
			return err
		}

		matched := []*matchedServiceAndWorkload{}
		pipeline := rivers.FromSlice(infos).Filter(func(t stream.T) bool {
			info := t.(*model.K8sInfo)
			if info.Kind == string(WORKLOAD_TYPE_DEPLOYMENT) || info.Kind == string(WORKLOAD_TYPE_STATEFULSET) {
				info.Namespace = projectname
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
				emitter.Emit(&matchedServiceAndWorkload{Service: svc, Workload: info, Type: WORKLOAD_TYPE(info.Kind)})
			}
		})
		if prepare != nil {
			pipeline = pipeline.Apply(&transformers.Observer{
				OnNext: func(data stream.T, emitter stream.Emitter) error {
					send, err := prepare(data.(*matchedServiceAndWorkload))
					if send {
						emitter.Emit(data)
					}
					return err
				},
			})
		}
		err = pipeline.CollectAs(&matched)
		if err != nil {
			return err
		}
		if process != nil {
			return process(matched)
		}
		return nil
	})
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

func ListAllReleases(projectName string) ([]model.Release, error) {
	var err error
	var models []model.ReleaseModel
	if projectName == "" {
		models, err = dao.GetAllHelmReleases()
	} else {
		models, err = dao.GetAllHelmReleasesByProjectName(projectName)
	}
	if err != nil {
		return nil, err
	}
	return generateModelReleases(models)
}

func ListReleasesByUserID(userid int64, projectName string) ([]model.Release, error) {
	var err error
	var models []model.ReleaseModel
	if projectName == "" {
		models, err = dao.GetHelmReleasesByUserID(userid)
	} else {
		models, err = dao.GetHelmReleasesByUserIDAndProjectName(userid, projectName)
	}
	if err != nil {
		return nil, err
	}
	return generateModelReleases(models)
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
	return DeleteReleaseByReleaseModel(release)
}

func DeleteReleaseByReleaseModel(release *model.ReleaseModel) error {
	helmhost, err := getHelmHost()
	if err != nil {
		return err
	}
	err = helmpkg.DeleteRelease(release.Name, helmhost)
	if err != nil {
		return err
	}
	_, err = dao.DeleteHelmRelease(model.ReleaseModel{
		ID: release.ID,
	})
	if err != nil {
		return err
	}
	return processBoardMatchedServiceAndWorkloads(release.Workloads, release.ProjectName, nil, func(matched []*matchedServiceAndWorkload) error {
		// delete the service from board
		services := []model.ServiceStatus{}
		rivers.FromSlice(matched).Map(func(t stream.T) stream.T {
			m := t.(*matchedServiceAndWorkload)
			return model.ServiceStatus{
				Name:        m.Service.Name,
				ProjectName: m.Service.Namespace,
			}
		}).CollectAs(&services)
		_, err = dao.DeleteServiceByNames(services)
		return err
	})
}

func DeleteReleaseByProjectName(projectName string) error {
	models, err := dao.GetAllHelmReleasesByProjectName(projectName)
	if err != nil {
		return err
	}
	for _, m := range models {
		err = DeleteReleaseByReleaseModel(&m)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteReleaseByUserIDAndProjectName(userid int64, projectName string) error {
	models, err := dao.GetHelmReleasesByUserIDAndProjectName(userid, projectName)
	if err != nil {
		return err
	}
	for _, m := range models {
		err = DeleteReleaseByReleaseModel(&m)
		if err != nil {
			return err
		}
	}
	return nil
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

	helmhost, err := getHelmHost()
	if err != nil {
		return nil, err
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
	podsChan := make(chan []model.PodMO)
	go func() {
		load, err := helmpkg.GetReleaseManifest(release.Name, helmhost)
		if err != nil {
			logs.Warning("Get release %s workloads from helm error:%+v", release.Name, err)
		}
		loadChan <- load
		// analysis the manifest and get the pods.
		var pods []model.PodMO
		model.NewK8sHelper().Visit(load, func(infos []*model.K8sInfo, err error) error {
			// add the kubernetes resources to board
			k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
				KubeConfigPath: kubeConfigPath(),
			})
			podlist, err := k8sclient.AppV1().Extend().ListSelectRelatePods(infos)
			if err != nil {
				logs.Warn("list release %s relate pods error:%+v", release.Name, err)
				return err
			}
			if podlist != nil {
				for i := range podlist.Items {
					var containers []model.ContainerMO
					for j := range podlist.Items[i].Spec.Containers {
						containers = append(containers, model.ContainerMO{
							Name:  podlist.Items[i].Spec.Containers[j].Name,
							Image: podlist.Items[i].Spec.Containers[j].Image,
						})
					}
					pods = append(pods, model.PodMO{
						Name:        podlist.Items[i].Name,
						ProjectName: podlist.Items[i].Namespace,
						Spec: model.PodSpecMO{
							Containers: containers,
						},
					})
				}
				// sort the pods by project and name.
				sort.SliceStable(pods, func(i, j int) bool {
					return strings.Compare(pods[i].ProjectName+"/"+pods[i].Name, pods[j].ProjectName+"/"+pods[j].Name) <= 0
				})
			}
			return nil
		})
		podsChan <- pods
	}()

	//get the result
	helmrelease := <-releaseChan

	notes := <-notesChan
	workloads := <-loadChan
	status := <-statusChan
	pods := <-podsChan
	detail := model.ReleaseDetail{
		Release:        generateModelRelease(release, helmrelease),
		Workloads:      workloads,
		Notes:          notes,
		WorkloadStatus: status,
		Pods:           pods,
	}
	return &detail, err
}

func generateModelReleases(models []model.ReleaseModel) ([]model.Release, error) {
	// get the releases from helm cmd
	helmhost, err := getHelmHost()
	if err != nil {
		return nil, err
	}
	list, err := helmpkg.ListAllReleases(helmhost)
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
	return ret, nil
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
	release, err := dao.GetHelmReleaseByName(name)
	if err != nil {
		return false, err
	}

	if release != nil {
		return true, nil
	}

	// get the releases from helm cmd
	helmhost, err := getHelmHost()
	if err != nil {
		return false, err
	}
	list, err := helmpkg.ListAllReleases(helmhost)
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
	repo, err := GetHelmRepository(defaultRepoID)
	if err != nil {
		return err
	}
	// get releases from kubernetes
	// get the releases from helm cmd
	helmhost, err := getHelmHost()
	if err != nil {
		return err
	}
	list, err := helmpkg.ListDeployedReleasesByNamespace(helmhost, projectname)
	if err != nil {
		return err
	}

	// get release from dao
	models, err := dao.GetAllHelmReleases()
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
			// add the release into database.
			err = initBoardReleaseFromKubernetesHelmRelease(repo, &model.Release{
				Name:           list.Releases[i].Name,
				ProjectID:      project.ID,
				ProjectName:    project.Name,
				RepositoryID:   repo.ID,
				RepositoryName: repo.Name,
				OwnerID:        int64(project.OwnerID),
				OwnerName:      project.OwnerName,
			})
			if err != nil {
				// ignore the initReleaseError
				logs.Warning("sync release from kubernetes release %s error: %+v", list.Releases[i].Name, err)
			}
		}
	}
	return nil
}

func getHelmHost() (string, error) {
	tillerport := utils.GetIntValue("TILLER_PORT")
	// add the kubernetes resources to board
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	nodes, err := k8sclient.AppV1().Node().List()
	if err != nil {
		return "", err
	}
	if nodes != nil {
		for i := range nodes.Items {
			logs.Info("ping the ip %s for tiller address", nodes.Items[i].NodeIP)
			ok, err := utils.PingIPAddr(nodes.Items[i].NodeIP)
			if err == nil && ok {
				return nodes.Items[i].NodeIP + ":" + strconv.Itoa(tillerport), nil
			}
		}
	}
	logs.Error("Can't find the available node used for tiller address")
	return "", fmt.Errorf("Can't find the available node used for tiller address")
}
