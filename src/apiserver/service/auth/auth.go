package auth

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"

	"errors"
	"time"

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
			logs.Debug("Failed times %n, Last Updated %v", user.FailedTimes, user.UpdateTime)
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
	_, err = service.UpdateUser(*user, "failed_times")
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
	_, err = service.UpdateUser(*user, "failed_times")
	if err != nil {
		logs.Error("Failed to udpated user in DB: %+v\n", err)
		return err
	}
	//TODO reset the request IP memory
	logs.Debug("Resetted the times user %s IP %s", principal, requestaddr)
	return nil
}
