package utils

import (
	"errors"
	"io"
	"io/ioutil"

	yaml "github.com/ghodss/yaml"
)

func UnmarshalYamlFile(file io.Reader, config interface{}) error {
	fileInfo, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("InternalError:" + err.Error())
	}
	return yaml.Unmarshal(fileInfo, config)
}
