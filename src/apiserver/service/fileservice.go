package service

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"os"
	"path/filepath"
)

const (
	metaFile       = "META.cfg"
	processImage   = "process-image"
	processService = "process-service"
	rollingUpdate  = "rolling-update"
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
	f.WriteString("[Configurations]\n")
	for key, value := range configurations {
		fmt.Fprintf(f, "%s=%s\n", key, value)
	}
	directories := []string{processImage, processService, rollingUpdate}
	for _, dir := range directories {
		targetDir := filepath.Join(targetPath, dir)
		tempFile := filepath.Join(targetDir, ".placehold.tmp")
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			return err
		}
		f, err := os.Create(tempFile)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}
