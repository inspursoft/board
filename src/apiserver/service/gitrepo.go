package service

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type repoHandler struct {
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

func getSSHAuth() (*gitssh.PublicKeys, error) {
	deployKey, err := ioutil.ReadFile("/root/.ssh/id_rsa")
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

func InitRepo(servePath, path string) (*repoHandler, error) {
	auth, err := getSSHAuth()
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:  "git@gitserver:/gitserver" + servePath,
		Auth: auth,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			return OpenRepo(path)
		}
		if err != transport.ErrEmptyRemoteRepository {
			return nil, err
		}
	}
	worktree, err := getWorktree(repo)
	if err != nil {
		return nil, err
	}
	return &repoHandler{repo: repo, worktree: worktree}, nil
}

func OpenRepo(path string) (*repoHandler, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	worktree, err := getWorktree(repo)
	if err != nil {
		return nil, err
	}
	return &repoHandler{repo: repo, worktree: worktree}, nil
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

func (r *repoHandler) Commit(message string, signature *object.Signature) (*repoHandler, error) {
	_, err := r.worktree.Commit(message, &git.CommitOptions{
		All:    true,
		Author: signature,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *repoHandler) Push() error {
	auth, err := getSSHAuth()
	if err != nil {
		return err
	}
	return r.repo.Push(&git.PushOptions{Auth: auth})
}

func (r *repoHandler) Pull() error {
	auth, err := getSSHAuth()
	if err != nil {
		return err
	}
	err = r.worktree.Pull(&git.PullOptions{Auth: auth})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}
