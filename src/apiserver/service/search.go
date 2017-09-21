package service

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"strings"

	modelK8s "k8s.io/client-go/pkg/api/v1"
)

type OriginImage struct {
	Repositories []string `json:"repositories"`
}
type SearchServiceResult struct {
	ServiceName string `json:"service_name"`
	ProjectName string `json:"project_name"`
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
		searchResult.ProjectResult = resProject
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

		resSvr, err = searchService(searchPara)
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
	var Node modelK8s.NodeList
	defer func() { recover() }()
	err = getFromRequest(kubeNodeURL(), &Node)
	if err != nil {
		return
	}
	for _, v := range Node.Items {
		if strings.Contains(v.Status.Addresses[1].Address, para) {
			res = append(res, SearchNodeResult{
				NodeName: v.Status.Addresses[1].Address,
				NodeIP:   v.Status.Addresses[1].Address,
			})
		}

	}
	return
}
func searchService(searchPara string) (res []SearchServiceResult, err error) {
	resSvr, err := dao.SearchService(searchPara)
	for _, val := range resSvr {
		var svr SearchServiceResult
		svr.ServiceName = val.Name
		svr.ProjectName = val.ProjectName
		res = append(res, svr)
	}
	return res, err
}
