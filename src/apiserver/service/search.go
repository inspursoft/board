package service

import (
	"strings"

	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
)

type OriginImage struct {
	Repositories []string `json:"repositories"`
}
type SearchServiceResult struct {
	ServiceName string `json:"service_name"`
	ProjectName string `json:"project_name"`
	IsPublic    bool   `json:"is_public"`
}
type SearchNodeResult struct {
	NodeName string `json:"node_name"`
	NodeIP   string `json:"node_ip"`
}
type SearchResult struct {
	ProjectResult []dao.SearchProjectResult `json:"project_result"`
	UserResult    []dao.SearchUserResult    `json:"user_result"`
	ImageResult   []SearchImageResult       `json:"images_name"`
	NodeResult    []SearchNodeResult        `json:"node_result"`
	ServiceResult []SearchServiceResult     `json:"service_result"`
}
type SearchImageResult struct {
	ImageName   string `json:"image_name"`
	ProjectName string `json:"project_name"`
}

var registryURL = utils.GetConfig("REGISTRY_URL")

func SearchSource(user *model.User, searchPara string) (searchResult SearchResult, err error) {
	var (
		resProject []dao.SearchProjectResult
		resUser    []dao.SearchUserResult
		resImages  []SearchImageResult
		resNode    []SearchNodeResult
		resSvr     []SearchServiceResult
	)
	if user == nil {
		resProject, err = dao.SearchPublicProject(searchPara)
		resSvr, err = searchPublicService(searchPara)
		searchResult.ProjectResult = resProject
		searchResult.ServiceResult = resSvr
	} else {

		resProject, err = dao.SearchPrivateProject(searchPara, user.Username)
		if err != nil {
			return searchResult, err
		}
		resUser, err = dao.SearchUser(user.Username, searchPara)
		if err != nil {
			return searchResult, err
		}
		currentProject, err := getProjectByUser(user.ID)
		if err != nil {
			return searchResult, err
		}
		resImages, err = searchImages(registryURL()+"/v2/_catalog", currentProject, searchPara)
		if err != nil {
			return searchResult, err
		}
		if user.SystemAdmin == 1 {
			resNode, err = searchNode(searchPara)
		}

		resSvr, err = searchService(searchPara, user.ID)
		if err != nil {
			return searchResult, err
		}
		searchResult = SearchResult{
			ProjectResult: resProject,
			UserResult:    resUser,
			ImageResult:   resImages,
			NodeResult:    resNode,
			ServiceResult: resSvr,
		}
	}
	return searchResult, nil
}
func searchImages(url string, projectNames []string, para string) (res []SearchImageResult, err error) {
	var oriImg OriginImage
	err = getFromRequest(url, &oriImg)
	if err != nil {
		return
	}
	for _, v := range oriImg.Repositories {
		temp := strings.Split(v, "/")
		if len(temp) == 0 {
			continue
		}
		for _, projectNameVal := range projectNames {
			if strings.EqualFold(temp[0], projectNameVal) {
				nameStr := strings.Join(temp[1:], "/")
				projectName := temp[0]
				if strings.Contains(nameStr, para) {
					res = append(res, SearchImageResult{
						ImageName:   nameStr,
						ProjectName: projectName})
				}
			}
		}

	}
	return
}

func getProjectByUser(userID int64) (projectName []string, err error) {
	var query model.Project
	projects, err := GetProjectsByUser(query, userID)
	if err != nil {
		return
	}
	for _, v := range projects {
		projectName = append(projectName, v.Name)
	}
	return
}

func searchNode(para string) (res []SearchNodeResult, err error) {
	defer func() { recover() }()
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	nodecli := k8sclient.AppV1().Node()

	nodes, err := nodecli.List()
	if err != nil {
		return
	}
	for _, v := range nodes.Items {
		if strings.Contains(v.Status.Addresses[1].Address, para) {
			res = append(res, SearchNodeResult{
				NodeName: v.Status.Addresses[1].Address,
				NodeIP:   v.Status.Addresses[1].Address,
			})
		}

	}
	return
}
func searchService(searchPara string, userID int64) (res []SearchServiceResult, err error) {
	serviceList, err := GetServiceList(searchPara, 0, userID, nil, nil)
	for _, val := range serviceList {
		var svr SearchServiceResult
		svr.ServiceName = val.Name
		svr.ProjectName = val.ProjectName
		svr.IsPublic = (val.Public == 1)
		res = append(res, svr)
	}
	return res, err
}

func searchPublicService(searchPara string) (res []SearchServiceResult, err error) {
	resSvr, err := dao.SearchPublicSvr(searchPara)
	for _, val := range resSvr {
		var svr SearchServiceResult
		svr.ServiceName = val.Name
		svr.ProjectName = val.ProjectName
		svr.IsPublic = (val.Public == 1)
		res = append(res, svr)
	}
	return res, err
}
