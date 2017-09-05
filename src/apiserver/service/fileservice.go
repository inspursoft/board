package service

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"path/filepath"
)

func ListUploadFiles(directory string) ([]model.FileInfo, error) {
	uploads := []model.FileInfo{}
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		upload := model.FileInfo{}
		upload.Path = path
		if !info.IsDir() {
			upload.FileName = info.Name()
			upload.Size = info.Size()
		}
		uploads = append(uploads, upload)
		return err
	})
	return uploads, nil
}
