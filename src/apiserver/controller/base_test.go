package controller_test

import (
	"git/inspursoft/board/src/apiserver/controller"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/gogs"
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
	adminUsername          = "admin"
	adminEmail             = "admin@inspur.com"
	defaultProject         = "library"
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

func initProjectRepo() {
	initialPassword := utils.GetStringValue("BOARD_ADMIN_PASSWORD")
	if initialPassword == "" {
		initialPassword = defaultInitialPassword
	}

	err := gogs.SignUp(model.User{Username: adminUsername, Email: adminEmail, Password: initialPassword})
	if err != nil {
		logs.Error("Failed to create admin user on Gogit: %+v", err)
	}

	token, err := gogs.CreateAccessToken(adminUsername, initialPassword)
	if err != nil {
		logs.Error("Failed to create access token for admin user: %+v", err)
	}
	user := model.User{ID: adminUserID, RepoToken: token.Sha1}
	service.UpdateUser(user, "repo_token")

	err = service.ConfigSSHAccess(adminUsername, token.Sha1)
	if err != nil {
		logs.Error("Failed to config SSH access for admin user: %+v", err)
	}
	logs.Info("Initialize serve repo ...")
	logs.Info("Init git repo for default project %s", defaultProject)

	// err = service.CreateRepoAndJob(adminUserID, defaultProject)
	// if err != nil {
	// 	logs.Error("Failed to create default repo %s: %+v", defaultProject, err)
	// }

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

	controller.InitController()
	os.Exit(m.Run())
}
