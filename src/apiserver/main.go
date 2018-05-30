package main

import (
	"bytes"
	"git/inspursoft/board/src/apiserver/controller"
	_ "git/inspursoft/board/src/apiserver/router"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/gogs"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"io/ioutil"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

const (
	defaultDBServer             = "db"
	defaultDBPort               = "3306"
	defaultAPIServerPort        = "8088"
	defaultTokenServer          = "tokenserver"
	defaultTokenServerPort      = "4000"
	defaultTokenCacheExpireTime = "1800"
	adminUserID                 = 1
	adminUsername               = "admin"
	adminEmail                  = "admin@inspur.com"
	defaultInitialPassword      = "123456a?"
	baseRepoPath                = "/repos"
	sshKeyPath                  = "/keys"
	defaultProject              = "library"
)

var gogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

func initBoardVersion() {
	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		logs.Error("Failed to read VERSION file: %+v", err)
		panic(err)
	}
	utils.SetConfig("BOARD_VERSION", string(bytes.TrimSpace(version)))
	err = service.SetSystemInfo("BOARD_VERSION", true)
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}
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

	err = service.CreateRepoAndJob(adminUserID, defaultProject)
	if err != nil {
		logs.Error("Failed to create default repo %s: %+v", defaultProject, err)
		panic(err)
	}

	utils.SetConfig("INIT_PROJECT_REPO", "created")
	err = service.SetSystemInfo("INIT_PROJECT_REPO", true)
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}
}

func initDefaultProjects() {
	logs.Info("Initialize default projects...")
	var err error
	// Sync namespace with specific project ownerID
	err = service.SyncNamespaceByOwnerID(adminUserID)
	if err != nil {
		logs.Error("Failed to sync namespace by userID: %d, err: %+v", adminUserID, err)
		panic(err)
	}
	logs.Info("Successful synchonized namespace for admin user.")
	// Sync projects from cluster namespaces
	err = service.SyncProjectsWithK8s()
	if err != nil {
		logs.Error("Failed to sync projects with K8s: %+v", err)
		panic(err)
	}
	logs.Info("Successful synchonized projects with Kubernetes.")
}

func syncServiceWithK8s() {
	service.SyncServiceWithK8s()
	utils.SetConfig("SYNC_K8S", "created")
	err := service.SetSystemInfo("SYNC_K8S", true)
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}
}
func updateSystemInfo() {
	var err error
	err = service.SetSystemInfo("BOARD_HOST_IP", true)
	if err != nil {
		logs.Error("Failed to set system config BOARD_HOST_IP: %+v", err)
		panic(err)
	}
	err = service.SetSystemInfo("AUTH_MODE", false)
	if err != nil {
		logs.Error("Failed to set system config AUTH_MODE: %+v", err)
		panic(err)
	}
	err = service.SetSystemInfo("REDIRECTION_URL", false)
	if err != nil {
		logs.Error("Failed to set system config REDIRECTION_URL: %+v", err)
	}
}

func main() {

	utils.InitializeDefaultConfig()

	utils.SetConfig("DB_IP", defaultDBServer)
	utils.SetConfig("DB_PORT", defaultDBPort)

	utils.SetConfig("TOKEN_SERVER_IP", defaultTokenServer)
	utils.SetConfig("TOKEN_SERVER_PORT", defaultTokenServerPort)
	utils.SetConfig("TOKEN_SERVER_URL", "http://%s:%s/tokenservice/token", "TOKEN_SERVER_IP", "TOKEN_SERVER_PORT")

	utils.SetConfig("BASE_REPO_PATH", baseRepoPath)
	utils.SetConfig("SSH_KEY_PATH", sshKeyPath)

	dao.InitDB()

	controller.InitController()
	updateSystemInfo()

	systemInfo, err := service.GetSystemInfo()
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}

	initBoardVersion()

	if systemInfo.SetAdminPassword == "" {
		updateAdminPassword()
	}

	if systemInfo.InitProjectRepo == "" {
		initProjectRepo()
	}

	if systemInfo.SyncK8s == "" || utils.GetStringValue("FORCE_INIT_SYNC") == "true" {
		initDefaultProjects()
		syncServiceWithK8s()
	}

	beego.Run(":" + defaultAPIServerPort)
}
