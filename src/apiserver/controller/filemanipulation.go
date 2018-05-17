package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"

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
	baseController
	toFilePath string
}

func (f *FileUploadController) Prepare() {
	f.resolveSignedInUser()
	f.resolveFilePath()
}

func (f *FileUploadController) resolveFilePath() {
	f.toFilePath = filepath.Join(baseRepoPath(), f.currentUser.Username, "upload")
	err := os.MkdirAll(f.toFilePath, 0755)
	if err != nil {
		logs.Error("Failed to make dir: %s, error: %+v", f.toFilePath, err)
	}
}

func (f *FileUploadController) Upload() {
	_, fh, err := f.GetFile("upload_file")
	if err != nil {
		f.internalError(err)
		return
	}
	targetFilePath := f.toFilePath

	logs.Info("User: %s uploaded file from %s to %s.", f.currentUser.Username, fh.Filename, targetFilePath)
	err = f.SaveToFile("upload_file", filepath.Join(targetFilePath, fh.Filename))
	if err != nil {
		f.internalError(err)
	}
}

func (f *FileUploadController) DownloadProbe() {
	if isEmpty, err := service.IsEmptyDirectory(f.toFilePath); isEmpty || err != nil {
		f.customAbort(http.StatusNotFound, "No uploaded file found.")
		return
	}
}

func (f *FileUploadController) Download() {
	fileName := f.GetString("file_name")
	if fileName == "" {
		logs.Info("Will zip files to be downloaded as no file name specified.")
		attachmentFilePath := filepath.Join(baseRepoPath(), f.currentUser.Username)
		err := service.ZipFiles(filepath.Join(attachmentFilePath, attachmentFile), f.toFilePath)
		if err != nil {
			f.customAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to zip file for attachment: %+v", err))
			return
		}
		f.toFilePath = attachmentFilePath
		fileName = attachmentFile
	}
	logs.Debug("Download file from path: %s", f.toFilePath)
	f.Ctx.Output.Download(filepath.Join(f.toFilePath, fileName), fileName)
}

func (f *FileUploadController) ListFiles() {
	uploads, err := service.ListUploadFiles(f.toFilePath)
	if err != nil {
		f.internalError(err)
		return
	}
	f.renderJSON(uploads)
}

func (f *FileUploadController) RemoveFile() {
	fileInfo := model.FileInfo{
		Path:     f.toFilePath,
		FileName: f.GetString("file_name"),
	}
	logs.Info("Removed file: %s", filepath.Join(fileInfo.Path, fileInfo.FileName))
	err := service.RemoveUploadFile(fileInfo)
	if err != nil {
		f.internalError(err)
	}
}
