package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	service "git/inspursoft/board/src/apiserver/service/dashboard"

	"github.com/astaxie/beego"
)

type JsonIn struct {
	TimeUnit      string `json:"time_unit"`
	TimeCount     string `json:"time_count"`
	TimestampBase string `json:"timestamp_base"`
}
type baseController struct {
	beego.Controller
}

func (b *baseController) resolveBody() (in JsonIn, err error) {
	data, err := ioutil.ReadAll(b.Ctx.Request.Body)
	json.Unmarshal(data, &in)
	if err != nil {
		return in, err
	}
	return in, nil
}

type DashboardServiceController struct {
	baseController
}

func (s *DashboardServiceController) GetService() {
	fmt.Println("asd")
	body, _ := s.resolveBody()
	serviceName := s.GetString("service_name")
	fmt.Println("servicename", serviceName)
	if body.TimeCount == "" {
		s.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	if body.TimestampBase == "" {
		s.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	if body.TimeUnit == "" {
		s.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	Json, err := service.GetService(body.TimeUnit, body.TimeCount, body.TimestampBase, serviceName)
	if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(500)
		fmt.Println(err)
		return
	}

	s.Ctx.ResponseWriter.Write(Json)
	fmt.Println(body, serviceName)
}
func (s *DashboardServiceController) GetList() {
	j := service.GetDashboardServiceList()
	s.Ctx.Output.Body(j)
}
