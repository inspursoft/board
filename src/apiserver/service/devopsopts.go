package service

import (
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
)

var devOpsOpt = utils.GetConfig("DEVOPS_OPT", "legacy")
var devOpsRegistries map[string]DevOps

type DevOps interface {
	CreateRepoAndJob(userID int64, projectName string) error
	ForkRepo(forkedUser *model.User, baseRepoName string) error
	CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, message string) error
}

func CurrentDevOps() DevOps {
	return devOpsRegistries[devOpsOpt()]
}

func init() {
	devOpsRegistries = make(map[string]DevOps)
	devOpsRegistries["legacy"] = LegacyDevOps{}
}
