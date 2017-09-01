package service

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
)

func SearchSource(user *model.User, projectName string) ([]dao.SearchResult, error) {
	if user == nil {
		return dao.SearchPublic(projectName)
	} else {
		return dao.SearchPrivite(projectName, user.Username)
	}
}
