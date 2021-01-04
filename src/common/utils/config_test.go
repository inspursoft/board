package utils_test

import (
	"github.com/inspursoft/board/src/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetConfig(t *testing.T) {
	assert := assert.New(t)
	utils.AddValue("URL", "test.inspursoft.com")
	utils.SetConfig("TEST_URL", "https://%s", "URL")
	testURL := utils.GetConfig("TEST_URL")
	assert.Equal("https://test.inspursoft.com", testURL(), "TestURLs are not equal.")
}

func TestSetDefaultConfig(t *testing.T) {
	assert := assert.New(t)
	value := utils.GetConfig("TEST_VALUE", "12345")
	assert.Equal("12345", value(), "Values get from config are not equal.")
}

func init() {
	utils.Initialize()
}
