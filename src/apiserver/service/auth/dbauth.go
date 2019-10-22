package auth

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
)

type DbAuth struct{}

func (auth DbAuth) DoAuth(principal, password string) (*model.User, error) {
	user, err := service.GetUserByName(principal)
	if err != nil {
		logs.Error("Failed to get user in SignIn: %+v\n", err)
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	query := model.User{Username: principal, Password: password}
	query.Password = utils.Encrypt(query.Password, user.Salt)
	return dao.GetUser(query, "username", "password")
}

func init() {
	registerAuth("db_auth", DbAuth{})
}
