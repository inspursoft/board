package service

import (
	"fmt"
	"github.com/inspursoft/board/src/common/utils"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/ssh"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var gogitsSSHPort = utils.GetConfig("GOGITS_SSH_PORT")
var gogitsHostIP = utils.GetConfig("GOGITS_HOST_IP")

type repoHandler struct {
	username string
	email    string
	repo     *git.Repository
	worktree *git.Worktree
	hash     plumbing.Hash
	repoPath string
	isRemove bool
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

	err := exec.Command("ssh", "-i", sshPrivateKeyPath, "-4", gogitsHostIP(), "-o", "StrictHostKeyChecking=no", "-p", gogitsSSHPort()).Run()
	if err != nil {
		logs.Warn("Failed to add Public key to known hosts: %+v", err)
	}
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

func InitRepo(serveURL, username, email, path string) (*repoHandler, error) {
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
			return OpenRepo(path, username, email)
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

func OpenRepo(path, username, email string) (*repoHandler, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	worktree, err := getWorktree(repo)
	if err != nil {
		return nil, err
	}
	return &repoHandler{username: username, email: email, repo: repo, worktree: worktree, repoPath: path}, nil
}

func (r *repoHandler) genSignature() *object.Signature {
	return &object.Signature{
		Name:  r.username,
		Email: r.email,
		When:  time.Now(),
	}
}

func getWorktree(repo *git.Repository) (*git.Worktree, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	return worktree, nil
}

func (r *repoHandler) ToRemove() *repoHandler {
	r.isRemove = true
	return r
}

func (r *repoHandler) Add(filename string) (*repoHandler, error) {
	_, err := r.worktree.Add(filename)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *repoHandler) Commit(message string) (*repoHandler, error) {
	var err error
	r.hash, err = r.worktree.Commit(message, &git.CommitOptions{
		All:    true,
		Author: r.genSignature(),
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *repoHandler) Tag(tagName, message string) (*repoHandler, error) {
	r, err := r.Commit(message)
	if err != nil {
		return nil, err
	}
	tag := object.Tag{
		Name:       tagName,
		Message:    message,
		Tagger:     *r.genSignature(),
		Target:     r.hash,
		TargetType: plumbing.CommitObject,
	}
	encodedObject := r.repo.Storer.NewEncodedObject()
	err = tag.Encode(encodedObject)
	if err != nil {
		return nil, err
	}
	hash, err := r.repo.Storer.SetEncodedObject(encodedObject)
	err = r.repo.Storer.SetReference(plumbing.NewReferenceFromStrings("refs/tags"+tagName, hash.String()))
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

func (r *repoHandler) SimplePush(items ...string) error {
	logs.Debug("Repo path for pushing objects: %s", r.repoPath)
	for _, item := range items {
		if r.isRemove {
			r.Remove(item)
		} else {
			r.Add(item)
		}
		logs.Debug(">>>>> pushed item: %s", item)
	}

	message := fmt.Sprintf("Push items: %s to repo: %s", strings.Join(items, ","), r.repoPath)
	if r.isRemove {
		message = "[DELETED]" + message
	}

	var err error
	_, err = r.Commit(message)
	if err != nil {
		return fmt.Errorf("failed to commit changes to user's repo: %+v", err)
	}
	err = r.Push()
	if err != nil {
		return fmt.Errorf("failed to push objects to git repo: %+v", err)
	}
	return nil
}
