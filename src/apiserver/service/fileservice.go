package service

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"os"
	"path/filepath"
)

const (
	metaFile = "META.cfg"
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

func CreateBaseDirectory(configurations map[string]string, targetPath string) error {
	if configurations == nil {
		return fmt.Errorf("configuration for generating base directory is nil")
	}
	f, err := os.OpenFile(filepath.Join(targetPath, metaFile), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create META.cfg file: %+v", err)
	}
	defer f.Close()
	f.WriteString("[para]\n")
	for key, value := range configurations {
		fmt.Fprintf(f, "%s=%s\n", key, value)
	}
	return nil
}
