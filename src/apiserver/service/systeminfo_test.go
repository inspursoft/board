package service_test

import (
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	AuthMode    = "db_auth"
	targetName  = "REDIRECTION_URL"
	targetValue = "http://test.domain.com"
)

func TestSetSystemInfo(t *testing.T) {
	utils.Initialize()
	utils.SetConfig(targetName, targetValue)
	err := service.SetSystemInfo(targetName, true)
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while setting system info: %+v", err)
	systemInfo, err := service.GetSystemInfo()
	assert.Equalf(targetValue, systemInfo.RedirectionURL, "System info %s is not as expected: %s", targetName, targetValue)
}

func TestGetSystemInfo(t *testing.T) {
	systemInfo, err := service.GetSystemInfo()
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while getting system info: %+v", err)
	assert.NotNilf(systemInfo, "Failed to get system info: %+v", err)
	assert.Equalf(targetValue, systemInfo.RedirectionURL, "System info auth_mode value is not as expected: %s", AuthMode)
}
