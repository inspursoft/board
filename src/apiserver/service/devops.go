package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/logs"
)

func ResolveRepoName(projectName, username string) (repoName string, err error) {
	project, err := GetProjectByName(projectName)
	if err != nil {
		return
	}
	if project == nil {
		err = errors.New("invalid project name")
		return
	}
	members, err := GetProjectMembers(project.ID)
	if err != nil {
		return
	}
	isMember := false
	for _, m := range members {
		if m.Username == username {
			isMember = true
		}
	}
	repoName = project.Name
	if isMember && project.OwnerName != username {
		repoName = username + "_" + project.Name
	}
	logs.Debug("Resolved repo name as: %s.", repoName)
	return
}

func ResolveRepoPath(repoName, username string) (repoPath string) {
	repoPath = filepath.Join(BaseRepoPath(), username, "contents", repoName)
	logs.Debug("Set repo path at file upload: %s", repoPath)
	return
}

func ResolveDockerfileName(imageName, tag string) string {
	imageName = imageName[strings.LastIndex(imageName, "/")+1:]
	return fmt.Sprintf("Dockerfile.%s_%s", imageName, tag)
}

func FetchFileContentByDevOpsOpt(branch string, filePath string) (io.Reader, error) {
	if devOpsOpt() == "legacy" {
		return os.Open(filePath)
	}
	relPath, err := filepath.Rel(BaseRepoPath(), filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve rel path base on: %s with raw path: %s, error: %+v", BaseRepoPath(), filePath, err)
	}
	logs.Debug("Current file path: %s and relative path is: %s", filePath, relPath)
	parts := strings.Split(relPath, "/")
	if len(parts) < 4 {
		logs.Error("Invalid path pattern: %s to resolve into parts.", filePath)
		return nil, fmt.Errorf("invalid path pattern: %s to resolve into parts", filePath)
	}
	username := parts[0]
	repoName := parts[2]
	filePathInRepo := strings.Join(parts[3:], "/")
	logs.Debug("Resolve raw file path with username: %s, repo name: %s with path in repo: %s", username, repoName, filePathInRepo)
	content, err := CurrentDevOps().GetRepoFile(username, repoName, branch, filePathInRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo file by username: %s, repo name: %s, branch: %s, file path in repo: %s", username, repoName, branch, filePathInRepo)
	}
	return bytes.NewBuffer(content), nil
}
