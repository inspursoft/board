package controller

import (
	"io/ioutil"
	"net/http"
	"time"

	"encoding/json"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"bytes"

	"git/inspursoft/board/src/apiserver/service"

	"strconv"

	"net/url"

	"strings"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var conf config.Configer
var tokenServerURL *url.URL
var tokenCacheExpireSeconds int
var memoryCache cache.Cache

var kubeMasterURL = utils.GetConfig("KUBE_MASTER_URL")
var registryURL = utils.GetConfig("REGISTRY_URL")
var registryBaseURI = utils.GetConfig("REGISTRY_BASE_URI")

type baseController struct {
	beego.Controller
	currentUser    *model.User
	isSysAdmin     bool
	isProjectAdmin bool
}

func (b *baseController) Render() error {
	return nil
}

func (b *baseController) resolveBody() ([]byte, error) {
	data, err := ioutil.ReadAll(b.Ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type messageStatus struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func (b *baseController) serveStatus(status int, message string) {
	ms := messageStatus{
		StatusCode: status,
		Message:    message,
	}
	b.Data["json"] = ms
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.ServeJSON()
}

func (b *baseController) internalError(err error) {
	logs.Error("Error occurred: %+v", err)
	b.CustomAbort(http.StatusInternalServerError, "Unexpected error occurred.")
}

func (b *baseController) customAbort(status int, body string) {
	logs.Error("Error occurred: %s", body)
	b.CustomAbort(status, body)
}

func (b *baseController) getCurrentUser() *model.User {
	token := b.Ctx.Request.Header.Get("token")
	if token == "" {
		token = b.GetString("token")
	}

	payload, err := verifyToken(token)
	if err != nil {
		return nil
	}

	newToken, err := signToken(payload)
	if err != nil {
		logs.Error("failed to re-assign token: %+v", err)
		return nil
	}

	if strID, ok := payload["id"].(string); ok {
		userID, err := strconv.Atoi(strID)
		if err != nil {
			logs.Error("Error occurred on converting userID: %+v\n", err)
			return nil
		}
		user, err := service.GetUserByID(int64(userID))
		if err != nil {
			logs.Error("Error occurred while getting user by ID: %d\n", err)
			return nil
		}
		if !memoryCache.IsExist(user.Username) {
			logs.Error("Token has been expired forcely.")
			return nil
		}
		memoryCache.Put(user.Username, newToken.TokenString, time.Second*time.Duration(tokenCacheExpireSeconds))
		b.Ctx.ResponseWriter.Header().Set("token", newToken.TokenString)

		return user
	}
	return nil
}

func (b *baseController) signOff() error {
	username := b.GetString("username")
	err := memoryCache.Delete(username)
	if err != nil {
		logs.Error("Failed to delete username from memory cache: %+v", err)
		return err
	}
	logs.Info("Successful sign off from API server.")
	return nil
}

func signToken(payload map[string]interface{}) (*model.Token, error) {
	var err error
	reqData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(tokenServerURL.String(), "application/json", bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var token model.Token
	err = json.Unmarshal(respData, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func verifyToken(tokenString string) (map[string]interface{}, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, fmt.Errorf("no token was provided")
	}
	resp, err := http.Get(tokenServerURL.String() + "?token=" + tokenString)
	if err != nil {
		return nil, err
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload map[string]interface{}
	err = json.Unmarshal(respData, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func InitController() {
	var err error
	conf, err = config.NewConfig("ini", "app.conf")
	if err != nil {
		logs.Error("Failed to load config file: %+v\n", err)
	}
	rawURL := conf.String("tokenServerURL")
	tokenServerURL, err = url.Parse(rawURL)
	if err != nil {
		logs.Error("Failed to parse token server URL: %+v\n", err)
	}
	tokenCacheExpireSeconds, err = conf.Int("tokenCacheExpireSeconds")
	if err != nil {
		logs.Error("Failed to parse token expire seconds: %+v\n", err)
	}

	logs.Info("Set token server URL as %s and will expiration time after %d second(s) in cache", tokenServerURL.String(), tokenCacheExpireSeconds)
	memoryCache, err = cache.NewCache("memory", `{"interval": 3600}`)
	if err != nil {
		logs.Error("Failed to initialize cache: %+v\n", err)
	}

	logs.Info("Initialize serve repo\n")
	_, err = service.InitBareRepo(repoServePath)
	if err != nil {
		logs.Error("Failed to initialize serve repo: %+v\n", err)
	}

	beego.BConfig.MaxMemory = 1 << 22

}
