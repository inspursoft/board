package service

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"path/filepath"
)

func ListUploadFiles(directory string) ([]model.FileInfo, error) {
	uploads := []model.FileInfo{}
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			uploads = append(uploads, model.FileInfo{
				Path:     filepath.Dir(path),
				FileName: info.Name(),
				Size:     info.Size(),
			})
		}
		return err
	})
	return uploads, nil
}

func RemoveUploadFile(file model.FileInfo) error {
	return os.Remove(filepath.Join(file.Path, file.FileName))
}
