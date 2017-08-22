package service

import (
	"os"
	"path/filepath"
)

func CreatePath(basePath string, directory string) error {
	path := filepath.Join(basePath, directory)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(filepath.Join(basePath, directory), 0640)
	}
	return nil
}
