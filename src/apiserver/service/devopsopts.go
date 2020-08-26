package service

import (
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
)

type key int

type yamlAction int

const (
	storeItem              key        = iota
	BuildDockerImageCIYAML yamlAction = iota
	PushDockerImageCIYAML
)

var devOpsOpt = utils.GetConfig("DEVOPS_OPT")
var devOpsRegistries map[string]DevOps
var boardAPIBaseURL = utils.GetConfig("BOARD_API_BASE_URL")

type CommitItem struct {
	PathWithName string
	Content      string
}

type DevOps interface {
	SignUp(user model.User) error
	CreateAccessToken(username string, password string) (string, error)
	CommitAndPush(repoName string, isRemoved bool, username string, email string, items ...CommitItem) error
	ConfigSSHAccess(username string, token string, publicKey string) error
	CreateRepoAndJob(userID int64, projectName string) error
	ForkRepo(forkedUser model.User, baseRepoName string) error
	CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, message string) error
	MergePullRequest(repoName, repoToken string) error
	DeleteRepo(username string, repoName string) error
	CustomHookPushPayload(rawPayload []byte, nodeSelection string) error
	GetRepoFile(username string, repoName string, branch string, filePath string) ([]byte, error)
	DeleteUser(username string) error
	CreateCIYAML(action yamlAction, configurations map[string]string) (yamlName string, err error)
}

func CurrentDevOps() DevOps {
	return devOpsRegistries[devOpsOpt()]
}

func init() {
	devOpsRegistries = make(map[string]DevOps)
	devOpsRegistries["legacy"] = LegacyDevOps{}
	devOpsRegistries["gitlab"] = GitlabDevOps{}
}
