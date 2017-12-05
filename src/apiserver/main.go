package main

import (
	"git/inspursoft/board/src/apiserver/controller"
	_ "git/inspursoft/board/src/apiserver/router"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"

	"github.com/astaxie/beego/logs"

	"path/filepath"

	"github.com/astaxie/beego"
)

const (
	adminUserID            = 1
	defaultInitialPassword = "123456a?"
	baseRepoPath           = "/repos"
	sshKeyPath             = "/root/.ssh/id_rsa"
	defaultProject         = "library"
)

var repoServePath = filepath.Join(baseRepoPath, "board_repo_serve")
var repoServeURL = filepath.Join("root@gitserver:", "gitserver", "repos", "board_repo_serve")
var repoPath = filepath.Join(baseRepoPath, "board_repo")

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
		err = service.SetSystemInfo("SET_ADMIN_PASSWORD")
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
	logs.Info("Initialize serve repo\n")
	_, err := service.InitBareRepo(repoServePath)
	if err != nil {
		logs.Error("Failed to initialize serve repo: %+v\n", err)
	}

	logs.Info("Init git repo for default project %s", defaultProject)
	_, err = service.InitRepo(repoServeURL, repoPath)
	if err != nil {
		logs.Error("Failed to initialize default user's repo: %+v\n", err)
		return
	}

	os.MkdirAll(filepath.Join(repoPath, defaultProject), 0755)
	if err != nil {
		logs.Error("Failed to make default user's repo: %+v\n", err)
	}

	utils.SetConfig("INIT_PROJECT_REPO", "created")
	err = service.SetSystemInfo("INIT_PROJECT_REPO")
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}
}

func syncServiceWithK8s() {
	service.SyncServiceWithK8s()
}

func main() {

	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(4)

	utils.Initialize()

	utils.AddEnv("BOARD_ADMIN_PASSWORD")
	utils.AddEnv("KUBE_MASTER_IP")
	utils.AddEnv("KUBE_MASTER_PORT")
	utils.AddEnv("REGISTRY_IP")
	utils.AddEnv("REGISTRY_PORT")

	utils.AddEnv("AUTH_MODE")

	utils.AddEnv("LDAP_URL")
	utils.AddEnv("LDAP_SEARCH_DN")
	utils.AddEnv("LDAP_SEARCH_PWD")
	utils.AddEnv("LDAP_BASE_DN")
	utils.AddEnv("LDAP_FILTER")
	utils.AddEnv("LDAP_UID")
	utils.AddEnv("LDAP_SCOPE")
	utils.AddEnv("LDAP_TIMEOUT")

	utils.SetConfig("REGISTRY_URL", "http://%s:%s", "REGISTRY_IP", "REGISTRY_PORT")
	utils.SetConfig("KUBE_MASTER_URL", "http://%s:%s", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")
	utils.SetConfig("KUBE_NODE_URL", "http://%s:%s/api/v1/nodes", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")

	utils.SetConfig("BASE_REPO_PATH", baseRepoPath)
	utils.SetConfig("REPO_SERVE_URL", repoServeURL)
	utils.SetConfig("REPO_SERVE_PATH", repoServePath)
	utils.SetConfig("REPO_PATH", repoPath)
	utils.SetConfig("SSH_KEY_PATH", sshKeyPath)

	utils.SetConfig("REGISTRY_BASE_URI", "%s:%s", "REGISTRY_IP", "REGISTRY_PORT")

	dao.InitDB()
	err := service.SetSystemInfo("AUTH_MODE")
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}

	utils.AddValue("IS_EXTERNAL_AUTH", (utils.GetStringValue("AUTH_MODE") != "db_auth"))
	utils.ShowAllConfigs()

	controller.InitController()

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
	syncServiceWithK8s()
	beego.Run(":8088")
}
