package gitlab

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"net/url"
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
		"content-type":  []string{"application/form-data", "application/json"},
		"Authorization": []string{"Bearer " + g.accessToken},
	}
}

func (g *gitlabHandler) defaultHeader(req *http.Request) error {
	req.Header = g.getAccessHeader()
	return nil
}

type UserCreation struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
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

type userStatus struct {
	Emoji       string `json:"emoji"`
	Message     string `json:"message"`
	MessageHTML string `json:"message_html"`
}

type userInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

type ProjectCreation struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Visibility string `json:"visibility"`
}

type Message struct {
	Fingerprint []string `json:"fingerprint"`
	Key         []string `json:"key"`
}

type AddSSHKeyResponse struct {
	AddSSHKeyMessage Message `json:message`
}

func (g *gitlabHandler) CreateUser(user model.User) (u UserCreation, err error) {
	_, err = g.getUserStatus(UserCreation{Username: user.Username})
	if err == utils.ErrNotFound {
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
	return
}

func (g *gitlabHandler) ImpersonationToken(user UserCreation) (token ImpersonationToken, err error) {
	userList, err := g.getUserInfo(user.Username)
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

func (g *gitlabHandler) getUserStatus(user UserCreation) (u userStatus, err error) {
	err = utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s/users/%s/status", g.gitlabAPIBaseURL, user.Username),
		g.defaultHeader, nil, func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode == http.StatusNotFound {
				return utils.ErrNotFound
			}
			return utils.UnmarshalToJSON(resp.Body, &u)
		})
	return
}

func (g *gitlabHandler) getUserInfo(username string) (userList []userInfo, err error) {
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

func (g *gitlabHandler) CreateRepo(user model.User, project model.Project) (p ProjectCreation, err error) {
	err = utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/projects", g.gitlabAPIBaseURL),
		func(req *http.Request) error {
			req.Header = g.getAccessHeader()
			formData := url.Values{}
			formData.Add("path", user.Username)
			formData.Add("name", project.Name)
			req.URL.RawQuery = formData.Encode()
			return nil
		}, nil, func(req *http.Request, resp *http.Response) error {
			return utils.UnmarshalToJSON(resp.Body, &p)
		})
	return
}

func (g *gitlabHandler) DeleteProject(projectID int) error {
	return utils.SimpleDeleteRequestHandle(fmt.Sprintf("%s/projects/%d", g.gitlabAPIBaseURL, projectID), g.getAccessHeader())
}

func (g *gitlabHandler) DeleteUser(userID int) error {
	return utils.SimpleDeleteRequestHandle(fmt.Sprintf("%s/users/%d", g.gitlabAPIBaseURL, userID), g.getAccessHeader())
}
