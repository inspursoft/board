package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"

	"net/http"
	"os"
	"path/filepath"
)

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
	user := f.getCurrentUser()
	if user == nil {
		f.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	f.currentUser = user
	f.resolveRepoPath()
	f.resolveFilePath()
}

func (f *FileUploadController) resolveFilePath() {
	projectName := f.GetString("project_name")
	serviceID, err := f.GetInt64("service_id", 0)
	if err != nil {
		f.internalError(err)
		return
	}
	imageName := f.GetString("image_name")
	tagName := f.GetString("tag_name")

	reqUploadFile := uploadFile{
		ProjectName: projectName,
		ServiceID:   serviceID,
		ImageName:   imageName,
		TagName:     tagName,
	}

	if reqUploadFile.ProjectName == "" && reqUploadFile.ServiceID == 0 {
		f.customAbort(http.StatusBadRequest, "No project name or service ID provided.")
		return
	}

	if reqUploadFile.ImageName == "" && reqUploadFile.TagName == "" {
		f.customAbort(http.StatusBadRequest, "No image name or tag name provided.")
		return
	}

	if reqUploadFile.ProjectName != "" {
		isMember, err := service.IsProjectMemberByName(reqUploadFile.ProjectName, f.currentUser.ID)
		if err != nil {
			f.internalError(err)
			return
		}
		if !isMember {
			f.customAbort(http.StatusForbidden, "Not member to the current project with provided ID.")
			return
		}
		f.toFilePath = filepath.Join(imageProcess, reqUploadFile.ImageName, reqUploadFile.TagName, "upload")
	}
}

func (f *FileUploadController) Upload() {
	_, fh, err := f.GetFile("upload_file")
	if err != nil {
		f.internalError(err)
		return
	}
	targetFilePath := filepath.Join(f.repoPath, f.toFilePath)
	os.MkdirAll(targetFilePath, 0755)

	logs.Info("User: %s uploaded file from %s to %s.", f.currentUser.Username, fh.Filename, targetFilePath)
	err = f.SaveToFile("upload_file", filepath.Join(targetFilePath, fh.Filename))
	if err != nil {
		f.internalError(err)
	}
}

func (f *FileUploadController) ListFiles() {
	uploads, err := service.ListUploadFiles(filepath.Join(f.repoPath, f.toFilePath))
	if err != nil {
		f.internalError(err)
		return
	}
	f.Data["json"] = uploads
	f.ServeJSON()
}

func (f *FileUploadController) RemoveFile() {
	fileInfo := model.FileInfo{
		Path:     filepath.Join(f.repoPath, f.toFilePath),
		FileName: f.GetString("file_name"),
	}
	logs.Info("Removed file: %s", filepath.Join(fileInfo.Path, fileInfo.FileName))
	err := service.RemoveUploadFile(fileInfo)
	if err != nil {
		f.internalError(err)
	}
}
