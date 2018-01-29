package service

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/ssh"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type repoHandler struct {
	username string
	repo     *git.Repository
	worktree *git.Worktree
}

func InitBareRepo(servePath string) (*repoHandler, error) {
	repo, err := git.PlainInit(servePath, true)
	if err != nil {
		return nil, err
	}
	return &repoHandler{repo: repo}, nil
}

func getSSHAuth(username string) (*gitssh.PublicKeys, error) {
	sshPrivateKeyPath := filepath.Join(sshKeyPath(), username, sshPrivateKey)
	logs.Debug("SSH private key path: %s", sshPrivateKeyPath)
	deployKey, err := ioutil.ReadFile(sshPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(deployKey)
	if err != nil {
		return nil, err
	}
	auth := &gitssh.PublicKeys{User: "git", Signer: signer}
	return auth, nil
}

func InitRepo(serveURL, username, path string) (*repoHandler, error) {
	auth, err := getSSHAuth(username)
	if err != nil {
		return nil, err
	}
	logs.Debug("Repo URL: %s", serveURL)
	logs.Debug("Repo path: %s", path)
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:  serveURL,
		Auth: auth,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			return OpenRepo(path, username)
		}
		if err == transport.ErrEmptyRemoteRepository {
			return nil, nil
		}
	}

	worktree, err := getWorktree(repo)
	if err != nil {
		return nil, err
	}
	return &repoHandler{username: username, repo: repo, worktree: worktree}, nil
}

func OpenRepo(path, username string) (*repoHandler, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	worktree, err := getWorktree(repo)
	if err != nil {
		return nil, err
	}
	return &repoHandler{username: username, repo: repo, worktree: worktree}, nil
}

func getWorktree(repo *git.Repository) (*git.Worktree, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	return worktree, nil
}

func (r *repoHandler) Add(filename string) (*repoHandler, error) {
	_, err := r.worktree.Add(filename)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *repoHandler) Commit(message, username, email string) (*repoHandler, error) {
	_, err := r.worktree.Commit(message, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  username,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *repoHandler) Push() error {
	auth, err := getSSHAuth(r.username)
	if err != nil {
		return err
	}
	return r.repo.Push(&git.PushOptions{Auth: auth})
}

func (r *repoHandler) Pull() error {
	auth, err := getSSHAuth(r.username)
	if err != nil {
		return err
	}
	err = r.worktree.Pull(&git.PullOptions{Auth: auth})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func (r *repoHandler) Remove(filename string) (*repoHandler, error) {
	_, err := r.worktree.Remove(filename)
	if err != nil {
		return nil, err
	}
	return r, nil
}
