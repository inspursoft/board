package service

import (
	"git/inspursoft/board/src/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	authMode    = "db_auth"
	targetName  = "REDIRECTION_URL"
	targetValue = "http://test.domain.com"
)

func TestGetSystemInfo(t *testing.T) {
	systemInfo, err := GetSystemInfo()
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while getting system info: %+v", err)
	assert.NotNilf(systemInfo, "Failed to get system info: %+v", err)
	assert.Equalf(authMode, systemInfo.AuthMode, "System info auth_mode value is not as expected: %s", authMode)
}

func TestSetSystemInfo(t *testing.T) {
	utils.Initialize()
	utils.SetConfig(targetName, targetValue)
	err := SetSystemInfo(targetName, true)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while setting system info: %+v", err)
	systemInfo, err := GetSystemInfo()
	assert.Equalf(targetValue, systemInfo.RedirectionURL, "System info %s is not as expected: %s", targetName, targetValue)
}
