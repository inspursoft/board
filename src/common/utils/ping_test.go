package utils_test

import (
	"github.com/inspursoft/board/src/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingIPAddr(t *testing.T) {
	status, err := utils.PingIPAddr("127.0.0.1")
	assert := assert.New(t)
	assert.Nilf(err, "Failed to ping IP address: %+v", err)
	assert.True(status, "IP Address cannot ping.")
}
