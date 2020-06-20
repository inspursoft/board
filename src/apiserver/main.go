package main

import (
	"bytes"

	c "git/inspursoft/board/src/apiserver/controllers/commons"
	v2routers "git/inspursoft/board/src/apiserver/routers"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/v1/controller"
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
	defaultAdminServer          = "adminserver"
	defaultAdminServerPort      = "8080"
	defaultKubeConfigPath       = "/root/kubeconfig"
	defaultSwaggerDoc           = "disabled"
	defaultAuthMode             = "db_auth"
	defaultMode                 = "normal"
	adminUserID                 = 1
	adminUsername               = "admin"
	adminEmail                  = "admin@inspur.com"
	defaultInitialPassword      = "123456a?"
	BaseRepoPath                = "/repos"
	sshKeyPath                  = "/keys"
	defaultProject              = "library"
	kvmToolsPath                = "/root/kvm"
	kvmRegistryPath             = "/root/kvmregistry"
)

var GogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var JenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

var apiServerPort = utils.GetConfig("API_SERVER_PORT", defaultAPIServerPort)
var swaggerDoc = utils.GetConfig("SWAGGER_DOC", defaultSwaggerDoc)
var devopsOpt = utils.GetConfig("DEVOPS_OPT")

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
	devops := service.CurrentDevOps()
	err := devops.SignUp(model.User{Username: adminUsername, Email: adminEmail, Password: initialPassword})
	if err != nil {
		logs.Error("Failed to create admin user on current DevOps: %+v", err)
	}

	token, err := devops.CreateAccessToken(adminUsername, initialPassword)
	if err != nil {
		logs.Error("Failed to create access token for admin user: %+v", err)
	}
	user := model.User{ID: adminUserID, RepoToken: token}
	service.UpdateUser(user, "repo_token")

	err = service.ConfigSSHAccess(adminUsername, token)
	if err != nil {
		logs.Error("Failed to config SSH access for admin user: %+v", err)
	}
	logs.Info("Initialize serve repo ...")
	logs.Info("Init git repo for default project %s", defaultProject)

	err = devops.CreateRepoAndJob(adminUserID, defaultProject)
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
	// err = service.SyncProjectsWithK8s()
	// if err != nil {
	// 	logs.Error("Failed to sync projects with K8s: %+v", err)
	// }
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

	utils.SetConfig("ADMINSERVER_IP", defaultAdminServer)
	utils.SetConfig("ADMINSERVER_PORT", defaultAdminServerPort)
	utils.SetConfig("ADMINSERVER_URL", "http://%s:%s/v1/admin", "ADMINSERVER_IP", "ADMINSERVER_PORT")

	utils.SetConfig("BASE_REPO_PATH", BaseRepoPath)
	utils.SetConfig("SSH_KEY_PATH", sshKeyPath)

	utils.SetConfig("KVM_TOOLS_PATH", kvmToolsPath)
	utils.SetConfig("KVM_REGISTRY_PATH", kvmRegistryPath)

	utils.SetConfig("KUBE_CONFIG_PATH", defaultKubeConfigPath)

	utils.SetConfig("AUTH_MODE", defaultAuthMode)

	dao.InitDB()

	c.InitController()
	controller.InitRouter()
	v2routers.InitRouterV2()
	err := v2routers.InitK8sRouter()
	if err != nil {
		logs.Error("Failed to init kubernetes api routes: %+v", err)
		panic(err)
	}
	if swaggerDoc() == "enabled" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
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
	service.SetSystemInfo("MODE", true)
	service.SetSystemInfo("BOARD_HOST_IP", true)
	service.SetSystemInfo("AUTH_MODE", false)
	service.SetSystemInfo("REDIRECTION_URL", false)
	service.SetSystemInfo("DEVOPS_OPT", false)

	if utils.GetStringValue("JENKINS_EXECUTION_MODE") != "single" {
		err = service.PrepareKVMHost()
		if err != nil {
			logs.Error("Failed to prepare KVM host: %+v", err)
			panic(err)
		}
	}

	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = "ILGOWezZZLeeDozS9Zg6xB2Ogyv1a2Ji"
	beego.BConfig.WebConfig.XSRFExpire = 1800

	beego.Run(":" + apiServerPort())
}
