package service

import (
	"errors"
	"fmt"
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

const (
	k8sAPIversion1 = "v1"
	adminUserID    = 1
	adminUserName  = "admin"
	projectPrivate = 0
)

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

func GetPaginatedProjectsByUser(query model.Project, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedProjects, error) {
	return dao.GetPaginatedProjectsByUser(query, userID, pageIndex, pageSize, orderField, orderAsc)
}

func GetProjectsByMember(query model.Project, userID int64) ([]*model.Project, error) {
	return dao.GetProjectsByMember(query, userID)
}

func DeleteProject(projectID int64) (bool, error) {
	project := model.Project{ID: projectID, Deleted: 1}
	_, err := dao.UpdateProject(project, "deleted")
	if err != nil {
		return false, err
	}
	return true, nil
}

func NamespaceExists(projectName string) (bool, error) {
	cli, err := K8sCliFactory("", kubeMasterURL(), k8sAPIversion1)
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return false, err
	}

	n := apiSet.Namespaces()
	var listOpt modelK8s.ListOptions
	namespaceList, err := n.List(listOpt)
	if err != nil {
		logs.Error("Failed to check namespace list in cluster", projectName)
		return false, err
	}

	for _, namespace := range (*namespaceList).Items {
		if projectName == namespace.Name {
			logs.Info("Namespace existing %+v", namespace)
			return true, nil
		}
	}
	return false, nil
}

func CreateNamespace(projectName string) (bool, error) {
	projectExists, err := NamespaceExists(projectName)
	if err != nil {
		return false, err
	}
	if projectExists {
		logs.Info("Project library already exists in cluster.")
		return true, nil
	}

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

func SyncNamespaceByOwnerID(userID int64) error {
	query := model.Project{OwnerID: int(userID)}
	projects, err := GetProjectsByUser(query, userID)
	if err != nil {
		return fmt.Errorf("Failed to get default projects: %+v", err)
	}

	for _, project := range projects {
		_, err = CreateNamespace((*project).Name)
		if err != nil {
			return fmt.Errorf("Failed to create namespace: %s", (*project).Name)
		}
	}
	return nil
}

func SyncProjectsWithK8s() error {
	cli, err := K8sCliFactory("", kubeMasterURL(), k8sAPIversion1)
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		logs.Error("Failed to get K8s cli")
		return err
	}

	n := apiSet.Namespaces()
	var listOpt modelK8s.ListOptions
	namespaceList, err := n.List(listOpt)
	if err != nil {
		logs.Error("Failed to check namespace list in cluster")
		return err
	}

	for _, namespace := range (*namespaceList).Items {
		existing, err := ProjectExists(namespace.Name)
		if err != nil {
			logs.Error("Failed to check prject existing %s %+v", namespace.Name, err)
			continue
		}
		if existing {
			logs.Info("Project existing %s", namespace.Name)
		} else {
			//Add it to projects
			var reqProject model.Project
			reqProject.Name = namespace.Name
			reqProject.OwnerID = adminUserID
			reqProject.OwnerName = adminUserName
			reqProject.Public = projectPrivate

			isSuccess, err := CreateProject(reqProject)
			if err != nil {
				logs.Error("Failed to create project %s %+v", namespace.Name, err)
				// Still can work
				continue
			}
			if !isSuccess {
				logs.Error("Failed to create project %s", namespace.Name)
				// Still can work
				continue
			}
		}
	}
	return err
}

func DeleteNamespace(nameSpace string) (bool, error) {
	namespaceExists, err := NamespaceExists(nameSpace)
	if err != nil {
		return false, err
	}
	if !namespaceExists {
		logs.Info("Namespace %s not exists in cluster.", nameSpace)
		return true, nil
	}

	cli, err := K8sCliFactory("", kubeMasterURL(), "v1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return false, err
	}

	n := apiSet.Namespaces()
	err = n.Delete(nameSpace, nil)
	if err != nil {
		logs.Error("Failed to delete namespace %s", nameSpace)
		return false, err
	}
	return true, nil
}
