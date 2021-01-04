package controller_test

import (
	"github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/v1/controller"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/model/dashboard"
	"github.com/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

const (
	adminUserID            = 1
	adminUsername          = "boardadmin"
	adminEmail             = "boardadmin@inspur.com"
	defaultProject         = "library"
	defaultInitialPassword = "123456a?"
)

var gitlabAccessToken = utils.GetConfig("GITLAB_ADMIN_TOKEN")

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

func initProjectRepo() {
	initialPassword := utils.GetStringValue("BOARD_ADMIN_PASSWORD")
	if initialPassword == "" {
		initialPassword = defaultInitialPassword
	}
	user := model.User{ID: adminUserID, RepoToken: gitlabAccessToken()}
	service.UpdateUser(user, "repo_token")

	logs.Info("Initialize serve repo ...")
	logs.Info("Init git repo for default project %s", defaultProject)

	utils.SetConfig("INIT_PROJECT_REPO", "created")
	service.SetSystemInfo("INIT_PROJECT_REPO", true)
	logs.Info("Finished to create initial project and repo.")
}

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	utils.AddEnv("NODE_IP")
	utils.AddEnv("TOKEN_SERVER_IP")
	utils.AddEnv("TOKEN_SERVER_PORT")
	utils.SetConfig("TOKEN_SERVER_URL", "http://%s:%s/tokenservice/token", "TOKEN_SERVER_IP", "TOKEN_SERVER_PORT")
	utils.SetConfig("SSH_KEY_PATH", "/tmp/ssh-keys")
	utils.SetConfig("AUDIT_DEBUG", "false")
	utils.SetConfig("INIT_STATUS", "READY")
	utils.AddEnv("GITLAB_ADMIN_TOKEN")

	utils.ShowAllConfigs()
	dao.InitDB()
	systemInfo, err := service.GetSystemInfo()
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}
	if systemInfo.SetAdminPassword == "" {
		updateAdminPassword()
	}

	if systemInfo.InitProjectRepo == "" {
		initProjectRepo()
	}
	commons.InitController()
	os.Exit(m.Run())
}
