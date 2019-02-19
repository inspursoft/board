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

func (hc *HelmController) GetHelmRepoDetailAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return
	}
	pageIndex, _ := hc.GetInt("page_index", 0)
	if pageIndex < 0 {
		hc.internalError(fmt.Errorf("The page index %d is not correct", pageIndex))
		return
	}
	pageSize, _ := hc.GetInt("page_size", 0)
	if pageSize < 0 {
		hc.internalError(fmt.Errorf("The page size %d is not correct", pageSize))
		return
	}
	nameRegex := hc.GetString("name_regex", "")
	logs.Info("get helm repository %d, filter by name regex %s and from index %d with size %d", id, nameRegex, pageIndex, pageSize)

	// do some check
	repo, err := service.GetRepository(int64(id))
	if err != nil {
		hc.internalError(err)
		return
	} else if repo == nil {
		hc.customAbort(http.StatusBadRequest, fmt.Sprintf("Helm repository %d does not exists.", int64(id)))
		return
	}

	detail, err := service.GetRepoDetail(repo, nameRegex, pageIndex, pageSize)
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

func (hc *HelmController) UploadHelmChartAction() {
	logs.Info("upload helm chart")
	if !hc.isSysAdmin {
		hc.customAbort(http.StatusForbidden, "Insufficient privileges to upload helm chart.")
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

func (hc *HelmController) DeleteHelmChartAction() {
	if !hc.isSysAdmin {
		hc.customAbort(http.StatusForbidden, "Insufficient privileges to delete helm chart.")
		return
	}
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
	logs.Info("install release %s in project %d with values %s", release.Name, release.ProjectId, release.Values)

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
	//Check name exist or not
	isExists, err := service.CheckReleaseNames(release.Name)
	if err != nil {
		hc.internalError(err)
		return
	}
	if isExists {
		hc.customAbort(http.StatusConflict, "release "+release.Name+" already exists.")
	}

	err = service.InstallChart(repo, release.Chart, release.ChartVersion, release.Name, release.ProjectId, project.Name, release.Values, hc.currentUser.ID, hc.currentUser.Username)
	if err != nil {
		hc.internalError(err)
		return
	}
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

	var userid int64 = -1
	if !hc.isSysAdmin {
		userid = hc.currentUser.ID
	}
	releases, err := service.ListReleases(repo, userid)
	if err != nil {
		hc.internalError(err)
		return
	}

	hc.renderJSON(releases)
}

func (hc *HelmController) DeleteHelmReleaseAction() {
	// get the release id
	m := hc.checkReleaseOperationPriviledges()
	if m == nil {
		return
	}
	err := service.DeleteRelease(m.ID)
	if err != nil {
		hc.internalError(err)
		return
	}

}

func (hc *HelmController) GetHelmReleaseAction() {
	// get the release id
	m := hc.checkReleaseOperationPriviledges()
	if m == nil {
		return
	}
	release, err := service.GetReleaseDetail(m.ID)
	if err != nil {
		hc.internalError(err)
		return
	}
	hc.renderJSON(release)
}

func (hc *HelmController) checkReleaseOperationPriviledges() *model.ReleaseModel {
	// get the release id
	releaseid, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.internalError(err)
		return nil
	}
	model, err := service.GetReleaseFromDB(int64(releaseid))
	if err != nil {
		hc.internalError(err)
		return nil
	}
	if model == nil {
		hc.customAbort(http.StatusNotFound, fmt.Sprintf("Can't find the release with id %d.", releaseid))
		return nil
	}
	if !hc.isSysAdmin && hc.currentUser.ID != model.OwnerID {
		hc.customAbort(http.StatusForbidden, fmt.Sprintf("Insufficient privileges to operate release %s.", model.Name))
		return nil
	}
	return model
}

func (hc *HelmController) ReleaseExists() {
	target := hc.GetString("release_name")
	isExists, err := service.CheckReleaseNames(target)
	if err != nil {
		hc.internalError(err)
		return
	}
	if isExists {
		hc.customAbort(http.StatusConflict, "release "+target+" already exists.")
	}
}
