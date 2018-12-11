package service

import (
	"fmt"
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
		Name:     repo.Name,
		URL:      repo.URL,
		Username: repo.Username,
		Password: repo.Password,
		Type:     repo.Type,
	}
	if repo.Cert != "" {
		entry.Cert = []byte(repo.Cert)
	}
	if repo.Key != "" {
		entry.Key = []byte(repo.Key)
	}
	if repo.CA != "" {
		entry.CA = []byte(repo.CA)
	}
	return entry
}

func GetRepoDetail(repo *model.Repository) (*model.RepositoryDetail, error) {
	chartrepo, err := helmpkg.NewChartRepository(toEntry(repo))
	if err != nil {
		return nil, err
	}
	var detail model.RepositoryDetail
	// TODO: replace ICON URL to ICON base64 code.
	detail.IndexFile = (*model.IndexFile)(chartrepo.IndexFile)
	detail.Repository = repo
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

type boardServiceInsertAction struct {
	ownerid     int64
	ownername   string
	projectname string
	release     *model.ReleaseModel
}

func (b *boardServiceInsertAction) PreInstall(templateInfo string) error {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: utils.GetConfig("KUBE_MASTER_URL")(),
	})
	//resolve the templateInfo into kubernetes service and deployments.....
	mapper := k8sclient.AppV1().Mapper()
	return mapper.Visit(templateInfo, func(infos []*model.Info) error {
		// board only support Service and Deployment yamls
		var serviceInfos, deploymentInfos []*model.Info
		for i := range infos {
			if infos[i].Namespace == "" {
				infos[i].Namespace = b.projectname
			}
			if infos[i].GroupVersionKind.Kind == "Service" && infos[i].Namespace == b.projectname {
				serviceInfos = append(serviceInfos, infos[i])
			} else if infos[i].GroupVersionKind.Kind == "Deployment" && infos[i].Namespace == b.projectname {
				deploymentInfos = append(deploymentInfos, infos[i])
			}
		}

		if len(serviceInfos) == 0 || len(deploymentInfos) == 0 {
			return fmt.Errorf("The helm chart %s-%s must have a service and deployment within namespace %s", b.release.Chart, b.release.ChartVersion, b.projectname)
		}
		updateBoard := false
		for i := range serviceInfos {
			var findDeploy *model.Info
			// find the deployment with service name
			for j := range deploymentInfos {
				if serviceInfos[i].Name == deploymentInfos[j].Name {
					findDeploy = deploymentInfos[j]
					break
				}
			}
			if findDeploy != nil {
				//add the service in board
				//TODO: add the service and deployment in board.....
				updateBoard = true
			}
		}
		if !updateBoard {
			return fmt.Errorf("The service and deployment in helm chart %s-%s must have same name", b.release.Chart, b.release.ChartVersion)
		}
		return nil
	})
}

func (b *boardServiceInsertAction) PostInstall(templateInfo string) error {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: utils.GetConfig("KUBE_MASTER_URL")(),
	})
	//resolve the templateInfo into kubernetes service and deployments.....
	mapper := k8sclient.AppV1().Mapper()
	err := mapper.Visit(templateInfo, func(infos []*model.Info) error {
		// board only support Service and Deployment yamls
		for i := range infos {
			if infos[i].Namespace == "" {
				infos[i].Namespace = b.projectname
			}
			if infos[i].GroupVersionKind.Kind == "Service" && infos[i].Namespace == b.projectname {
				// install the service into service_status table.
				status, err := CreateServiceConfig(model.ServiceStatus{
					Name:          infos[i].Name,
					ProjectID:     b.release.RepositoryId,
					ProjectName:   b.projectname,
					Comment:       "helm chart created service",
					OwnerID:       b.ownerid,
					OwnerName:     b.ownername,
					Status:        defaultStatus,
					Public:        projectPrivate,
					Deleted:       defaultDeleted,
					CreationTime:  time.Now(),
					UpdateTime:    time.Now(),
					Source:        helm,
					ServiceConfig: "",
				})
				if err != nil {
					return err
				}
				b.release.Services = append(b.release.Services, model.ReleaseServiceModel{
					ServiceId: status.ID,
				})
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	b.release.Workloads = templateInfo
	_, err = dao.AddHelmRelease(*b.release)

	return err
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
	_, err = chartrepo.InstallChart(chartname, chartversion, name, projectname, values, helmhost, &boardServiceInsertAction{
		ownerid:     ownerid,
		ownername:   ownername,
		projectname: projectname,
		release: &model.ReleaseModel{
			Name:         name,
			ProjectId:    projectid,
			RepositoryId: repo.ID,
			Chart:        chartname,
			ChartVersion: chartversion,
			Value:        values,
		},
	})
	return err
}

func ListReleases(repo *model.Repository) ([]model.Release, error) {
	if repo == nil {
		return dao.GetHelmReleases()
	}
	return dao.GetHelmReleasesByRepositoryId(repo.ID)
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
	for i := range release.Services {
		//delete the service
		_, err := DeleteService(release.Services[i].ServiceId)
		if err != nil {
			return err
		}
	}
	_, err = dao.DeleteHelmRelease(model.ReleaseModel{
		ID: releaseid,
	})
	return err
}
