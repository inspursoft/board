package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeNow(t *testing.T) {
	assert := assert.New(t)
	timeNow := GetServerTime()
	assert.NotZero(timeNow.TimeNow, "Error occurred while testing GetServerTime.")
}
