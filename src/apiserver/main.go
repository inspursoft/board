package main

import (
	_ "git/inspursoft/board/src/apiserver/router"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

var adminUserID int64 = 1
var defaultInitialPassword = "123456a?"

func updateAdminPassword(initialPassword string) {
	if initialPassword == "" {
		initialPassword = defaultInitialPassword
	}
	salt := utils.GenerateRandomString()
	encryptedPassword := utils.Encrypt(initialPassword, salt)
	user := model.User{ID: adminUserID, Password: encryptedPassword, Salt: salt}
	isSuccess, err := service.UpdateUser(user, "password", "salt")
	if err != nil {
		logs.Error("Failed to update user password: %+v", err)
	}
	if isSuccess {
		logs.Info("Admin password has been updated successfully.")
	} else {
		logs.Info("Failed to update admin initial password.")
	}
}

func main() {
	initialPassword := os.Getenv("BOARD_ADMIN_PASSWORD")
	updateAdminPassword(initialPassword)
	beego.Run(":8088")
}
