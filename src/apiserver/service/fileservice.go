package service

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"path/filepath"
)

func CreatePath(basePath string, directory string) (string, error) {
	path := filepath.Join(basePath, directory)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", os.MkdirAll(filepath.Join(basePath, directory), 0640)
	}
	return path, nil
}

func ListUploadFiles(directory string) ([]model.FileUploaded, error) {
	file, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	files, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}
	uploads := []model.FileUploaded{}
	for _, f := range files {
		uploaded := model.FileUploaded{
			Path:     filepath.Dir(f.Name()),
			FileName: filepath.FromSlash(f.Name()),
			Size:     f.Size(),
		}
		uploads = append(uploads, uploaded)
	}
	return uploads, nil
}
