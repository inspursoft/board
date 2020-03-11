package main

import (
	"bytes"
	"git/inspursoft/board/src/apiserver/controller"
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
	defaultKubeConfigPath       = "/root/kubeconfig"
	adminUserID                 = 1
	adminUsername               = "admin"
	adminEmail                  = "admin@inspur.com"
	defaultInitialPassword      = "123456a?"
	baseRepoPath                = "/repos"
	sshKeyPath                  = "/keys"
	defaultProject              = "library"
	kvmToolsPath                = "/root/kvm"
	kvmRegistryPath             = "/root/kvmregistry"
)

var gogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var jenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

func initBoardVersion() {
	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		logs.Error("Failed to read VERSION file: %+v", err)
	}
	utils.SetConfig("BOARD_VERSION", string(bytes.TrimSpace(version)))
	service.SetSystemInfo("BOARD_VERSION", true)
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
		service.SetSystemInfo("SET_ADMIN_PASSWORD", false)
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
	}

	utils.SetConfig("INIT_PROJECT_REPO", "created")
	service.SetSystemInfo("INIT_PROJECT_REPO", true)
	logs.Info("Finished to create initial project and repo.")
}

func syncUpWithK8s() {
	logs.Info("Initialize to sync up with K8s status ...")
	defer func() {
		utils.SetConfig("SYNC_K8S", "finished")
		service.SetSystemInfo("SYNC_K8S", false)
	}()
	var err error
	// Sync namespace with specific project ownerID
	err = service.SyncNamespaceByOwnerID(adminUserID)
	if err != nil {
		logs.Error("Failed to sync namespace by userID: %d, err: %+v", adminUserID, err)
	}
	logs.Info("Successful sync up with namespaces for admin user.")
	// Sync projects from cluster namespaces
	err = service.SyncProjectsWithK8s()
	if err != nil {
		logs.Error("Failed to sync projects with K8s: %+v", err)
	}
	logs.Info("Successful sync up with projects with K8s.")
}

func initKubernetesInfo() {
	logs.Info("Initialize kubernetes info")
	info, err := service.GetKubernetesInfo()
	if err != nil {
		logs.Error("Failed to initialize kubernetes info, err: %+v", err)
		utils.SetConfig("KUBERNETES_VERSION", "NA")
	} else {
		utils.SetConfig("KUBERNETES_VERSION", info.GitVersion)
	}
	service.SetSystemInfo("KUBERNETES_VERSION", true)
	logs.Info("Finished to initialize kubernetes info.")
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

	utils.SetConfig("KVM_TOOLS_PATH", kvmToolsPath)
	utils.SetConfig("KVM_REGISTRY_PATH", kvmRegistryPath)

	utils.SetConfig("KUBE_CONFIG_PATH", defaultKubeConfigPath)

	dao.InitDB()

	controller.InitController()
	controller.InitRouter()

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

	if systemInfo.KubernetesVersion == "" || systemInfo.KubernetesVersion == "NA" {
		initKubernetesInfo()
	}

	if systemInfo.SyncK8s == "" || utils.GetStringValue("FORCE_INIT_SYNC") == "true" {
		syncUpWithK8s()
	}

	service.SetSystemInfo("DNS_SUFFIX", true)
	service.SetSystemInfo("BOARD_HOST_IP", true)
	service.SetSystemInfo("AUTH_MODE", false)
	service.SetSystemInfo("REDIRECTION_URL", false)

	if utils.GetStringValue("JENKINS_EXECUTION_MODE") != "single" {
		err = service.PrepareKVMHost()
		if err != nil {
			logs.Error("Failed to prepare KVM host: %+v", err)
			panic(err)
		}
	}

	beego.Run(":" + defaultAPIServerPort)
}
