package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/astaxie/beego/logs"
)

var configStorage map[string]interface{}

func addValue(name string, value interface{}) {
	configStorage[name] = value
}

func getValue(name string) (interface{}, bool) {
	if val, exists := configStorage[name]; exists {
		return val, true
	}
	return nil, false
}

func AddEnv(name string) {
	addValue(name, os.Getenv(name))
}

func GetIntValue(name string) int {
	if v, ok := configStorage[name].(int); ok {
		return v
	}
	panic(fmt.Sprintf("Failed to get value for key: %s", name))
}

func GetStringValue(name string) string {
	if s, ok := configStorage[name].(string); ok {
		return s
	}
	panic(fmt.Sprintf("Failed to get value for key: %s", name))
}

func SetConfig(name, formatter string, keys ...string) {
	configStorage[name] = fmt.Sprintf(formatter,
		func() (values []interface{}) {
			for _, key := range keys {
				values = append(values, GetStringValue(key))
			}
			return
		}()...)
	return
}

func GetConfig(name string) func() string {
	return func() string { return GetStringValue(name) }
}

func Initialize() {
	configStorage = make(map[string]interface{})
}

func ShowAllConfigs() {
	logs.Info("Current configurations in storage:\n")
	for k, v := range configStorage {
		if strings.Contains(strings.ToUpper(k), "PASSWORD") {
			continue
		}
		logs.Info("\t%s: %s", k, v)
	}
}
