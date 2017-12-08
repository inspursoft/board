package service

import (
	"fmt"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
)

func GetSystemInfo() (*model.SystemInfo, error) {
	configs, err := dao.GetAllConfigs()
	if err != nil {
		return nil, err
	}
	var systemInfo model.SystemInfo
	for _, config := range configs {
		switch config.Name {
		case "BOARD_HOST":
			systemInfo.BoardHost = config.Value
		case "AUTH_MODE":
			systemInfo.AuthMode = config.Value
		case "SET_ADMIN_PASSWORD":
			systemInfo.SetAdminPassword = config.Value
		case "INIT_PROJECT_REPO":
			systemInfo.InitProjectRepo = config.Value
		}
	}
	return &systemInfo, nil
}

func SetSystemInfo(name string, reconfigurable bool) error {
	config, err := dao.GetConfig(name)
	if err != nil {
		return err
	}
	if config.Name == "" || reconfigurable {
		value := utils.GetStringValue(name)
		if value == "" {
			return fmt.Errorf("Has not set config %s yet", name)
		}
		_, err := dao.AddOrUpdateConfig(model.Config{Name: name, Value: value, Comment: fmt.Sprintf("Set config %s as %s.", name, value)})
		return err
	}
	utils.SetConfig(name, config.Value)
	return nil
}
