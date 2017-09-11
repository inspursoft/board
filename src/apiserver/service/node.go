package service

import "git/inspursoft/board/src/common/dao"

func GetNode(nodeName string) (node dao.Node, err error) {
	return dao.GetNode(nodeName)
}
