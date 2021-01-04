package vm

import (
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
)

type User struct {
	ID          int64  `json:"user_id"`
	Username    string `json:"user_name"`
	Password    string `json:"user_password"`
	Email       string `json:"user_email"`
	Realname    string `json:"user_realname"`
	Comment     string `json:"user_comment"`
	SystemAdmin int    `json:"user_system_admin"`
}

func (u User) ToMO() (m model.User) {
	err := utils.Adapt(u, &m)
	if err != nil {
		logs.Error("Failed to convert VM to MO: %+v", err)
		return
	}
	logs.Debug("Converted to model: %+v", m)
	return
}
