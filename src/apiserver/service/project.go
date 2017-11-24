package service

import (
	"errors"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"

	modelK8s "k8s.io/client-go/pkg/api/v1"

	"k8s.io/client-go/kubernetes"
)

var repoServeURL = utils.GetConfig("REPO_SERVE_URL")
var repoPath = utils.GetConfig("REPO_PATH")

func CreateProject(project model.Project) (bool, error) {
	projectID, err := dao.AddProject(project)
	if err != nil {
		return false, err
	}

	projectMember := model.ProjectMember{
		ProjectID: projectID,
		UserID:    int64(project.OwnerID),
		RoleID:    model.ProjectAdmin,
	}
	projectMemberID, err := dao.InsertOrUpdateProjectMember(projectMember)
	if err != nil {
		return false, errors.New("failed to create project member")
	}
	if projectID == 0 || projectMemberID == 0 {
		return false, errors.New("failed to create projectID memberID")
	}

	// Setup git repo for this project
	logs.Info("Initializing project %s repo", project.Name)
	_, err = InitRepo(repoServeURL(), repoPath())
	if err != nil {
		return false, errors.New("Initialize Project repo failed.")
	}

	subPath := project.Name
	if subPath != "" {
		os.MkdirAll(filepath.Join(repoPath(), subPath), 0755)
		if err != nil {
			return false, errors.New("Initialize Project path failed.")
		}
	}
	return true, nil
}

func GetProject(project model.Project, selectedFields ...string) (*model.Project, error) {
	p, err := dao.GetProject(project, selectedFields...)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func ProjectExists(projectName string) (bool, error) {
	query := model.Project{Name: projectName}
	project, err := dao.GetProject(query, "name")
	if err != nil {
		return false, err
	}
	return (project != nil && project.ID != 0), nil
}

func ProjectExistsByID(projectID int64) (bool, error) {
	query := model.Project{ID: projectID, Deleted: 0}
	project, err := dao.GetProject(query, "id", "deleted")
	if err != nil {
		return false, err
	}
	return (project != nil && project.Name != ""), nil
}

func UpdateProject(project model.Project, fieldNames ...string) (bool, error) {
	if project.ID == 0 {
		return false, errors.New("no Project ID provided")
	}
	_, err := dao.UpdateProject(project, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetProjectsByUser(query model.Project, userID int64) ([]*model.Project, error) {
	return dao.GetProjectsByUser(query, userID)
}

func DeleteProject(projectID int64) (bool, error) {
	project := model.Project{ID: projectID, Deleted: 1}
	_, err := dao.UpdateProject(project, "deleted")
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateNamespace(projectName string) (bool, error) {
	cli, err := K8sCliFactory("", kubeMasterURL(), "v1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return false, err
	}

	n := apiSet.Namespaces()
	var namespace modelK8s.Namespace
	namespace.ObjectMeta.Name = projectName
	_, err = n.Create(&namespace)
	if err != nil {
		logs.Error("Failed to creat namespace", projectName)
		return false, err
	}
	logs.Info(namespace)
	return true, nil
}
