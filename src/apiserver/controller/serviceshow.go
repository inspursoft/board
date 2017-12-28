package controller

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

type ServiceShowController struct {
	baseController
}

func (s *ServiceShowController) Get() {

	ownerName := s.Ctx.Input.Param(":owner_name")
	projectName := s.Ctx.Input.Param(":project_name")
	serviceName := s.Ctx.Input.Param(":service_name")
	serviceIdentity := strings.ToLower(ownerName + "_" + projectName + "_" + serviceName)

	if serviceURL, ok := memoryCache.Get(serviceIdentity).(string); ok {
		logs.Debug("Service URL: %s", serviceURL)
		http.Redirect(s.Ctx.ResponseWriter, s.Ctx.Request, serviceURL, http.StatusFound)
	}
}
