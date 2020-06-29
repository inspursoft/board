package main

import (
	"bytes"
	"context"
	"sync"
	"time"

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

type systemInfoKey int

const (
	systemInfo                  systemInfoKey = iota
	defaultDBServer                           = "db"
	defaultDBPort                             = "3306"
	defaultAPIServerPort                      = "8088"
	defaultTokenServer                        = "tokenserver"
	defaultTokenServerPort                    = "4000"
	defaultTokenCacheExpireTime               = "1800"
	defaultAdminServer                        = "adminserver"
	defaultAdminServerPort                    = "8080"
	defaultKubeConfigPath                     = "/root/kubeconfig"
	defaultSwaggerDoc                         = "disabled"
	defaultAuthMode                           = "db_auth"
	defaultMode                               = "normal"
	adminUserID                               = 1
	adminUsername                             = "admin"
	adminEmail                                = "admin@inspur.com"
	defaultInitialPassword                    = "123456a?"
	BaseRepoPath                              = "/repos"
	sshKeyPath                                = "/keys"
	defaultProject                            = "library"
	kvmToolsPath                              = "/root/kvm"
	kvmRegistryPath                           = "/root/kvmregistry"
)

var GogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var JenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

var apiServerPort = utils.GetConfig("API_SERVER_PORT", defaultAPIServerPort)
var swaggerDoc = utils.GetConfig("SWAGGER_DOC", defaultSwaggerDoc)
var devopsOpt = utils.GetConfig("DEVOPS_OPT")
var jenkinsExecutionMode = utils.GetConfig("JENKINS_EXECUTION_MODE")
var k8sForceInitSync = utils.GetConfig("FORCE_INIT_SYNC")
var initTimeExpiry = time.Minute * 10

func setConfigurations() {

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
}

func initBoardVersion(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		logs.Error("Error occurred from context: %+v", err)
		return
	}
	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		logs.Error("Failed to read VERSION file: %+v", err)
	}
	utils.SetConfig("BOARD_VERSION", string(bytes.TrimSpace(version)))
	service.SetSystemInfo("BOARD_VERSION", true)
}

func updateAdminPassword(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		logs.Error("Error occurred from context: %+v", err)
		return
	}
	if ctx.Value(systemInfo).(*model.SystemInfo).SetAdminPassword == "updated" {
		logs.Info("Skip updating admin user as it has been updated.")
		return
	}
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

func initProjectRepo(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		logs.Error("Error occurred from context: %+v", err)
		return
	}
	if ctx.Value(systemInfo).(*model.SystemInfo).InitProjectRepo == "created" {
		logs.Info("Skip initializing project repo as it has been created.")
		return
	}

	initialPassword := utils.GetStringValue("BOARD_ADMIN_PASSWORD")
	if initialPassword == "" {
		initialPassword = defaultInitialPassword
	}
	service.SetSystemInfo("DEVOPS_OPT", false)
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
	service.SetSystemInfo("INIT_PROJECT_REPO", false)
	logs.Info("Finished to create initial project and repo.")
}

func prepareKVMHost(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		logs.Error("Error occurred at: %s from context: %+v", "prepare KVM host", err)
		return
	}
	if jenkinsExecutionMode() == "single" {
		logs.Info("Skip preparing KVM host as it set as single slave node.")
		return
	}
	if err := service.PrepareKVMHost(); err != nil {
		logs.Error("Failed to prepare KVM host, error: %+v", err)
	}
}

func initKubernetesInfo(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		logs.Error("Error occurred from context: %+v", err)
		return
	}
	logs.Info("Initialize Kubernetes info")
	info, err := service.GetKubernetesInfo()
	if err != nil {
		logs.Error("Failed to initialize kubernetes info, err: %+v", err)
		utils.SetConfig("KUBERNETES_VERSION", "NA")
	} else {
		utils.SetConfig("KUBERNETES_VERSION", info.GitVersion)
	}
	service.SetSystemInfo("KUBERNETES_VERSION", true)
	logs.Info("Finished to initialize Kubernetes info.")
}

func syncUpWithK8s(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		logs.Error("Error occurred from context: %+v", err)
		return
	}
	if ctx.Value(systemInfo).(*model.SystemInfo).SyncK8s == "created" {
		logs.Info("Skip initializing project repo as it has been created.")
		return
	}
	if k8sForceInitSync() == "false" {
		logs.Info("Skip sync up with K8s forcely. ")
		return
	}
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

func main() {
	setConfigurations()
	c.InitController()
	controller.InitRouter()
	v2routers.InitRouterV2()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.SetConfig("INIT_STATUS", "NOT_READY")
		dao.InitDB()
		service.SetSystemInfo("DNS_SUFFIX", true)
		service.SetSystemInfo("MODE", true)
		service.SetSystemInfo("BOARD_HOST_IP", true)
		service.SetSystemInfo("AUTH_MODE", false)
		service.SetSystemInfo("REDIRECTION_URL", false)
		service.SetSystemInfo("DEVOPS_OPT", false)

		info, err := service.GetSystemInfo()
		if err != nil {
			logs.Error("Failed to set system config: %+v", err)
			panic(err)
		}
		ctx := context.WithValue(context.Background(), systemInfo, info)
		initBoardVersion(ctx)
		utils.SetConfig("INIT_STATUS", "UPDATE_ADMIN_PASSWORD")
		updateAdminPassword(ctx)
		utils.SetConfig("INIT_STATUS", "INIT_PROJECT_REPO")
		initProjectRepo(ctx)
		utils.SetConfig("INIT_STATUS", "PREPARE_KVM_HOST")
		prepareKVMHost(ctx)
		utils.SetConfig("INIT_STATUS", "INIT_KUBERNETES_INFO")
		initKubernetesInfo(ctx)
		utils.SetConfig("INIT_STATUS", "SYNC_UP_K8S")
		syncUpWithK8s(ctx)
		utils.SetConfig("INIT_STATUS", "READY")
	}()

	if swaggerDoc() == "enabled" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = "ILGOWezZZLeeDozS9Zg6xB2Ogyv1a2Ji"
	beego.BConfig.WebConfig.XSRFExpire = 1800

	beego.Run(":" + apiServerPort())
}
