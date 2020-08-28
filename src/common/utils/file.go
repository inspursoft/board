package utils

import (
	"io/ioutil"

	"github.com/astaxie/beego/logs"
)

func GetContentFromFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error("Failed to read file from path: %s, error: %+v", path, err)
		return ""
	}
	return string(data)
}
