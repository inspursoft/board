package auth

import (
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type indataAccount struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"name"`
}

type postParam struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type InDataAuth struct{}

func (auth InDataAuth) DoAuth(principal, password string) (*model.User, error) {

	verificationURL := utils.GetStringValue("VERIFICATION_URL")
	logs.Debug("Verification URL: %s", verificationURL)
	logs.Debug("External token: %s", principal)

	param := postParam{
		Token: principal,
		Type:  "id_token",
	}

	var account indataAccount
	err := utils.RequestHandle(http.MethodPost, verificationURL, func(req *http.Request) error {
		req.Header = http.Header{
			"content-type": []string{"application/json"},
		}
		return nil
	}, &param,
		func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &account)
		})
	if err != nil {
		logs.Error("Failed to create request: %+v", err)
		return nil, nil
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
	user, err := service.GetUserByName(u.Username)
	if err != nil {
		logs.Error("Failed to get user in SignIn: %+v\n", err)
		return nil, err
	}
	return dao.GetUser(*user, "username", "password")
}

func init() {
	registerAuth("indata_auth", InDataAuth{})
}
