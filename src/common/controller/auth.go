package controller

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
	return token.TokenString, true
}
