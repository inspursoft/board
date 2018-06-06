package controller

import (
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/dao"
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
	utils.InitializeDefaultConfig()
	utils.AddEnv("NODE_IP")
	utils.AddEnv("TOKEN_SERVER_IP")
	utils.AddEnv("TOKEN_SERVER_PORT")
	utils.SetConfig("TOKEN_SERVER_URL", "http://%s:%s/tokenservice/token", "TOKEN_SERVER_IP", "TOKEN_SERVER_PORT")
	utils.SetConfig("SSH_KEY_PATH", "/tmp/ssh-keys")

	dao.InitDB()
	updateAdminPassword()
	InitController()
	os.Exit(m.Run())
}
