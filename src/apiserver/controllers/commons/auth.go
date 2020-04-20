package commons

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/auth"
	"net/http"
	"strconv"
	"time"
)

var ReservedUsernames = [...]string{"explore", "create", "assets", "css", "img", "js", "less", "plugins", "debug", "raw", "install", "api", "avatar", "user", "org", "help", "stars", "issues", "pulls", "commits", "repo", "template", "new", ".", ".."}

func (ca *BaseController) ProcessAuth(principal, password string) (string, bool) {
	var currentAuth *auth.Auth
	var err error
	//Check in MemoryCache
	failedtimes, quickdeny := auth.CacheCheckAuthFailedTimes(principal)
	if quickdeny {
		ca.Ctx.SetCookie("failedtimes", string(failedtimes))
		ca.ServeStatus(http.StatusNotAcceptable, "NotAcceptable")
		return "", false
	}

	//Check signin failed times
	failedtimes, deny, _ := auth.CheckAuthFailedTimes(principal)
	if deny {
		ca.Ctx.SetCookie("failedtimes", string(failedtimes))
		ca.ServeStatus(http.StatusNotAcceptable, "NotAcceptable.")
		return "", false
	}

	if principal == "admin" {
		currentAuth, err = auth.GetAuth("db_auth")
	} else {
		currentAuth, err = auth.GetAuth(AuthMode())
	}
	if err != nil {
		ca.InternalError(err)
		return "", false
	}
	user, err := (*currentAuth).DoAuth(principal, password)
	if err != nil {
		ca.InternalError(err)
		return "", false
	}

	if user == nil {
		// Signin failed, update the user access info
		failedtimes, err = auth.UpdateAuthFailedTimes(principal, ca.Ctx.Request.RemoteAddr)
		if err != nil {
			ca.InternalError(err)
			return "", false
		}
		ca.Ctx.SetCookie("failedtimes", string(failedtimes))
		ca.ServeStatus(http.StatusBadRequest, "Incorrect username or password.")
		return "", false
	}
	payload := make(map[string]interface{})
	payload["id"] = strconv.Itoa(int(user.ID))
	payload["username"] = user.Username
	payload["email"] = user.Email
	payload["realname"] = user.Realname
	payload["is_system_admin"] = user.SystemAdmin
	token, err := ca.SignToken(payload)
	if err != nil {
		ca.InternalError(err)
		return "", false
	}
	MemoryCache.Put(user.Username, token.TokenString, time.Second*time.Duration(TokenCacheExpireSeconds))
	MemoryCache.Put(token.TokenString, payload, time.Second*time.Duration(TokenCacheExpireSeconds))
	ca.AuditUser, _ = service.GetUserByName(user.Username)

	//Reset the user failed times
	err = auth.ResetAuthFailedTimes(principal, ca.Ctx.Request.RemoteAddr)
	if err != nil {
		ca.InternalError(err)
		return "", false
	}

	return token.TokenString, true
}
