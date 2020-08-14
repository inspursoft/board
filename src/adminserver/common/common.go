package common

import (
	"errors"

	"github.com/alyu/configparser"
)

var ErrAdminLogin = errors.New("another admin user has signed in other place")
var ErrForbidden = errors.New("Forbidden")
var ErrWrongPassword = errors.New("Wrong password")
var ErrTokenServer = errors.New("tokenserver is down")
var ErrNoData = errors.New("Board already uninstalled")

func ReadCfgItem(items ...string) (string, error) {
	var cfgPath string
	if len(items) == 1 {
		cfgPath = "/go/cfgfile/board.cfg"
	} else {
		cfgPath = "/go/cfgfile/board.cfg.tmp"
	}

	config, err := configparser.Read(cfgPath)
	if err != nil {
		return "", err
	}
	//section sensitive, global refers to all sections.
	section, err := config.Section("global")
	if err != nil {
		return "", err
	}
	return section.ValueOf(items[0]), nil
}
