package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"git/inspursoft/board/src/common/model"

	"encoding/json"

	"bytes"

	"git/inspursoft/board/src/apiserver/service"

	"time"

	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
)

var memoryCache cache.Cache

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

func (b *baseController) getCurrentUser() (*model.User, error) {
	userID, err := b.GetInt("user_id")
	if err != nil {
		return nil, err
	}
	payload, err := verifyToken(strconv.Itoa(userID))
	if err != nil {
		memoryCache.Delete(strconv.Itoa(userID))
		return nil, err
	}

	if strID, ok := payload["id"].(string); ok {
		userID, err := strconv.Atoi(strID)
		if err != nil {
			return nil, err
		}
		return service.GetUserByID(int64(userID))
	}
	return nil, err
}

func (b *baseController) signToken(key interface{}, payload map[string]interface{}) error {
	var err error
	reqData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post("http://localhost:4000/tokenservice/token", "application/json", bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var token model.Token
	err = json.Unmarshal(respData, &token)
	if err != nil {
		return err
	}
	log.Printf("Get token from server: %s\n", token.TokenString)
	memoryCache.Put(key.(string), token.TokenString, time.Second*1800)
	return nil
}

func verifyToken(key string) (map[string]interface{}, error) {
	token := memoryCache.Get(key)
	if tokenString, ok := token.(string); ok {
		resp, err := http.Get("http://localhost:4000/tokenservice/token?token=" + tokenString)
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
	return nil, fmt.Errorf("invalid token in cache")
}

func init() {
	var err error
	memoryCache, err = cache.NewCache("memory", "")
	if err != nil {
		log.Fatalf("Failed to init memory cache: %+v", err)
	}
}
