package controller

import (
	"net/http"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
)

type StorageController struct {
	baseController
}

func (p *StorageController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
}

func (s *StorageController)Storage()  {
	name := s.GetString("name")
	server := s.GetString("server")
	path := s.GetString("path")
	storageCap,err := s.GetInt64("cap")
	if err !=nil {
		s.customAbort(http.StatusNotImplemented, fmt.Sprint(err))
		return
	}
	err = service.SetNFSVol(name,server, path, storageCap)
	if err != nil {
		s.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	return
}