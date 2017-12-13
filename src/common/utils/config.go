package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

var configStorage map[string]interface{}

func add(name string, value interface{}) {
	configStorage[name] = value
}

func AddEnv(name string, defaultValue ...string) {
	value := os.Getenv(name)
	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	add(name, value)
}

func AddValue(name string, value interface{}) {
	add(name, value)
}

func GetIntValue(name string) int {
	if v, ok := configStorage[name].(int); ok {
		return v
	}
	panic(fmt.Sprintf("Failed to get int value for key: %s", name))
}

func GetBoolValue(name string) bool {
	s := configStorage[name]
	v, err := strconv.ParseBool(s.(string))
	if err != nil {
		panic(fmt.Sprintf("Failed to get bool value for key: %s", name))
	}
	return v
}

func GetStringValue(name string, defaultValue ...string) string {
	if defaultValue != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
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

func GetConfig(name string, defaultValue ...string) func() string {
	return func() string { return GetStringValue(name, defaultValue...) }
}

func Initialize() {
	configStorage = make(map[string]interface{})
}

func ShowAllConfigs() {
	logs.Info("Current configurations in storage:\n")
	for k, v := range configStorage {
		if strings.Contains(strings.ToUpper(k), "PASSWORD") || strings.Contains(strings.ToUpper(k), "PWD") {
			continue
		}
		logs.Info("\t%s: %v", k, v)
	}
}
