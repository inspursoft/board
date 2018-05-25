package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockEmailIdentity = ""
var mockEmailHost = "smtp.inspur.com"
var mockEmailUsername = "admin"
var mockEmailPassword = "123456a?"
var mockEmailPort = 25
var mockEmailTLS = false
var mockEmailFrom = "admin@inspur.com"
var mockEmailTo = "tester@inspur.com"

func TestSendingEmail(t *testing.T) {
	err := NewEmailHandler(mockEmailIdentity, mockEmailUsername, mockEmailPassword, mockEmailHost, mockEmailPort).
		IsTLS(mockEmailTLS).
		Send(mockEmailFrom, []string{mockEmailTo}, "Testing", "This is a test mail.")
	assert := assert.New(t)
	assert.Nilf(err, "Failed to send email: %+v", err)
}

func TestPingEmail(t *testing.T) {
	err := NewEmailHandler(mockEmailIdentity, mockEmailUsername, mockEmailPassword, mockEmailHost, mockEmailPort).
		IsTLS(mockEmailTLS).
		Ping()
	assert := assert.New(t)
	assert.Nilf(err, "Failed to ping email server: %+v", err)
}
