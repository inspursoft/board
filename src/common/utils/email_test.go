package utils_test

import (
	"github.com/inspursoft/board/src/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockEmailIdentity = ""
var mockEmailHost = "smtp.myserver.com"
var mockEmailUsername = "tester@myserver.com"
var mockEmailPassword = ""
var mockEmailPort = 12225
var mockEmailTLS = false
var mockEmailFrom = "tester@myserver.com"
var mockEmailTo = "target@myserver.com"

func TestSendEmail(t *testing.T) {
	err := utils.NewEmailHandler(mockEmailIdentity, mockEmailUsername, mockEmailPassword, mockEmailHost, mockEmailPort).
		IsTLS(mockEmailTLS).
		Send(mockEmailFrom, []string{mockEmailTo}, "Testing", "This is a <b>test12345</b> mail.")
	assert := assert.New(t)
	assert.Nilf(err, "Failed to send email: %+v", err)
}

func TestPingEmail(t *testing.T) {
	err := utils.NewEmailHandler(mockEmailIdentity, mockEmailUsername, mockEmailPassword, mockEmailHost, mockEmailPort).
		IsTLS(mockEmailTLS).
		Ping()
	assert := assert.New(t)
	assert.Nilf(err, "Failed to ping email server: %+v", err)
}
