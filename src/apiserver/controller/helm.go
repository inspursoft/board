package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
)

type HelmController struct {
	BaseController
}

func (hc *HelmController) ListHelmReposAction() {
	logs.Info("list all helm repos")

	// list the repos from storage
	repo, err := service.ListRepositories()
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(repo)
}

func (hc *HelmController) CreateHelmRepoAction() {
	if hc.isSysAdmin == false {
		hc.customAbort(http.StatusForbidden, "Insufficient privileges to delete image.")
		return
	}
	// resolve the repository
	repo := new(model.Repository)
	err := hc.resolveBody(repo)
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("Added repository %s: %+v", repo.Name, repo)

	// do some check
	exist, err := service.CheckRepoNameNotExist(repo.Name)
	if err != nil {
		hc.internalError(err)
		return
	} else if exist {
		hc.customAbort(http.StatusConflict, fmt.Sprintf("Helm Repository %s already exists in cluster.", repo.Name))
		return
	}

	// add the repo to k8s
	repoid, err := service.AddRepository(*repo)
	if err != nil {
		hc.internalError(err)
		return
	}
	repo.ID = repoid
	hc.renderJSON(repo)
}

func (hc *HelmController) UpdateHelmRepoAction() {
	if hc.isSysAdmin == false {
		hc.customAbort(http.StatusForbidden, "Insufficient privileges to delete image.")
		return
	}
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}

	// resolve the repo
	repo := new(model.Repository)
	err = hc.resolveBody(repo)
	if err != nil {
		hc.internalError(err)
		return
	}

	logs.Info("update helm repository %d to %+v", id, repo)
	// override the fields
	repo.ID = int64(id)

	// do some check
	oldrepo, err := service.GetRepository(repo.ID)
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

	err = service.UpdateRepository(*repo)
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(repo)
}

func (hc *HelmController) DeleteHelmRepoAction() {
	if hc.isSysAdmin == false {
		hc.customAbort(http.StatusForbidden, "Insufficient privileges to delete image.")
		return
	}
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("delete helm repository %d", id)

	// do some check
	oldrepo, err := service.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if oldrepo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	// delete the repo
	err = service.DeleteRepository(int64(id))
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
	repo, err := service.GetRepository(int64(id))
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
	repo, err := service.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	detail, err := service.GetRepoDetail(repo)
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
	repo, err := service.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	detail, err := service.GetChartDetail(repo, chartname, chartversion)
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(detail)
}

func (hc *HelmController) DeleteHelmChartAction() {
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
	logs.Info("delete helm chart %s with version %s from repository %d", chartname, chartversion, id)

	// do some check
	repo, err := service.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	err = service.DeleteChart(repo, chartname, chartversion)
	if err != nil {
		hc.internalError(err)
		return
	}
}

func (hc *HelmController) InstallHelmChartAction() {
	// resolve the release info
	release := new(model.Release)
	err := hc.resolveBody(release)
	if err != nil {
		hc.internalError(err)
		return
	}

	logs.Info("install helm chart %s with version %s from repository %d", release.Chart, release.ChartVersion, release.RepositoryId)
	logs.Info("install release %s in project %d with value %s", release.Name, release.ProjectId, release.Value)

	// do some check
	repo, err := service.GetRepository(release.RepositoryId)
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", release.RepositoryId))
		return
	}

	//Judge authority
	project := hc.resolveUserPrivilegeByID(release.ProjectId)
	err = service.InstallChart(repo, release.Chart, release.ChartVersion, release.Name, release.ProjectId, project.Name, release.Value, hc.currentUser.ID, hc.currentUser.Username)
	if err != nil {
		hc.internalError(err)
		return
	}
	//	hc.renderJSON(detail)
}

func (hc *HelmController) ListHelmReleaseAction() {
	// get the repo id
	repoStr := hc.GetString("repository_id")
	var repo *model.Repository
	var err error
	if repoStr != "" {
		repoid, err := strconv.Atoi(repoStr)
		if err != nil {
			hc.internalError(err)
			return
		}
		// do some check
		repo, err = service.GetRepository(int64(repoid))
		if err != nil {
			hc.internalError(err)
			return
		} else if repo == nil {
			hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", repoid))
			return
		}
	}

	releases, err := service.ListReleases(repo)
	if err != nil {
		hc.internalError(err)
		return
	}

	hc.renderJSON(releases)
}

func (hc *HelmController) DeleteHelmReleaseAction() {
	// get the release id
	releaseid, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}

	err = service.DeleteRelease(int64(releaseid))
	if err != nil {
		hc.internalError(err)
		return
	}

}

func (hc *HelmController) UploadHelmChartAction() {
	logs.Info("upload helm chart")
	if hc.isSysAdmin == false {
		hc.customAbort(http.StatusForbidden, "Insufficient privileges to delete image.")
		return
	}

	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	logs.Info("get helm repository %d", id)

	// do some check
	repo, err := service.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	_, fileHeader, err := hc.GetFile("upload_file")
	if err != nil {
		hc.internalError(err)
		return
	}
	if !strings.HasSuffix(fileHeader.Filename, "tgz") && !strings.HasSuffix(fileHeader.Filename, "tar.gz") {
		hc.internalError(fmt.Errorf("the upload file must be a gzip tar file"))
		return
	}
	// save to the museum chart directory
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
	tempFile := filepath.Join(tempDir, fileHeader.Filename)
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	err = hc.SaveToFile("upload_file", tempFile)
	if err != nil {
		hc.internalError(err)
	}
	err = service.UploadChart(repo, tempFile)
	if err != nil {
		hc.internalError(err)
	}
}
