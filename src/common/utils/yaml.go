package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	yaml "github.com/ghodss/yaml"
)

func UnmarshalYamlFile(file io.Reader, config interface{}) error {
	fileInfo, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("InternalError:" + err.Error())
	}
	return yaml.Unmarshal(fileInfo, config)
}

func GenerateFile(fileInfo []byte, loadPath string, fileName string) error {
	err := CheckFilePath(loadPath)
	if err != nil {
		logs.Error("Check yaml file path error, err:%+v\n", err)
		return err
	}
	absFileName := filepath.Join(loadPath, fileName)
	err = ioutil.WriteFile(absFileName, fileInfo, 0644)
	if err != nil {
		logs.Error("Generate yaml file failed, err:%+v\n", err)
		return err
	}
	return nil
}

func CheckFilePath(loadPath string) error {
	if fi, err := os.Stat(loadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(loadPath, 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return errors.New("ERR_DEPLOYMENT_PATH_NOT_DIRECTORY")
	}
	return nil
}

func UnmarshalYamlData(data []byte, config interface{}, callback func(in interface{}) error) error {
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return callback(config)
}
