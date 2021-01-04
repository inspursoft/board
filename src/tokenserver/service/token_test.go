package service_test

import (
	"fmt"
	"github.com/inspursoft/board/src/tokenserver/service"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	TOKEN_BAD              = "abcd"
	TOKEN_OTHER_SIGNMETHOD = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA"
)

var PAYLOAD = map[string]interface{}{
	"id":               "1",
	"username":         "zhangsan",
	"email":            "zhangsan@inspur.com",
	"realname":         "zhangsan",
	"is_project_admin": float64(1),
	"is_system_admin":  float64(0),
}

func TestSign(t *testing.T) {
	assert := assert.New(t)

	// nil condition
	token, err := service.Sign(nil)
	assert.Nil(err, fmt.Sprintf("Sign nil map error: %+v", err))
	assert.NotEmpty(token, "The sign token is empty")

	// empty map
	token, err = service.Sign(make(map[string]interface{}))
	assert.Nil(err, fmt.Sprintf("Sign empty map error: %+v", err))
	assert.NotEmpty(token, "The sign token is empty")
}

func TestTokenWithInvalidPayload(t *testing.T) {
	assert := assert.New(t)

	_, err := service.Verify(TOKEN_BAD)
	assert.NotNil(err, "Verify bad token should failed")

	_, err = service.Verify(TOKEN_OTHER_SIGNMETHOD)
	assert.NotNil(err, "Verify ECDSASHA256 signed token should failed")
}

func TestTokenWithValidPayload(t *testing.T) {
	assert := assert.New(t)

	// test timeout token
	os.Setenv("TOKEN_EXPIRE_TIME", "1")
	service.InitService()
	token, err := service.Sign(PAYLOAD)
	assert.Nil(err, fmt.Sprintf("Sign payload error: %+v", err))
	assert.NotEmpty(token, "The sign token is empty")
	time.Sleep(2 * time.Second)
	v, err := service.Verify(token)
	assert.NotNil(err, fmt.Sprintf("Verify token should timeout error: %+v", err))

	// normal test
	os.Setenv("TOKEN_EXPIRE_TIME", "1200")
	service.InitService()
	token, err = service.Sign(PAYLOAD)
	assert.Nil(err, fmt.Sprintf("Sign payload error: %+v", err))
	assert.NotEmpty(token, "The sign token is empty")
	v, err = service.Verify(token)
	assert.Nil(err, fmt.Sprintf("Verify token error: %+v", err))
	assert.Equal(PAYLOAD, v, "Verify origin payload error")

}
