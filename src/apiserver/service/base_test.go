package service_test

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var (
	boardHostIP         = utils.GetConfig("BOARD_HOST_IP")
	repoPath            = utils.GetConfig("BASE_REPO_PATH")
	gogitsBaseURL       = utils.GetConfig("GOGITS_BASE_URL")
	gogitsRepoURL       = utils.GetConfig("GOGITS_SSH_URL")
	sshKeyPath          = utils.GetConfig("SSH_KEY_PATH")
	jenkinsNodeIP       = utils.GetConfig("JENKINS_NODE_IP")
	jenkinsNodeSSHPort  = utils.GetConfig("JENKINS_NODE_SSH_PORT")
	jenkinsNodeUsername = utils.GetConfig("JENKINS_NODE_USERNAME")
	jenkinsNodePassword = utils.GetConfig("JENKINS_NODE_PASSWORD")
	jenkinsNodeVolume   = utils.GetConfig("JENKINS_NODE_VOLUME")
)

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	utils.AddValue("BASE_REPO_PATH", "/tmp/test-repos")
	utils.AddValue("SSH_KEY_PATH", "/tmp/test-keys")
	utils.AddValue("KUBE_CONFIG_PATH", "/tmp/config")

	dao.InitDB()
	os.Exit(func() int {
		r := m.Run()
		defer func() {
			cleanUpProject()
			cleanUpUser()
			cleanUpProjectMember()
		}()
		return r
	}())
}

func cleanUpProject() {
	o := orm.NewOrm()
	affectedCount, err := o.Delete(&project)
	if err != nil {
		logs.Error("Failed to clean up project: %+v", err)
	}
	logs.Info("Deleted in project %d row(s) affected.", affectedCount)
}

func cleanUpUser() {
	o := orm.NewOrm()
	affectedCount, err := o.Delete(&user)
	if err != nil {
		logs.Error("Failed to clean up user: %+v", err)
	}
	logs.Info("Deleted  in user %d row(s) affected.", affectedCount)
}

func cleanUpProjectMember() {
// 	o := orm.NewOrm()
// 	var affectedCount int64
// 	var err error
// 	affectedCount, err = o.Delete(&testMemberProject)
// 	if err != nil {
// 		logs.Error("Failed to delete project: %+v", err)
// 	}
// 	logs.Info("Deleted in project %d row(s) affected.", affectedCount)
// 	affectedCount, err = o.Delete(&testMember)
// 	if err != nil {
// 		logs.Error("Failed to delete member: %+v", err)
// 	}
// 	logs.Info("Deleted in member %d row(s) affected.", affectedCount)
}
