package service

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service/devops/gitlab"
	"git/inspursoft/board/src/apiserver/service/devops/jenkins"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

var gitlabAdminToken = utils.GetConfig("GITLAB_ADMIN_TOKEN")

type gitlabJenkinsPushProjectPayload struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	GitHTTPURL string `json:"git_http_url"`
}

type gitlabJenkinsPushPayload struct {
	Project      gitlabJenkinsPushProjectPayload `json:"project"`
	NodeSelector string                          `json:"node_selector"`
}

type GitlabDevOps struct{}

func (g GitlabDevOps) SignUp(user model.User) error {
	userCreation, err := gitlab.NewGitlabHandler(gitlabAdminToken()).CreateUser(user)
	if err != nil {
		logs.Error("Failed to sign up via Gitlab API, error: %+v", err)
		return err
	}
	logs.Debug("Successful signed up user: %+v", userCreation)
	return nil
}

func (g GitlabDevOps) CreateAccessToken(username string, password string) (string, error) {
	userCreation := gitlab.UserInfo{Name: username, Username: username}
	token, err := gitlab.NewGitlabHandler(gitlabAdminToken()).ImpersonationToken(userCreation)
	if err != nil {
		logs.Error("Failed to create access token via Gitlab API, error %+v", err)
		return "", err
	}
	return token.Token, nil
}

func generateCommitActionInfo(repoUser model.User, repoProject model.Project, action string, items ...CommitItem) (commitActionInfos []gitlab.CommitActionInfo, commitMessage string) {
	for i, item := range items {
		fi := gitlab.FileInfo{Path: item.PathWithName}
		_, err := gitlab.NewGitlabHandler(repoUser.RepoToken).ManipulateFile("detect", repoUser, repoProject, "master", fi)
		if err == nil {
			logs.Debug("Update file: %s as it already exist.", item.PathWithName)
			action = "update"
		}
		if err == gitlab.ErrFileDoesNotExists {
			logs.Debug("Create file: %s as it does not exist.", item.PathWithName)
			action = "create"
		}
		commitActionInfos = append(commitActionInfos, gitlab.CommitActionInfo{
			Action:   action,
			FilePath: item.PathWithName,
			Content:  item.Content,
		})
		if i == len(items)-1 {
			commitMessage += fmt.Sprintf(" %s", item.PathWithName)
		} else {
			commitMessage += fmt.Sprintf(" %s,", item.PathWithName)
		}
	}
	return
}

func (g GitlabDevOps) CommitAndPush(repoName string, isRemoved bool, username string, email string, items ...CommitItem) error {
	user, err := GetUserByName(username)
	if err != nil {
		return fmt.Errorf("failed to get project owner by username: %s, error: %+v", username, err)
	}
	repoUser, err := g.GetUser(user.RepoToken, user.Username)
	if err != nil {
		return fmt.Errorf("failed to get user from repo by name: %s, error: %+v", username, err)
	}
	repoUser.RepoToken = user.RepoToken
	repoProject, err := g.GetRepo(user.RepoToken, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repo project by name: %s, error: %+v", repoName, err)
	}
	logs.Debug("Got repo: %+v to commit and push.", repoProject)
	action := "create"
	if isRemoved {
		action = "delete"
	}
	commitActionInfos, commitMessage := generateCommitActionInfo(repoUser, repoProject, action, items...)
	logs.Debug("Commit action info: %+v", commitActionInfos)
	gitlab.NewGitlabHandler(user.RepoToken).CommitMultiFiles(repoUser, repoProject, "master", commitMessage, isRemoved, commitActionInfos)
	return nil
}

func (g GitlabDevOps) ConfigSSHAccess(username string, token string, publicKey string) error {
	addSSHKeyResponse, err := gitlab.NewGitlabHandler(token).AddSSHKey(fmt.Sprintf("%s's SSH access.", username), publicKey)
	if err != nil {
		logs.Error("Failed to config SSH access via Gitlab API, error: %+v", err)
		return err
	}
	logs.Debug("Successful configured SSH access: %+v", addSSHKeyResponse)
	return nil
}

func (g GitlabDevOps) CreateRepoAndJob(userID int64, projectName string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		logs.Error("Failed to get user: %+v", err)
		return err
	}
	if user == nil {
		return fmt.Errorf("user with ID: %d is nil", userID)
	}
	username := user.Username
	accessToken := user.RepoToken
	logs.Info("Create repo and job with username: %s, project name: %s.", username, projectName)
	repoName, err := ResolveRepoName(projectName, username)
	if err != nil {
		return err
	}
	logs.Info("Initialize serve repo with name: %s ...", repoName)

	gitlabHandler := gitlab.NewGitlabHandler(accessToken)
	if gitlabHandler == nil {
		return fmt.Errorf("failed to create Gitlab handler")
	}
	userInfo := model.User{Username: user.Username, Email: user.Email, RepoToken: user.RepoToken}

	projectInfo := model.Project{Name: repoName}
	projectCreation, err := gitlabHandler.CreateRepo(userInfo, projectInfo)
	if err != nil {
		logs.Error("Failed to create repo via Gitlab API, error %+v", err)
		return err
	}
	logs.Debug("Successful created Gitlab project: %+v", projectCreation)
	projectInfo.ID = int64(projectCreation.ID)

	hookURL := fmt.Sprintf("%s/jenkins-job/invoke", boardAPIBaseURL())
	hookCreation, err := gitlabHandler.CreateHook(projectInfo, hookURL)
	if err != nil {
		logs.Error("Failed to create hook: %s to the repo: %s, error: %+v", hookURL, projectInfo.Name, err)
	}
	logs.Debug("Successful created hook: %+v to Gitlab repository: %s", hookCreation, projectInfo.Name)

	fileInfo := gitlab.FileInfo{
		Name:    "README.md",
		Path:    "README.md",
		Content: "README file created by Board.",
	}

	fileCreation, err := gitlabHandler.ManipulateFile("create", userInfo, projectInfo, "master", fileInfo)
	if err != nil {
		logs.Error("Failed to create file: %+v to the repo: %s, error: %+v", fileInfo, projectInfo.Name, err)
		return err
	}
	logs.Debug("Successful created file: %+v to Gitlab repository: %s", fileCreation, projectInfo.Name)

	jenkinsHandler := jenkins.NewJenkinsHandler()
	err = jenkinsHandler.CreateJobWithParameter(repoName)
	if err != nil {
		logs.Error("Failed to create Jenkins' job with repo name: %s, error: %+v", repoName, err)
		return err
	}
	logs.Debug("Waiting for services to be created...")
	time.Sleep(time.Second * 12)
	return nil
}

func (g GitlabDevOps) GetRepo(token string, repoName string) (project model.Project, err error) {
	foundProjectList, err := gitlab.NewGitlabHandler(token).GetRepoInfo(model.Project{Name: repoName})
	if err != nil {
		logs.Error("Failed to get repo for name: %s with error: %+v", repoName, err)
		return
	}
	if len(foundProjectList) == 0 {
		logs.Error("Repo: %s not found.", repoName)
		return
	}
	project.ID = int64(foundProjectList[0].ID)
	project.Name = foundProjectList[0].Name
	project.OwnerName = foundProjectList[0].Owner.Name
	return
}

func (g GitlabDevOps) GetUser(token string, username string) (user model.User, err error) {
	foundUserList, err := gitlab.NewGitlabHandler(token).GetUserInfo(username)
	if err != nil {
		logs.Error("Failed to get user by name: %s with error: %+v", username, err)
		return
	}
	if len(foundUserList) == 0 {
		logs.Error("User: %s not found.", username)
		return
	}
	user.ID = int64(foundUserList[0].ID)
	user.Username = foundUserList[0].Name
	user.Email = foundUserList[0].Email
	return
}

func (g GitlabDevOps) ForkRepo(forkedUser model.User, baseRepoName string) error {
	project, err := GetProjectByName(baseRepoName)
	if err != nil {
		return fmt.Errorf("failed to get project by name: %s, error: %+v", baseRepoName, err)
	}
	projectOwner, err := GetUserByName(project.OwnerName)
	if err != nil {
		return fmt.Errorf("failed to get project owner by username: %s, error: %+v", project.OwnerName, err)
	}
	baseRepo, err := g.GetRepo(projectOwner.RepoToken, baseRepoName)
	if err != nil {
		return fmt.Errorf("failed to get repo info name: %s, error: %+v", baseRepoName, err)
	}
	forkedRepoUser, err := g.GetUser(forkedUser.RepoToken, forkedUser.Username)
	if err != nil {
		return fmt.Errorf("failed to get repo user: %s, error: %+v", forkedUser.Username, err)
	}
	memberUser, err := gitlab.NewGitlabHandler(projectOwner.RepoToken).AddMemberToRepo(forkedRepoUser, baseRepo)
	if err != nil {
		return fmt.Errorf("failed to add member: %s to project: %+v, error: %+v", forkedRepoUser.Username, baseRepo, err)
	}
	logs.Debug("Successful added member: %+v to project ID: %d", memberUser, baseRepo.ID)

	gitlabHandler := gitlab.NewGitlabHandler(forkedUser.RepoToken)
	if gitlabHandler == nil {
		return fmt.Errorf("failed to create Gitlab handler")
	}
	forkedRepoName, err := ResolveRepoName(baseRepoName, forkedUser.Username)
	if err != nil {
		return fmt.Errorf("failed to resolve repo name via base repo name: %s, error: %+v", baseRepoName, err)
	}
	forkedCreation, err := gitlabHandler.ForkRepo(int(baseRepo.ID), forkedRepoName)
	if err != nil {
		return fmt.Errorf("failed to fork repo with name: %s from base repo ID: %d", baseRepoName, baseRepo.ID)
	}
	logs.Debug("Successful forked repo with name: %s, with detail: %+v", baseRepoName, forkedCreation)

	projectInfo := model.Project{ID: int64(forkedCreation.ID)}
	hookURL := fmt.Sprintf("%s/jenkins-job/invoke", boardAPIBaseURL())
	hookCreation, err := gitlabHandler.CreateHook(projectInfo, hookURL)
	if err != nil {
		logs.Error("Failed to create hook: %s to the repo: %s, error: %+v", hookURL, projectInfo.Name, err)
	}
	logs.Debug("Successful created hook: %+v to Gitlab repository: %s", hookCreation, projectInfo.Name)

	jenkinsHandler := jenkins.NewJenkinsHandler()
	err = jenkinsHandler.CreateJobWithParameter(forkedRepoName)
	if err != nil {
		logs.Error("Failed to create Jenkins' job with project name: %s, error: %+v", forkedRepoName, err)
		return err
	}
	return nil
}

func (g GitlabDevOps) CreatePullRequestAndComment(username, ownerName, repoName, repoToken, compareInfo, title, message string) error {
	assignee, err := g.GetUser(repoToken, username)
	if err != nil {
		return fmt.Errorf("failed to get assignee by name: %s, error: %+v", username, err)
	}
	sourceProject, err := g.GetRepo(repoToken, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repo by name: %s, error: %+v", repoName, err)
	}
	foundRepoList, err := gitlab.NewGitlabHandler(repoToken).GetRepoInfo(model.Project{Name: repoName})
	if err != nil {
		return fmt.Errorf("failed to list repo info by name: %s, error: %+v", repoName, err)
	}
	if len(foundRepoList) == 0 {
		return fmt.Errorf("repo: %s not found", repoName)
	}
	targetRepo := foundRepoList[0].ForkedFromProject
	targetProject := model.Project{ID: int64(targetRepo.ID)}
	mergeInfo := strings.Split(compareInfo, "...")
	sourceBranch := mergeInfo[0]
	subMergeInfo := strings.Split(mergeInfo[1], ":")
	targetBranch := subMergeInfo[1]
	logs.Debug("Resolve merge request info by compareInfo: %s - sourceBranch: %s, targetBranch: %s", compareInfo, sourceBranch, targetBranch)

	mrCreation, err := gitlab.NewGitlabHandler(repoToken).CreateMR(assignee, sourceProject, targetProject, sourceBranch, targetBranch, title, message)
	if err != nil {
		return fmt.Errorf("failed to create MR by repo name: %s with source branch: %s, target branch: %s, to the target project: %s", repoName, sourceBranch, targetBranch, targetProject.Name)
	}
	logs.Debug("Successful created MR with detail: %+v", mrCreation)
	return nil
}

func (g GitlabDevOps) MergePullRequest(repoName, repoToken string) error {
	sourceProject, err := g.GetRepo(repoToken, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repo by name: %s, error: %+v", repoName, err)
	}
	foundMRList, err := gitlab.NewGitlabHandler(repoToken).ListMR(sourceProject)
	if err != nil {
		return fmt.Errorf("failed to list merge request by name: %s, error: %+v", repoName, err)
	}
	if len(foundMRList) == 0 {
		return fmt.Errorf("repo: %s has no merge request", repoName)
	}
	mrIID := foundMRList[0].IID
	mrAcceptance, err := gitlab.NewGitlabHandler(repoToken).AcceptMR(sourceProject, mrIID)
	if err != nil {
		return fmt.Errorf("failed to accept MR by repo name: %s, error: %+v", repoName, err)
	}
	logs.Debug("Successful accepted MR with detail: %+v", mrAcceptance)
	return nil
}

func (g GitlabDevOps) DeleteRepo(username string, repoName string) error {
	user, err := GetUserByName(username)
	if err != nil {
		return fmt.Errorf("failed to get user by name: %s, error: %+v", username, err)
	}
	gitlabHandler := gitlab.NewGitlabHandler(user.RepoToken)
	if gitlabHandler == nil {
		return fmt.Errorf("failed to create Gitlab handler")
	}
	project, err := g.GetRepo(user.RepoToken, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repo by name: %s, error: %+v", repoName, err)
	}
	err = gitlabHandler.DeleteProject(int(project.ID))
	if err != nil {
		return fmt.Errorf("failed to delete project by ID: %d, error: %+v", project.ID, err)
	}
	logs.Debug("Successful deleted project by ID: %d", project.ID)
	return nil
}

func (g GitlabDevOps) CustomHookPushPayload(rawPayload []byte, nodeSelection string) error {
	var cp gitlabJenkinsPushPayload
	err := json.Unmarshal(rawPayload, &cp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON custom push payload: %+v", err)
	}
	cp.NodeSelector = nodeSelection
	logs.Debug("Resolve for push event hook payload: %+v", cp)
	header := http.Header{
		"content-type":   []string{"application/json"},
		"X-Gitlab-Event": []string{"Push Hook"},
	}
	return utils.SimplePostRequestHandle(fmt.Sprintf("%s/generic-webhook-trigger/invoke", JenkinsBaseURL()), header, cp)
}
