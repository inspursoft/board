package controller_test

import (
	"git/inspursoft/board/src/apiserver/controller"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/model/dashboard"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

const (
	adminUserID            = 1
	defaultInitialPassword = "123456a?"
)

func init() {
	controller.InitRouter()
	orm.RegisterModel(new(dashboard.NodeDashboardMinute), new(dashboard.NodeDashboardHour), new(dashboard.NodeDashboardDay))
}

func updateAdminPassword() {
	initialPassword := utils.GetStringValue("BOARD_ADMIN_PASSWORD")
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
		utils.SetConfig("SET_ADMIN_PASSWORD", "updated")
		err = service.SetSystemInfo("SET_ADMIN_PASSWORD", false)
		if err != nil {
			logs.Error("Failed to set system config: %+v", err)
			panic(err)
		}
		logs.Info("Admin password has been updated successfully.")
	} else {
		logs.Info("Failed to update admin initial password.")
	}
}

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	utils.AddEnv("NODE_IP")
	utils.AddEnv("TOKEN_SERVER_IP")
	utils.AddEnv("TOKEN_SERVER_PORT")
	utils.SetConfig("TOKEN_SERVER_URL", "http://%s:%s/tokenservice/token", "TOKEN_SERVER_IP", "TOKEN_SERVER_PORT")
	utils.SetConfig("SSH_KEY_PATH", "/tmp/ssh-keys")

	dao.InitDB()
	updateAdminPassword()
	controller.InitController()
	os.Exit(m.Run())
}
