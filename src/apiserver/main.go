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
	defaultAPIServerPort   = "8088"
	adminUserID            = 1
	adminUsername          = "admin"
	adminEmail             = "admin@inspur.com"
	defaultInitialPassword = "123456a?"
	baseRepoPath           = "/repos"
	sshKeyPath             = "/keys"
	defaultProject         = "library"
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

/*
func syncServiceWithK8s() {
	service.SyncServiceWithK8s()
	utils.SetConfig("SYNC_K8S", "created")
	err := service.SetSystemInfo("SYNC_K8S", true)
	if err != nil {
		logs.Error("Failed to set system config: %+v", err)
		panic(err)
	}
}
*/

func updateSystemInfo() {
	var err error
	err = service.SetSystemInfo("BOARD_HOST", true)
	if err != nil {
		logs.Error("Failed to set system config BOARD_HOST: %+v", err)
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

	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(4)

	utils.Initialize()

	utils.AddEnv("BOARD_HOST")
	utils.AddValue("BOARD_API_PORT", defaultAPIServerPort)
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
	utils.AddEnv("FORCE_INIT_SYNC")
	utils.AddEnv("VERIFICATION_URL")
	utils.AddEnv("REDIRECTION_URL")

	utils.SetConfig("BOARD_API_BASE_URL", "http://%s:%s/api/v1", "BOARD_HOST", "BOARD_API_PORT")

	utils.AddEnv("GOGITS_HOST_IP")
	utils.AddEnv("GOGITS_HOST_PORT")
	utils.SetConfig("GOGITS_BASE_URL", "http://%s:%s", "GOGITS_HOST_IP", "GOGITS_HOST_PORT")

	utils.AddEnv("GOGITS_SSH_PORT")
	utils.SetConfig("GOGITS_SSH_URL", "ssh://git@%s:%s", "GOGITS_HOST_IP", "GOGITS_SSH_PORT")

	utils.AddEnv("JENKINS_HOST_IP")
	utils.AddEnv("JENKINS_HOST_PORT")
	utils.SetConfig("JENKINS_BASE_URL", "http://%s:%s", "JENKINS_HOST_IP", "JENKINS_HOST_PORT")

	utils.SetConfig("REGISTRY_URL", "http://%s:%s", "REGISTRY_IP", "REGISTRY_PORT")
	utils.SetConfig("KUBE_MASTER_URL", "http://%s:%s", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")
	utils.SetConfig("KUBE_NODE_URL", "http://%s:%s/api/v1/nodes", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")

	utils.SetConfig("BASE_REPO_PATH", baseRepoPath)
	utils.SetConfig("SSH_KEY_PATH", sshKeyPath)
	utils.SetConfig("API_SERVER_PORT", defaultAPIServerPort)

	utils.SetConfig("API_SERVER_URL", "http://%s:%s", "BOARD_HOST", "API_SERVER_PORT")

	utils.SetConfig("REGISTRY_BASE_URI", "%s:%s", "REGISTRY_IP", "REGISTRY_PORT")

	dao.InitDB()

	utils.AddValue("IS_EXTERNAL_AUTH", (utils.GetStringValue("AUTH_MODE") != "db_auth"))

	utils.ShowAllConfigs()

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
		//already do sync service in sync project
		//syncServiceWithK8s()
	}

	beego.Run(":" + defaultAPIServerPort)
}
