package auth

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"

	"errors"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
)

const (
	defaultFailedTimes  = 5
	defaultDenyDuration = 600 //Seconds
)

type Auth interface {
	DoAuth(principal, password string) (*model.User, error)
}

var registry map[string]Auth

var SignInCache cache.Cache

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

// Check the failed times by user name, TODO check by request IP
func CheckAuthFailedTimes(principal string) (int, bool, error) {
	user, err := service.GetUserByName(principal)
	if err != nil {
		logs.Error("Failed to get user in SignIn: %+v\n", err)
		return 0, false, err
	}
	if user == nil {
		// a new user name, pass, TODO check request IP failed times
		return 0, false, nil
	}
	if user.FailedTimes >= defaultFailedTimes {
		//Failed times more than the limit, check access deny duration
		if time.Since(user.UpdateTime).Seconds() < defaultDenyDuration {
			logs.Debug("Failed times %d, Last Updated %v", user.FailedTimes, user.UpdateTime)
			// Add a record in SignInCache for quick check
			err = SignInCache.Put(principal,
				model.UserSignInCache{FailedTimes: user.FailedTimes, UpdateTime: user.UpdateTime},
				time.Second*time.Duration(defaultFailedTimes))
			if err != nil {
				logs.Error("Failed SignIn cached %s %v", principal, err)
			}
			return user.FailedTimes, true, nil
		}
	}
	return user.FailedTimes, false, nil
}

// Update the failed times by user name
func UpdateAuthFailedTimes(principal string, requestaddr string) (int, error) {
	user, err := service.GetUserByName(principal)
	if err != nil {
		logs.Error("Failed to get user in SignIn: %+v\n", err)
		return 0, err
	}
	if user == nil {
		// TODO not a user name, update the request IP failed times in memorycache
		logs.Info("Deny no user %s request IP %s", principal, requestaddr)
		return 0, nil
	}
	user.FailedTimes = user.FailedTimes + 1
	_, err = service.UpdateUser(*user, "failed_times", "update_time")
	if err != nil {
		logs.Error("Failed to udpated user in DB: %+v\n", err)
		return user.FailedTimes, err
	}
	return user.FailedTimes, nil
}

//Reset the access check times
func ResetAuthFailedTimes(principal string, requestaddr string) error {
	user, err := service.GetUserByName(principal)
	if err != nil {
		logs.Error("Failed to get user in SignIn: %+v\n", err)
		return err
	}
	if user == nil {
		logs.Error("Failed to get user in DB %s", principal)
		return errors.New("Failed to get user in DB")
	}
	user.FailedTimes = 0
	_, err = service.UpdateUser(*user, "failed_times", "update_time")
	if err != nil {
		logs.Error("Failed to udpated user in DB: %+v\n", err)
		return err
	}
	//Remove the failed user record from SignInCache
	SignInCache.Delete(principal)

	//TODO reset the request IP memory
	logs.Debug("Resetted the times user %s IP %s", principal, requestaddr)
	return nil
}

// Check the failed times in Cache
func CacheCheckAuthFailedTimes(principal string) (int, bool) {
	//TODO remove this Debug
	logs.Debug("SignInCache: %v user: %s Exist:%v", SignInCache, principal, SignInCache.IsExist(principal))
	if record, ok := SignInCache.Get(principal).(model.UserSignInCache); ok {
		//Check access deny duration
		if time.Since(record.UpdateTime).Seconds() < defaultDenyDuration {
			logs.Debug("Cache check Failed times %d, Last Updated %v", record.FailedTimes, record.UpdateTime)
			return record.FailedTimes, true
		}
	}
	return 0, false
}

func init() {
	var err error
	logs.Debug("Init SignInCache")
	SignInCache, err = cache.NewCache("memory", `{"interval": 3600}`)
	if err != nil {
		logs.Error("Failed to initialize SignIn cache: %+v", err)
	}
}
