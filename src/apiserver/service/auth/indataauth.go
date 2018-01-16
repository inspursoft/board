package auth

import (
	"bytes"
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"

	"net/http"

	"github.com/astaxie/beego/logs"
)

type indataAccount struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"name"`
}

type InDataAuth struct{}

func (auth InDataAuth) DoAuth(principal, password string) (*model.User, error) {

	verificationURL := utils.GetStringValue("VERIFICATION_URL")
	logs.Debug("Verification URL: %s", verificationURL)
	logs.Debug("External token: %s", principal)

	params := make(map[string]string)
	params["token"] = principal
	reqData, err := json.Marshal(params)
	if err != nil {
		logs.Error("Failed to marshal token from request: %+v", err)
		return nil, nil
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", verificationURL, bytes.NewReader(reqData))
	if err != nil {
		logs.Error("Failed to create request: %+v", err)
		return nil, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		logs.Error("Failed request remote endpoint: %+v", err)
	}
	if resp == nil {
		logs.Error("Failed to get response from request.")
		return nil, nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("Failed to read from response body: %+v", err)
	}

	var account indataAccount
	err = json.Unmarshal(data, &account)
	if err != nil {
		logs.Error("Failed to unmarshal response data: %+v", err)
	}
	if account.Username == "" {
		logs.Error("Invalid token for request verification.")
		return nil, nil
	}

	var u model.User
	u.Username = account.Username
	u.Email = account.Email

	logs.Debug("username:", u.Username, ",email:", u.Email)

	exist, err := service.UserExists("username", u.Username, 0)
	if err != nil {
		logs.Debug("err: %+v", err)
		return nil, err
	}

	if !exist {
		u.Realname = account.FullName
		u.Password = "12345678AbC"
		u.Comment = "registered from InData platform."
		if u.Email == "" {
			u.Email = u.Username + "@placeholder.com"
		}
		_, err := service.SignUp(u)
		if err != nil {
			return nil, err
		}
	}
	user, err := dao.GetUser(u, "username", "deleted")
	if err != nil {
		logs.Error("Failed to get user in SignIn: %+v\n", err)
		return nil, err
	}
	return dao.GetUser(*user, "username", "password")
}

func init() {
	registerAuth("indata_auth", InDataAuth{})
}
