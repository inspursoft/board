package controller

import (
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"
	"git/inspursoft/board/src/common/model"

	"bytes"

	"git/inspursoft/board/src/apiserver/service"

	"strconv"

	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var conf config.Configer
var tokenServerURL *url.URL

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
	tokenString := b.GetString("token")
	payload, err := verifyToken(tokenString)
	if err != nil {
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
	log.Printf("Get token from server: %s\n", token.TokenString)
	return &token, nil
}

func verifyToken(tokenString string) (map[string]interface{}, error) {
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

func init() {
	var err error
	conf, err = config.NewConfig("ini", "app.conf")
	if err != nil {
		log.Fatalf("Failed to load config file: %+v\n", err)
	}
	rawURL := conf.String("tokenServerURL")
	tokenServerURL, err = url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Failed to parse token server URL: %+v\n", err)
	}
	log.Printf("Set tokenservice URL as %s", tokenServerURL.String())
}
