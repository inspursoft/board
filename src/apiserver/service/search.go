package service

import "git/inspursoft/board/src/common/dao"

func SearchSource(usrName string, pjName string) ([]dao.SearchResult, error) {
	switch usrName {
	case "":
		return dao.SearchPublic(pjName)
	default:
		return dao.SearchPrivite(pjName, usrName)
	}
}
