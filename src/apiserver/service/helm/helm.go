package helm

import (
	"fmt"

	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
)

const (
	helmName = "helm"
)

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
	return dao.GetHelmRepository(repo)
}

func UpdateRepository(repo model.Repository) error {
	_, err := dao.UpdateHelmRepository(repo)
	return err
}

func DeleteRepository(id int64) error {
	_, err := dao.DeleteHelmRepository(model.Repository{
		ID: id,
	})
	return err
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

func toEntry(repo *model.Repository) *Entry {
	entry := &Entry{
		Name:     repo.Name,
		URL:      repo.URL,
		Username: repo.Username,
		Password: repo.Password,
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
	chartrepo, err := NewChartRepository(toEntry(repo))
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
	chartrepo, err := NewChartRepository(toEntry(repo))
	if err != nil {
		return nil, err
	}
	chart, err := chartrepo.FetchChart(chartname, chartversion)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

func InstallChart(repo *model.Repository, chartname, chartversion, name, namespace, values string) error {
	chartrepo, err := NewChartRepository(toEntry(repo))
	if err != nil {
		return err
	}
	helmhost := utils.GetStringValue("HELM_HOST")
	if helmhost == "" {
		return fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
	return chartrepo.InstallChart(chartname, chartversion, name, namespace, values, helmhost, func(templateInfo string) error {
		//resolve the templateInfo into kubernetes service and deployments.....
		mapper := NewMapper()
		return mapper.Visit(templateInfo, func(infos []*Info) error {
			// board only support Service and Deployment yamls
			var serviceInfos, deploymentInfos []*Info
			for i := range infos {
				if infos[i].Namespace == "" {
					infos[i].Namespace = namespace
				}
				if infos[i].GroupVersionKind.Kind == "Service" && infos[i].Namespace == namespace {
					serviceInfos = append(serviceInfos, infos[i])
				} else if infos[i].GroupVersionKind.Kind == "Deployment" && infos[i].Namespace == namespace {
					deploymentInfos = append(deploymentInfos, infos[i])
				}
			}

			if len(serviceInfos) == 0 || len(deploymentInfos) == 0 {
				return fmt.Errorf("The helm chart %s-%s must have a service and deployment within namespace %s", chartname, chartversion, namespace)
			}
			updateBoard := false
			for i := range serviceInfos {
				var findDeploy *Info
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
				return fmt.Errorf("The service and deployment in helm chart %s-%s must have same name", chartname, chartversion)
			}
			return nil
		})
	})
}
