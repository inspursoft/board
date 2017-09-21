package main

import (
	_ "git/inspursoft/board/src/apiserver/controller"
	_ "git/inspursoft/board/src/apiserver/router"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"

	"path/filepath"

	"github.com/astaxie/beego"
)

var adminUserID int64 = 1
var defaultInitialPassword = "123456a?"

var baseRepoPath = `/repos`
var repoServeURL = filepath.Join("root@gitserver:", "gitserver", "repos", "board_repo_serve")
var repoPath = filepath.Join(baseRepoPath, "board_repo")

func updateAdminPassword(initialPassword string) {
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
		logs.Info("Admin password has been updated successfully.")
	} else {
		logs.Info("Failed to update admin initial password.")
	}
}

func main() {
	utils.Initialize()

	utils.AddEnv("BOARD_ADMIN_PASSWORD")
	utils.AddEnv("KUBE_MASTER_HOST")
	utils.AddEnv("KUBE_MASTER_PORT")
	utils.AddEnv("REGISTRY_HOST")
	utils.AddEnv("REGISTRY_PORT")

	utils.SetConfig("REGISTRY_URL", "%s:%s", "REGISTRY_HOST", "REGISTRY_PORT")
	utils.SetConfig("KUBE_MASTER_URL", "%s:%s", "KUBE_MASTER_HOST", "KUBE_MASTER_PORT")
	utils.SetConfig("KUBE_NODE_URL", "%s:%s/api/v1/nodes", "KUBE_MASTER_HOST", "KUBE_MASTER_PORT")

	utils.SetConfig("REPO_SERVE_URL", repoServeURL)
	utils.SetConfig("REPO_PATH", repoPath)

	utils.SetConfig("REGISTRY_BASE_URI", "http://%s:%s", "REGISTRY_HOST", "REGISTRY_PORT")

	utils.ShowAllConfigs()

	updateAdminPassword(utils.GetStringValue("BOARD_ADMIN_PASSWORD"))

	beego.Run(":8088")
}
