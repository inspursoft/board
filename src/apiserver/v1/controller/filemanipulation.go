package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"

	"net/http"
	"os"
	"path/filepath"
)

var attachmentFile = "attachment.zip"

type uploadFile struct {
	ProjectName string `json:"project_name"`
	ServiceID   int64  `json:"service_id"`
	ImageName   string `json:"image_name"`
	TagName     string `json:"tag_name"`
}

type FileUploadController struct {
	c.BaseController
	ToFilePath string
}

func (f *FileUploadController) Prepare() {
	f.EnableXSRF = false
	f.ResolveSignedInUser()
	f.RecordOperationAudit()
	f.resolveFilePath()
}

func (f *FileUploadController) resolveFilePath() {
	f.ToFilePath = filepath.Join(c.BaseRepoPath(), f.CurrentUser.Username, "upload")
	err := os.MkdirAll(f.ToFilePath, 0755)
	if err != nil {
		logs.Error("Failed to make dir: %s, error: %+v", f.ToFilePath, err)
	}
}

func (f *FileUploadController) Upload() {
	_, fh, err := f.GetFile("upload_file")
	if err != nil {
		f.InternalError(err)
		return
	}
	targetFilePath := f.ToFilePath

	logs.Info("User: %s uploaded file from %s to %s.", f.CurrentUser.Username, fh.Filename, targetFilePath)
	err = f.SaveToFile("upload_file", filepath.Join(targetFilePath, fh.Filename))
	if err != nil {
		f.InternalError(err)
	}
}

func (f *FileUploadController) DownloadProbe() {
	if isEmpty, err := service.IsEmptyDirectory(f.ToFilePath); isEmpty || err != nil {
		f.CustomAbortAudit(http.StatusNotFound, "No uploaded file found.")
		return
	}
}

func (f *FileUploadController) Download() {
	fileName := f.GetString("file_name")
	if fileName == "" {
		logs.Info("Will zip files to be downloaded as no file name specified.")
		attachmentFilePath := filepath.Join(c.BaseRepoPath(), f.CurrentUser.Username)
		err := service.ZipFiles(filepath.Join(attachmentFilePath, attachmentFile), f.ToFilePath)
		if err != nil {
			f.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprintf("Failed to zip file for attachment: %+v", err))
			return
		}
		f.ToFilePath = attachmentFilePath
		fileName = attachmentFile
	}
	logs.Debug("Download file from path: %s", f.ToFilePath)
	f.Ctx.Output.Download(filepath.Join(f.ToFilePath, fileName), fileName)
}

func (f *FileUploadController) ListFiles() {
	uploads, err := service.ListUploadFiles(f.ToFilePath)
	if err != nil {
		f.InternalError(err)
		return
	}
	f.RenderJSON(uploads)
}

func (f *FileUploadController) RemoveFile() {
	fileInfo := model.FileInfo{
		Path:     f.ToFilePath,
		FileName: f.GetString("file_name"),
	}
	logs.Info("Removed file: %s", filepath.Join(fileInfo.Path, fileInfo.FileName))
	err := service.RemoveUploadFile(fileInfo)
	if err != nil {
		f.InternalError(err)
	}
}
