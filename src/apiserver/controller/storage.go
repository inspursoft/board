package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
)

type StorageController struct {
	BaseController
}

func (s *StorageController) Storage() {
	name := s.GetString("name")
	server := s.GetString("server")
	path := s.GetString("path")
	storageCap, err := s.GetInt64("cap")
	if err != nil {
		s.customAbort(http.StatusNotImplemented, fmt.Sprint(err))
		return
	}
	err = service.SetNFSVol(name, server, path, storageCap)
	if err != nil {
		s.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
	}
}
