package service

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service/devops/gitlab"
	"git/inspursoft/board/src/apiserver/service/devops/gitlabci"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"golang.org/x/net/context"
)

const (
	gitlabBuildConsoleTemplateURL = "%s/{{.JobName}}/-/jobs/{{.BuildSerialID}}/raw"
)

var gitlabAdminToken = utils.GetConfig("GITLAB_ADMIN_TOKEN")
var gitlabBaseURL = utils.GetConfig("GITLAB_BASE_URL")
var registryBaseURI = utils.GetConfig("REGISTRY_BASE_URI")
var kanikoImage = "kaniko-project/executor:dev"

type gitlabJenkinsPushProjectPayload struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	GitHTTPURL string `json:"git_http_url"`
}

type gitlabJenkinsPushPayload struct {
	Project      gitlabJenkinsPushProjectPayload `json:"project"`
	NodeSelector string                          `json:"node_selector"`
}

type objectAttr struct {
	ID             int      `json:"id"`
	Ref            string   `json:"ref"`
	Tag            bool     `json:"tag"`
	Sha            string   `json:"sha"`
	BeforeSha      string   `json:"before_sha"`
	Source         string   `json:"source"`
	Status         string   `json:"status"`
	DetailedStatus string   `json:"detailed_status"`
	Stages         []string `json:"stages"`
}

type build struct {
	ID    int    `json:"id"`
	Stage string `json:"stage"`
}

type gitlabPipelinePayload struct {
	ObjectAttr objectAttr `json:"object_attributes"`
	Builds     []build    `json:"builds"`
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan bool)

	handle := func() <-chan error {
		e := make(chan error)
		go func() {
			defer close(e)
			ctx = context.WithValue(ctx, storeItem, "Gitlab")
			gitlabHandler := gitlab.NewGitlabHandler(accessToken)
			if gitlabHandler == nil {
				logs.Error("failed to create Gitlab handler")
				cancel()
			}
			userInfo := model.User{Username: user.Username, Email: user.Email, RepoToken: user.RepoToken}

			projectInfo := model.Project{Name: repoName}
			projectCreation, err := gitlabHandler.CreateRepo(userInfo, projectInfo)
			if err != nil {
				logs.Error("Failed to create repo via Gitlab API, error %+v", err)
				close(e)
			}
			logs.Debug("Successful created Gitlab project: %+v", projectCreation)
			hookURL := fmt.Sprintf("%s/jenkins-job/pipeline?user_id=%d", boardAPIBaseURL(), userID)
			projectInfo.ID = int64(projectCreation.ID)

			gitlabHookCreation, err := gitlabHandler.CreateHook(projectInfo, hookURL)
			if err != nil {
				logs.Error("Failed to create hook to repo: %s, error: %+v", repoName, err)
				cancel()
			}
			logs.Debug("Successful created Gitlab hook: %+v", gitlabHookCreation)

			projectInfo.ID = int64(projectCreation.ID)

			fileInfo := gitlab.FileInfo{
				Name:    "README.md",
				Path:    "README.md",
				Content: "README file created by Board.",
			}

			fileCreation, err := gitlabHandler.ManipulateFile("create", userInfo, projectInfo, "master", fileInfo)
			if err != nil {
				logs.Error("Failed to create file: %+v to the repo: %s, error: %+v", fileInfo, projectInfo.Name, err)
				cancel()
			}
			logs.Debug("Successful created file: %+v to Gitlab repository: %s", fileCreation, projectInfo.Name)

			for {
				select {
				case <-ctx.Done():
					if err := ctx.Err(); err != nil {
						logs.Error("Canceled context in creation %s, error: %+v", ctx.Value(storeItem), err)
						e <- err
					}
					logs.Debug("Execution for %s in context has done.", ctx.Value(storeItem))
				case <-done:
					logs.Debug("Finished executing Gitlab process.")
					e <- nil
				}
			}
		}()
		return e
	}
	close(done)
	return <-handle()
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

func (g GitlabDevOps) CustomHookPipelinePayload(rawPayload []byte) (pipelineID int, buildNumber int, err error) {
	var cp gitlabPipelinePayload
	err = json.Unmarshal(rawPayload, &cp)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal JSON custom pipeline payload: %+v", err)
		return
	}
	pipelineID = cp.ObjectAttr.ID
	buildNumber = cp.Builds[0].ID
	logs.Debug("Resolved pipeline: %d payload for build number: %d", pipelineID, buildNumber)
	return
}

func (g GitlabDevOps) GetRepoFile(username string, repoName string, branch string, filePath string) ([]byte, error) {
	user, err := GetUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by name: %s, error: %+v", username, err)
	}
	gitlabHandler := gitlab.NewGitlabHandler(user.RepoToken)
	if gitlabHandler == nil {
		return nil, fmt.Errorf("failed to create Gitlab handler")
	}
	project, err := g.GetRepo(user.RepoToken, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo by name: %s, error: %+v", repoName, err)
	}
	content, err := gitlabHandler.GetFileRawContent(project, branch, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %s from branch: %s in repo: %s", filePath, branch, repoName)
	}
	logs.Debug("Got file: %s with content: %s", filePath, string(content))
	return content, nil
}

func (g GitlabDevOps) DeleteUser(username string) error {
	user, err := g.GetUser(gitlabAdminToken(), username)
	if err != nil {
		return fmt.Errorf("failed to get repo user: %s, error: %+v", user.Username, err)
	}
	return gitlab.NewGitlabHandler(gitlabAdminToken()).DeleteUser(int(user.ID))
}

func generateBuildingImageGitlabCIYAML(configurations map[string]string) error {
	token := configurations["token"]
	imageURI := configurations["image_uri"]
	dockerfileName := configurations["dockerfile"]
	repoPath := configurations["repo_path"]
	ciImage := gitlabci.Image{Name: fmt.Sprintf("%s/%s", registryBaseURI(), kanikoImage)}
	ciJobs := make(map[string]gitlabci.Job)
	var ci gitlabci.GitlabCI
	ciJobs["build-image"] = gitlabci.Job{
		Image: &ciImage,
		Stage: "build-image",
		Tags:  []string{"docker-ci"},
		Script: []string{
			ci.WriteMultiLine("CI_REGISTRY=%s", registryBaseURI()),
			ci.WriteMultiLine("CI_REGISTRY_USER=%s", "admin"),
			ci.WriteMultiLine("CI_REGISTRY_PASSWORD=%s", "$(echo -n 123456a? | base64)"),
			"if [ -d 'upload' ]; then rm -rf upload; fi",
			"if [ -e 'attachment.zip' ]; then rm -f attachment.zip; fi",
			ci.WriteMultiLine("token=%s", token),
			ci.WriteMultiLine("status=`curl -I \"%s/files/download?token=$token\" 2>/dev/null | head -n 1 | awk '{print $2}'`", boardAPIBaseURL()),
			ci.WriteMultiLine("bash -c \"if [ $status == '200' ]; then curl -o attachment.zip \"%s/files/download?token=$token\" && mkdir -p upload && unzip attachment.zip -d upload; fi\"", boardAPIBaseURL()),
			"export PATH=/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin",
			ci.WriteMultiLine("/kaniko/executor --context $CI_PROJECT_DIR --dockerfile $CI_PROJECT_DIR/containers/%s --destination %s --cache-repo $CI_REGISTRY --cache=true", dockerfileName, imageURI),
		},
	}
	return ci.GenerateGitlabCI(ciJobs, repoPath)
}

func generatePushingImageGitlabCIYAML(configurations map[string]string) error {
	token := configurations["token"]
	imagePackageName := configurations["image_package_name"]
	imageURI := configurations["image_uri"]
	repoPath := configurations["repo_path"]
	ciJobs := make(map[string]gitlabci.Job)
	var ci gitlabci.GitlabCI
	ciJobs["push-image"] = gitlabci.Job{
		Stage: "push-image",
		Tags:  []string{"shell-ci"},
		Script: []string{
			"if [ -d 'upload' ]; then rm -rf upload; fi",
			"if [ -e 'attachment.zip' ]; then rm -f attachment.zip; fi",
			ci.WriteMultiLine("token=%s", token),
			ci.WriteMultiLine("status=`curl -I \"%s/files/download?token=$token\" 2>/dev/null | head -n 1 | awk '{print $2}'`", boardAPIBaseURL()),
			ci.WriteMultiLine("bash -c \"if [ $status == '200' ]; then curl -o attachment.zip \"%s/files/download?token=$token\" && mkdir -p upload && unzip attachment.zip -d upload; fi\"", boardAPIBaseURL()),
			"export PATH=/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin",
			ci.WriteMultiLine("image_name_tag=$(docker load -i upload/%s |grep 'Loaded image'|awk '{print $NF}')", imagePackageName),
			ci.WriteMultiLine("image_name_tag=${image_name_tag#sha256:}"),
			ci.WriteMultiLine("docker tag $image_name_tag %s", imageURI),
			ci.WriteMultiLine("docker push %s", imageURI),
			ci.WriteMultiLine("docker rmi %s", imageURI),
			ci.WriteMultiLine("if [[ $image_name_tag =~ ':' ]]; then docker rmi $image_name_tag; fi"),
		},
	}
	return ci.GenerateGitlabCI(ciJobs, repoPath)
}

func (g GitlabDevOps) CreateCIYAML(action yamlAction, configurations map[string]string) (yamlName string, err error) {
	yamlName = gitlabci.GitlabCIFilename
	switch action {
	case BuildDockerImageCIYAML:
		err = generateBuildingImageGitlabCIYAML(configurations)
	case PushDockerImageCIYAML:
		err = generatePushingImageGitlabCIYAML(configurations)
	}
	return
}

func (g GitlabDevOps) ResolveHandleURL(configurations map[string]string) (consoleURL string, stopURL string, err error) {
	jobName := configurations["project_name"]
	repoToken := configurations["repo_token"]
	pipelineID, _ := strconv.Atoi(configurations["pipeline_id"])
	buildSerialID := configurations["build_serial_id"]
	query := CIConsole{JobName: jobName, BuildSerialID: buildSerialID}
	consoleURL, err = utils.GenerateURL(fmt.Sprintf(gitlabBuildConsoleTemplateURL, gitlabBaseURL()), query)

	gitlabHandler := gitlab.NewGitlabHandler(repoToken)
	if gitlabHandler == nil {
		return
	}
	project, err := g.GetRepo(repoToken, jobName)
	if err != nil {
		err = fmt.Errorf("failed to get repo by name: %s, error: %+v", jobName, err)
		return
	}
	stopURL = fmt.Sprintf("%s/api/v4/projects/%d/pipelines/%d/cancel?private_token=%s", gitlabBaseURL(), project.ID, pipelineID, repoToken)
	return
}
