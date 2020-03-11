package auth

import (
	"fmt"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
)

type Auth interface {
	DoAuth(principal, password string) (*model.User, error)
}

var registry map[string]Auth

func GetAuth(authMode string) (*Auth, error) {
	if auth, ok := registry[authMode]; ok {
		return &auth, nil
	}
	return nil, fmt.Errorf("unsupported auth mode: %s", authMode)
}

func registerAuth(authMode string, auth Auth) {
	if registry == nil {
		registry = make(map[string]Auth)
	}
	if _, ok := registry[authMode]; !ok {
		registry[authMode] = auth
		logs.Debug("Auth mode: %s has been registered.", authMode)
	}
}
