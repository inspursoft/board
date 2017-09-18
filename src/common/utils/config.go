package utils

import (
	"fmt"
	"os"

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
	logs.Error("Failed to get value for key: %s", name)
	return 0
}

func GetStringValue(name string) string {
	if s, ok := configStorage[name].(string); ok {
		return s
	}
	logs.Error("Failed to get value for key: %s", name)
	return ""
}

func SetConfig(name, formatter string, keys ...string) {
	configStorage[name] = fmt.Sprintf(formatter, func() (values []interface{}) {
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
		logs.Info("%s: %s\n", k, v)
	}
}
