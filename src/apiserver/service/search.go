package service

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
)

type SearchResult struct {
	ProjectResult []dao.SearchProjectResult `json:"project_result"`
	UserResult    []dao.SearchUserResult    `json:"user_result"`
}

func SearchSource(user *model.User, searchPara string) (searchResult SearchResult, err error) {
	var (
		resProject []dao.SearchProjectResult
		resUser    []dao.SearchUserResult
	)
	if user == nil {
		resProject, err = dao.SearchPublicProject(searchPara)
		searchResult.ProjectResult = resProject
	} else {
		resProject, err = dao.SearchPrivateProject(searchPara, user.Username)
		resUser, err = dao.SearchUser(user.Username, searchPara)
		searchResult.ProjectResult = resProject
		searchResult.UserResult = resUser
	}
	if err != nil {
		return
	}

	return searchResult, nil
}
