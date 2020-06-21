package service

import (
	"errors"
	"fmt"
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
