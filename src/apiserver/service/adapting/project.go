package adapting

import (
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
)

func GetProjectsByMember(v vm.Project, userID int64) (results []*vm.Project, err error) {
	projects, err := service.GetProjectsByMember(v.ToMO(), userID)
	if err != nil {
		return
	}
	utils.Adapt(projects, &results)
	return
}

func GetProjectsByUser(v vm.Project, userID int64) (results []*vm.Project, err error) {
	projects, err := service.GetProjectsByUser(v.ToMO(), userID)
	if err != nil {
		return
	}
	for _, v := range projects {
		logs.Debug("%+v", v)
	}
	utils.Adapt(projects, &results)
	return
}

func GetPaginatedProjectsByUser(v vm.Project, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (pp *vm.PaginatedProjects, err error) {
	paginagedProjects, err := service.GetPaginatedProjectsByUser(v.ToMO(), userID, pageIndex, pageSize, orderField, orderAsc)
	if err != nil {
		return
	}
	utils.Adapt(paginagedProjects, pp)
	return
}

func CreateProject(project vm.Project) (bool, error) {
	return service.CreateProject(project.ToMO())
}

func CreateNamespace(project *vm.Project) (status bool, err error) {
	var p *model.Project
	err = utils.Adapt(*project, &p)
	if err != nil {
		return
	}
	status, err = service.CreateNamespace(p)
	return
}

func GetProject(project vm.Project, selectedFields ...string) (target *vm.Project, err error) {
	p, err := service.GetProject(project.ToMO(), selectedFields...)
	err = utils.Adapt(p, &target)
	if err != nil {
		return
	}
	return
}
