package gitlab

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

var gitlabBaseURL = utils.GetConfig("GITLAB_BASE_URL")
var gitlabAPIPrefix = "/api/v4"
var gitlabDefaultPassword = "123456a?"
var maxRetryCount = 30

type gitlabHandler struct {
	accessToken      string
	gitlabAPIBaseURL string
}

var defaultAccessLevel = 30

func NewGitlabHandler(accessToken string) *gitlabHandler {
	pingURL := fmt.Sprintf("%s", gitlabBaseURL())
	for i := 0; i < maxRetryCount; i++ {
		logs.Debug("Ping Gitlab server %d time(s)...", i+1)
		if i == maxRetryCount-1 {
			logs.Warn("Failed to ping Gitlab due to exceed max retry count.")
			break
		}
		err := utils.RequestHandle(http.MethodGet, pingURL, nil, nil,
			func(req *http.Request, resp *http.Response) error {
				if resp.StatusCode <= 400 {
					return nil
				}
				return fmt.Errorf("Requested URL %s with unexpected response: %d", pingURL, resp.StatusCode)
			})
		if err == nil {
			logs.Info("Successful connected to the Gitlab service.")
			break
		}
		time.Sleep(time.Second)
	}
	return &gitlabHandler{
		accessToken:      accessToken,
		gitlabAPIBaseURL: fmt.Sprintf("%s%s", gitlabBaseURL(), gitlabAPIPrefix),
	}
}

func (g *gitlabHandler) getAccessHeader() http.Header {
	return http.Header{
		"content-type":  []string{"application/json"},
		"Authorization": []string{"Bearer " + g.accessToken},
	}
}

func (g *gitlabHandler) defaultHeader(req *http.Request) error {
	req.Header = g.getAccessHeader()
	return nil
}

type ImpersonationToken struct {
	ID        int       `json:"id"`
	Active    bool      `json:"active"`
	Scopes    []string  `json:"scopes"`
	Token     string    `json:"token"`
	Revoked   bool      `json:"revoked"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

type ProjectCreation struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	PathWithNamespace string        `json:"path_with_namespace"`
	Visibility        string        `json:"visibility"`
	ForkedFromProject ForkedProject `json:"forked_from_project"`
	Owner             UserInfo      `json:"owner"`
}

type HookCreation struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	ProjectID  int       `json:"project_id"`
	PushEvents bool      `json:"push_events"`
	CreatedAt  time.Time `json:"created_at"`
}

type MRCreation struct {
	ID              int      `json:"id"`
	IID             int      `json:"iid"`
	ProjectID       int      `json:"project_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	State           string   `json:"state"`
	TargetBranch    string   `json:"target_branch"`
	SourceBranch    string   `json:"source_branch"`
	Author          UserInfo `json:"author"`
	Assignee        UserInfo `json:"assignee"`
	SourceProjectID int      `json:"source_project_id"`
	TargetProjectID int      `json:"target_project_id"`
}

type ForkedProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	Visibility        string `json:"visibility"`
}

type Message struct {
	Fingerprint []string `json:"fingerprint"`
	Key         []string `json:"key"`
}

type FileInfo struct {
	Name    string
	Path    string
	Content string
}

type CommitInfo struct {
	Branch        string             `json:"branch"`
	CommitMessage string             `json:"commit_message"`
	AuthorName    string             `json:"author_name"`
	AuthorEmail   string             `json:"author_email"`
	Actions       []CommitActionInfo `json:"actions"`
}

type CommitActionInfo struct {
	Action   string `json:"action"`
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

type CommitCreation struct {
	ID             string      `json:"id"`
	ShortID        string      `json:"short_id"`
	Title          string      `json:"title"`
	AuthorName     string      `json:"author_name"`
	AuthorEmail    string      `json:"author_email"`
	CommitterName  string      `json:"committer_name"`
	CommitterEmail string      `json:"committer_email"`
	CreatedAt      string      `json:"created_at"`
	Message        string      `json:"message"`
	ParentIDs      []string    `json:"parent_ids"`
	Stats          CommitStats `json:"stats"`
	WebURL         string      `json:"web_url"`
}

type CommitStats struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Total     int `json:"total"`
}

type commonRespMessage struct {
	Message string `json:"message"`
}

type PipelineStatus struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	Sha       string `json:"sha"`
	BeforeSha string `json:"before_sha"`
	Tag       bool   `json:"tag"`
}

var ErrFileAlreadyExists = errors.New("A file with this name already exists")
var ErrFileDoesNotExists = errors.New("A file with this name doesn't exist")

func (f FileInfo) EscapedPath() string {
	return strings.ReplaceAll(url.PathEscape(f.Path), ".", "%2E")
}

type CommitRepoData struct {
	Branch        string `json:"branch"`
	AuthorEmail   string `json:"author_email"`
	AuthorName    string `json:"author_name"`
	Content       string `json:"content"`
	CommitMessage string `json:"commit_message"`
	FilePath      string `json:"file_path"`
}

type FileCreation struct {
	FilePath string `json:"file_path"`
	Branch   string `json:"branch"`
	Ref      string `json:"ref"`
	Content  string `json:"content"`
}

type AddSSHKeyResponse struct {
	AddSSHKeyMessage Message `json:message`
}

func (g *gitlabHandler) CreateUser(user model.User) (u UserInfo, err error) {
	userList, err := g.GetUserInfo(user.Username)
	if len(userList) == 0 {
		err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/users", g.gitlabAPIBaseURL),
			func(req *http.Request) error {
				req.Header = g.getAccessHeader()
				formData := url.Values{}
				formData.Add("name", user.Username)
				formData.Add("email", user.Email)
				formData.Add("password", user.Password)
				formData.Add("username", user.Username)
				formData.Add("skip_confirmation", "true")
				req.URL.RawQuery = formData.Encode()
				return nil
			}, nil, func(req *http.Request, resp *http.Response) error {
				return utils.UnmarshalToJSON(resp.Body, &u)
			})
		return
	}
	logs.Debug("User: %s already exists bypassing to create.", user.Username)
	if len(userList) == 1 {
		logs.Debug("Found user from Gitlab: %+v", u)
		u = userList[0]
		return
	}
	return
}

func (g *gitlabHandler) ImpersonationToken(user UserInfo) (token ImpersonationToken, err error) {
	userList, err := g.GetUserInfo(user.Username)
	if err != nil {
		logs.Error("Failed to get user info via Gitlab API by username: %s, error: %+v", user.Username, err)
		return
	}
	if len(userList) == 0 {
		logs.Error("No user found by searching name: %s", user.Username)
		return
	}
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/users/%d/impersonation_tokens", g.gitlabAPIBaseURL, userList[0].ID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("name", fmt.Sprintf("%s's token", user.Name))
			formData.Add("scopes[]", "api")
			formData.Add("scopes[]", "read_repository")
			formData.Add("scopes[]", "write_repository")
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &token)
		})
	return
}

func (g *gitlabHandler) getUserStatus(user UserInfo) (err error) {
	err = utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/users/%s/status", g.gitlabAPIBaseURL, user.Username),
		g.defaultHeader, nil, func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode == http.StatusNotFound {
				return utils.ErrNotFound
			}
			return nil
		})
	return
}

func (g *gitlabHandler) GetUserInfo(username string) (userList []UserInfo, err error) {
	err = utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/users?search=%s", g.gitlabAPIBaseURL, username),
		g.defaultHeader, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &userList)
		})
	return
}

func (g *gitlabHandler) AddSSHKey(title string, key string) (a AddSSHKeyResponse, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/user/keys", g.gitlabAPIBaseURL),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("title", title)
			formData.Add("key", key)
			formData.Add("expires_at", "")
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &a)
		})
	return
}

func (g *gitlabHandler) CreateHook(project model.Project, hookURL string) (h HookCreation, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects/%d/hooks", g.gitlabAPIBaseURL, project.ID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("url", hookURL)
			formData.Add("push_events", "false")
			formData.Add("pipeline_events", "true")
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &h)
		})
	return
}

func (g *gitlabHandler) GetRepoInfo(project model.Project) (p []ProjectCreation, err error) {
	err = utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/projects?search=%s", g.gitlabAPIBaseURL, project.Name),
		g.defaultHeader, nil, func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode == http.StatusNotFound {
				return utils.ErrNotFound
			}
			return utils.UnmarshalToJSON(resp.Body, &p)
		})
	return
}

func (g *gitlabHandler) CreateRepo(user model.User, project model.Project) (p ProjectCreation, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects", g.gitlabAPIBaseURL),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("path", fmt.Sprintf("%s", project.Name))
			formData.Add("visibility", "public")
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &p)
		})
	return
}

func (g *gitlabHandler) AddMemberToRepo(user model.User, project model.Project) (u UserInfo, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects/%d/members", g.gitlabAPIBaseURL, project.ID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("user_id", fmt.Sprintf("%d", user.ID))
			formData.Add("access_level", fmt.Sprintf("%d", defaultAccessLevel))
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &u)
		})
	return
}

func (g *gitlabHandler) ForkRepo(forkedFromProjectID int, repoName string) (p ProjectCreation, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects/%d/fork", g.gitlabAPIBaseURL, forkedFromProjectID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("path", repoName)
			formData.Add("name", repoName)
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &p)
		})
	return
}

func (g *gitlabHandler) GetFileRawContent(project model.Project, branch string, filePath string) (content []byte, err error) {
	err = utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/projects/%d/repository/files/%s/raw", g.gitlabAPIBaseURL, project.ID, url.PathEscape(filePath)),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			queryParam := url.Values{}
			queryParam.Add("ref", branch)
			req.URL.RawQuery = queryParam.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode == http.StatusNotFound {
				return ErrFileDoesNotExists
			}
			content, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read from response body with error: %+v", err)
			}
			return err
		})
	return
}

func (g *gitlabHandler) ManipulateFile(action string, user model.User, project model.Project, branch string, fileInfo FileInfo) (f FileCreation, err error) {
	requestMethod := http.MethodPost
	if action == "create" {
		requestMethod = http.MethodPost
	} else if action == "update" {
		requestMethod = http.MethodPut
	} else if action == "detect" {
		requestMethod = http.MethodGet
	}
	err = utils.RequestHandle(requestMethod, fmt.Sprintf("%s/projects/%d/repository/files/%s", g.gitlabAPIBaseURL, project.ID, fileInfo.EscapedPath()),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			if req.Method == http.MethodGet {
				queryParam := url.Values{}
				queryParam.Add("ref", branch)
				req.URL.RawQuery = queryParam.Encode()
			}
			return nil
		}, CommitRepoData{
			Branch:        branch,
			AuthorEmail:   user.Email,
			AuthorName:    user.Username,
			Content:       fileInfo.Content,
			CommitMessage: fmt.Sprintf("Add file: %s", fileInfo.Name),
		}, func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode == http.StatusNotFound {
				return ErrFileDoesNotExists
			}
			utils.UnmarshalToJSON(resp.Body, &f)
			return nil
		})
	return
}

func (g *gitlabHandler) CommitMultiFiles(user model.User, project model.Project, branch string, commitMessage string, isRemoved bool, commitActionInfos []CommitActionInfo) (c CommitCreation, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects/%d/repository/commits", g.gitlabAPIBaseURL, project.ID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			return nil
		}, CommitInfo{
			Branch:        branch,
			CommitMessage: commitMessage,
			AuthorName:    user.Username,
			AuthorEmail:   user.Email,
			Actions:       commitActionInfos,
		}, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &c)
		})
	return
}
func (g *gitlabHandler) ListMR(sourceProject model.Project) (mrList []MRCreation, err error) {
	err = utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/projects/%d/merge_requests", g.gitlabAPIBaseURL, sourceProject.ID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &mrList)
		})
	return
}

func (g *gitlabHandler) CreateMR(assignee model.User, sourceProject model.Project, targetProject model.Project, sourceBranch string, targetBranch string, title string, description string) (m MRCreation, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects/%d/merge_requests", g.gitlabAPIBaseURL, sourceProject.ID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("source_branch", sourceBranch)
			formData.Add("target_branch", targetBranch)
			formData.Add("title", title)
			formData.Add("description", description)
			formData.Add("assignee_id", fmt.Sprintf("%d", assignee.ID))
			formData.Add("target_project_id", fmt.Sprintf("%d", targetProject.ID))
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &m)
		})
	return
}

func (g *gitlabHandler) AcceptMR(sourceProject model.Project, mergeRequestID int) (m MRCreation, err error) {
	err = utils.RequestHandle(http.MethodPut, fmt.Sprintf("%s/projects/%d/%d/merge", g.gitlabAPIBaseURL, sourceProject.ID, mergeRequestID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &m)
		})
	return
}

func (g *gitlabHandler) DeleteProject(projectID int) error {
	return utils.SimpleDeleteRequestHandle(fmt.Sprintf("%s/projects/%d", g.gitlabAPIBaseURL, projectID), g.getAccessHeader())
}

func (g *gitlabHandler) DeleteUser(userID int) error {
	return utils.SimpleDeleteRequestHandle(fmt.Sprintf("%s/users/%d", g.gitlabAPIBaseURL, userID), g.getAccessHeader())
}

func (g *gitlabHandler) CancelPipeline(projectID int, pipelineID int) (p PipelineStatus, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects/%d/pipelines/%d/cancel", g.gitlabAPIBaseURL, projectID, pipelineID),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &p)
		})
	return
}
