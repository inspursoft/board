package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"

	"git/inspursoft/board/src/apiserver/service/helm"
	"git/inspursoft/board/src/common/model"
)

type HelmController struct {
	BaseController
}

func (hc *HelmController) ListHelmReposAction() {
	logs.Info("list all helm repos")

	// list the repos from storage
	repo, err := helm.ListRepositories()
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(repo)
}

func (hc *HelmController) CreateHelmRepoAction() {
	// resolve the repository
	repo := new(model.Repository)
	err := hc.resolveBody(repo)
	if err != nil {
		return
	}
	logs.Info("Added repository %s: %+v", repo.Name, repo)

	// do some check
	exist, err := helm.CheckRepoNameNotExist(repo.Name)
	if err != nil {
		hc.internalError(err)
		return
	} else if exist {
		hc.customAbort(http.StatusConflict, fmt.Sprintf("Helm Repository %s already exists in cluster.", repo.Name))
		return
	}

	// add the hpa to k8s
	repoid, err := helm.AddRepository(*repo)
	if err != nil {
		hc.internalError(err)
		return
	}
	repo.ID = repoid
	hc.renderJSON(repo)
}

func (hc *HelmController) UpdateHelmRepoAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}

	// resolve the hpa
	repo := new(model.Repository)
	err = hc.resolveBody(repo)
	if err != nil {
		return
	}

	logs.Info("update helm repository %d to %+v", id, repo)
	// override the fields
	repo.ID = int64(id)

	// do some check
	oldrepo, err := helm.GetRepository(repo.ID)
	if err != nil {
		hc.internalError(err)
		return
	} else if oldrepo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm Repository %d does not exists.", repo.ID))
		return
	} else if oldrepo.Name != repo.Name {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("can't change Helm Repository %s's name to %s", oldrepo.Name, repo.Name))
		return
	}

	err = helm.UpdateRepository(*repo)
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(repo)
}

func (hc *HelmController) DeleteHelmRepoAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("delete helm repository %d", id)

	// do some check
	oldrepo, err := helm.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if oldrepo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	// delete the repo
	err = helm.DeleteRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	}
}

func (hc *HelmController) GetHelmRepoAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("get helm repository %d", id)

	// do some check
	repo, err := helm.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	hc.renderJSON(repo)
}

func (hc *HelmController) GetHelmRepoDetailAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("get helm repository %d", id)

	// do some check
	repo, err := helm.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	detail, err := helm.GetRepoDetail(repo)
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(detail)
}

func (hc *HelmController) GetHelmChartDetailAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	chartname := hc.Ctx.Input.Param(":chartname")
	if err != nil {
		hc.internalError(err)
		return
	}
	chartversion := hc.Ctx.Input.Param(":chartversion")
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("get helm chart %s with version %s from repository %d", chartname, chartversion, id)

	// do some check
	repo, err := helm.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	detail, err := helm.GetChartDetail(repo, chartname, chartversion)
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(detail)
}

func (hc *HelmController) InstallHelmChartAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	chartname := hc.Ctx.Input.Param(":chartname")
	if err != nil {
		hc.internalError(err)
		return
	}
	chartversion := hc.Ctx.Input.Param(":chartversion")
	if err != nil {
		hc.internalError(err)
		return
	}
	name := hc.Ctx.Input.Param(":name")
	if err != nil {
		hc.internalError(err)
		return
	}
	namespace := hc.Ctx.Input.Param(":namespace")
	if err != nil {
		hc.internalError(err)
		return
	}
	defer hc.Ctx.Request.Body.Close()
	value, err := ioutil.ReadAll(hc.Ctx.Request.Body)
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("install helm chart %s with version %s from repository %d", chartname, chartversion, id)
	logs.Info("install release %s in namespace %s with value %s", name, namespace, value)

	// do some check
	repo, err := helm.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	err = helm.InstallChart(repo, chartname, chartversion, name, namespace, string(value))
	if err != nil {
		hc.internalError(err)
		return
	}
	//	hc.renderJSON(detail)
}
