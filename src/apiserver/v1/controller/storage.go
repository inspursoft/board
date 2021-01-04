package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"

	//"github.com/inspursoft/board/src/apiserver/service"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type StorageController struct {
	c.BaseController
}

func (s *StorageController) Storage() {

	name := s.GetString("name")
	server := s.GetString("server")
	path := s.GetString("path")
	storageCap, err := s.GetInt64("cap")
	logs.Debug(name, server, path, storageCap)
	if err != nil {
		s.CustomAbortAudit(http.StatusNotImplemented, fmt.Sprint(err))
		return
	}
	//TODO NFS PV in the next version
	//err = service.SetNFSVol(name, server, path, storageCap)
	if err != nil {
		s.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
	}

}
