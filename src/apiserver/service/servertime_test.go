package service_test

import (
	"github.com/inspursoft/board/src/apiserver/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeNow(t *testing.T) {
	assert := assert.New(t)
	timeNow := service.GetServerTime()
	assert.NotZero(timeNow.TimeNow, "Error occurred while testing GetServerTime.")
}
