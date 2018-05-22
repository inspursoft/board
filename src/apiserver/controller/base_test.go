package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/model/dashboard"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

const (
	adminUserID            = 1
	defaultInitialPassword = "123456a?"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/sign-in",
				&AuthController{},
				"post:SignInAction"),
			beego.NSRouter("/ext-auth",
				&AuthController{},
				"get:ExternalAuthAction"),
			beego.NSRouter("/sign-up",
				&AuthController{},
				"post:SignUpAction"),
			beego.NSRouter("/log-out",
				&AuthController{},
				"get:LogOutAction"),
			beego.NSRouter("/user-exists",
				&AuthController{},
				"get:UserExists"),
			beego.NSRouter("/users/current",
				&AuthController{},
				"get:CurrentUserAction"),
			beego.NSRouter("/systeminfo",
				&AuthController{},
				"get:GetSystemInfo"),
			beego.NSRouter("/users",
				&UserController{},
				"get:GetUsersAction"),
			beego.NSRouter("/users/:id([0-9]+)/password",
				&UserController{},
				"put:ChangePasswordAction"),
			beego.NSRouter("/users/changeaccount",
				&UserController{},
				"put:ChangeUserAccount"),
			beego.NSRouter("/adduser",
				&SystemAdminController{},
				"post:AddUserAction"),
			beego.NSRouter("/users/:id([0-9]+)",
				&SystemAdminController{},
				"get:GetUserAction;put:UpdateUserAction;delete:DeleteUserAction"),
			beego.NSRouter("/users/:id([0-9]+)/systemadmin",
				&SystemAdminController{},
				"put:ToggleSystemAdminAction"),
			beego.NSRouter("/projects",
				&ProjectController{},
				"head:ProjectExists;get:GetProjectsAction;post:CreateProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/publicity",
				&ProjectController{},
				"put:ToggleProjectPublicAction"),
			beego.NSRouter("/projects/:id([0-9]+)",
				&ProjectController{},
				"get:GetProjectAction;delete:DeleteProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/members",
				&ProjectMemberController{},
				"get:GetProjectMembersAction;post:AddOrUpdateProjectMemberAction"),
			beego.NSRouter("/projects/:projectId([0-9]+)/members/:userId([0-9]+)",
				&ProjectMemberController{},
				"delete:DeleteProjectMemberAction"),
			beego.NSRouter("/images",
				&ImageController{},
				"get:GetImagesAction;delete:DeleteImageAction"),
			beego.NSRouter("/images/:imagename(.*)",
				&ImageController{},
				"get:GetImageDetailAction;delete:DeleteImageTagAction"),
			beego.NSRouter("/images/reset-temp",
				&ImageController{},
				"put:ResetBuildImageTempAction"),
			beego.NSRouter("/images/building",
				&ImageController{},
				"post:BuildImageAction"),
			beego.NSRouter("/images/dockerfilebuilding",
				&ImageController{},
				"post:DockerfileBuildImageAction"),
			beego.NSRouter("/images/dockerfile",
				&ImageController{},
				"get:GetImageDockerfileAction"),
			beego.NSRouter("/images/registry",
				&ImageController{},
				"get:GetImageRegistryAction"),
			beego.NSRouter("/images/preview",
				&ImageController{},
				"post:DockerfilePreviewAction"),
			beego.NSRouter("/images/configclean",
				&ImageController{},
				"delete:ConfigCleanAction"),
			beego.NSRouter("/images/:imagename(.*)/existing",
				&ImageController{},
				"get:CheckImageTagExistingAction"),
			beego.NSRouter("/search",
				&SearchSourceController{}, "get:Search"),
			beego.NSRouter("/node",
				&NodeController{}, "get:GetNode"),
			beego.NSRouter("/nodes",
				&NodeController{}, "get:NodeList"),
			beego.NSRouter("/node/toggle",
				&NodeController{}, "get:NodeToggle"),
			beego.NSRouter("/node/:id([0-9]+)/group",
				&NodeController{},
				"get:GetGroupsOfNodeAction;post:AddNodeToGroupAction;delete:RemoveNodeFromGroupAction"),
			beego.NSRouter("/nodegroup",
				&NodeGroupController{},
				"get:GetNodeGroupsAction;post:AddNodeGroupAction;delete:DeleteNodeGroupAction"),
			beego.NSRouter("/nodegroup/existing",
				&NodeGroupController{},
				"get:CheckNodeGroupNameExistingAction"),
			beego.NSNamespace("/storage",
				beego.NSRouter("/setnfs", &StorageController{}, "post:Storage"),
			),
			beego.NSNamespace("/dashboard", beego.NSRouter("/service",
				&DashboardServiceController{},
				"post:GetServiceData"),
				beego.NSRouter("/node",
					&DashboardNodeController{}, "post:GetNodeData"),
				beego.NSRouter("/data",
					&Dashboard{}, "post:GetData"),
				beego.NSRouter("/time",
					&DashboardServiceController{}, "get:GetServerTime"),
			),
			beego.NSRouter("/services",
				&ServiceController{},
				"post:CreateServiceConfigAction;get:GetServiceListAction"),
			beego.NSRouter("/services/exists",
				&ServiceController{},
				"get:ServiceExists"),
			beego.NSRouter("/services/rollingupdate/image",
				&ServiceRollingUpdateController{},
				"get:GetRollingUpdateServiceImageConfigAction;patch:PatchRollingUpdateServiceImageAction"),
			beego.NSRouter("/services/rollingupdate/nodegroup",
				&ServiceRollingUpdateController{},
				"get:GetRollingUpdateServiceNodeGroupConfigAction;patch:PatchRollingUpdateServiceNodeGroupAction"),
			beego.NSRouter("/services/:id([0-9]+)",
				&ServiceController{},
				"delete:DeleteServiceAction"),
			beego.NSRouter("/services/deployment",
				&ServiceController{},
				"post:DeployServiceAction"),
			beego.NSRouter("/services/config",
				&ServiceConfigController{},
				"post:SetConfigServiceStepAction;get:GetConfigServiceStepAction;delete:DeleteServiceStepAction"),
			beego.NSRouter("/services/reconfig",
				&ServiceConfigController{},
				"get:GetConfigServiceFromDBAction"),
			beego.NSRouter("/services/:id([0-9]+)/status",
				&ServiceController{},
				"get:GetServiceStatusAction"),
			beego.NSRouter("/services/selectservices",
				&ServiceController{},
				"get:GetSelectableServicesAction"),
			beego.NSRouter("/services/yaml/upload",
				&ServiceController{},
				"post:UploadYamlFileAction"),
			beego.NSRouter("/services/yaml/download",
				&ServiceController{},
				"get:DownloadDeploymentYamlFileAction"),
			beego.NSRouter("/images/dockerfile/upload",
				&ImageController{},
				"post:UploadDockerfileFileAction"),
			beego.NSRouter("/images/dockerfile/download",
				&ImageController{},
				"get:DownloadDockerfileFileAction"),
			beego.NSRouter("/services/:id([0-9]+)/info",
				&ServiceController{},
				"get:GetServiceInfoAction"),
			beego.NSRouter("/services/info",
				&ServiceController{},
				"post:StoreServiceRoute"),
			beego.NSRouter("/services/:id([0-9]+)/test",
				&ServiceController{},
				"post:DeployServiceTestAction"),
			beego.NSRouter("/services/:id([0-9]+)/toggle",
				&ServiceController{},
				"put:ToggleServiceAction"),
			beego.NSRouter("/services/:id([0-9]+)/scale",
				&ServiceController{},
				"put:ScaleServiceAction;get:GetScaleStatusAction"),
			beego.NSRouter("/services/:id([0-9]+)/publicity",
				&ServiceController{},
				"put:ServicePublicityAction"),
			beego.NSRouter("/files/upload",
				&FileUploadController{},
				"post:Upload"),
			beego.NSRouter("/files/download",
				&FileUploadController{},
				"get:Download"),
			beego.NSRouter("/files/list",
				&FileUploadController{},
				"post:ListFiles"),
			beego.NSRouter("/files/remove",
				&FileUploadController{},
				"post:RemoveFile"),
			beego.NSRouter("/jenkins-job/:userID([0-9]+)/:buildNumber([0-9]+)",
				&JenkinsJobCallbackController{},
				"get:BuildNumberCallback"),
			beego.NSRouter("/jenkins-job/console",
				&JenkinsJobController{},
				"get:Console"),
			beego.NSRouter("/jenkins-job/stop",
				&JenkinsJobController{},
				"get:Stop"),
		),
	)

	beego.AddNamespace(ns)
	beego.Router("/deploy/:owner_name/:project_name/:service_name", &ServiceShowController{})
	beego.SetStaticPath("/swagger", "swagger")

	orm.RegisterModel(new(dashboard.NodeDashboardMinute), new(dashboard.NodeDashboardHour), new(dashboard.NodeDashboardDay))
}

func connectToDB() {
	hostIP := os.Getenv("HOST_IP")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:root123@tcp(%s:3306)/board?charset=utf8", hostIP))
	if err != nil {
		logs.Error("Failed to connect to DB.")
	}
}

func createAppConf() {
	var content []byte
	hostIP := os.Getenv("HOST_IP")
	tokenServer := fmt.Sprintf("tokenServerURL=http://%s:4000/tokenservice/token\n", hostIP)
	f, err := os.OpenFile("app.conf", os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		logs.Error("Open app.conf fail")
		panic(err)
	}
	_, err = f.Write([]byte(tokenServer))
	if err != nil {
		logs.Error("Write app.conf fail")
		panic(err)
	}
	content = []byte(`tokenCacheExpireSeconds=1800
dbPassword=root123
dbHost=db`)
	_, err = f.Write(content)
	if err != nil {
		logs.Error("write app.conf fail.")
		panic(err)
	}
}

func removeAppConf() {
	if err := os.Remove("app.conf"); err != nil {
		logs.Error("remove app.conf fail.")
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

func TestMain(m *testing.M) {
	utils.Initialize()
	utils.AddEnv("KUBE_MASTER_IP")
	utils.AddEnv("KUBE_MASTER_PORT")
	utils.SetConfig("KUBE_MASTER_URL", "http://%s:%s", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")
	utils.AddEnv("NODE_IP")
	utils.AddEnv("REGISTRY_IP")
	utils.AddEnv("REGISTRY_PORT")
	utils.SetConfig("REGISTRY_URL", "http://%s:%s", "REGISTRY_IP", "REGISTRY_PORT")
	utils.SetConfig("REGISTRY_BASE_URI", "%s:%s", "REGISTRY_IP", "REGISTRY_PORT")
	utils.AddValue("IS_EXTERNAL_AUTH", false)
	utils.AddValue("AUTH_MODE", "db_auth")
	utils.AddValue("BOARD_ADMIN_PASSWORD", "123456a?")
	utils.AddEnv("JENKINS_HOST_IP")
	utils.AddEnv("JENKINS_HOST_PORT")
	utils.SetConfig("JENKINS_BASE_URL", "http://%s:%s", "JENKINS_HOST_IP", "JENKINS_HOST_PORT")
	connectToDB()
	createAppConf()
	updateAdminPassword()
	defer removeAppConf()

	InitController()
	//os.Exit(m.Run())
	m.Run()
}
