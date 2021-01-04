package auth

import (
	"fmt"
	"github.com/inspursoft/board/src/common/model"

	"time"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
)

const (
	defaultFailedTimes  = 10
	defaultDenyDuration = 120 //Seconds
)

type Auth interface {
	DoAuth(principal, password string) (*model.User, error)
}

var registry map[string]Auth

func GetAuth(AuthMode string) (*Auth, error) {
	if auth, ok := registry[AuthMode]; ok {
		return &auth, nil
	}
	return nil, fmt.Errorf("unsupported auth mode: %s", AuthMode)
}

func registerAuth(AuthMode string, auth Auth) {
	if registry == nil {
		registry = make(map[string]Auth)
	}
	if _, ok := registry[AuthMode]; !ok {
		registry[AuthMode] = auth
		logs.Debug("Auth mode: %s has been registered.", AuthMode)
	}
}

var temporaryStore cache.Cache

type AuthCheckResult struct {
	TimeNow              time.Time
	FailedTimes          int
	IsTemporarilyBlocked bool
	TimeElapsed          int
	TimeRemain           int
}

// Check the failed times by user name, TODO check by request IP
func CheckAuthFailedTimes(principal string) (result AuthCheckResult) {
	// a new user name, pass, TODO check request IP failed times
	if temporaryStore.IsExist(principal) {
		if r, ok := temporaryStore.Get(principal).(AuthCheckResult); ok {
			r.FailedTimes++
			logs.Debug("Increase auth check failed times with principal: %s, with times: %d", principal, r.FailedTimes)
			if r.FailedTimes >= defaultFailedTimes {
				if !r.IsTemporarilyBlocked {
					r.TimeNow = time.Now()
				}
				if r.TimeElapsed < defaultDenyDuration {
					r.IsTemporarilyBlocked = true
					r.TimeElapsed = int(time.Since(r.TimeNow).Seconds())
					logs.Debug("Blocking auth principal: %s as it has not reached the deny duration: %d", principal, r.TimeRemain)
				}
				timeRemain := defaultDenyDuration - r.TimeElapsed
				if timeRemain >= 0 {
					r.TimeRemain = timeRemain
				} else {
					r.FailedTimes = 0
					r.TimeNow = time.Now()
					r.TimeElapsed = 0
					r.TimeRemain = 0
					r.IsTemporarilyBlocked = false
					logs.Debug("Unlocking auth principal: %s as it has reached deny duration", principal)
				}
			}
			temporaryStore.Put(principal, r, 300*time.Second)
			result = r
		}
	} else {
		result.FailedTimes = 0
		logs.Debug("Store auth check with principal: %s.", principal)
		temporaryStore.Put(principal, result, 300*time.Second)
	}
	return
}

//Reset the access check times
func ResetAuthFailedTimes(principal string, requestaddr string) error {
	return temporaryStore.Delete(principal)
}

func init() {
	temporaryStore = cache.NewMemoryCache()
}
