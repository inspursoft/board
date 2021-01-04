package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"

	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
)

type HelmController struct {
	c.BaseController
}

func (hc *HelmController) ListHelmReposAction() {
	logs.Info("list all helm repos")

	// list the repos from storage
	repo, err := service.ListHelmRepositories()
	if err != nil {
		hc.InternalError(err)
		return
	}
	hc.RenderJSON(repo)
}

func (hc *HelmController) GetHelmRepoDetailAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.InternalError(err)
		return
	}
	pageIndex, _ := hc.GetInt("page_index", 0)
	if pageIndex < 0 {
		hc.InternalError(fmt.Errorf("The page index %d is not correct", pageIndex))
		return
	}
	pageSize, _ := hc.GetInt("page_size", 0)
	if pageSize < 0 {
		hc.InternalError(fmt.Errorf("The page size %d is not correct", pageSize))
		return
	}
	nameRegex := hc.GetString("name_regex", "")
	logs.Info("get helm repository %d, filter by name regex %s and from index %d with size %d", id, nameRegex, pageIndex, pageSize)

	// do some check
	repo := hc.resolveHelmRepositoryByID(int64(id))

	var detail interface{}
	if pageSize == 0 {
		detail, err = service.GetRepoDetail(repo, nameRegex)
	} else {
		detail, err = service.GetPaginatedRepoDetail(repo, nameRegex, pageIndex, pageSize)
	}
	if err != nil {
		hc.InternalError(err)
		return
	}
	hc.RenderJSON(detail)
}

func (hc *HelmController) GetHelmChartDetailAction() {
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.InternalError(err)
		return
	}
	chartname := hc.Ctx.Input.Param(":chartname")
	if err != nil {
		hc.InternalError(err)
		return
	}
	chartversion := hc.Ctx.Input.Param(":chartversion")
	if err != nil {
		hc.InternalError(err)
		return
	}
	logs.Info("get helm chart %s with version %s from repository %d", chartname, chartversion, id)

	detail, err := service.GetChartDetail(hc.resolveHelmRepositoryByID(int64(id)), chartname, chartversion)
	if err != nil {
		hc.InternalError(err)
		return
	}
	hc.RenderJSON(detail)
}

func (hc *HelmController) UploadHelmChartAction() {
	logs.Info("upload helm chart")
	if !hc.IsSysAdmin {
		hc.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to upload helm chart.")
		return
	}

	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.InternalError(err)
		return
	}
	logs.Info("get helm repository %d", id)

	// do some check
	repo := hc.resolveHelmRepositoryByID(int64(id))

	_, fileHeader, err := hc.GetFile("upload_file")
	if err != nil {
		hc.InternalError(err)
		return
	}
	if !strings.HasSuffix(fileHeader.Filename, "tgz") && !strings.HasSuffix(fileHeader.Filename, "tar.gz") {
		hc.InternalError(fmt.Errorf("the upload file must be a gzip tar file"))
		return
	}
	// save to the museum chart directory
	tempDir, err := ioutil.TempDir("", "upload")
	if err != nil {
		hc.InternalError(err)
		return
	}
	defer os.RemoveAll(tempDir)
	tempFile := filepath.Join(tempDir, fileHeader.Filename)
	err = hc.SaveToFile("upload_file", tempFile)
	if err != nil {
		hc.InternalError(err)
	}
	err = service.UploadChart(repo, tempFile)
	if err != nil {
		hc.InternalError(err)
	}
}

func (hc *HelmController) DeleteHelmChartAction() {
	if !hc.IsSysAdmin {
		hc.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to delete helm chart.")
		return
	}
	// get the repo id
	id, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.InternalError(err)
		return
	}
	chartname := hc.Ctx.Input.Param(":chartname")
	if err != nil {
		hc.InternalError(err)
		return
	}
	chartversion := hc.Ctx.Input.Param(":chartversion")
	if err != nil {
		hc.InternalError(err)
		return
	}
	logs.Info("delete helm chart %s with version %s from repository %d", chartname, chartversion, id)

	err = service.DeleteChart(hc.resolveHelmRepositoryByID(int64(id)), chartname, chartversion)
	if err != nil {
		hc.InternalError(err)
		return
	}
}

func (hc *HelmController) InstallHelmChartAction() {
	// resolve the release info
	release := new(model.Release)
	err := hc.ResolveBody(release)
	if err != nil {
		hc.InternalError(err)
		return
	}
	logs.Info("install helm chart %s with version %s from repository %d", release.Chart, release.ChartVersion, release.RepositoryID)
	logs.Info("install release %s in project %d with values %s", release.Name, release.ProjectID, release.Values)
	//Judge authority
	project := hc.ResolveUserPrivilegeByID(release.ProjectID)

	//Check name exist or not
	isExists, err := service.CheckReleaseNames(release.Name)
	if err != nil {
		hc.InternalError(err)
		return
	}
	if isExists {
		hc.CustomAbortAudit(http.StatusConflict, "release "+release.Name+" already exists.")
	}
	// update the release attributes
	release.ProjectName = project.Name
	release.OwnerID = hc.CurrentUser.ID
	release.OwnerName = hc.CurrentUser.Username

	err = service.InstallChart(hc.resolveHelmRepositoryByID(release.RepositoryID), release)
	if err != nil {
		hc.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprintf("%+v", err))
		return
	}
}

func (hc *HelmController) ListHelmReleaseAction() {
	var err error
	projectName := hc.GetString("project_name")
	var releases interface{}
	if hc.IsSysAdmin {
		releases, err = service.ListAllReleases(projectName)
	} else {
		releases, err = service.ListReleasesByUserID(hc.CurrentUser.ID, projectName)
	}
	if err != nil {
		hc.InternalError(err)
		return
	}

	hc.RenderJSON(releases)
}

func (hc *HelmController) DeleteHelmReleaseAction() {
	// get the release id
	r := hc.resolveRelease()
	if !hc.IsSysAdmin && hc.CurrentUser.ID != r.OwnerID {
		hc.CustomAbortAudit(http.StatusForbidden, fmt.Sprintf("Insufficient privileges to operate release %s.", r.Name))
		return
	}
	err := service.DeleteRelease(r.ID)
	if err != nil {
		hc.InternalError(err)
		return
	}

}

func (hc *HelmController) DeleteHelmReleaseByProjectAction() {
	// get the project name
	projectName := hc.GetString("project_name")
	var err error
	if hc.IsSysAdmin {
		err = service.DeleteReleaseByProjectName(projectName)
	} else {
		err = service.DeleteReleaseByUserIDAndProjectName(hc.CurrentUser.ID, projectName)
	}

	if err != nil {
		hc.InternalError(err)
		return
	}

}

func (hc *HelmController) GetHelmReleaseAction() {
	// get the release id
	r := hc.resolveRelease()
	if !hc.IsSysAdmin && hc.CurrentUser.ID != r.OwnerID {
		hc.CustomAbortAudit(http.StatusForbidden, fmt.Sprintf("Insufficient privileges to operate release %s.", r.Name))
		return
	}
	release, err := service.GetReleaseDetail(r.ID)
	if err != nil {
		hc.InternalError(err)
		return
	}
	hc.RenderJSON(release)
}

func (hc *HelmController) GetReleasePodLogsAction() {
	var err error
	r := hc.resolveRelease()
	//Judge authority
	if !hc.IsSysAdmin && hc.CurrentUser.ID != r.OwnerID {
		hc.CustomAbortAudit(http.StatusForbidden, fmt.Sprintf("Insufficient privileges to operate release %s.", r.Name))
		return
	}
	podName := hc.Ctx.Input.Param(":podname")
	readCloser, err := service.GetK8sPodLogs(r.ProjectName, podName, hc.GeneratePodLogOptions())
	if err != nil {
		hc.ParseError(err, c.ParseGetK8sError)
		return
	}
	defer readCloser.Close()
	_, err = io.Copy(&utils.FlushResponseWriter{hc.Ctx.Output.Context.ResponseWriter}, readCloser)
	if err != nil {
		logs.Error("get job logs error:%+v", err)
	}
}

func (hc *HelmController) resolveRelease() *model.ReleaseModel {
	// get the release id
	releaseid, err := strconv.Atoi(hc.Ctx.Input.Param(":id"))
	if err != nil {
		hc.InternalError(err)
		return nil
	}
	model, err := service.GetReleaseFromDB(int64(releaseid))
	if err != nil {
		hc.InternalError(err)
		return nil
	}
	if model == nil {
		hc.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("Can't find the release with id %d.", releaseid))
		return nil
	}
	return model
}

func (hc *HelmController) resolveHelmRepositoryByID(repositoryID int64) *model.HelmRepository {
	// do some check
	repo, err := service.GetHelmRepository(repositoryID)
	if err != nil {
		hc.InternalError(err)
		return nil
	} else if repo == nil {
		hc.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("Helm repository %d does not exists.", repositoryID))
		return nil
	}
	return repo
}

func (hc *HelmController) ReleaseExists() {
	target := hc.GetString("release_name")
	isExists, err := service.CheckReleaseNames(target)
	if err != nil {
		hc.InternalError(err)
		return
	}
	if isExists {
		hc.CustomAbortAudit(http.StatusConflict, "release "+target+" already exists.")
	}
}
