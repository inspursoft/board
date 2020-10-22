package commons

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/auth"
	t "git/inspursoft/board/src/common/token"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
)

var ReservedUsernames = [...]string{"explore", "create", "assets", "css", "img", "js", "less", "plugins", "debug", "raw", "install", "api", "avatar", "user", "org", "help", "stars", "issues", "pulls", "commits", "repo", "template", "new", ".", ".."}
var defaultFailedTimes = 10
var defaultFailedTimesForCaptcha = 3

type resolveType int

const (
	failedRetries resolveType = iota + 1
	invalidCaptcha
	temporaryBlock
	invalidPassword
)

type resolveInfo struct {
	Retries     int         `json:"resolve_sign_in_retries"`
	Type        resolveType `json:"resolve_sign_in_type"`
	Description string      `json:"resolve_sign_in_description"`
	Value       int         `json:"resolve_sign_in_value"`
}

func (ca *BaseController) ProcessAuth(principal, password string) (string, bool) {
	var currentAuth *auth.Auth
	var err error
	var message string
	var info resolveInfo
	//Check signin failed times
	authCheck := auth.CheckAuthFailedTimes(principal)
	info.Retries = authCheck.FailedTimes
	if authCheck.IsTemporarilyBlocked {
		message = "Temporarily blocked."
		info.Type = temporaryBlock
		info.Description = "Temporarily blocked."
		info.Value = authCheck.TimeRemain
		ca.ServeJSONOutput(http.StatusBadRequest, info)
		return "", false
	}
	var validateCaptcha bool
	if v, ok := MemoryCache.Get("validate_captcha").(bool); ok {
		validateCaptcha = v
	}
	if authCheck.FailedTimes >= defaultFailedTimesForCaptcha && authCheck.FailedTimes < defaultFailedTimes || validateCaptcha {
		MemoryCache.Put("validate_captcha", true, 300*time.Second)
		captchaID := ca.GetString("captcha_id")
		challenge := ca.GetString("captcha")
		if !Cpt.Verify(captchaID, challenge) {
			message = "Invalid captcha."
			info.Type = invalidCaptcha
			info.Description = message
			ca.ServeJSONOutput(http.StatusBadRequest, info)
			return "", false
		}
	}

	if principal == "boardadmin" {
		currentAuth, err = auth.GetAuth("db_auth")
	} else {
		currentAuth, err = auth.GetAuth(AuthMode())
	}
	if err != nil {
		ca.InternalError(err)
		return "", false
	}
	password, err = service.DecodeUserPassword(password)
	if err != nil {
		message = "Invalid password content."
		info.Type = invalidPassword
		info.Description = "Invalid password content."
		ca.ServeJSONOutput(http.StatusBadRequest, info)
	}

	user, err := (*currentAuth).DoAuth(principal, password)
	if err != nil {
		ca.InternalError(err)
		return "", false
	}

	if user == nil {
		message = "Incorrect username or password."
		info.Type = failedRetries
		info.Description = message
		ca.ServeJSONOutput(http.StatusBadRequest, info)
		return "", false
	}
	MemoryCache.Delete("validate_captcha")
	if existing := MemoryCache.IsExist(user.Username); existing {
		if lastToken, ok := MemoryCache.Get(user.Username).(string); ok {
			logs.Info("Found last token stored in cache, will be removing it ...")
			MemoryCache.Delete(lastToken)
		}
		logs.Info("Found last user stored in cache, will be removing it ...")
		MemoryCache.Delete(user.Username)
	}
	payload := make(map[string]interface{})
	payload["id"] = strconv.Itoa(int(user.ID))
	payload["username"] = user.Username
	payload["email"] = user.Email
	payload["realname"] = user.Realname
	payload["is_system_admin"] = user.SystemAdmin

	token, err := t.SignToken(TokenServerURL(), payload)
	if err != nil {
		ca.InternalError(err)
		return "", false
	}
	MemoryCache.Put(user.Username, token.TokenString, DefaultCacheDuration)
	MemoryCache.Put(token.TokenString, payload, DefaultCacheDuration)
	ca.AuditUser, _ = service.GetUserByName(user.Username)

	//Reset the user failed times
	err = auth.ResetAuthFailedTimes(principal, ca.Ctx.Request.RemoteAddr)
	if err != nil {
		ca.InternalError(err)
		return "", false
	}

	return token.TokenString, true
}
