package controller

import (
	"git/inspursoft/board/src/apiserver/service"

	"github.com/astaxie/beego/logs"

	"net/http"
	"os"
	"path/filepath"
)

type uploadFile struct {
	ProjectName string `json:"project_name"`
	ServiceID   int64  `json:"service_id"`
}

const maxFileuploadSize = 1 << 22

type FileUploadController struct {
	baseController
	toFilePath string
}

func (f *FileUploadController) Prepare() {
	user := f.getCurrentUser()
	if user == nil {
		f.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	f.currentUser = user
	f.isProjectAdmin = (user.ProjectAdmin == 1)
	if !f.isProjectAdmin {
		f.CustomAbort(http.StatusForbidden, "Insufficient privileges.")
		return
	}
}

func (f *FileUploadController) resolveFilePath() {
	projectName := f.GetString("project_name")
	serviceID, err := f.GetInt64("service_id", 0)
	if err != nil {
		f.internalError(err)
		return
	}
	reqUploadFile := uploadFile{
		ProjectName: projectName,
		ServiceID:   serviceID,
	}

	if reqUploadFile.ProjectName == "" && reqUploadFile.ServiceID == 0 {
		f.CustomAbort(http.StatusBadRequest, "No project name or service ID provided.")
		return
	}
	if reqUploadFile.ProjectName != "" {
		isMember, err := service.IsProjectMemberByName(reqUploadFile.ProjectName)
		if err != nil {
			f.internalError(err)
			return
		}
		if !isMember {
			f.CustomAbort(http.StatusForbidden, "Not member to the current project with provided ID.")
			return
		}
		f.toFilePath = reqUploadFile.ProjectName
	}
}

func (f *FileUploadController) Upload() {
	f.resolveFilePath()
	_, fh, err := f.GetFile("uploadFile")
	if err != nil {
		f.internalError(err)
		return
	}
	targetFilePath := filepath.Join(repoPath, f.toFilePath, "upload")
	os.MkdirAll(targetFilePath, 0755)

	logs.Info("User: %s uploaded file from %s to %s.", f.currentUser.Username, fh.Filename, targetFilePath)
	err = f.SaveToFile("uploadFile", filepath.Join(targetFilePath, fh.Filename))
	if err != nil {
		f.internalError(err)
	}
}

func (f *FileUploadController) ListFiles() {
	f.resolveFilePath()
	uploads, err := service.ListUploadFiles(filepath.Join(repoPath, f.toFilePath, "upload"))
	if err != nil {
		f.internalError(err)
		return
	}
	f.Data["json"] = uploads
	f.ServeJSON()
}
