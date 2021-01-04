package service

import (
	"fmt"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"strconv"
)

const (
	K8SPROXY_ENABLE_KEY     = "k8s_proxy_enable"
	DEFAULT_K8SPROXY_ENABLE = false
)

func SetK8SProxyConfig(config model.K8SProxyConfig) error {
	_, err := dao.AddOrUpdateConfig(model.Config{Name: K8SPROXY_ENABLE_KEY, Value: fmt.Sprint(config.Enable), Comment: fmt.Sprintf("Set config %s.", K8SPROXY_ENABLE_KEY)})
	return err
}

func GetK8SProxyConfig() (*model.K8SProxyConfig, error) {
	config, err := dao.GetConfig(K8SPROXY_ENABLE_KEY)
	if err != nil {
		return nil, err
	}
	if config.Value == "" {
		return &model.K8SProxyConfig{Enable: DEFAULT_K8SPROXY_ENABLE}, nil
	}
	enable, err := strconv.ParseBool(config.Value)
	if err != nil {
		return nil, err
	}
	return &model.K8SProxyConfig{Enable: enable}, nil
}
