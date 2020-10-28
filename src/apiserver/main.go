package main

import (
	"bytes"
	"context"
	"fmt"
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
	adminUsername                             = "boardadmin"
	adminEmail                                = "boardadmin@inspur.com"
	defaultInitialPassword                    = "123456a?"
	BaseRepoPath                              = "/repos"
	sshKeyPath                                = "/keys"
	defaultProject                            = "library"
)

var GogitsSSHURL = utils.GetConfig("GOGITS_SSH_URL")
var JenkinsBaseURL = utils.GetConfig("JENKINS_BASE_URL")

var apiServerPort = utils.GetConfig("API_SERVER_PORT", defaultAPIServerPort)
var swaggerDoc = utils.GetConfig("SWAGGER_DOC", defaultSwaggerDoc)
var devopsOpt = utils.GetConfig("DEVOPS_OPT")
var jenkinsExecutionMode = utils.GetConfig("JENKINS_EXECUTION_MODE")
var k8sForceInitSync = utils.GetConfig("FORCE_INIT_SYNC")

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

	utils.SetConfig("KUBE_CONFIG_PATH", defaultKubeConfigPath)

	utils.SetConfig("AUTH_MODE", defaultAuthMode)
}

func initBoardVersion(ctx context.Context, cancel context.CancelFunc, e chan error) {
	if err := detectContextStatus(ctx); err != nil {
		e <- err
		return
	}
	utils.SetConfig("INIT_STATUS", "5_1_INIT_BOARD_VERSION")
	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		logs.Error("Failed to read VERSION file: %+v", err)
		e <- fmt.Errorf("Failed to read VERSION file: %+v", err)
		cancel()
	}
	utils.SetConfig("BOARD_VERSION", string(bytes.TrimSpace(version)))
	service.SetSystemInfo("BOARD_VERSION", true)
}

func updateAdminPassword(ctx context.Context, cancel context.CancelFunc, e chan error) {
	if err := detectContextStatus(ctx); err != nil {
		e <- err
		return
	}
	utils.SetConfig("INIT_STATUS", "5_2_UPDATE_ADMIN_PASSWORD")
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
		e <- fmt.Errorf("Failed to update user password: %+v", err)
		cancel()
	}
	if isSuccess {
		utils.SetConfig("SET_ADMIN_PASSWORD", "updated")
		service.SetSystemInfo("SET_ADMIN_PASSWORD", false)
		logs.Info("Admin password has been updated successfully.")
	} else {
		logs.Info("Failed to update admin initial password.")
	}
}

func initProjectRepo(ctx context.Context, cancel context.CancelFunc, e chan error) {
	if err := detectContextStatus(ctx); err != nil {
		e <- err
		return
	}
	utils.SetConfig("INIT_STATUS", "5_3_INIT_PROJECT_REPO")
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
		e <- fmt.Errorf("Failed to create admin user on current DevOps: %+v", err)
		cancel()
	}
	token, err := devops.CreateAccessToken(adminUsername, initialPassword)
	if err != nil {
		logs.Error("Failed to create access token for admin user: %+v", err)
		e <- fmt.Errorf("Failed to create access token for admin user: %+v", err)
		cancel()
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
		e <- fmt.Errorf("Failed to create default repo %s: %+v", defaultProject, err)
		cancel()
	}
	utils.SetConfig("INIT_PROJECT_REPO", "created")
	service.SetSystemInfo("INIT_PROJECT_REPO", false)
	logs.Info("Finished to create initial project and repo.")
}

func initKubernetesInfo(ctx context.Context, cancel context.CancelFunc, e chan error) {
	if err := detectContextStatus(ctx); err != nil {
		e <- err
		return
	}
	utils.SetConfig("INIT_STATUS", "5_4_INIT_KUBERNETES_INFO")
	logs.Info("Initializing Kubernetes info")
	info, err := service.GetKubernetesInfo()
	if err != nil {
		logs.Error("Failed to initialize kubernetes info, err: %+v", err)
		utils.SetConfig("KUBERNETES_VERSION", "NA")
		e <- fmt.Errorf("Failed to initialize kubernetes info, err: %+v", err)
		cancel()
	}
	utils.SetConfig("KUBERNETES_VERSION", info.GitVersion)
	service.SetSystemInfo("KUBERNETES_VERSION", true)
	logs.Info("Finished to initialize Kubernetes info.")
}

func syncUpWithK8s(ctx context.Context, cancel context.CancelFunc, e chan error) {
	if err := detectContextStatus(ctx); err != nil {
		e <- err
		return
	}
	utils.SetConfig("INIT_STATUS", "5_5_SYNC_UP_K8S")
	if ctx.Value(systemInfo).(*model.SystemInfo).SyncK8s == "created" {
		logs.Info("Skip initializing project repo as it has been created.")
		if k8sForceInitSync() == "false" {
			logs.Info("Skip sync up with K8s forcely. ")
			return
		}
	}
	logs.Info("Initialize to sync up with K8s status ...")
	defer func() {
		utils.SetConfig("SYNC_K8S", "finished")
		service.SetSystemInfo("SYNC_K8S", false)
	}()
	// Sync namespace with specific project ownerID
	err := service.SyncNamespaceByOwnerID(adminUserID)
	if err != nil {
		logs.Error("Failed to sync namespace by userID: %d, err: %+v", adminUserID, err)
		e <- fmt.Errorf("Failed to sync namespace by userID: %d, err: %+v", adminUserID, err)
		cancel()
	}
	logs.Info("Successful sync up with namespaces for admin user.")
	// Sync projects from cluster namespaces
	err = service.SyncProjectsWithK8s()
	if err != nil {
		logs.Error("Failed to sync projects with K8s: %+v", err)
		e <- fmt.Errorf("Failed to sync projects with K8s: %+v", err)
		cancel()
	}
	logs.Info("Successful sync up with projects with K8s.")
}

func detectContextStatus(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Millisecond * 500):
		return nil
	}
}

func main() {
	setConfigurations()
	c.InitController()
	controller.InitRouter()
	v2routers.InitRouterV2()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	initializer := func(done chan bool) chan error {
		e := make(chan error)
		go func() {
			defer close(e)

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
				e <- fmt.Errorf("Failed to set system config: %+v", err)
				cancel()
			}
			utils.SetConfig("INIT_STATUS", "NOT_READY")
			ctx = context.WithValue(ctx, systemInfo, info)
			initBoardVersion(ctx, cancel, e)
			updateAdminPassword(ctx, cancel, e)
			initProjectRepo(ctx, cancel, e)
			initKubernetesInfo(ctx, cancel, e)
			syncUpWithK8s(ctx, cancel, e)
			close(done)

			select {
			case <-done:
				return
			case <-e:
				return
			}
		}()
		return e
	}
	done := make(chan bool)
	go func() {
		select {
		case err := <-initializer(done):
			if err != nil {
				utils.SetConfig("INIT_STATUS", err.Error())
				logs.Error("Failed to execute initialization with error: %+v", err)
			} else {
				utils.SetConfig("INIT_STATUS", "READY")
				logs.Info("Finished executing Board initialization process.")
			}
		}
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
